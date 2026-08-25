package mysql

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEnsureTLS_AddsTLSWhenMissing(t *testing.T) {
	dsn, err := ensureTLS("user:pass@tcp(db:3306)/app?parseTime=true", &Options{})
	if err != nil {
		t.Fatalf("ensureTLS: %v", err)
	}
	if !strings.Contains(dsn, "tls="+tlsConfigName) {
		t.Fatalf("expected tls=%s in DSN, got %q", tlsConfigName, dsn)
	}
}

func TestEnsureTLS_RespectsExistingTLSParam(t *testing.T) {
	in := "user:pass@tcp(db:3306)/app?parseTime=true&tls=false"
	dsn, err := ensureTLS(in, &Options{})
	if err != nil {
		t.Fatalf("ensureTLS: %v", err)
	}
	if dsn != in {
		t.Fatalf("expected DSN unchanged, got %q", dsn)
	}
}

func TestEnsureTLS_DisableTLS(t *testing.T) {
	in := "user:pass@tcp(db:3306)/app?parseTime=true"
	dsn, err := ensureTLS(in, &Options{DisableTLS: true})
	if err != nil {
		t.Fatalf("ensureTLS: %v", err)
	}
	if dsn != in {
		t.Fatalf("expected DSN unchanged, got %q", dsn)
	}
}

func TestEnsureMultiStatements(t *testing.T) {
	in := "user:pass@tcp(db:3306)/app?parseTime=true"
	dsn := ensureMultiStatements(in)
	if !strings.Contains(dsn, "multiStatements=true") {
		t.Fatalf("expected multiStatements=true in DSN, got %q", dsn)
	}

	already := "user:pass@tcp(db:3306)/app?parseTime=true&multiStatements=false"
	if got := ensureMultiStatements(already); got != already {
		t.Fatalf("expected DSN unchanged, got %q", got)
	}
}

func TestHasTLSParam(t *testing.T) {
	cases := []struct {
		dsn  string
		want bool
	}{
		{"user@tcp(localhost:3306)/db", false},
		{"user@tcp(localhost:3306)/db?parseTime=true", false},
		{"user@tcp(localhost:3306)/db?tls=true", true},
		{"user@tcp(localhost:3306)/db?parseTime=true&tls=false", true},
		{"user@tcp(localhost:3306)/db?TLS=skip-verify", true},
	}
	for _, tc := range cases {
		if got := hasTLSParam(tc.dsn); got != tc.want {
			t.Fatalf("hasTLSParam(%q)=%v, want %v", tc.dsn, got, tc.want)
		}
	}
}

func TestSetDSNParam_ReplacesExistingValue(t *testing.T) {
	in := "user:pass@tcp(db:3306)/app?parseTime=true&tls=true"
	got := setDSNParam(in, "tls", tlsConfigName)
	if strings.Contains(got, "tls=true") {
		t.Fatalf("expected tls=true to be replaced, got %q", got)
	}
	if !strings.Contains(got, "tls="+tlsConfigName) {
		t.Fatalf("expected tls=%s in DSN, got %q", tlsConfigName, got)
	}
}

func TestIsIPHost(t *testing.T) {
	if !isIPHost("34.65.29.144") {
		t.Fatal("expected IPv4 host to be detected as IP")
	}
	if !isIPHost("2001:db8::1") {
		t.Fatal("expected IPv6 host to be detected as IP")
	}
	if isIPHost("db.example.com") {
		t.Fatal("expected hostname not to be detected as IP")
	}
}

func TestBuildTLSConfig_WithoutCASkipsVerification(t *testing.T) {
	dsn := "user:pass@tcp(34.65.29.144:3306)/app?parseTime=true"
	cfg, err := buildTLSConfig(dsn, &Options{})
	if err != nil {
		t.Fatalf("buildTLSConfig: %v", err)
	}
	if !cfg.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify when no CA file is configured")
	}
	if cfg.VerifyConnection != nil {
		t.Fatal("expected no VerifyConnection when no CA file is configured")
	}
	if cfg.RootCAs != nil {
		t.Fatal("expected no RootCAs when no CA file is configured")
	}
}

func TestBuildTLSConfig_IPHostWithCAUsesVerifyConnection(t *testing.T) {
	caFile := writeTempCAFile(t)
	dsn := "user:pass@tcp(34.65.29.144:3306)/app?parseTime=true"
	cfg, err := buildTLSConfig(dsn, &Options{TLSCAFile: caFile})
	if err != nil {
		t.Fatalf("buildTLSConfig: %v", err)
	}
	if !cfg.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify for IP host so default IP SAN checks are skipped")
	}
	if cfg.VerifyConnection == nil {
		t.Fatal("expected VerifyConnection for IP host with CA")
	}
	if cfg.RootCAs == nil {
		t.Fatal("expected RootCAs when CA file is configured")
	}
	if cfg.ServerName != "" {
		t.Fatalf("expected empty ServerName, got %q", cfg.ServerName)
	}
}

func TestBuildTLSConfig_ServerNameOverride(t *testing.T) {
	caFile := writeTempCAFile(t)
	dsn := "user:pass@tcp(34.65.29.144:3306)/app?parseTime=true"
	cfg, err := buildTLSConfig(dsn, &Options{TLSCAFile: caFile, TLSServerName: "db.example.com"})
	if err != nil {
		t.Fatalf("buildTLSConfig: %v", err)
	}
	if cfg.ServerName != "db.example.com" {
		t.Fatalf("expected ServerName db.example.com, got %q", cfg.ServerName)
	}
	if cfg.VerifyConnection != nil {
		t.Fatal("expected VerifyConnection to be unset when ServerName is provided")
	}
	if cfg.InsecureSkipVerify {
		t.Fatal("expected hostname verification when ServerName is provided")
	}
}

func TestDSNHost(t *testing.T) {
	if got := dsnHost("user:pass@tcp(34.65.29.144:3306)/app"); got != "34.65.29.144" {
		t.Fatalf("expected 34.65.29.144, got %q", got)
	}
	if got := dsnHost("user:pass@tcp(db.example.com:3306)/app"); got != "db.example.com" {
		t.Fatalf("expected db.example.com, got %q", got)
	}
}

func writeTempCAFile(t *testing.T) string {
	t.Helper()
	cert, err := generateTestCA()
	if err != nil {
		t.Fatalf("generate test ca: %v", err)
	}
	path := t.TempDir() + "/server-ca.pem"
	if err := os.WriteFile(path, cert, 0o600); err != nil {
		t.Fatalf("write temp ca file: %v", err)
	}
	return path
}

func generateTestCA() ([]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), nil
}
