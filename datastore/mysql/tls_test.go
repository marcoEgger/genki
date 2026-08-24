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
