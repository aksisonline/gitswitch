package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// The first-run setup wizard: Welcome → Scan → Import → Accounts → Done.
// Every screen shares one skeleton — header with a step indicator, a short
// headline, the content, then the key hints — and is drawn on the same canvas
// as the main tabs, so it shrinks with the terminal the same way.

// wizardSteps is the number of dots in the step indicator.
const wizardSteps = 5

// scanRows is the number of sources the Scan step checks.
const scanRows = 3

// stepper renders the step indicator ("● ● ○ ○ ○"). Empty outside first-run
// setup — the "add account" screen is reused from the main list, where a
// progress indicator would be meaningless.
func (m Model) stepper(step int) string {
	if !m.onboarding {
		return ""
	}
	done := lipgloss.NewStyle().Foreground(colorGreen)
	here := lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
	todo := lipgloss.NewStyle().Foreground(colorDim)
	var dots []string
	for i := 0; i < wizardSteps; i++ {
		switch {
		case i < step:
			dots = append(dots, done.Render("●"))
		case i == step:
			dots = append(dots, here.Render("●"))
		default:
			dots = append(dots, todo.Render("○"))
		}
	}
	return strings.Join(dots, " ")
}

func identities(n int) string {
	if n == 1 {
		return "1 identity"
	}
	return fmt.Sprintf("%d identities", n)
}

// headline is the bold one-line statement at the top of a wizard step.
func headline(s string) string {
	return "  " + lipgloss.NewStyle().Foreground(colorWhite).Bold(true).Render(s)
}

// drawWizardWelcome — step 0: what this is, what's needed, start or skip.
func (m Model) drawWizardWelcome(c *canvas) {
	m.header(c, "Setup", m.stepper(0))
	c.gap()
	ok := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	missing := lipgloss.NewStyle().Foreground(colorYellow).Render("!")

	if m.arcadeMode {
		c.add("  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("★  PLAYER 1  —  READY TO SET UP?  ★"))
	} else {
		c.add(headline("Welcome — let's set up your git identities."))
		if c.d != dTight {
			c.add("  " + styleBrand.Render("So every repo commits as the right person."))
		}
		c.gap()
		row := func(found bool, name, yes, no string) {
			mark, note := ok, yes
			if !found {
				mark, note = missing, no
			}
			c.add("  " + mark + "  " + styleCurrentVal.Render(padTo(name, 4)) + styleBrand.Render(note))
		}
		row(m.hasGit, "git", "installed", "not found — run: gitswitch doctor")
		row(m.hasGH, "gh", "installed (GitHub CLI, used to log in)", "not found — needed to log in with GitHub")
	}
	c.gap()

	start, skip := "Start setup", "Skip for now"
	if m.arcadeMode {
		start, skip = "INSERT COIN  (Start Setup)", "SKIP SETUP"
	}
	play := lipgloss.NewStyle().Foreground(colorGreen).Render("▶")
	later := lipgloss.NewStyle().Foreground(colorDim).Render("→")
	m.drawChoices(c, m.wizardStep, []choiceView{{play, start}, {later, skip}})

	c.gap()
	c.footer(m, [][2]string{{"↑/↓", "move"}, {"enter", "select"}, {"q", "quit"}})
}

// drawWizardDetect — step 1: scan the machine for identities already in use.
// Rows tick off one at a time while the real scan runs in the background.
func (m Model) drawWizardDetect(c *canvas) {
	m.header(c, "Setup", m.stepper(1))
	c.gap()
	if m.arcadeMode {
		c.add("  " + styleTitle.Render("SCANNING FOR EXISTING IDENTITIES..."))
	} else {
		c.add(headline("Looking for identities you already use…"))
	}
	c.gap()

	sources := [scanRows][2]string{
		{"~/.gitconfig", "name, email, signing key"},
		{"gh CLI accounts", "GitHub logins on this machine"},
		{"~/.ssh keys", "private keys"},
	}
	nameW := 0
	for _, s := range sources {
		nameW = max(nameW, lipgloss.Width(s[0]))
	}
	step := m.scanStep
	if !m.scanning {
		step = scanRows // no animation running: show the finished state
	}
	for i, s := range sources {
		var mark string
		name := styleCurrentVal
		switch {
		case i < step:
			mark = lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
		case i == step:
			mark = lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render(m.spinner())
		default:
			mark = styleItemDim.Render("·")
			name = styleItemDim
		}
		c.add("  " + mark + "  " + name.Render(padTo(s[0], nameW)) + "  " + styleBrand.Render(s[1]))
	}

	c.gap()
	// Reserved result row — filled in when the scan finishes.
	result := ""
	if step >= scanRows && (m.scanReady || !m.scanning) {
		if n := len(m.detectedProfiles); n > 0 {
			result = "  " + lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("Found "+identities(n))
		} else {
			result = "  " + styleBrand.Render("Nothing found — you can add an account next.")
		}
	}
	c.add(result)
	c.gap()
	c.footer(m, [][2]string{{"enter", "continue"}, {"esc", "back"}})
}

// drawWizardImport — step 2: pick which detected identities to import.
func (m Model) drawWizardImport(c *canvas) {
	m.header(c, "Setup", m.stepper(2))
	c.gap()
	count := len(m.detectedProfiles)
	if count == 0 {
		c.add("  " + styleBrand.Render("No existing identities found."))
	} else {
		c.add(headline("Found " + identities(count) + " — choose what to import:"))
	}
	c.gap()

	for i, p := range m.detectedProfiles {
		on := i < len(m.importSelected) && m.importSelected[i]
		box := styleItemDim.Render("[ ]")
		if on {
			box = lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("[✓]")
		}
		lead, nick := "  ", styleItemInactive
		if i == m.wizardStep {
			lead = lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("❯") + " "
			nick = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
		}
		who := p.UserName + "  ·  " + p.Email
		if c.d == dTight {
			c.mark(zImport, i, lead+box+" "+nick.Render(truncate(p.Nickname, 16))+"  "+styleBrand.Render(p.Email))
			continue
		}
		c.mark(zImport, i, lead+box+" "+nick.Render(p.Nickname)+"\n"+"      "+styleBrand.Render(who))
		if i < count-1 {
			c.gap()
		}
	}

	c.gap()
	c.add(buttonRow("Import selected", "Skip"))
	c.gap()
	c.footer(m, [][2]string{{"↑/↓", "move"}, {"space", "toggle"}, {"enter", "import"}, {"s", "skip"}, {"esc", "back"}})
}

// drawWizardAddMore — step 3: add accounts. Also reached with `a` from the
// main list, where it is simply the "add an account" menu.
func (m Model) drawWizardAddMore(c *canvas) {
	title := "Add Account"
	if m.onboarding {
		title = "Setup"
	}
	m.header(c, title, m.stepper(3))
	c.gap()

	if n := len(m.profiles); n > 0 {
		c.add("  " + lipgloss.NewStyle().Foreground(colorGreen).Render("✓") + "  " + styleCurrentVal.Render(identities(n)+" ready"))
	} else {
		c.add("  " + styleBrand.Render("No identities yet — add your first one."))
	}
	if m.arcadeMode {
		c.add("  " + styleTitle.Render("ADD ANOTHER ACCOUNT?"))
	} else if len(m.profiles) > 0 {
		c.add(headline("Add another account?"))
	} else {
		c.add(headline("How do you want to add it?"))
	}
	c.gap()

	login, manual, done := "Log in with GitHub", "Add manually", "Done"
	if m.onboarding {
		done = "Done, open gitswitch"
	}
	if m.arcadeMode {
		login, manual, done = strings.ToUpper(login), strings.ToUpper(manual), strings.ToUpper(done)
	}
	m.drawChoices(c, m.wizardStep, []choiceView{
		{lipgloss.NewStyle().Foreground(colorGreen).Render("◉"), login},
		{lipgloss.NewStyle().Foreground(colorPurple).Render("✎"), manual},
		{lipgloss.NewStyle().Foreground(colorYellow).Render("→"), done},
	})

	c.add(m.statusLine()) // reserved: result of an add (or its error)
	c.footer(m, [][2]string{{"↑/↓", "move"}, {"enter", "select"}, {"esc", "back"}})
}

// drawWizardDone — step 4: setup complete, with what to try next.
func (m Model) drawWizardDone(c *canvas) {
	m.header(c, "Setup complete", m.stepper(4))
	c.gap()
	check := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	n := len(m.profiles)
	if m.arcadeMode {
		players := "PLAYERS"
		if n == 1 {
			players = "PLAYER"
		}
		c.add("  " + check + "  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("+500 EXP  ·  %d %s LOADED", n, players)))
	} else {
		c.add("  " + check + "  " + lipgloss.NewStyle().Foreground(colorWhite).Bold(true).Render("You're set — "+identities(n)+" ready."))
	}

	if cur := m.active; cur != nil {
		label := "Current"
		if m.arcadeMode {
			label = "PLAYER 1"
		}
		c.add("     " + styleCurrent.Render(label+"  ") + styleCheckmark.Render(cur.Nickname) + styleBrand.Render("  ·  "+cur.Email))
	} else {
		c.add("     " + styleBrand.Render("None is active yet — pick one on the next screen."))
	}

	if c.d != dTight {
		c.gap()
		c.add("  " + styleItemDim.Render("Try next"))
		key := styleFooterKey.Render
		for _, tip := range [][2]string{
			{"enter", "switch to the highlighted account"},
			{"2", "Utilities: shell integration, commit guard"},
			{"a", "add another account any time"},
		} {
			c.add("  " + key(padTo(tip[0], 6)) + styleBrand.Render(tip[1]))
		}
	}
	c.gap()

	open := "Open gitswitch"
	if m.arcadeMode {
		open = "OPEN GITSWITCH"
	}
	m.drawChoices(c, 0, []choiceView{{lipgloss.NewStyle().Foreground(colorGreen).Render("▶"), open}})
	c.gap()
	c.footer(m, [][2]string{{"enter", "open"}})
}
