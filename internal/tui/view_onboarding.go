package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// viewWizardWelcome — new user step 0: welcome screen.
func (m Model) viewWizardWelcome(pw int) string {
	header := m.viewHeader("Setup")
	sep := "\n\n"

	var greeting, prereqs string
	var ctaBorderStyle lipgloss.Style
	if m.arcadeMode {
		greeting = sep + "  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("★  PLAYER 1  —  READY TO SET UP?  ★")
		prereqs = ""
		ctaBorderStyle = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
	} else {
		greeting = sep + "  First time here. Let's get you set up."
		prereqs = "\n\n  " + styleBrand.Render("You'll need:") +
			"\n  " + styleBrand.Render("  ·  git   (you have it — you're reading this)") +
			"\n  " + styleBrand.Render("  ·  gh    (GitHub CLI, optional but recommended)")
		ctaBorderStyle = lipgloss.NewStyle().Bold(true).Foreground(colorPurple)
	}

	divLine := "\n\n  " + divider(pw-4)

	var ctaLabel string
	if m.arcadeMode {
		ctaLabel = "▶  INSERT COIN  (Start Setup)"
	} else {
		ctaLabel = "▶  Start Setup"
	}

	cta := divLine +
		"\n\n  " + ctaBorderStyle.Render("╔══════════════════════════════════════╗") +
		"\n  " + ctaBorderStyle.Render("║") + "  " + lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render(fmt.Sprintf("%-38s", ctaLabel)) + ctaBorderStyle.Render("║") +
		"\n  " + ctaBorderStyle.Render("╚══════════════════════════════════════╝")

	var secondary string
	if m.arcadeMode {
		secondary = "\n\n  " + styleBrand.Render("[ SKIP SETUP ]")
	} else {
		secondary = "\n\n  " + styleBrand.Render("[ Skip for now ]")
	}

	body := header + greeting + prereqs + cta + secondary
	return stylePanelBorder(pw).Render(body)
}

// viewWizardDetect — new user step 1: scanning for existing configs.
func (m Model) viewWizardDetect(pw int) string {
	header := m.viewHeader("Setup")
	sep := "\n\n"

	var title string
	if m.arcadeMode {
		title = sep + "  " + styleTitle.Render("SCANNING FOR EXISTING IDENTITIES...")
	} else {
		title = sep + "  Scanning your machine for existing git identities…"
	}

	divLine := "\n\n  " + divider(pw-4)

	spin := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	checks := divLine +
		"\n\n  " + spin + "  " + styleCurrentVal.Render("~/.gitconfig") + "          " + styleBrand.Render("name, email, keys") +
		"\n  " + spin + "  " + styleCurrentVal.Render("gh CLI accounts") + "       " + styleBrand.Render("logged-in users") +
		"\n  " + spin + "  " + styleCurrentVal.Render("SSH keys in ~/.ssh/") + "   " + styleBrand.Render("private keys")

	hint := "\n\n  " + styleBrand.Render("press ") + styleFooterKey.Render("enter") + styleBrand.Render(" or click to continue")

	body := header + title + checks + hint
	return stylePanelBorder(pw).Render(body)
}

// viewWizardImport — new user step 2: show detected profiles, let user select.
func (m Model) viewWizardImport(pw int) string {
	header := m.viewHeader("Setup")
	sep := "\n\n"

	count := len(m.detectedProfiles)
	pluralSuffix := func() string {
		if count == 1 {
			return "y"
		}
		return "ies"
	}()

	var title string
	if m.arcadeMode {
		title = sep + "  Found " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("%d", count)) + " existing identit" + pluralSuffix + ":"
	} else {
		title = sep + "  Found " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("%d", count)) + " existing identit" + pluralSuffix + " to import:"
	}

	var list string
	for i, p := range m.detectedProfiles {
		checked := "✓"
		checkColor := colorGreen
		if i < len(m.importSelected) && !m.importSelected[i] {
			checked = " "
			checkColor = colorDim
		}
		checkStyle := lipgloss.NewStyle().Foreground(checkColor)
		ribbonColor := colorDim
		if i == m.wizardStep {
			ribbonColor = colorPurple
		}
		ribbon := lipgloss.NewStyle().Foreground(ribbonColor).Render("▐")

		line := "\n\n  " + ribbon + " " + checkStyle.Render("["+checked+"]") + " " + styleCurrentVal.Render(p.Nickname) +
			"\n       " + styleBrand.Render(p.UserName+"  ·  "+p.Email)
		list += line
	}

	if count == 0 {
		list = "\n\n  " + styleBrand.Render("No existing identities found.")
	}

	divLine := "\n\n  " + divider(pw-4)
	confirmStyle := lipgloss.NewStyle().Bold(true).Foreground(colorPurple).Render
	confirmBtn := divLine +
		"\n\n  " + confirmStyle("╔══════════════════════════════╗") + " " + styleBrand.Render("╔════════════════╗") +
		"\n  " + confirmStyle("║") + " " + lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("✓  Import selected            ") + confirmStyle("║") + " " + styleBrand.Render("║  ×  Skip       ║") +
		"\n  " + confirmStyle("╚══════════════════════════════╝") + " " + styleBrand.Render("╚════════════════╝")

	footerPairs := [][2]string{
		{"↑/↓", "move"},
		{"space", "toggle"},
		{"enter", "import"},
		{"esc", "back"},
	}
	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, footerPairs)

	body := header + title + list + confirmBtn + footer
	return stylePanelBorder(pw).Render(body)
}

// viewWizardAddMore — new user step 3: offer to add more accounts.
func (m Model) viewWizardAddMore(pw int) string {
	header := m.viewHeader("Setup")
	sep := "\n\n"

	imported := len(m.profiles)
	var status string
	if imported > 0 {
		importedSuffix := func() string {
			if imported == 1 {
				return "y"
			}
			return "ies"
		}()
		status = sep + "  " + lipgloss.NewStyle().Foreground(colorGreen).Render("✓") + "  " + styleCurrentVal.Render(fmt.Sprintf("%d identit%s imported.", imported, importedSuffix))
	} else {
		status = sep + "  " + styleBrand.Render("No identities imported yet.")
	}

	divLine := "\n\n  " + divider(pw-4)

	var question string
	if m.arcadeMode {
		question = divLine + sep + "  " + styleTitle.Render("ADD ANOTHER ACCOUNT?")
	} else {
		question = divLine + sep + "  Add another GitHub account?"
	}

	btn := func(focused bool, icon, label string) string {
		var border lipgloss.Style
		if focused {
			if m.arcadeMode {
				border = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
			} else {
				border = lipgloss.NewStyle().Foreground(colorPurple).Bold(true)
			}
		} else {
			border = styleBrand
		}
		top := border.Render("┌──────────────────────────────────────┐")
		mid := border.Render("│") + "  " + icon + "  " + styleCurrentVal.Render(fmt.Sprintf("%-36s", label)) + border.Render("│")
		bot := border.Render("└──────────────────────────────────────┘")
		return "\n\n  " + top + "\n  " + mid + "\n  " + bot
	}

	oauthIcon := lipgloss.NewStyle().Foreground(colorGreen).Render("◉")
	manualIcon := lipgloss.NewStyle().Foreground(colorPurple).Render("✎")
	doneIcon := lipgloss.NewStyle().Foreground(colorYellow).Render("→")

	var oauthLabel, manualLabel, doneLabel string
	if m.arcadeMode {
		oauthLabel = "LOG IN WITH GITHUB  [OAuth]"
		manualLabel = "ADD MANUALLY"
		doneLabel = "DONE, OPEN GITSWITCH"
	} else {
		oauthLabel = "Log in with GitHub  [OAuth device flow]"
		manualLabel = "Add manually"
		doneLabel = "Done, open gitswitch"
	}

	buttons := question +
		btn(m.wizardStep == 0, oauthIcon, oauthLabel) +
		btn(m.wizardStep == 1, manualIcon, manualLabel) +
		btn(m.wizardStep == 2, doneIcon, doneLabel)

	footerPairs := [][2]string{
		{"↑/↓", "move"},
		{"enter", "select"},
		{"esc", "back"},
	}
	footer := "\n\n" + divider(pw) + "\n" + m.footerKeys(pw, footerPairs)

	body := header + status + buttons + footer
	return stylePanelBorder(pw).Render(body)
}

// viewWizardDone — new user step 4: setup complete.
func (m Model) viewWizardDone(pw int) string {
	header := m.viewHeader("Setup complete")
	sep := "\n\n"

	imported := len(m.profiles)
	check := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	importedSuffix := func() string {
		if imported == 1 {
			return "y"
		}
		return "ies"
	}()
	var doneMsg string
	if m.arcadeMode {
		playerSuffix := func() string {
			if imported == 1 {
				return ""
			}
			return "S"
		}()
		doneMsg = sep + "  " + check + "  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("+500 EXP  ·  %d PLAYER%s LOADED", imported, playerSuffix))
	} else {
		doneMsg = sep + "  " + check + "  All done! " + styleCurrentVal.Render(fmt.Sprintf("%d identit%s ready.", imported, importedSuffix))
	}

	divLine := "\n\n  " + divider(pw-4)

	// Show first profile as active
	var activeInfo string
	if len(m.profiles) > 0 {
		p := m.profiles[0]
		label := "Current"
		if m.arcadeMode {
			label = "PLAYER 1"
		}
		activeInfo = divLine + sep + "  " + styleCurrent.Render(label+"  ") + check + "  " + styleCurrentVal.Render(p.Nickname) +
			"\n  " + styleBrand.Render("               "+p.UserName+" · "+p.Email)
	}

	ctaBorderStyle := lipgloss.NewStyle().Bold(true).Foreground(colorPurple)
	var ctaLabel string
	if m.arcadeMode {
		ctaBorderStyle = lipgloss.NewStyle().Bold(true).Foreground(colorYellow)
		ctaLabel = "▶  OPEN GITSWITCH"
	} else {
		ctaLabel = "▶  Open gitswitch"
	}
	cta := divLine +
		"\n\n  " + ctaBorderStyle.Render("╔══════════════════════════════════════╗") +
		"\n  " + ctaBorderStyle.Render("║") + "  " + lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render(fmt.Sprintf("%-38s", ctaLabel)) + ctaBorderStyle.Render("║") +
		"\n  " + ctaBorderStyle.Render("╚══════════════════════════════════════╝")

	var tip string
	if m.arcadeMode {
		tip = "\n\n  " + styleBrand.Render("tip: press ") + styleFooterKey.Render("[A]") + styleBrand.Render(" to add more players anytime")
	} else {
		tip = "\n\n  " + styleBrand.Render("tip: press ") + styleFooterKey.Render("[a]") + styleBrand.Render(" to add more accounts anytime")
	}

	body := header + doneMsg + activeInfo + cta + tip
	return stylePanelBorder(pw).Render(body)
}
