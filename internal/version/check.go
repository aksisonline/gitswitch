package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	githubReleaseURL    = "https://api.github.com/repos/aksisonline/gitswitch/releases/latest"
	githubReleaseTagURL = "https://api.github.com/repos/aksisonline/gitswitch/releases/tags/%s"
	installScriptURL    = "https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh"
	cacheTTL            = 24 * time.Hour
	seenVersionFile     = "version-seen.txt"
)

// semverRe matches stable (v0.1.22) and pre-release (v0.2.0-beta.1, v0.2.0-canary.3) tags.
var semverRe = regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[a-zA-Z][a-zA-Z0-9]*(\.\d+)?)?$`)

type cache struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// CachedLatestVersion returns the latest version using a 24-hour disk cache.
// Returns "" on any error (graceful degradation).
func CachedLatestVersion(configDir string) string {
	cachePath := filepath.Join(configDir, "version-check.json")

	if data, err := os.ReadFile(cachePath); err == nil {
		var c cache
		if json.Unmarshal(data, &c) == nil && time.Since(c.CheckedAt) < cacheTTL {
			// Validate the cached tag still exists — guards against deleted/yanked releases.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if tagExists(ctx, c.LatestVersion) {
				return c.LatestVersion
			}
			// Tag gone: bust cache and fall through to re-fetch.
			_ = os.Remove(cachePath)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latest, err := fetchLatestVersion(ctx)
	if err != nil {
		return ""
	}

	_ = saveCache(cachePath, latest)
	return latest
}

// tagExists returns true if the given version tag has a release on GitHub.
// Returns true on network error to avoid false-positive cache busts.
func tagExists(ctx context.Context, version string) bool {
	url := fmt.Sprintf(githubReleaseTagURL, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return true
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitswitch-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ClearVersionCache removes the on-disk version cache so the next check always hits GitHub.
func ClearVersionCache(configDir string) error {
	return os.Remove(filepath.Join(configDir, "version-check.json"))
}

// CurrentVersionExists checks whether the given version has a published release on GitHub.
// Returns true on network error to avoid spurious forced upgrades.
func CurrentVersionExists(v string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return tagExists(ctx, v)
}

// FetchLatestVersionFresh always fetches from GitHub API, bypassing the cache.
// Used by the upgrade command so it always targets the true latest release.
func FetchLatestVersionFresh() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return fetchLatestVersion(ctx)
}

func fetchLatestVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubReleaseURL, nil)
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

	var result struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.TagName, nil
}

func saveCache(cachePath, version string) error {
	data, err := json.Marshal(cache{LatestVersion: version, CheckedAt: time.Now()})
	if err != nil {
		return err
	}
	// Atomic write: write to temp file, then rename.
	tmp := cachePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, cachePath)
}

// IsUpdateAvailable compares two semver strings.
// Returns false if current is not a valid semver (e.g. "dev" builds).
func IsUpdateAvailable(current, latest string) bool {
	if !semverRe.MatchString(current) || !semverRe.MatchString(latest) {
		return false
	}
	c := parseSemver(current)
	l := parseSemver(latest)
	if l[0] != c[0] {
		return l[0] > c[0]
	}
	if l[1] != c[1] {
		return l[1] > c[1]
	}
	return l[2] > c[2]
}

func parseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		fmt.Sscanf(parts[i], "%d", &result[i])
	}
	return result
}

// UpgradeCommand returns a configured exec.Cmd for upgrading to targetVersion.
// Returns an error if targetVersion fails semver validation (prevents shell injection).
func UpgradeCommand(targetVersion string) (*exec.Cmd, error) {
	if !semverRe.MatchString(targetVersion) {
		return nil, fmt.Errorf("invalid version format: %q", targetVersion)
	}
	script := fmt.Sprintf(`curl -fsSL %s | bash -s -- %s`, installScriptURL, targetVersion)
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

// RunUpgrade downloads and runs the install script to upgrade to the given version.
func RunUpgrade(targetVersion string) error {
	cmd, err := UpgradeCommand(targetVersion)
	if err != nil {
		return err
	}
	return cmd.Run()
}

// ShouldShowWhatsNew returns (show, releaseNotes).
// Shows the splash on first install, or when the current version has a higher
// major or minor than the last seen version. Silently skips on any error.
// Does NOT mark the version as seen — call MarkVersionSeen after the user dismisses.
func ShouldShowWhatsNew(configDir, currentVersion string) (bool, string) {
	if !semverRe.MatchString(currentVersion) {
		return false, ""
	}
	path := filepath.Join(configDir, seenVersionFile)
	data, err := os.ReadFile(path)

	firstTime := err != nil
	if !firstTime {
		lastSeen := strings.TrimSpace(string(data))
		if semverRe.MatchString(lastSeen) {
			cur := parseSemver(currentVersion)
			seen := parseSemver(lastSeen)
			// Recurring update: only show for minor/major bumps.
			if cur[0] <= seen[0] && cur[1] <= seen[1] {
				return false, ""
			}
		}
		// Invalid lastSeen falls through — treat as first time.
	}

	notes, fetchErr := FetchReleaseNotes(currentVersion)
	if fetchErr != nil || notes == "" {
		// Can't fetch notes — mark seen silently so we don't retry every launch.
		_ = os.WriteFile(path, []byte(currentVersion), 0600)
		return false, ""
	}
	return true, notes
}

// FetchReleaseNotes fetches the release body for the given version tag from GitHub.
func FetchReleaseNotes(version string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return fetchReleaseNotes(ctx, version)
}

// MarkVersionSeen writes the current version as "seen" so the splash won't show again.
func MarkVersionSeen(configDir, version string) {
	_ = os.WriteFile(filepath.Join(configDir, seenVersionFile), []byte(version), 0600)
}

func fetchReleaseNotes(ctx context.Context, version string) (string, error) {
	url := fmt.Sprintf(githubReleaseTagURL, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	var result struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Body, nil
}
