// Package install provides the full-screen interactive gitswitch setup wizard.
//
// Launches a Bubble Tea TUI when called from a real terminal with no flags.
// Follows the same clean visual style as the main gitswitch TUI — rounded
// borders, purple/green palette, no arcade effects here.
//
// Non-interactive paths (piped stdin, --yes, --shell flag) bypass the TUI
// and apply defaults immediately.
package install

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// ── palette — mirrors main TUI normal/default theme ──────────────────────────

var (
	cPrimary   = lipgloss.Color("141") // purple  — borders, labels, keys
	cAccent    = lipgloss.Color("84")  // green   — checkmarks, YES, success
	cHighlight = lipgloss.Color("226") // yellow  — selected option
	cDim       = lipgloss.Color("241") // muted   — secondary text
	cWhite     = lipgloss.Color("255") // bright  — body text
	cRed       = lipgloss.Color("196") // error / NO
)

// ── Options / Config ─────────────────────────────────────────────────────────

// Options carries the resolved choices after a wizard run.
type Options struct {
	Shell            shell.Shell
	Framework        shell.Framework
	InstallShell     bool
	InstallHTTPS     bool
	InstallGHWrapper bool
	InstallGuard     bool
}

// Config controls wizard behaviour.
type Config struct {
	ShellOverride string // non-empty → non-interactive
	Yes           bool   // --yes → accept all defaults
	HTTPSDefault  bool   // default for HTTPS step
	GuardDefault  bool   // default for Pre-commit Guard step
}

// ── Entry point ───────────────────────────────────────────────────────────────

// Run runs the wizard and returns the user's choices.
func Run(cfg Config) (Options, error) {
	sh, fw := resolveShell(cfg.ShellOverride)
	opts := Options{
		Shell:            sh,
		Framework:        fw,
		InstallShell:     true,
		InstallHTTPS:     cfg.HTTPSDefault,
		InstallGHWrapper: true,
		InstallGuard:     cfg.GuardDefault,
	}

	stat, _ := os.Stdin.Stat()
	isTerminal := stat.Mode()&os.ModeCharDevice != 0
	interactive := isTerminal && !cfg.Yes && cfg.ShellOverride == ""
	if !interactive {
		return opts, nil
	}

	m := newModel(sh, fw, cfg.HTTPSDefault, cfg.GuardDefault)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return opts, err
	}
	final, ok := result.(model)
	if !ok || final.cancelled {
		return Options{Shell: sh, Framework: fw}, nil
	}
	return Options{
		Shell:            sh,
		Framework:        fw,
		InstallShell:     final.wantShell,
		InstallHTTPS:     final.wantHTTPS,
		InstallGHWrapper: final.wantGHWrapper,
		InstallGuard:     final.wantGuard,
	}, nil
}

// PrintSummary prints the post-install result (plain text, outside TUI).
func PrintSummary(w io.Writer, shellResult string, shellDone, httpsDone, ghWrapperDone bool, httpsErr error, guardDone ...bool) {
	check := lipgloss.NewStyle().Foreground(cAccent).Render("✓")
	skip := lipgloss.NewStyle().Foreground(cDim).Render("–")
	warn := lipgloss.NewStyle().Foreground(cHighlight).Render("⚠")
	dim := func(s string) string { return lipgloss.NewStyle().Foreground(cDim).Render(s) }

	fmt.Fprintln(w)
	// Cancelling the wizard (or declining every step) lands here too — don't
	// call that "complete".
	title := "  Setup complete"
	anyGuard := len(guardDone) > 0 && guardDone[0]
	if !shellDone && !httpsDone && !ghWrapperDone && !anyGuard && httpsErr == nil {
		title = "  Nothing was changed"
	}
	fmt.Fprintln(w, lipgloss.NewStyle().Bold(true).Foreground(cWhite).Render(title))
	fmt.Fprintln(w, dim("  ────────────────────────────────────────────────"))
	fmt.Fprintln(w)

	if shellDone {
		fmt.Fprintf(w, "  %s  Shell integration    %s\n", check, dim(shellResult))
	} else {
		fmt.Fprintf(w, "  %s  Shell integration    %s\n", skip, dim("skipped"))
	}

	switch {
	case httpsErr != nil:
		fmt.Fprintf(w, "  %s  HTTPS routing        %s\n", warn, dim(httpsErr.Error()))
	case httpsDone:
		fmt.Fprintf(w, "  %s  HTTPS routing        %s\n", check, dim("credential helper registered"))
	default:
		fmt.Fprintf(w, "  %s  HTTPS routing        %s\n", skip, dim("skipped"))
	}

	if ghWrapperDone {
		fmt.Fprintf(w, "  %s  Session Isolation    %s\n", check, dim("enabled"))
	} else {
		fmt.Fprintf(w, "  %s  Session Isolation    %s\n", skip, dim("skipped"))
	}

	guardEnabled := len(guardDone) > 0 && guardDone[0]
	if len(guardDone) > 0 {
		if guardEnabled {
			fmt.Fprintf(w, "  %s  Pre-commit guard     %s\n", check, dim("enabled (blocks wrong-author commits)"))
		} else {
			fmt.Fprintf(w, "  %s  Pre-commit guard     %s\n", skip, dim("skipped"))
		}
	}

	fmt.Fprintln(w)
	if shellDone || httpsDone || ghWrapperDone || guardEnabled {
		fmt.Fprintln(w, dim("  Reload your shell (or open a new terminal) to activate."))
		fmt.Fprintln(w)
		fmt.Fprintln(w, lipgloss.NewStyle().Foreground(cWhite).Render("  Next:"))
		key := lipgloss.NewStyle().Foreground(cPrimary).Bold(true)
		fmt.Fprintln(w, "  "+key.Render("gitswitch current")+"   see your active profile")
		fmt.Fprintln(w, "  "+key.Render("gitswitch pin")+"       pin a repo to always use a specific profile")
		fmt.Fprintln(w, "  "+key.Render("gitswitch guard")+"     check or manage the pre-commit identity guard")
		fmt.Fprintln(w)
	}
}

// ── Bubble Tea model ──────────────────────────────────────────────────────────

type wizStep int

const (
	stepShell wizStep = iota
	stepHTTPS
	stepGHWrapper
	stepGuard
	stepDone
)

type model struct {
	width, height int
	step          wizStep
	cursor        int // 0 = YES, 1 = NO
	wantShell     bool
	wantHTTPS     bool
	wantGHWrapper bool
	wantGuard     bool
	cancelled     bool

	sh               shell.Shell
	fw               shell.Framework
	httpsDefault     bool
	guardDefault     bool
	alreadyShell     bool
	alreadyHTTPS     bool
	alreadyGHWrapper bool
	alreadyGuard     bool
	ghInstalled      bool
}

func newModel(sh shell.Shell, fw shell.Framework, httpsDefault, guardDefault bool) model {
	m := model{
		wantShell:        true,
		wantHTTPS:        httpsDefault,
		wantGHWrapper:    true,
		wantGuard:        guardDefault,
		sh:               sh,
		fw:               fw,
		httpsDefault:     httpsDefault,
		guardDefault:     guardDefault,
		alreadyShell:     shell.IsInstalled(shell.RCFile(sh)),
		alreadyHTTPS:     git.IsCredentialHelperInstalled(),
		alreadyGHWrapper: shell.IsGHWrapperInstalled(shell.RCFile(sh)),
		alreadyGuard:     git.IsGuardHookInstalled(),
		ghInstalled:      git.IsGHInstalled(),
	}
	m.step = m.nextStep(stepShell)
	return m
}

// done reports whether a step's feature is already set up.
func (m model) done(s wizStep) bool {
	switch s {
	case stepShell:
		return m.alreadyShell
	case stepHTTPS:
		return m.alreadyHTTPS
	case stepGHWrapper:
		return m.alreadyGHWrapper
	case stepGuard:
		return m.alreadyGuard
	}
	return false
}

// nextStep returns the first step at or after from that still needs an
// answer — steps whose feature is already set up are never asked about.
func (m model) nextStep(from wizStep) wizStep {
	for s := from; s < stepDone; s++ {
		if !m.done(s) {
			return s
		}
	}
	return stepDone
}

// progress returns this step's position among the steps being asked.
func (m model) progress() (n, total int) {
	for s := stepShell; s < stepDone; s++ {
		if m.done(s) {
			continue
		}
		total++
		if s <= m.step {
			n++
		}
	}
	return n, total
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		case "left", "right", "h", "l", "tab":
			m.cursor = 1 - m.cursor
		case "up", "k":
			m.cursor = 0
		case "down", "j":
			m.cursor = 1
		case "y", "Y":
			m.cursor = 0
			return m.commit()
		case "n", "N":
			m.cursor = 1
			return m.commit()
		case "enter", " ":
			return m.commit()
		}
	}
	return m, nil
}

func (m model) commit() (tea.Model, tea.Cmd) {
	yes := m.cursor == 0
	switch m.step {
	case stepShell:
		m.wantShell = yes
	case stepHTTPS:
		m.wantHTTPS = yes
	case stepGHWrapper:
		m.wantGHWrapper = yes
	case stepGuard:
		m.wantGuard = yes
	case stepDone:
		return m, tea.Quit
	}
	m.cursor = 0
	if m.step = m.nextStep(m.step + 1); m.step == stepDone {
		return m, tea.Quit
	}
	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────────────

// Density: how much vertical room a step is drawn with. View tries roomy
// first and tightens until the step fits the window, so a 24-row terminal
// shows the whole step instead of scrolling its title off the top.
const (
	dRoomy   = iota // padded panel, boxed before/after example
	dCompact        // no padding, example inline
	dTight          // title, one-line bullets, question
)

// stepCopy is everything a step says.
type stepCopy struct {
	title    string
	detail   string   // styled line under the title (what was detected)
	bullets  []string // what turning this on does
	status   string   // styled prerequisite line (gh CLI found / missing)
	before   string
	after    string // may be styled
	note     string
	question string
}

func (m model) copy() stepCopy {
	dim := lipgloss.NewStyle().Foreground(cDim)
	good := lipgloss.NewStyle().Foreground(cAccent).Bold(true)
	gh := lipgloss.NewStyle().Foreground(cAccent).Render("✓") + dim.Render("  gh CLI found")
	if !m.ghInstalled {
		gh = lipgloss.NewStyle().Foreground(cHighlight).Render("!") + dim.Render("  gh CLI not found — stays inactive until gh is installed")
	}

	switch m.step {
	case stepShell:
		detected := lipgloss.NewStyle().Foreground(cWhite).Bold(true).Render(shellName(m.sh))
		if fw := frameworkName(m.fw); fw != "" {
			detected += dim.Render("  ·  " + fw)
		}
		return stepCopy{
			title:  "Shell integration",
			detail: dim.Render("Detected  ") + detected,
			bullets: []string{
				"Shows the active profile in your prompt",
				"Offers to switch when you cd into a repo you use a different identity in",
				"Learns which identity you use in each repo",
			},
			before:   "user@machine ~/code $",
			after:    "user@machine ~/code " + lipgloss.NewStyle().Foreground(cPrimary).Bold(true).Render("[work]") + " $",
			note:     `On cd: "You usually work here as 'work' — switch? [y/N]"`,
			question: "Install shell integration?",
		}
	case stepHTTPS:
		return stepCopy{
			title: "HTTPS credential routing",
			bullets: []string{
				"Makes gitswitch the helper git asks for HTTPS logins",
				"Hands git the right account's token per repo — via gh, stores nothing",
				"Works alongside existing helpers (osxkeychain and others)",
			},
			status:   gh,
			before:   "git push  →  wrong account, or a password prompt",
			after:    "git push  →  " + good.Render("right account, no prompt"),
			question: "Register the HTTPS credential helper?",
		}
	case stepGHWrapper:
		return stepCopy{
			title: "Session Isolation",
			bullets: []string{
				"Each repo gets its own git identity and gh account",
				"Required for 'gitswitch pin' to take effect",
				"Wraps the gh command so it follows the repo you are in",
			},
			status:   gh,
			before:   "gh pr view  →  always your one global gh account",
			after:    "gh pr view  →  " + good.Render("this repo's own gh account"),
			question: "Enable Session Isolation?",
		}
	case stepGuard:
		return stepCopy{
			title: "Pre-commit identity guard",
			bullets: []string{
				"Stops a commit made with the wrong name or email",
				"Checks repo pins, your usual identity per repo, and unknown emails",
				"Your existing hooks (Husky, lefthook, pre-commit) still run",
			},
			before:   "git commit  →  wrong email quietly lands in git log",
			after:    "git commit  →  " + good.Render("blocked, with the one-line fix"),
			note:     "Bypass once with: git commit --no-verify",
			question: "Enable the pre-commit identity guard?",
		}
	}
	return stepCopy{}
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}
	// Panel outer width, with a column of margin each side.
	outer := m.width - 2
	if outer > 92 {
		outer = 92
	}
	if outer < 24 {
		outer = 24
	}

	var lines []string
	for d := dRoomy; d <= dTight; d++ {
		lines = m.render(outer, d)
		if m.height <= 0 || len(lines) <= m.height {
			break
		}
	}
	if m.height > 0 && len(lines) > m.height {
		lines = lines[:m.height]
	}

	left := (m.width - outer) / 2
	if left < 0 {
		left = 0
	}
	top := 0
	if m.height > len(lines) {
		top = (m.height - len(lines)) / 3
	}
	pad := strings.Repeat(" ", left)
	var b strings.Builder
	b.WriteString(strings.Repeat("\n", top))
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(pad + clipLine(l, m.width-left))
	}
	return b.String()
}

// render draws the whole screen at density d as a list of terminal rows,
// none wider than outer.
func (m model) render(outer, d int) []string {
	inner := outer - 4 // border (2) + padding (2)
	var lines []string
	lines = append(lines, m.headerLines(outer)...)
	if d == dRoomy {
		lines = append(lines, "")
	}

	var body []string
	if m.step == stepDone {
		body = []string{
			lipgloss.NewStyle().Foreground(cAccent).Bold(true).Render("✓  Already fully set up"),
			"",
		}
		body = append(body, wrapLines("Nothing to do. Run 'gitswitch current' to see your active profile.", inner, lipgloss.NewStyle().Foreground(cDim))...)
	} else {
		body = m.stepLines(inner, d)
	}
	padV := 0
	if d == dRoomy {
		padV = 1
	}
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPrimary).
		Width(outer-2).
		Padding(padV, 1).
		Render(strings.Join(body, "\n"))
	lines = append(lines, strings.Split(panel, "\n")...)

	if d == dRoomy {
		lines = append(lines, "")
	}
	return append(lines, m.footerLines(outer)...)
}

// stepLines is the inside of the panel for the current step.
func (m model) stepLines(w, d int) []string {
	c := m.copy()
	white := lipgloss.NewStyle().Foreground(cWhite)
	dim := lipgloss.NewStyle().Foreground(cDim)
	gap := func(out []string) []string {
		if d == dRoomy {
			return append(out, "")
		}
		return out
	}

	n, total := m.progress()
	badge := lipgloss.NewStyle().Foreground(cPrimary).Bold(true).Render(fmt.Sprintf("%d / %d", n, total))
	out := []string{clipLine(badge+"  "+white.Bold(true).Render(c.title), w)}
	if c.detail != "" && d != dTight {
		out = gap(out)
		out = append(out, clipLine(c.detail, w))
	}

	out = gap(out)
	dot := lipgloss.NewStyle().Foreground(cPrimary).Render("·")
	for _, bl := range c.bullets {
		if d == dTight {
			out = append(out, clipLine(" "+dot+" "+white.Render(bl), w))
			continue
		}
		// Hanging indent: wrapped lines line up under the text, not the dot.
		for i, l := range wrapLines(bl, w-3, white) {
			lead := " " + dot + " "
			if i > 0 {
				lead = "   "
			}
			out = append(out, lead+l)
		}
	}

	if c.status != "" && d != dTight {
		out = gap(out)
		out = append(out, clipLine(" "+c.status, w))
	}

	if d != dTight {
		out = gap(out)
		label := func(s string, col lipgloss.Color) string {
			return lipgloss.NewStyle().Foreground(col).Bold(true).Width(8).Render(s)
		}
		example := []string{
			clipLine(label("before", cDim)+dim.Render(c.before), w-4),
			clipLine(label("after", cAccent)+c.after, w-4),
		}
		if d == dRoomy {
			if c.note != "" {
				example = append(example, clipLine(dim.Italic(true).Render("        "+c.note), w-4))
			}
			box := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(cDim).
				Width(w-2).
				Padding(0, 1).
				Render(strings.Join(example, "\n"))
			out = append(out, strings.Split(box, "\n")...)
		} else {
			for _, l := range example {
				out = append(out, " "+l)
			}
		}
	}

	if d != dTight {
		out = append(out, "")
	}
	// The buttons are never clipped: if they don't fit beside the question,
	// they get a line of their own.
	question, buttons := white.Bold(true).Render(c.question), m.buttons()
	if lipgloss.Width(question)+2+lipgloss.Width(buttons) <= w {
		return append(out, question+"  "+buttons)
	}
	return append(out, clipLine(question, w), buttons)
}

// buttons renders the Yes / No choice with the focused one highlighted.
func (m model) buttons() string {
	on := lipgloss.NewStyle().Background(cAccent).Foreground(lipgloss.Color("0")).Bold(true).Padding(0, 2)
	onNo := on.Background(cRed).Foreground(cWhite)
	off := lipgloss.NewStyle().Foreground(cDim).Padding(0, 2)
	if m.cursor == 0 {
		return on.Render("Yes") + " " + off.Render("No")
	}
	return off.Render("Yes") + " " + onNo.Render("No")
}

func (m model) headerLines(outer int) []string {
	brand := lipgloss.NewStyle().Bold(true).Foreground(cPrimary).Render("gitswitch")
	left := " " + brand + lipgloss.NewStyle().Foreground(cDim).Render(" — setup")

	// Step dots at the right: done, current, to do.
	right := ""
	if n, total := m.progress(); m.step != stepDone && total > 1 {
		var dots []string
		for i := 1; i <= total; i++ {
			switch {
			case i < n:
				dots = append(dots, lipgloss.NewStyle().Foreground(cAccent).Render("●"))
			case i == n:
				dots = append(dots, lipgloss.NewStyle().Foreground(cPrimary).Bold(true).Render("●"))
			default:
				dots = append(dots, lipgloss.NewStyle().Foreground(cDim).Render("○"))
			}
		}
		right = strings.Join(dots, " ") + " "
	}
	if gap := outer - lipgloss.Width(left) - lipgloss.Width(right); gap > 0 && right != "" {
		left += strings.Repeat(" ", gap) + right
	}
	rule := " " + lipgloss.NewStyle().Foreground(cDim).Render(strings.Repeat("─", outer-2))
	return []string{clipLine(left, outer), rule}
}

// footerLines lists the keys, wrapping onto a second line on narrow windows.
func (m model) footerLines(outer int) []string {
	key := lipgloss.NewStyle().Foreground(cHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(cDim)
	pairs := [][2]string{{"←/→", "choose"}, {"enter", "confirm"}, {"y/n", "answer"}, {"q", "quit"}}
	if m.step == stepDone {
		pairs = [][2]string{{"enter", "close"}}
	}
	var lines []string
	cur, curW := " ", 1
	for i, p := range pairs {
		item := key.Render(p[0]) + dim.Render(" "+p[1])
		itemW := lipgloss.Width(p[0]) + 1 + lipgloss.Width(p[1])
		sepW := 0
		if i > 0 {
			sepW = 3
		}
		if curW+sepW+itemW > outer && curW > 1 {
			lines = append(lines, cur)
			cur, curW, sepW = " ", 1, 0
		}
		cur += strings.Repeat(" ", sepW) + item
		curW += sepW + itemW
	}
	return append(lines, cur)
}

// wrapLines word-wraps plain text to w cells and styles each line.
func wrapLines(text string, w int, style lipgloss.Style) []string {
	if w < 8 {
		w = 8
	}
	var out []string
	for _, l := range strings.Split(ansi.Wordwrap(text, w, ""), "\n") {
		out = append(out, style.Render(l))
	}
	return out
}

// clipLine shortens a (possibly styled) line to w cells.
func clipLine(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= w {
		return s
	}
	return ansi.Truncate(s, w, "…")
}

// ── helpers ───────────────────────────────────────────────────────────────────

func resolveShell(override string) (shell.Shell, shell.Framework) {
	var sh shell.Shell
	switch override {
	case "zsh":
		sh = shell.ShellZsh
	case "bash":
		sh = shell.ShellBash
	case "fish":
		sh = shell.ShellFish
	default:
		sh = shell.DetectShell()
	}
	return sh, shell.DetectFramework()
}

func shellName(sh shell.Shell) string {
	switch sh {
	case shell.ShellZsh:
		return "zsh"
	case shell.ShellBash:
		return "bash"
	case shell.ShellFish:
		return "fish"
	default:
		return "bash"
	}
}

func frameworkName(fw shell.Framework) string {
	switch fw {
	case shell.FrameworkOMZ:
		return "oh-my-zsh"
	case shell.FrameworkP10k:
		return "Powerlevel10k"
	case shell.FrameworkStarship:
		return "Starship"
	default:
		return ""
	}
}
