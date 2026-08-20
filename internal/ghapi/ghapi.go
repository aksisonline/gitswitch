// Package ghapi is a minimal GitHub REST client used to fill in profile
// details (name, verified email) for accounts already authenticated via the
// gh CLI — it never authenticates anything itself, callers hand it a token
// obtained elsewhere (gh auth token).
package ghapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// httpClient bounds every GitHub API call — ghGet previously used
// http.DefaultClient with no timeout, so an unreachable/slow API could hang
// a caller indefinitely (the onboarding identity cross-reference, which
// must never stall the wizard).
var httpClient = &http.Client{Timeout: 3 * time.Second}

// FetchName best-effort resolves the GitHub profile "name" field for the
// account owning token on host. Returns "" on any failure — never an
// error, same contract as FetchVerifiedEmails.
func FetchName(token, host string) string {
	data, err := ghGet(apiBaseURL(host)+"/user", token)
	if err != nil {
		return ""
	}
	var u struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(data, &u) != nil {
		return ""
	}
	return u.Name
}

// FetchVerifiedEmails best-effort resolves every verified email address for
// the account owning token on host, primary email first (so callers that
// just want "the" email — like fetchUser, below — can take out[0]). Tries
// /user/emails first (all verified emails — needs the user:email scope,
// which a plain `gh auth login` token commonly lacks); falls back to
// /user's public profile email, which GitHub returns to any authenticated
// caller regardless of scope (null if the account has none set public).
// Returns nil on total failure — never an error, callers never branch on why.
func FetchVerifiedEmails(token, host string) []string {
	apiBase := apiBaseURL(host)

	if data, err := ghGet(apiBase+"/user/emails", token); err == nil {
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if json.Unmarshal(data, &emails) == nil {
			var out []string
			primaryIdx := -1
			for _, e := range emails {
				if !e.Verified || e.Email == "" {
					continue
				}
				if e.Primary {
					primaryIdx = len(out)
				}
				out = append(out, e.Email)
			}
			if primaryIdx > 0 {
				out[0], out[primaryIdx] = out[primaryIdx], out[0]
			}
			if len(out) > 0 {
				return out
			}
		}
	}

	if data, err := ghGet(apiBase+"/user", token); err == nil {
		var u struct {
			Email string `json:"email"`
		}
		if json.Unmarshal(data, &u) == nil && u.Email != "" {
			return []string{u.Email}
		}
	}

	return nil
}

func ghGet(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
}

// apiBaseURL returns the REST API base for a given GitHub host.
func apiBaseURL(host string) string {
	if host == "github.com" || host == "" {
		return "https://api.github.com"
	}
	return "https://" + host + "/api/v3"
}
