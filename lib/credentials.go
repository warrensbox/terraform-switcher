package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// credentialsFile mirrors the subset of Terraform's credentials.tfrc.json
// (written by `terraform login`) that we need:
//
//	{"credentials": {"host.example.com": {"token": "..."}}}
type credentialsFile struct {
	Credentials map[string]struct {
		Token string `json:"token"`
	} `json:"credentials"`
}

// DefaultCredentialsPath returns the location of Terraform's credentials file:
// ~/.terraform.d/credentials.tfrc.json (%APPDATA%\terraform.d on Windows).
func DefaultCredentialsPath() (string, error) {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "terraform.d", "credentials.tfrc.json"), nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".terraform.d", "credentials.tfrc.json"), nil
}

// LookupCredentialToken returns the token stored for host (host or host:port)
// in the given credentials file.
func LookupCredentialToken(credentialsPath, host string) (string, error) {
	data, err := os.ReadFile(credentialsPath) // nolint:gosec // path is user-controlled by design
	if err != nil {
		return "", fmt.Errorf("could not read credentials file %q: %w", credentialsPath, err)
	}
	var creds credentialsFile
	if err = json.Unmarshal(data, &creds); err != nil {
		return "", fmt.Errorf("could not parse credentials file %q: %w", credentialsPath, err)
	}
	for key, entry := range creds.Credentials {
		if strings.EqualFold(key, host) && strings.TrimSpace(entry.Token) != "" {
			return strings.TrimSpace(entry.Token), nil
		}
	}
	return "", fmt.Errorf("no token found for host %q in credentials file %q", host, credentialsPath)
}

// ResolveMirrorTokens builds the host -> token map for the given mirror URLs.
//   - If useCredentialsFile is set, each URL's host is looked up in
//     credentialsPath (empty means the default path); a missing entry is an error
//     because the user explicitly asked for authentication.
//   - Otherwise the TF_MIRROR_TOKEN value (envToken), if any, is used for all hosts.
//
// With useCredentialsFile set, envToken still takes precedence when non-empty.
func ResolveMirrorTokens(envToken string, useCredentialsFile bool, credentialsPath string, mirrorURLs ...string) (map[string]string, error) {
	tokens := map[string]string{}
	envToken = strings.TrimSpace(envToken)

	for _, raw := range mirrorURLs {
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			continue
		}
		if envToken != "" {
			tokens[strings.ToLower(u.Host)] = envToken
			continue
		}
		if !useCredentialsFile {
			continue
		}
		if credentialsPath == "" {
			if credentialsPath, err = DefaultCredentialsPath(); err != nil {
				return nil, err
			}
		}
		token, err := LookupCredentialToken(credentialsPath, u.Host)
		if err != nil && u.Port() != "" {
			// Terraform keys may omit the port
			token, err = LookupCredentialToken(credentialsPath, u.Hostname())
		}
		if err != nil {
			return nil, err
		}
		tokens[strings.ToLower(u.Host)] = token
	}

	if useCredentialsFile && len(tokens) == 0 {
		return nil, errors.New("--auth requires a valid mirror URL")
	}
	return tokens, nil
}
