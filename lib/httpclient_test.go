package lib

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func hostOf(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u.Host
}

func TestMirrorAuthSentOnlyToMirrorHost(t *testing.T) {
	var got string
	mirror := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer mirror.Close()
	other := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer other.Close()

	SetMirrorAuth(map[string]string{hostOf(t, mirror.URL): "s3cret"})
	defer SetMirrorAuth(nil)

	resp, err := httpGet(mirror.URL + "/x")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got != "Bearer s3cret" {
		t.Errorf("mirror host: expected Bearer header, got %q", got)
	}

	got = ""
	resp, err = httpGet(other.URL + "/x")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got != "" {
		t.Errorf("non-mirror host: expected no Authorization header, got %q", got)
	}
}

func TestMirrorAuthDisabledWithoutToken(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer srv.Close()

	SetMirrorAuth(map[string]string{hostOf(t, srv.URL): ""})
	resp, err := httpGet(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got != "" {
		t.Errorf("expected no Authorization header, got %q", got)
	}
}

func TestMirrorAuthDroppedOnCrossHostRedirect(t *testing.T) {
	var got string
	target := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
	}))
	defer target.Close()
	// Use a different host name for the redirect target so Go treats it as cross-host
	tu, err := url.Parse(target.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	redirectTo := strings.Replace(target.URL, tu.Hostname(), "localhost", 1)
	mirror := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}))
	defer mirror.Close()

	SetMirrorAuth(map[string]string{hostOf(t, mirror.URL): "s3cret"})
	defer SetMirrorAuth(nil)

	resp, err := httpGet(mirror.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got != "" {
		t.Errorf("Authorization leaked across redirect: %q", got)
	}
}

func TestStatusErrorHint(t *testing.T) {
	for _, code := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		if err := statusError("https://m/x", code); !strings.Contains(err.Error(), MirrorTokenEnvVar) {
			t.Errorf("HTTP %d: expected hint mentioning %s, got %v", code, MirrorTokenEnvVar, err)
		}
	}
	if err := statusError("https://m/x", http.StatusNotFound); strings.Contains(err.Error(), MirrorTokenEnvVar) {
		t.Errorf("404 should not hint at auth: %v", err)
	}
}
