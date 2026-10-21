// Package skill embeds gitswitch's Agent Skills package (SKILL.md), living
// at this repo-root conventional path so it's discoverable both by Go's
// go:embed (which can't reach outside this directory) and by skills.sh's
// repo crawl (npx skills add aksisonline/gitswitch).
package skill

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed SKILL.md
var Content []byte

// StandardSkillBases returns the conventional base directories for agent skills:
// ~/.agents, ~/.claude, ~/.codex, and ~/.hermes.
func StandardSkillBases(home string) []string {
	names := []string{".agents", ".claude", ".codex", ".hermes"}
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = filepath.Join(home, n)
	}
	return out
}

// InstalledLocations returns the paths of existing SKILL.md files on the machine.
func InstalledLocations() []string {
	var found []string
	home, err := os.UserHomeDir()
	if err != nil {
		return found
	}

	for _, base := range StandardSkillBases(home) {
		p := filepath.Join(base, "skills", "gitswitch", "SKILL.md")
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}

	// Also check project-level if running inside a repo with .agents or .claude
	for _, base := range []string{".agents", ".claude"} {
		p := filepath.Join(base, "skills", "gitswitch", "SKILL.md")
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}

	return found
}

// IsOutdated reports whether any installed skill file differs from Content.
func IsOutdated() (outdated []string, upToDate []string) {
	for _, path := range InstalledLocations() {
		existing, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(existing, Content) {
			outdated = append(outdated, path)
		} else {
			upToDate = append(upToDate, path)
		}
	}
	return outdated, upToDate
}

// UpdateAll rewrites all installed SKILL.md files with Content.
func UpdateAll() ([]string, error) {
	outdated, _ := IsOutdated()
	var updated []string
	for _, path := range outdated {
		if err := os.WriteFile(path, Content, 0o644); err != nil {
			return updated, err
		}
		updated = append(updated, path)
	}
	return updated, nil
}
