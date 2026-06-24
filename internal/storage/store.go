package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	Nickname string `json:"nickname"  yaml:"nickname"`
	UserName string `json:"user_name"  yaml:"user_name"`
	Email    string `json:"email"      yaml:"email"`
	SignKey  string `json:"sign_key,omitempty"  yaml:"sign_key,omitempty"`
	SSHKey   string `json:"ssh_key,omitempty"   yaml:"ssh_key,omitempty"`
	GHUser   string `json:"gh_user,omitempty"   yaml:"gh_user,omitempty"`
	TokenRef string `json:"token_ref,omitempty" yaml:"token_ref,omitempty"`
	Active   bool   `json:"active"     yaml:"active"`
}

// legacyProfile handles migration from older JSON formats.
type legacyProfile struct {
	Nickname string `json:"nickname"`
	UserName string `json:"user_name"`
	Name     string `json:"name"` // v1: was both label and git user.name
	Email    string `json:"email"`
	SignKey  string `json:"sign_key,omitempty"`
	SSHKey   string `json:"ssh_key,omitempty"`
	GHUser   string `json:"gh_user,omitempty"`
	Active   bool   `json:"active"`
}

// configFile is the on-disk YAML schema (v2+).
type configFile struct {
	Version  int       `yaml:"version"`
	Profiles []Profile `yaml:"profiles"`
}

type Store struct {
	path        string
	wasMigrated bool
}

// WasMigrated reports whether Load() performed a first-run migration from
// profiles.json to config.yaml during this session.
func (s *Store) WasMigrated() bool { return s.wasMigrated }

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

func (s *Store) yamlPath() string {
	return filepath.Join(s.path, "config.yaml")
}

func (s *Store) legacyJSONPath() string {
	return filepath.Join(s.path, "profiles.json")
}

func (s *Store) ConfigDir() string {
	return s.path
}

type Prefs struct {
	ColorTheme    int  `json:"color_theme"`
	SplashSeen020 bool `json:"splash_seen_020"`
	ShellEnabled  bool `json:"shell_enabled"`
	// ShowUsername toggles the Accounts secondary column between email
	// (zero value / default) and the GitHub username.
	ShowUsername bool `json:"show_username"`
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
	if p.ColorTheme < 0 || p.ColorTheme >= 12 {
		p.ColorTheme = 0
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
	// Primary: config.yaml
	if data, err := os.ReadFile(s.yamlPath()); err == nil {
		var cf configFile
		if err := yaml.Unmarshal(data, &cf); err != nil {
			// config.yaml corrupt — try the v1 backup before giving up
			if bakData, bakErr := os.ReadFile(s.legacyJSONPath() + ".v1.bak"); bakErr == nil {
				return s.migrateFromJSON(bakData)
			}
			return nil, fmt.Errorf("parse config.yaml: %w (no backup found)", err)
		}
		return cf.Profiles, nil
	}

	// Migration path: profiles.json exists, config.yaml does not
	data, err := os.ReadFile(s.legacyJSONPath())
	if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil
		}
		return nil, err
	}
	profiles, err := s.migrateFromJSON(data)
	if err != nil {
		return nil, err
	}
	// Write config.yaml atomically and rename old file
	if saveErr := s.Save(profiles); saveErr == nil {
		_ = os.Rename(s.legacyJSONPath(), s.legacyJSONPath()+".v1.bak")
		s.wasMigrated = true
	}
	return profiles, nil
}

func (s *Store) migrateFromJSON(data []byte) ([]Profile, error) {
	var raw []legacyProfile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	profiles := make([]Profile, len(raw))
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
		}
		if p.Nickname == "" && p.UserName != "" {
			p.Nickname = p.UserName
		}
		profiles[i] = p
	}
	return profiles, nil
}

func (s *Store) Save(profiles []Profile) error {
	cf := configFile{Version: 2, Profiles: profiles}
	data, err := yaml.Marshal(cf)
	if err != nil {
		return err
	}
	tmp := s.yamlPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.yamlPath())
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
