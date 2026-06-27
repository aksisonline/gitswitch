package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// per_page=100 is the GitHub API maximum; avoids missing a canary that sits
// below many stable releases in the list without needing full pagination.
const githubReleasesURL = "https://api.github.com/repos/aksisonline/gitswitch/releases?per_page=100"

// IsCanaryBuild reports whether the given version string is a canary / pre-release build.
func IsCanaryBuild(v string) bool {
	v = strings.ToLower(v)
	return strings.Contains(v, "canary") || strings.Contains(v, "beta") || strings.Contains(v, "rc")
}

// FetchLatestCanaryVersion fetches the most recent pre-release tag from GitHub.
// Returns an error if no pre-release exists.
func FetchLatestCanaryVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubReleasesURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitswitch-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var releases []struct {
		TagName    string `json:"tag_name"`
		Prerelease bool   `json:"prerelease"`
		Draft      bool   `json:"draft"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}

	for _, r := range releases {
		if r.Prerelease && !r.Draft {
			return r.TagName, nil
		}
	}
	return "", fmt.Errorf("no canary release found on GitHub — check back soon")
}
