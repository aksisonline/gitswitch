package tui

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// profileFormData is heap-allocated and pointed to by Model so that huh's
// value bindings survive the value-copy of Model on every Update.
type profileFormData struct {
	nickname string
	userName string
	email    string
	signKey  string
	sshKey   string
	ghUser   string
}

func required(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errRequired{label}
		}
		return nil
	}
}

type errRequired struct{ label string }

func (e errRequired) Error() string { return e.label + " is required" }

// huhTheme builds a huh theme from the current (already theme-applied) color
// vars so forms match the active palette / arcade skin.
func huhTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.
		BorderForeground(colorPurple)
	t.Focused.Title = t.Focused.Title.Foreground(colorPurple).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(colorDim)
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(colorYellow)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(colorGreen)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(colorWhite)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(colorDim)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(colorRed)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(colorRed)

	t.Blurred.Title = t.Blurred.Title.Foreground(colorDim)
	t.Blurred.Description = t.Blurred.Description.Foreground(colorDim)
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(colorWhite)
	t.Blurred.TextInput.Placeholder = t.Blurred.TextInput.Placeholder.Foreground(colorDim)
	t.Blurred.Base = t.Blurred.Base.BorderForeground(colorDim)

	t.Help.ShortKey = lipgloss.NewStyle().Foreground(colorYellow)
	t.Help.ShortDesc = lipgloss.NewStyle().Foreground(colorDim)
	t.Help.ShortSeparator = lipgloss.NewStyle().Foreground(colorDim)
	return t
}

// newProfileForm builds the add/edit form bound to d. arcade only tweaks copy.
func newProfileForm(d *profileFormData, edit bool, width int) *huh.Form {
	if width < 32 {
		width = 32
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("nickname").
				Title("Nickname").
				Description("A label for this identity — not written to git.").
				Placeholder("personal").
				Value(&d.nickname).
				Validate(required("Nickname")),
			huh.NewInput().
				Key("userName").
				Title("Name").
				Description("git user.name — the author name on your commits.").
				Placeholder("Ada Lovelace").
				Value(&d.userName).
				Validate(required("Name")),
			huh.NewInput().
				Key("email").
				Title("Email").
				Description("git user.email — the author email on your commits.").
				Placeholder("ada@example.com").
				Value(&d.email).
				Validate(required("Email")),
			huh.NewInput().
				Key("ghUser").
				Title("GitHub username").
				Description("Optional — used to switch the active gh CLI account.").
				Placeholder("optional").
				Value(&d.ghUser),
			huh.NewInput().
				Key("sshKey").
				Title("SSH key path").
				Description("Optional — sets core.sshCommand, e.g. ~/.ssh/id_work.").
				Placeholder("optional").
				Value(&d.sshKey),
			huh.NewInput().
				Key("signKey").
				Title("GPG signing key").
				Description("Optional — git user.signingkey for signed commits.").
				Placeholder("optional").
				Value(&d.signKey),
		),
	).
		WithTheme(huhTheme()).
		WithWidth(width).
		WithShowHelp(true).
		WithShowErrors(true)

	return form
}
