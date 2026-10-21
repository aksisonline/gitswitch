package skill

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestIsOutdatedAndUpdateAll(t *testing.T) {
	tmp := t.TempDir()
	skillDir := filepath.Join(tmp, ".agents", "skills", "gitswitch")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")

	// Write old content
	if err := os.WriteFile(skillPath, []byte("old content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Mock home dir
	t.Setenv("HOME", tmp)

	locs := InstalledLocations()
	if len(locs) != 1 || locs[0] != skillPath {
		t.Fatalf("InstalledLocations = %v, want [%s]", locs, skillPath)
	}

	outdated, upToDate := IsOutdated()
	if len(outdated) != 1 || outdated[0] != skillPath {
		t.Fatalf("outdated = %v, want [%s]", outdated, skillPath)
	}
	if len(upToDate) != 0 {
		t.Fatalf("upToDate = %v, want empty", upToDate)
	}

	updated, err := UpdateAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(updated) != 1 || updated[0] != skillPath {
		t.Fatalf("updated = %v, want [%s]", updated, skillPath)
	}

	got, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, Content) {
		t.Fatal("file content was not updated to Content")
	}

	outdatedAfter, upToDateAfter := IsOutdated()
	if len(outdatedAfter) != 0 {
		t.Fatalf("outdatedAfter = %v, want empty", outdatedAfter)
	}
	if len(upToDateAfter) != 1 || upToDateAfter[0] != skillPath {
		t.Fatalf("upToDateAfter = %v, want [%s]", upToDateAfter, skillPath)
	}
}
