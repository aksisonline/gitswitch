package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Nickname string `json:"nickname"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	SignKey  string `json:"sign_key,omitempty"`
	SSHKey   string `json:"ssh_key,omitempty"`  // path to SSH private key, e.g. ~/.ssh/id_work
	GHUser   string `json:"gh_user,omitempty"`  // GitHub CLI username for gh auth switch
	Active   bool   `json:"active"`
}

// legacyProfile handles migration from old formats.
type legacyProfile struct {
	Nickname string `json:"nickname"`
	UserName string `json:"user_name"`
	Name     string `json:"name"` // v1 format: was both label and git user.name
	Email    string `json:"email"`
	SignKey  string `json:"sign_key,omitempty"`
	SSHKey   string `json:"ssh_key,omitempty"`
	GHUser   string `json:"gh_user,omitempty"`
	Active   bool   `json:"active"`
}

type Store struct {
	path string
}

func New() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".config", "gitswitch")
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

func (s *Store) filePath() string {
	return filepath.Join(s.path, "profiles.json")
}

func (s *Store) ConfigDir() string {
	return s.path
}

type Prefs struct {
	ColorTheme int `json:"color_theme"`
}

func (s *Store) prefsPath() string {
	return filepath.Join(s.path, "config.json")
}

func (s *Store) LoadPrefs() (Prefs, error) {
	data, err := os.ReadFile(s.prefsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return Prefs{}, nil
		}
		return Prefs{}, err
	}
	var p Prefs
	if err := json.Unmarshal(data, &p); err != nil {
		return Prefs{}, err
	}
	return p, nil
}

func (s *Store) SavePrefs(p Prefs) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.prefsPath(), data, 0600)
}

func (s *Store) Load() ([]Profile, error) {
	data, err := os.ReadFile(s.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil
		}
		return nil, err
	}

	var raw []legacyProfile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	profiles := make([]Profile, len(raw))
	needsMigration := false
	for i, r := range raw {
		p := Profile{
			Nickname: r.Nickname,
			UserName: r.UserName,
			Email:    r.Email,
			SignKey:  r.SignKey,
			SSHKey:   r.SSHKey,
			GHUser:   r.GHUser,
			Active:   r.Active,
		}
		if p.UserName == "" && r.Name != "" {
			p.UserName = r.Name
			needsMigration = true
		}
		if p.Nickname == "" && p.UserName != "" {
			p.Nickname = p.UserName
			needsMigration = true
		}
		profiles[i] = p
	}

	if needsMigration {
		_ = s.Save(profiles)
	}

	return profiles, nil
}

func (s *Store) Save(profiles []Profile) error {
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(), data, 0600)
}

func (s *Store) Add(nickname, userName, email, signKey, sshKey, ghUser string) error {
	profiles, err := s.Load()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		if p.Nickname == nickname {
			return fmt.Errorf("profile '%s' already exists", nickname)
		}
	}
	profiles = append(profiles, Profile{
		Nickname: nickname,
		UserName: userName,
		Email:    email,
		SignKey:  signKey,
		SSHKey:   sshKey,
		GHUser:   ghUser,
	})
	return s.Save(profiles)
}

func (s *Store) Update(nickname string, updated Profile) error {
	profiles, err := s.Load()
	if err != nil {
		return err
	}
	for i, p := range profiles {
		if p.Nickname == nickname {
			updated.Active = p.Active
			profiles[i] = updated
			return s.Save(profiles)
		}
	}
	return fmt.Errorf("profile '%s' not found", nickname)
}

func (s *Store) Remove(nickname string) error {
	profiles, err := s.Load()
	if err != nil {
		return err
	}
	filtered := []Profile{}
	found := false
	for _, p := range profiles {
		if p.Nickname != nickname {
			filtered = append(filtered, p)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("profile '%s' not found", nickname)
	}
	return s.Save(filtered)
}

func (s *Store) Get(nickname string) (*Profile, error) {
	profiles, err := s.Load()
	if err != nil {
		return nil, err
	}
	for _, p := range profiles {
		if p.Nickname == nickname {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("profile '%s' not found", nickname)
}

func (s *Store) SetActive(nickname string) error {
	profiles, err := s.Load()
	if err != nil {
		return err
	}
	found := false
	for i := range profiles {
		profiles[i].Active = false
		if profiles[i].Nickname == nickname {
			profiles[i].Active = true
			found = true
		}
	}
	if !found {
		return fmt.Errorf("profile '%s' not found", nickname)
	}
	return s.Save(profiles)
}

func (s *Store) GetActive() (*Profile, error) {
	profiles, err := s.Load()
	if err != nil {
		return nil, err
	}
	for _, p := range profiles {
		if p.Active {
			return &p, nil
		}
	}
	return nil, nil
}
