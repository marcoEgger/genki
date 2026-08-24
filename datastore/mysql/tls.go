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

// ensureTLS enables TLS on the DSN unless it already specifies a tls parameter
// or DisableTLS was set. Uses a registered tls.Config with MinVersion TLS 1.2.
func ensureTLS(dsn string, opts *Options) (string, error) {
	if opts.DisableTLS || hasTLSParam(dsn) {
		return dsn, nil
	}

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

	return appendDSNParam(dsn, "tls", tlsConfigName), nil
}

func hasTLSParam(dsn string) bool {
	idx := strings.Index(dsn, "?")
	if idx < 0 || idx == len(dsn)-1 {
		return false
	}
	values, err := url.ParseQuery(dsn[idx+1:])
	if err != nil {
		return strings.Contains(strings.ToLower(dsn[idx+1:]), "tls=")
	}
	for key := range values {
		if strings.EqualFold(key, "tls") {
			return true
		}
	}
	return false
}

func appendDSNParam(dsn, key, value string) string {
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	return dsn + sep + key + "=" + url.QueryEscape(value)
}
