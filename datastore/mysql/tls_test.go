package mysql

import (
	"strings"
	"testing"
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

func TestBuildTLSConfig_IPHostUsesVerifyConnection(t *testing.T) {
	dsn := "user:pass@tcp(34.65.29.144:3306)/app?parseTime=true"
	cfg, err := buildTLSConfig(dsn, &Options{})
	if err != nil {
		t.Fatalf("buildTLSConfig: %v", err)
	}
	if cfg.VerifyConnection == nil {
		t.Fatal("expected VerifyConnection for IP host")
	}
	if cfg.ServerName != "" {
		t.Fatalf("expected empty ServerName, got %q", cfg.ServerName)
	}
}

func TestBuildTLSConfig_ServerNameOverride(t *testing.T) {
	dsn := "user:pass@tcp(34.65.29.144:3306)/app?parseTime=true"
	cfg, err := buildTLSConfig(dsn, &Options{TLSServerName: "db.example.com"})
	if err != nil {
		t.Fatalf("buildTLSConfig: %v", err)
	}
	if cfg.ServerName != "db.example.com" {
		t.Fatalf("expected ServerName db.example.com, got %q", cfg.ServerName)
	}
	if cfg.VerifyConnection != nil {
		t.Fatal("expected VerifyConnection to be unset when ServerName is provided")
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
