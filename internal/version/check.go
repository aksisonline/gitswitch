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
	githubReleaseURL = "https://api.github.com/repos/aksisonline/gitswitch/releases/latest"
	installScriptURL = "https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh"
	cacheTTL         = 24 * time.Hour
)

var semverRe = regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

type cache struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// CachedLatestVersion returns the latest version using a 24-hour disk cache.
// Returns "" on any error (graceful degradation)
func CachedLatestVersion(configDir string) string {
	cachePath := filepath.Join(configDir, "version-check.json")

	if data, err := os.ReadFile(cachePath); err == nil {
		var c cache
		if json.Unmarshal(data, &c) == nil && time.Since(c.CheckedAt) < cacheTTL {
			return c.LatestVersion
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
