//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// MirrorTokenEnvVar is the environment variable holding the Bearer token used
// to authenticate against a private mirror. It is deliberately env-only (no
// CLI flag, no TOML key) so the secret does not leak via `ps`, shell history
// or committed config files.
const MirrorTokenEnvVar = "TF_MIRROR_TOKEN" //nolint:gosec // G101 false positive: this is the environment-variable name, not the token.

const httpTimeout = 5 * time.Minute

var (
	mirrorAuthMu     sync.RWMutex
	mirrorAuthTokens = map[string]string{} // lowercase host[:port] -> token
)

// httpClient has a timeout (unlike http.DefaultClient). Go's client strips the
// Authorization header on redirects to a different host, which is what we want
// when Artifactory redirects to a CDN / pre-signed storage URL.
var httpClient = &http.Client{Timeout: httpTimeout}

// SetMirrorAuth configures Bearer authentication from a host -> token map. A
// token is only ever sent to its own host (never to e.g. PGP key hosts). An
// empty map disables authentication.
func SetMirrorAuth(tokens map[string]string) {
	mirrorAuthMu.Lock()
	defer mirrorAuthMu.Unlock()
	mirrorAuthTokens = map[string]string{}
	for host, token := range tokens {
		if token = strings.TrimSpace(token); token != "" {
			mirrorAuthTokens[strings.ToLower(host)] = token
		}
	}
}

// newGetRequest builds a GET request, adding the Authorization header only
// when the target host has a configured token.
func newGetRequest(rawURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	mirrorAuthMu.RLock()
	defer mirrorAuthMu.RUnlock()
	if token, ok := mirrorAuthTokens[strings.ToLower(req.URL.Host)]; ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

// httpGet performs an (optionally authenticated) GET request.
func httpGet(rawURL string) (*http.Response, error) {
	req, err := newGetRequest(rawURL)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req) // nolint:gosec // URL is expected to be variable
}

// statusError describes a non-200 response, with a hint for auth failures.
func statusError(rawURL string, status int) error {
	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return fmt.Errorf("Unable to download from %s: HTTP %d (authentication may be required; set the %s environment variable or use --auth with a credentials.tfrc.json entry for the host)", rawURL, status, MirrorTokenEnvVar)
	default:
		return fmt.Errorf("Unable to download from %s: HTTP %d", rawURL, status)
	}
}
