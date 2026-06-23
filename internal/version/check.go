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
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	githubReleaseURL  = "https://api.github.com/repos/aksisonline/gitswitch/releases/latest"
	githubReleasesURL = "https://api.github.com/repos/aksisonline/gitswitch/releases?per_page=20"
	installScriptURL  = "https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.sh"
	installScriptPS1  = "https://raw.githubusercontent.com/aksisonline/gitswitch/main/.github/install.ps1"
	cacheTTL          = 24 * time.Hour
)

// semverRe matches stable versions (v0.1.22) and beta versions (v0.2.0-beta.1).
var semverRe = regexp.MustCompile(`^v?\d+\.\d+\.\d+(-beta\.\d+)?$`)

type cache struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// IsBeta reports whether version is a pre-release beta build.
func IsBeta(version string) bool {
	return strings.Contains(version, "-beta.")
}

// CachedLatestVersion returns the best available update for the given current
// version using a 24-hour disk cache. Beta builds check for newer betas AND
// stable releases; stable builds check stable releases only.
func CachedLatestVersion(configDir, currentVersion string) string {
	cacheFile := "version-check.json"
	if IsBeta(currentVersion) {
		cacheFile = "version-check-beta.json"
	}
	cachePath := filepath.Join(configDir, cacheFile)

	if data, err := os.ReadFile(cachePath); err == nil {
		var c cache
		if json.Unmarshal(data, &c) == nil && time.Since(c.CheckedAt) < cacheTTL {
			return c.LatestVersion
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var latest string
	var err error
	if IsBeta(currentVersion) {
		latest, err = fetchLatestForBeta(ctx, currentVersion)
	} else {
		latest, err = fetchLatestStable(ctx)
	}
	if err != nil {
		return ""
	}

	_ = saveCache(cachePath, latest)
	return latest
}

// FetchLatestVersionFresh always fetches from GitHub API, bypassing the cache.
func FetchLatestVersionFresh() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return fetchLatestStable(ctx)
}

func fetchLatestStable(ctx context.Context) (string, error) {
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

// fetchLatestForBeta returns the best available version for a beta user:
// the highest beta of the same target release, or the stable release if it
// has shipped (stable always wins over any beta of the same base).
func fetchLatestForBeta(ctx context.Context, currentVersion string) (string, error) {
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

	// Base version we're targeting (e.g. v0.2.0 from v0.2.0-beta.1)
	base := strings.SplitN(currentVersion, "-", 2)[0]

	var bestBeta, bestStable string
	for _, r := range releases {
		if r.Draft || !semverRe.MatchString(r.TagName) {
			continue
		}
		if !r.Prerelease {
			// Stable release for our base version → promote immediately
			if strings.HasPrefix(r.TagName, base) || compareVersions(r.TagName, base) >= 0 {
				if bestStable == "" || compareVersions(r.TagName, bestStable) > 0 {
					bestStable = r.TagName
				}
			}
		} else if strings.HasPrefix(r.TagName, base+"-beta.") {
			if bestBeta == "" || compareVersions(r.TagName, bestBeta) > 0 {
				bestBeta = r.TagName
			}
		}
	}

	// Stable beats beta — if the stable is out, recommend it
	if bestStable != "" {
		return bestStable, nil
	}
	if bestBeta != "" {
		return bestBeta, nil
	}
	return currentVersion, nil
}

func saveCache(cachePath, version string) error {
	data, err := json.Marshal(cache{LatestVersion: version, CheckedAt: time.Now()})
	if err != nil {
		return err
	}
	tmp := cachePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, cachePath)
}

// IsUpdateAvailable compares two semver strings, including beta versions.
// Returns false for dev/unversioned builds.
func IsUpdateAvailable(current, latest string) bool {
	if !semverRe.MatchString(current) || !semverRe.MatchString(latest) {
		return false
	}
	return compareVersions(latest, current) > 0
}

// compareVersions returns >0 if a > b, 0 if equal, <0 if a < b.
// Handles beta suffixes: stable > beta of same base; higher beta.N > lower.
func compareVersions(a, b string) int {
	aParts, aBeta := splitVersion(a)
	bParts, bBeta := splitVersion(b)

	for i := 0; i < 3; i++ {
		if aParts[i] != bParts[i] {
			return aParts[i] - bParts[i]
		}
	}
	// Same base version — stable beats beta; higher beta.N beats lower
	switch {
	case aBeta < 0 && bBeta < 0:
		return 0 // both stable
	case aBeta < 0:
		return 1 // a is stable, b is beta → a wins
	case bBeta < 0:
		return -1 // b is stable, a is beta → b wins
	default:
		return aBeta - bBeta // both beta, higher N wins
	}
}

// splitVersion parses "v0.2.0-beta.3" into ([0,2,0], 3).
// Returns betaN=-1 for stable releases.
func splitVersion(v string) (parts [3]int, betaN int) {
	v = strings.TrimPrefix(v, "v")
	betaN = -1
	if idx := strings.Index(v, "-beta."); idx >= 0 {
		n, err := strconv.Atoi(v[idx+6:])
		if err == nil {
			betaN = n
		}
		v = v[:idx]
	}
	segs := strings.SplitN(v, ".", 3)
	for i := 0; i < len(segs) && i < 3; i++ {
		fmt.Sscanf(segs[i], "%d", &parts[i])
	}
	return
}

// UpgradeCommand returns a configured exec.Cmd for upgrading to targetVersion.
func UpgradeCommand(targetVersion string) (*exec.Cmd, error) {
	if !semverRe.MatchString(targetVersion) {
		return nil, fmt.Errorf("invalid version format: %q", targetVersion)
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// PowerShell: set version env var then run install script
		script := fmt.Sprintf(
			`$env:GS_VERSION='%s'; Invoke-Expression (Invoke-WebRequest -Uri '%s' -UseBasicParsing).Content`,
			targetVersion, installScriptPS1,
		)
		cmd = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	} else {
		script := fmt.Sprintf(`curl -fsSL %s | bash -s -- %s`, installScriptURL, targetVersion)
		cmd = exec.Command("sh", "-c", script)
	}
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
