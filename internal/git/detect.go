package git

import (
	"github.com/aksisonline/gitswitch/internal/storage"
)

// DetectActive returns the profile whose user.name + user.email match
// the current global git config. Falls back to the JSON active flag if
// no config match is found.
func DetectActive(profiles []storage.Profile) *storage.Profile {
	cfg := New(true)
	currentName, currentEmail, err := cfg.GetUser()
	if err != nil || (currentName == "" && currentEmail == "") {
		return flaggedActive(profiles)
	}

	for i, p := range profiles {
		if p.UserName == currentName && p.Email == currentEmail {
			return &profiles[i]
		}
	}

	// No exact match — fall back to the JSON flag
	return flaggedActive(profiles)
}

func flaggedActive(profiles []storage.Profile) *storage.Profile {
	for i, p := range profiles {
		if p.Active {
			return &profiles[i]
		}
	}
	return nil
}
