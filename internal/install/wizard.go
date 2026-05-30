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
	"time"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// ── palette — mirrors main TUI normal/default theme ──────────────────────────

var (
	cPrimary  = lipgloss.Color("141") // purple  — borders, labels, keys
	cAccent   = lipgloss.Color("84")  // green   — checkmarks, YES, success
	cHighlight = lipgloss.Color("226") // yellow  — selected option
	cDim      = lipgloss.Color("241") // muted   — secondary text
	cWhite    = lipgloss.Color("255") // bright  — body text
	cRed      = lipgloss.Color("196") // error / NO
	cBgHover  = lipgloss.Color("237") // selected option background
)

// ── Options / Config ─────────────────────────────────────────────────────────

// Options carries the resolved choices after a wizard run.
type Options struct {
	Shell        shell.Shell
	Framework    shell.Framework
	InstallShell bool
	InstallHTTPS bool
}

// Config controls wizard behaviour.
type Config struct {
	ShellOverride string // non-empty → non-interactive
	Yes           bool   // --yes → accept all defaults
	HTTPSDefault  bool   // default for HTTPS step
}

// ── Entry point ──────────────────────────────────────────────────────────────

// Run runs the wizard and returns the user's choices.
func Run(cfg Config, _ io.Writer) (Options, error) {
	sh, fw := resolveShell(cfg.ShellOverride)
	opts := Options{Shell: sh, Framework: fw, InstallShell: true, InstallHTTPS: cfg.HTTPSDefault}

	interactive := isatty.IsTerminal(os.Stdin.Fd()) && !cfg.Yes && cfg.ShellOverride == ""
	if !interactive {
		return opts, nil
	}

	m := newModel(sh, fw, cfg.HTTPSDefault)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return opts, err
	}
	final, ok := result.(model)
	if !ok || final.cancelled {
		return Options{Shell: sh, Framework: fw, InstallShell: false, InstallHTTPS: false}, nil
	}
	return Options{
		Shell:        sh,
		Framework:    fw,
		InstallShell: final.wantShell,
		InstallHTTPS: final.wantHTTPS,
	}, nil
}

// PrintSummary prints the post-install result (plain text, outside TUI).
func PrintSummary(w io.Writer, shellResult string, shellDone, httpsDone bool, httpsErr error) {
	check := lipgloss.NewStyle().Foreground(cAccent).Render("✓")
	skip  := lipgloss.NewStyle().Foreground(cDim).Render("–")
	warn  := lipgloss.NewStyle().Foreground(cHighlight).Render("⚠")
	dim   := func(s string) string { return lipgloss.NewStyle().Foreground(cDim).Render(s) }

	fmt.Fprintln(w)
	fmt.Fprintln(w, lipgloss.NewStyle().Bold(true).Foreground(cWhite).Render("  Setup complete"))
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

	fmt.Fprintln(w)
	if shellDone || httpsDone {
		fmt.Fprintln(w, dim("  Reload your shell (or open a new terminal) to activate."))
		fmt.Fprintln(w)
		fmt.Fprintln(w, lipgloss.NewStyle().Foreground(cWhite).Render("  Next:"))
		key := lipgloss.NewStyle().Foreground(cPrimary).Bold(true)
		fmt.Fprintln(w, "  "+key.Render("gitswitch current")+"   see your active profile")
		fmt.Fprintln(w, "  "+key.Render("gitswitch pin")+"       pin a repo to always use a specific profile")
		fmt.Fprintln(w)
	}
}

// ── Bubble Tea model ─────────────────────────────────────────────────────────

type wizStep int

const (
	stepShell wizStep = iota
	stepHTTPS
	stepDone
)

type tickMsg time.Time

type model struct {
	width, height int
	step          wizStep
	cursor        int // 0 = YES, 1 = NO
	wantShell     bool
	wantHTTPS     bool
	cancelled     bool
	blink         bool

	sh           shell.Shell
	fw           shell.Framework
	httpsDefault bool
	alreadyShell bool
	alreadyHTTPS bool
	ghInstalled  bool
}

func newModel(sh shell.Shell, fw shell.Framework, httpsDefault bool) model {
	alreadyShell := shell.IsInstalled(shell.RCFile(sh))
	alreadyHTTPS := git.IsCredentialHelperInstalled()

	start := stepShell
	if alreadyShell && !alreadyHTTPS {
		start = stepHTTPS
	} else if alreadyShell && alreadyHTTPS {
		start = stepDone
	}

	return model{
		step:         start,
		cursor:       0,
		wantShell:    true,
		wantHTTPS:    httpsDefault,
		sh:           sh,
		fw:           fw,
		httpsDefault: httpsDefault,
		alreadyShell: alreadyShell,
		alreadyHTTPS: alreadyHTTPS,
		ghInstalled:  git.IsGHInstalled(),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(600*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd { return tickCmd() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tickMsg:
		m.blink = !m.blink
		return m, tickCmd()
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
		m.step = stepHTTPS
		m.cursor = 0
	case stepHTTPS:
		m.wantHTTPS = yes
		return m, tea.Quit
	case stepDone:
		return m, tea.Quit
	}
	return m, nil
}

// ── View ─────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.width == 0 {
		return ""
	}
	pw := m.width - 6
	if pw > 90 {
		pw = 90
	}
	if pw < 40 {
		pw = 40
	}

	padLeft := (m.width - pw - 4) / 2
	if padLeft < 0 {
		padLeft = 0
	}
	pad := strings.Repeat(" ", padLeft)

	var sections []string
	sections = append(sections, m.viewHeader(pw))

	switch m.step {
	case stepShell:
		sections = append(sections, m.viewShell(pw))
	case stepHTTPS:
		sections = append(sections, m.viewHTTPS(pw))
	case stepDone:
		sections = append(sections, m.viewAlreadyDone(pw))
	}

	sections = append(sections, m.viewFooter(pw))

	// vertical centering
	totalLines := 0
	for _, s := range sections {
		totalLines += strings.Count(s, "\n") + 1
	}
	topPad := (m.height - totalLines) / 3
	if topPad < 0 {
		topPad = 0
	}

	var b strings.Builder
	b.WriteString(strings.Repeat("\n", topPad))
	for i, s := range sections {
		lines := strings.Split(s, "\n")
		for _, l := range lines {
			b.WriteString(pad + l + "\n")
		}
		if i < len(sections)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// ── header ────────────────────────────────────────────────────────────────────

func (m model) viewHeader(pw int) string {
	brand := lipgloss.NewStyle().
		Bold(true).
		Foreground(cPrimary).
		Render("gitswitch")
	sub := lipgloss.NewStyle().Foreground(cDim).Render(" — setup")

	line := lipgloss.NewStyle().Foreground(cDim).Render(strings.Repeat("─", pw-2))

	return lipgloss.NewStyle().
		PaddingLeft(1).
		Render(brand+sub) + "\n" +
		lipgloss.NewStyle().PaddingLeft(1).Render(line)
}

// ── shell step ────────────────────────────────────────────────────────────────

func (m model) viewShell(pw int) string {
	inner := pw - 4

	stepLine := m.stepBadge(1, m.totalSteps()) + "  " +
		lipgloss.NewStyle().Bold(true).Foreground(cWhite).Render("Shell integration")

	shName := shellName(m.sh)
	fwName := frameworkName(m.fw)
	detectedVal := lipgloss.NewStyle().Foreground(cWhite).Bold(true).Render(shName)
	if fwName != "" {
		detectedVal += lipgloss.NewStyle().Foreground(cDim).Render("  ·  "+fwName)
	}
	detected := lipgloss.NewStyle().Foreground(cDim).Render("Detected  ") + detectedVal

	bullets := m.bullets([]string{
		"Active profile visible in your prompt at all times",
		"Nudge when you cd into a repo usually worked on as a different identity",
		"Usage learning — gitswitch builds per-repo identity history automatically",
	})

	ba := m.beforeAfter(inner,
		"user@machine ~/code $",
		"user@machine ~/code "+
			lipgloss.NewStyle().Foreground(cPrimary).Bold(true).Render("[work]")+" $",
		`On cd: "You usually work here as 'work' — switch? [y/N]"`,
	)

	prompt := m.choicePrompt("Install shell integration?")

	body := strings.Join([]string{stepLine, "", detected, "", bullets, "", ba, "", prompt}, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPrimary).
		Width(pw).
		PaddingLeft(1).PaddingRight(1).PaddingTop(1).PaddingBottom(1).
		Render(body)
}

// ── HTTPS step ────────────────────────────────────────────────────────────────

func (m model) viewHTTPS(pw int) string {
	inner := pw - 4
	stepN := 2
	if m.alreadyShell {
		stepN = 1
	}

	stepLine := m.stepBadge(stepN, m.totalSteps()) + "  " +
		lipgloss.NewStyle().Bold(true).Foreground(cWhite).Render("HTTPS credential routing")

	upgradeNote := ""
	if m.alreadyShell {
		upgradeNote = lipgloss.NewStyle().Foreground(cDim).
			Render("Shell integration already active — this is the one new step.") + "\n\n"
	}

	bullets := m.bullets([]string{
		"Registers gitswitch as git's HTTPS credential helper",
		"Routes the right token per repo — delegates to gh CLI, stores nothing",
		"Works alongside existing helpers (osxkeychain etc.)",
	})

	ghStatus := ""
	if m.ghInstalled {
		ghStatus = lipgloss.NewStyle().Foreground(cAccent).Render("✓") +
			lipgloss.NewStyle().Foreground(cDim).Render("  gh CLI found")
	} else {
		ghStatus = lipgloss.NewStyle().Foreground(cHighlight).Render("⚠") +
			lipgloss.NewStyle().Foreground(cDim).Render("  gh CLI not found — installs but stays inert until gh is set up")
	}

	ba := m.beforeAfter(inner,
		"git push  →  wrong account / password prompt",
		"git push  →  "+lipgloss.NewStyle().Foreground(cAccent).Bold(true).Render("right token routed silently"),
		"",
	)

	prompt := m.choicePrompt("Register HTTPS credential helper?")

	body := strings.Join([]string{stepLine, "", upgradeNote + bullets, "", ghStatus, "", ba, "", prompt}, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPrimary).
		Width(pw).
		PaddingLeft(1).PaddingRight(1).PaddingTop(1).PaddingBottom(1).
		Render(body)
}

// ── already done ──────────────────────────────────────────────────────────────

func (m model) viewAlreadyDone(pw int) string {
	msg := lipgloss.NewStyle().Foreground(cAccent).Bold(true).Render("✓  Already fully set up")
	sub := lipgloss.NewStyle().Foreground(cDim).
		Render("Nothing to do. Run `gitswitch current` to see your active profile.")
	body := lipgloss.JoinVertical(lipgloss.Left, msg, "", sub)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cPrimary).
		Width(pw).
		PaddingLeft(2).PaddingRight(2).PaddingTop(2).PaddingBottom(2).
		Render(body)
}

// ── components ────────────────────────────────────────────────────────────────

func (m model) stepBadge(n, total int) string {
	return lipgloss.NewStyle().
		Foreground(cPrimary).
		Bold(true).
		Render(fmt.Sprintf("%d / %d", n, total))
}

func (m model) totalSteps() int {
	if m.alreadyShell && !m.alreadyHTTPS {
		return 1
	}
	return 2
}

func (m model) bullets(items []string) string {
	dim := lipgloss.NewStyle().Foreground(cDim)
	white := lipgloss.NewStyle().Foreground(cWhite)
	dot := lipgloss.NewStyle().Foreground(cPrimary).Render("·")
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = dim.Render("  ") + dot + " " + white.Render(item)
	}
	return strings.Join(lines, "\n")
}

func (m model) beforeAfter(w int, before, after, note string) string {
	label := func(s string, c lipgloss.Color) string {
		return lipgloss.NewStyle().Foreground(c).Bold(true).Width(7).Render(s)
	}
	lines := label("before", cDim) + "  " +
		lipgloss.NewStyle().Foreground(cDim).Render(before) + "\n" +
		label("after", cAccent) + "  " + after
	if note != "" {
		lines += "\n\n" + lipgloss.NewStyle().Foreground(cDim).Italic(true).Render("        "+note)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cDim).
		Width(w-2).
		PaddingLeft(1).PaddingRight(1).
		Render(lines)
}

func (m model) choicePrompt(question string) string {
	q := lipgloss.NewStyle().Foreground(cWhite).Bold(true).Render(question)

	var yes, no string
	if m.cursor == 0 {
		yes = lipgloss.NewStyle().Background(cBgHover).Foreground(cHighlight).Bold(true).Padding(0, 2).Render("Yes")
		no = lipgloss.NewStyle().Foreground(cDim).Padding(0, 2).Render("No")
	} else {
		yes = lipgloss.NewStyle().Foreground(cDim).Padding(0, 2).Render("Yes")
		no = lipgloss.NewStyle().Background(cBgHover).Foreground(cRed).Bold(true).Padding(0, 2).Render("No")
	}

	cursor := "  "
	if m.blink {
		cursor = lipgloss.NewStyle().Foreground(cPrimary).Render("  ▌")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		q,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, yes, "  ", no, cursor),
	)
}

// ── footer ────────────────────────────────────────────────────────────────────

func (m model) viewFooter(pw int) string {
	key := lipgloss.NewStyle().Foreground(cHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(cDim)

	pairs := []string{
		key.Render("←/→") + dim.Render(" toggle"),
		key.Render("Enter") + dim.Render(" confirm"),
		key.Render("Y/N") + dim.Render(" quick answer"),
		key.Render("Q") + dim.Render(" quit"),
	}

	return lipgloss.NewStyle().
		Foreground(cDim).
		Width(pw).
		Render(strings.Join(pairs, "   "))
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
