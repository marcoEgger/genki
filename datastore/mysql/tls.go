package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
)

const tlsConfigName = "genki"

// ensureTLS enables TLS on the DSN unless disabled or explicitly set to tls=false.
// When TLSCAFile is configured, a registered tls.Config verifies the server against
// that CA. Without a CA file, TLS is still enabled but certificate verification is
// skipped (equivalent to tls=skip-verify).
func ensureTLS(dsn string, opts *Options) (string, error) {
	if opts.DisableTLS {
		return dsn, nil
	}

	if opts.TLSCAFile == "" && hasTLSParamValue(dsn, "false") {
		return dsn, nil
	}

	if opts.TLSCAFile != "" || !hasTLSParam(dsn) {
		tlsCfg, err := buildTLSConfig(dsn, opts)
		if err != nil {
			return "", err
		}

		if err := mysqldriver.RegisterTLSConfig(tlsConfigName, tlsCfg); err != nil {
			return "", fmt.Errorf("register mysql tls config: %w", err)
		}

		return setDSNParam(dsn, "tls", tlsConfigName), nil
	}

	return dsn, nil
}

func buildTLSConfig(dsn string, opts *Options) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Without a CA file, enable encrypted TLS but skip certificate verification
	// (equivalent to tls=skip-verify). Useful for local development.
	if opts.TLSCAFile == "" {
		tlsCfg.InsecureSkipVerify = true
		return tlsCfg, nil
	}

	pem, err := os.ReadFile(opts.TLSCAFile)
	if err != nil {
		return nil, fmt.Errorf("read mysql tls ca file: %w", err)
	}
	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
		return nil, fmt.Errorf("append mysql tls ca certificates from %s", opts.TLSCAFile)
	}
	tlsCfg.RootCAs = rootCAs

	if opts.TLSServerName != "" {
		tlsCfg.ServerName = opts.TLSServerName
		return tlsCfg, nil
	}

	if isIPHost(dsnHost(dsn)) {
		// The MySQL driver sets ServerName to the DSN host. For IP hosts that
		// causes hostname verification against IP SANs and fails before any
		// custom VerifyConnection runs. Skip default hostname checks and verify
		// the certificate chain against the configured CA instead.
		tlsCfg.InsecureSkipVerify = true
		tlsCfg.VerifyConnection = verifyPeerCertChain(rootCAs)
	}

	return tlsCfg, nil
}

// verifyPeerCertChain validates the server certificate chain without requiring
// the connection host to appear in the certificate SANs. This allows connecting
// via IP while still verifying the server certificate against a trusted CA.
func verifyPeerCertChain(rootCAs *x509.CertPool) func(tls.ConnectionState) error {
	return func(state tls.ConnectionState) error {
		if len(state.PeerCertificates) == 0 {
			return fmt.Errorf("mysql tls: server did not present a certificate")
		}

		verifyOpts := x509.VerifyOptions{Roots: rootCAs}
		if len(state.PeerCertificates) > 1 {
			verifyOpts.Intermediates = x509.NewCertPool()
			for _, cert := range state.PeerCertificates[1:] {
				verifyOpts.Intermediates.AddCert(cert)
			}
		}

		if _, err := state.PeerCertificates[0].Verify(verifyOpts); err != nil {
			return fmt.Errorf("mysql tls: verify server certificate: %w", err)
		}
		return nil
	}
}

func dsnHost(dsn string) string {
	cfg, err := mysqldriver.ParseDSN(dsn)
	if err != nil {
		return ""
	}
	host, _, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return cfg.Addr
	}
	return host
}

func isIPHost(host string) bool {
	return net.ParseIP(host) != nil
}

// ensureMultiStatements adds multiStatements=true when missing. Required by
// golang-migrate when executing migration files.
func ensureMultiStatements(dsn string) string {
	if hasDSNParam(dsn, "multiStatements") {
		return dsn
	}
	return appendDSNParam(dsn, "multiStatements", "true")
}

func hasTLSParam(dsn string) bool {
	return hasDSNParam(dsn, "tls")
}

func hasTLSParamValue(dsn, value string) bool {
	values := dsnQueryValues(dsn)
	for _, v := range values["tls"] {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func hasDSNParam(dsn, key string) bool {
	values := dsnQueryValues(dsn)
	for k := range values {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

func dsnQueryValues(dsn string) url.Values {
	idx := strings.Index(dsn, "?")
	if idx < 0 || idx == len(dsn)-1 {
		return nil
	}
	values, err := url.ParseQuery(dsn[idx+1:])
	if err != nil {
		return nil
	}
	return values
}

func setDSNParam(dsn, key, value string) string {
	base, query := splitDSNQuery(dsn)
	if query == "" {
		return base + "?" + key + "=" + url.QueryEscape(value)
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return appendDSNParam(dsn, key, value)
	}

	escapedKey := url.QueryEscape(key)
	values.Del(escapedKey)
	for existingKey := range values {
		if strings.EqualFold(existingKey, key) {
			values.Del(existingKey)
		}
	}
	values.Set(key, value)

	return base + "?" + values.Encode()
}

func splitDSNQuery(dsn string) (base, query string) {
	idx := strings.Index(dsn, "?")
	if idx < 0 {
		return dsn, ""
	}
	return dsn[:idx], dsn[idx+1:]
}

func appendDSNParam(dsn, key, value string) string {
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	return dsn + sep + key + "=" + url.QueryEscape(value)
}
