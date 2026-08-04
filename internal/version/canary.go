package version

import (
	"context"
	"fmt"
	"time"
)

// FetchLatestCanaryVersion returns the most recent pre-release tag from GitHub.
// Returns an error if no pre-release exists.
func FetchLatestCanaryVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := fetchLatestReleaseForBeta(ctx)
	if err != nil {
		return "", err
	}
	if !IsBeta(r.Version) {
		return "", fmt.Errorf("no canary release found on GitHub — check back soon")
	}
	return r.Version, nil
}
