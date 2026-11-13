package version

import "testing"

// The install method decides who may replace the binary, so reading it wrong
// either breaks a Homebrew install (self-update over a brew-owned file) or
// strands a script install (told to run brew). Pin the path parsing.
func TestInstallFromPath(t *testing.T) {
	cases := []struct {
		path    string
		brew    bool
		formula string
	}{
		{"/opt/homebrew/Cellar/gitswitch/0.5.0/bin/gitswitch", true, "gitswitch"},
		{"/usr/local/Cellar/gitswitch/0.5.0/bin/gitswitch", true, "gitswitch"},
		{"/opt/homebrew/Cellar/gitswitch-beta/0.5.1-beta.2/bin/gitswitch", true, "gitswitch-beta"},
		{"/home/linuxbrew/.linuxbrew/Cellar/gitswitch/0.5.0/bin/gitswitch", true, "gitswitch"},
		{"/home/linuxbrew/.linuxbrew/bin/gitswitch", true, "gitswitch"},
		{"/usr/local/bin/gitswitch", false, ""},
		{"/Users/me/go/bin/gitswitch", false, ""},
		{`C:\Users\me\AppData\Local\gitswitch\gitswitch.exe`, false, ""},
	}
	for _, c := range cases {
		got := installFromPath(c.path)
		if got.IsBrew() != c.brew || got.Formula != c.formula {
			t.Errorf("installFromPath(%q) = {brew:%v formula:%q}, want {brew:%v formula:%q}",
				c.path, got.IsBrew(), got.Formula, c.brew, c.formula)
		}
	}
}

func TestUpgradeHint(t *testing.T) {
	if got := (Install{Method: MethodBrew, Formula: BetaFormula}).UpgradeHint(); got != "brew upgrade gitswitch-beta" {
		t.Errorf("brew beta hint = %q", got)
	}
	if got := (Install{}).UpgradeHint(); got != "gitswitch upgrade" {
		t.Errorf("script hint = %q", got)
	}
	if got := (Install{Method: MethodBrew, Formula: StableFormula}).FormulaRef(); got != "aksisonline/tap/gitswitch" {
		t.Errorf("FormulaRef = %q", got)
	}
}

// A Homebrew install must never be handed the install-script command.
func TestUpgradeCommandForRejectsBadVersion(t *testing.T) {
	for _, inst := range []Install{{}, {Method: MethodBrew, Formula: StableFormula}} {
		if _, err := UpgradeCommandFor(inst, "v1; rm -rf ~"); err == nil {
			t.Errorf("UpgradeCommandFor(%+v) accepted a malformed version", inst)
		}
	}
}
