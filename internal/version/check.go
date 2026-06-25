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
	ReleaseNotes  string    `json:"release_notes"`
	CheckedAt     time.Time `json:"checked_at"`
}

// Release bundles a version tag with its GitHub release body.
type Release struct {
	Version string
	Notes   string
}

// IsBeta reports whether version is a pre-release beta build.
func IsBeta(version string) bool {
	return strings.Contains(version, "-beta.")
}

// CachedLatestVersion returns the best available update for the given current
// version using a 24-hour disk cache. Beta builds check for newer betas AND
// stable releases; stable builds check stable releases only.
func CachedLatestVersion(configDir, currentVersion string) string {
	r := CachedLatestRelease(configDir, currentVersion)
	return r.Version
}

// CachedLatestRelease returns the best available release (version + notes) for
// the given current version using a 24-hour disk cache.
func CachedLatestRelease(configDir, currentVersion string) Release {
	cacheFile := "version-check.json"
	if IsBeta(currentVersion) {
		cacheFile = "version-check-beta.json"
	}
	cachePath := filepath.Join(configDir, cacheFile)

	if data, err := os.ReadFile(cachePath); err == nil {
		var c cache
		if json.Unmarshal(data, &c) == nil && time.Since(c.CheckedAt) < cacheTTL {
			// Discard stale cache: if cached latest is behind current, the user
			// upgraded and the cache is from the previous install (e.g. 0.1.x left
			// version-check.json with v0.1.22; we're now running v0.2.0).
			if compareVersions(c.LatestVersion, currentVersion) >= 0 {
				return Release{Version: c.LatestVersion, Notes: c.ReleaseNotes}
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rel Release
	var err error
	if IsBeta(currentVersion) {
		rel, err = fetchLatestReleaseForBeta(ctx, currentVersion)
	} else {
		rel, err = fetchLatestReleaseStable(ctx)
	}
	if err != nil {
		return Release{}
	}

	_ = saveCacheRelease(cachePath, rel)
	return rel
}

// FetchLatestVersionFresh always fetches from GitHub API, bypassing the cache.
func FetchLatestVersionFresh() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return fetchLatestStable(ctx)
}

// FetchLatestVersionFreshFor fetches the best available version for the given
// current version, bypassing the cache. Beta builds use the full releases list
// so users can discover newer betas; stable builds use the latest-release endpoint.
func FetchLatestVersionFreshFor(currentVersion string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if IsBeta(currentVersion) {
		return fetchLatestForBeta(ctx, currentVersion)
	}
	return fetchLatestStable(ctx)
}

func fetchLatestStable(ctx context.Context) (string, error) {
	r, err := fetchLatestReleaseStable(ctx)
	return r.Version, err
}

func fetchLatestReleaseStable(ctx context.Context) (Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubReleaseURL, nil)
	if err != nil {
		return Release{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitswitch-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var result struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Release{}, err
	}
	return Release{Version: result.TagName, Notes: result.Body}, nil
}

// fetchLatestForBeta returns the best available version for a beta user:
// the highest beta of the same target release, or the stable release if it
// has shipped (stable always wins over any beta of the same base).
func fetchLatestForBeta(ctx context.Context, currentVersion string) (string, error) {
	r, err := fetchLatestReleaseForBeta(ctx, currentVersion)
	return r.Version, err
}

func fetchLatestReleaseForBeta(ctx context.Context, currentVersion string) (Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubReleasesURL, nil)
	if err != nil {
		return Release{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gitswitch-updater")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var releases []struct {
		TagName    string `json:"tag_name"`
		Prerelease bool   `json:"prerelease"`
		Draft      bool   `json:"draft"`
		Body       string `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return Release{}, err
	}

	// Base version we're targeting (e.g. v0.2.0 from v0.2.0-beta.1)
	base := strings.SplitN(currentVersion, "-", 2)[0]

	var bestBeta, bestStable Release
	for _, r := range releases {
		if r.Draft || !semverRe.MatchString(r.TagName) {
			continue
		}
		if !r.Prerelease {
			if strings.HasPrefix(r.TagName, base) || compareVersions(r.TagName, base) >= 0 {
				if bestStable.Version == "" || compareVersions(r.TagName, bestStable.Version) > 0 {
					bestStable = Release{Version: r.TagName, Notes: r.Body}
				}
			}
		} else if strings.HasPrefix(r.TagName, base+"-beta.") {
			if bestBeta.Version == "" || compareVersions(r.TagName, bestBeta.Version) > 0 {
				bestBeta = Release{Version: r.TagName, Notes: r.Body}
			}
		}
	}

	// Stable beats beta — if the stable is out, recommend it
	if bestStable.Version != "" {
		return bestStable, nil
	}
	if bestBeta.Version != "" {
		return bestBeta, nil
	}
	return Release{Version: currentVersion}, nil
}

func saveCache(cachePath, version string) error {
	return saveCacheRelease(cachePath, Release{Version: version})
}

func saveCacheRelease(cachePath string, rel Release) error {
	data, err := json.Marshal(cache{LatestVersion: rel.Version, ReleaseNotes: rel.Notes, CheckedAt: time.Now()})
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
