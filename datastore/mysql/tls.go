package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
)

const tlsConfigName = "genki"

// ensureTLS enables TLS on the DSN unless disabled or explicitly set to tls=false.
// When TLSCAFile is configured, a registered tls.Config is always applied and
// overrides any existing tls= DSN parameter so the CA file is actually used.
func ensureTLS(dsn string, opts *Options) (string, error) {
	if opts.DisableTLS {
		return dsn, nil
	}

	if opts.TLSCAFile == "" && hasTLSParamValue(dsn, "false") {
		return dsn, nil
	}

	if opts.TLSCAFile != "" || !hasTLSParam(dsn) {
		tlsCfg := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		if opts.TLSCAFile != "" {
			pem, err := os.ReadFile(opts.TLSCAFile)
			if err != nil {
				return "", fmt.Errorf("read mysql tls ca file: %w", err)
			}
			rootCAs := x509.NewCertPool()
			if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
				return "", fmt.Errorf("append mysql tls ca certificates from %s", opts.TLSCAFile)
			}
			tlsCfg.RootCAs = rootCAs
		}

		if err := mysqldriver.RegisterTLSConfig(tlsConfigName, tlsCfg); err != nil {
			return "", fmt.Errorf("register mysql tls config: %w", err)
		}

		return setDSNParam(dsn, "tls", tlsConfigName), nil
	}

	return dsn, nil
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
