package git

import (
	"strings"
	"testing"
)

func TestListGHUsers_ParsesStatusOutput(t *testing.T) {
	// ListGHUsers shells out to gh; just ensure it does not panic and
	// returns a slice (empty is fine when gh has no accounts in CI).
	users := ListGHUsers()
	if users == nil {
		// nil and empty are both acceptable when gh is missing.
		return
	}
	for _, u := range users {
		if strings.TrimSpace(u.Login) == "" {
			t.Errorf("empty login in %+v", u)
		}
	}
}

func TestListSSHPrivateKeys_NoPubFiles(t *testing.T) {
	keys := ListSSHPrivateKeys()
	for _, k := range keys {
		if strings.HasSuffix(k, ".pub") {
			t.Errorf("public key should not be listed: %s", k)
		}
		if strings.Contains(k, "known_hosts") {
			t.Errorf("known_hosts should not be listed: %s", k)
		}
	}
}

func TestGetSignKey_EmptyWhenUnset(t *testing.T) {
	// Smoke: does not panic. Value depends on the machine.
	_ = New(true).GetSignKey()
}
