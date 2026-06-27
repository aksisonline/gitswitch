package version

import (
	"context"
	"fmt"
	"time"
)

// IsCanaryBuild reports whether the given version is a pre-release (beta/canary/rc) build.
func IsCanaryBuild(v string) bool {
	return IsBeta(v)
}

// FetchLatestCanaryVersion returns the most recent pre-release tag from GitHub.
// Returns an error if no pre-release exists.
func FetchLatestCanaryVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Reuse the releases list endpoint already declared in check.go.
	// fetchLatestReleaseForBeta expects a base version to match against; we
	// pass a synthetic v0.2.0-beta.0 so any beta tag is a candidate.
	r, err := fetchLatestReleaseForBeta(ctx, "v0.2.0-beta.0")
	if err != nil {
		return "", err
	}
	if !IsBeta(r.Version) {
		return "", fmt.Errorf("no canary release found on GitHub — check back soon")
	}
	return r.Version, nil
}
