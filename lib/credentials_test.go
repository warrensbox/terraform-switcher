package lib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCreds(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "credentials.tfrc.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

const testCreds = `{"credentials": {
  "art.example.com": {"token": " abc123 "},
  "other.example.com:8443": {"token": "def"},
  "empty.example.com": {"token": ""}
}}`

func TestLookupCredentialToken(t *testing.T) {
	p := writeCreds(t, testCreds)

	if tok, err := LookupCredentialToken(p, "ART.example.com"); err != nil || tok != "abc123" {
		t.Errorf("expected abc123 (case-insensitive, trimmed), got %q, %v", tok, err)
	}
	if _, err := LookupCredentialToken(p, "missing.example.com"); err == nil {
		t.Error("expected error for missing host")
	}
	if _, err := LookupCredentialToken(p, "empty.example.com"); err == nil {
		t.Error("expected error for empty token")
	}
	if _, err := LookupCredentialToken(filepath.Join(t.TempDir(), "nope"), "x"); err == nil {
		t.Error("expected error for missing file")
	}
	if _, err := LookupCredentialToken(writeCreds(t, "not json"), "x"); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestResolveMirrorTokens(t *testing.T) {
	p := writeCreds(t, testCreds)

	// --auth: token from credentials file, keyed by mirror host
	got, err := ResolveMirrorTokens("", true, p, "https://art.example.com/artifactory/x/index.json", "https://art.example.com/artifactory/x")
	if err != nil || got["art.example.com"] != "abc123" {
		t.Errorf("unexpected result %v, %v", got, err)
	}

	// host:port key present in the file
	got, err = ResolveMirrorTokens("", true, p, "https://other.example.com:8443/x")
	if err != nil || got["other.example.com:8443"] != "def" {
		t.Errorf("unexpected result %v, %v", got, err)
	}

	// --auth with no matching entry is an error
	if _, err = ResolveMirrorTokens("", true, p, "https://missing.example.com/x"); err == nil || !strings.Contains(err.Error(), "missing.example.com") {
		t.Errorf("expected missing-host error, got %v", err)
	}

	// env token takes precedence over the credentials file
	got, err = ResolveMirrorTokens("envtok", true, p, "https://art.example.com/x")
	if err != nil || got["art.example.com"] != "envtok" {
		t.Errorf("unexpected result %v, %v", got, err)
	}

	// no --auth, no env token: nothing configured
	got, err = ResolveMirrorTokens("", false, p, "https://art.example.com/x")
	if err != nil || len(got) != 0 {
		t.Errorf("unexpected result %v, %v", got, err)
	}
}
