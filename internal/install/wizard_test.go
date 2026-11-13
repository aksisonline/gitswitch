package install

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// TestEveryStepFitsTheTerminal renders every wizard step from a roomy window
// down to a small one and asserts nothing exceeds the terminal. Before this,
// every step was 26–29 rows tall and lost its title off the top of a standard
// 80×24 window.
func TestEveryStepFitsTheTerminal(t *testing.T) {
	var sheet strings.Builder
	for _, size := range [][2]int{{120, 40}, {80, 24}, {60, 18}, {44, 12}, {30, 10}} {
		for st := stepShell; st <= stepDone; st++ {
			for _, gh := range []bool{true, false} {
				m := model{width: size[0], height: size[1], step: st, sh: shell.ShellZsh, fw: shell.FrameworkOMZ, ghInstalled: gh}
				view := m.View()
				lines := strings.Split(view, "\n")
				if len(lines) > size[1] {
					t.Errorf("step %d at %v: %d rows > height %d", st, size, len(lines), size[1])
				}
				for i, l := range lines {
					if w := lipgloss.Width(l); w > size[0] {
						t.Errorf("step %d at %v line %d: width %d > %d: %q", st, size, i, w, size[0], ansi.Strip(l))
					}
				}
				// From the smallest supported size up, the step's title and both
				// answers must be on screen — that is the whole point of the step.
				if plain := ansi.Strip(view); st != stepDone && size[0] >= 44 && size[1] >= 12 {
					for _, want := range []string{m.copy().title, "Yes", "No"} {
						if !strings.Contains(plain, want) {
							t.Errorf("step %d at %v: %q is not on screen", st, size, want)
						}
					}
				}
				if gh {
					fmt.Fprintf(&sheet, "\n━━ step %d · %d×%d · %d rows ━━\n%s\n", st, size[0], size[1], len(lines), ansi.Strip(view))
				}
			}
		}
	}
	// GS_PREVIEW=<file> writes a contact sheet for eyeballing the layouts.
	if out := os.Getenv("GS_PREVIEW"); out != "" {
		if err := os.WriteFile(out, []byte(sheet.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// Steps whose feature is already set up are never asked about, and the step
// counter only counts the ones that are.
func TestWizardSkipsWhatIsAlreadySetUp(t *testing.T) {
	m := model{alreadyHTTPS: true, alreadyGuard: true}
	m.step = m.nextStep(stepShell)
	if m.step != stepShell {
		t.Fatalf("first step = %d, want shell", m.step)
	}
	if n, total := m.progress(); n != 1 || total != 2 {
		t.Errorf("progress = %d/%d, want 1/2", n, total)
	}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = next.(model)
	if m.step != stepGHWrapper {
		t.Errorf("after shell the wizard went to step %d, want Session Isolation (HTTPS is already set up)", m.step)
	}
	if n, total := m.progress(); n != 2 || total != 2 {
		t.Errorf("progress = %d/%d, want 2/2", n, total)
	}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if cmd == nil {
		t.Error("answering the last pending step should finish the wizard")
	}
}
