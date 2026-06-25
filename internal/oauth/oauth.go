package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	ghOAuth "github.com/cli/oauth"
)

const defaultClientID = "Ov23li75EL8QeU1iIfkR"

// Scopes requested during device flow.
var defaultScopes = []string{"repo", "read:user", "user:email", "gist", "workflow"}

// GHUser holds the GitHub user fields fetched after auth.
type GHUser struct {
	Login string
	Name  string
	Email string
}

// Login runs the GitHub device flow and returns the access token and user info.
// host is "github.com" for github.com or a GHES hostname.
// clientID overrides the baked-in default when non-empty.
func Login(host, clientID string) (token string, user GHUser, err error) {
	if clientID == "" {
		clientID = defaultClientID
	}
	if host == "" {
		host = "github.com"
	}

	ghHost, err := ghOAuth.NewGitHubHost("https://" + host)
	if err != nil {
		return "", GHUser{}, fmt.Errorf("invalid host %q: %w", host, err)
	}
	flow := &ghOAuth.Flow{
		Host:     ghHost,
		ClientID: clientID,
		Scopes:   defaultScopes,
		DisplayCode: func(code, verificationURL string) error {
			fmt.Println()
			fmt.Printf("  Open this URL in your browser:\n\n    %s\n\n", verificationURL)
			fmt.Printf("  Then enter the code:\n\n    %s\n\n", code)
			fmt.Println("  Waiting for authorization...")
			return nil
		},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return "", GHUser{}, fmt.Errorf("authorization failed: %w", err)
	}
	token = accessToken.Token

	user, err = fetchUser(token, host)
	if err != nil {
		return token, GHUser{}, fmt.Errorf("token obtained but could not fetch user info: %w", err)
	}
	return token, user, nil
}

func fetchUser(token, host string) (GHUser, error) {
	apiBase := apiBaseURL(host)

	userJSON, err := ghGet(apiBase+"/user", token)
	if err != nil {
		return GHUser{}, err
	}
	var u struct {
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(userJSON, &u); err != nil {
		return GHUser{}, err
	}

	email, err := primaryEmail(apiBase, token)
	if err != nil {
		// Non-fatal: fall back to noreply address.
		email = u.Login + "@users.noreply.github.com"
	}

	return GHUser{Login: u.Login, Name: u.Name, Email: email}, nil
}

func primaryEmail(apiBase, token string) (string, error) {
	data, err := ghGet(apiBase+"/user/emails", token)
	if err != nil {
		return "", err
	}
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(data, &emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", fmt.Errorf("no primary verified email found")
}

func ghGet(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
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
