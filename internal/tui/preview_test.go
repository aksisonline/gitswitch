package tui

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// TestPreviewScreens writes a plain-text contact sheet of every screen at
// several terminal sizes, so layouts can be reviewed without clicking through
// the app (the setup wizard in particular only ever shows on a fresh machine):
//
//	GS_PREVIEW=/tmp/screens.txt go test ./internal/tui -run TestPreviewScreens
//
// Set GS_PREVIEW_ARCADE=1 for arcade mode. Skipped unless GS_PREVIEW is set.
func TestPreviewScreens(t *testing.T) {
	out := os.Getenv("GS_PREVIEW")
	if out == "" {
		t.Skip("set GS_PREVIEW=<file> to write the contact sheet")
	}
	arcade := os.Getenv("GS_PREVIEW_ARCADE") != ""
	screens := []struct {
		name string
		st   State
		tab  int
	}{
		{"accounts", StateList, 0}, {"utilities", StateList, 1}, {"settings", StateList, 2},
		{"setup: welcome", StateWizardWelcome, 0}, {"setup: scan", StateWizardDetect, 0},
		{"setup: import", StateWizardImport, 0}, {"setup: add account", StateWizardAddMore, 0},
		{"setup: done", StateWizardDone, 0}, {"add profile", StateAdd, 0}, {"edit profile", StateEdit, 0},
		{"delete confirm", StateDeleteConfirm, 0}, {"update prompt", StateUpdatePrompt, 0},
		{"shell confirm", StateShellConfirm, 0}, {"cli tips", StateTips, 0}, {"what's new", StateWhatsNew, 0},
	}
	var b strings.Builder
	for _, size := range [][2]int{{100, 34}, {80, 24}, {60, 18}, {minTermW, minTermH}} {
		for _, sc := range screens {
			m := screenModel(t, sc.st, sc.tab, size[0], size[1], arcade)
			m.hasGit, m.hasGH = true, true
			m.shellEnabled, m.guardEnabled = true, true
			lines := strings.Split(ansi.Strip(m.View()), "\n")
			fmt.Fprintf(&b, "\n━━ %s · %d×%d · %d rows ━━\n", sc.name, size[0], size[1], len(lines))
			for _, l := range lines {
				b.WriteString("│" + l + "\n")
			}
		}
	}
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}
