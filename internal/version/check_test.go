package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsBeta(t *testing.T) {
	cases := map[string]bool{
		"v0.2.2-beta.1":  true,
		"v0.2.2":         false,
		"v0.2.2-beta.42": true,
	}
	for v, want := range cases {
		if got := IsBeta(v); got != want {
			t.Errorf("IsBeta(%q) = %v, want %v", v, got, want)
		}
	}
}

// Canary releases only ever bump PATCH by construction, so ShouldShowWhatsNew
// treats any distinct canary tag as notable rather than gating on minor/major
// like it does for stable releases. Same-tag re-launches must still stay
// quiet, and that path needs no network fetch, so it's safe to exercise here.
func TestShouldShowWhatsNew_CanarySameVersionStaysQuiet(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, seenVersionFile), []byte("v0.2.2-beta.1"), 0600); err != nil {
		t.Fatal(err)
	}
	show, notes := ShouldShowWhatsNew(dir, "v0.2.2-beta.1")
	if show || notes != "" {
		t.Errorf("expected no splash when canary version is unchanged, got show=%v notes=%q", show, notes)
	}
}
