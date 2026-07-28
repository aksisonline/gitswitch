package tui

import (
	"strings"
	"testing"
)

// TestProfileFormFieldTitlesMatchForm guards against the mouse map (in
// update.go) drifting from the actual field titles (in form.go) — a rename
// there would silently break click-to-focus without this.
func TestProfileFormFieldTitlesMatchForm(t *testing.T) {
	d := &profileFormData{}
	form := newProfileForm(d, false, 60)
	form.Init()
	rendered := form.View()
	for _, f := range profileFormFieldTitles {
		if !strings.Contains(rendered, f.title) {
			t.Errorf("expected form to render title %q for key %q, got:\n%s", f.title, f.key, rendered)
		}
	}
}

// TestFormFieldAt pins the add/edit form's click-to-field mapping: huh has no
// mouse support of its own, so a click anywhere in a field's block (title,
// description, input line) must resolve to that field's key.
func TestFormFieldAt(t *testing.T) {
	lines := []string{
		"  Nickname",
		"  A label for this identity",
		"  > personal",
		"",
		"  Name",
		"  git user.name",
		"  > Ada Lovelace",
		"",
		"  Signing key",
		"  optional",
		"  > ",
	}
	cases := []struct {
		relY int
		want string
		ok   bool
	}{
		{0, "nickname", true},
		{1, "nickname", true},
		{2, "nickname", true},
		{3, "nickname", true}, // blank spacer still belongs to the field above
		{4, "userName", true},
		{6, "userName", true},
		{8, "signKey", true},
		{10, "signKey", true},
	}
	for _, c := range cases {
		got, ok := formFieldAt(lines, c.relY)
		if ok != c.ok || got != c.want {
			t.Errorf("formFieldAt(relY=%d) = (%q, %v), want (%q, %v)", c.relY, got, ok, c.want, c.ok)
		}
	}

	if _, ok := formFieldAt(lines, -1); ok {
		t.Error("formFieldAt(relY=-1) should report not-found")
	}
	if _, ok := formFieldAt(nil, 0); ok {
		t.Error("formFieldAt(nil lines) should report not-found")
	}
}
