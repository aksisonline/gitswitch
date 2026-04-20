package storage

import (
	"fmt"
	"os/exec"
	"strings"
)

func (s *Store) ImportCurrent() error {
	profiles, err := s.Load()
	if err != nil {
		return err
	}
	if len(profiles) > 0 {
		return fmt.Errorf("profiles already exist")
	}

	userName, email, err := getGitConfig()
	if err != nil {
		return err
	}
	if userName == "" || email == "" {
		return fmt.Errorf("no existing git config found")
	}

	signKey, _ := getGitSignKey()

	p := Profile{
		Nickname: "default",
		UserName: userName,
		Email:    email,
		SignKey:  signKey,
		Active:   true,
	}

	return s.Save([]Profile{p})
}

func getGitConfig() (userName, email string, err error) {
	nameOut, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		return "", "", nil
	}
	emailOut, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return "", "", nil
	}
	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut)), nil
}

func getGitSignKey() (string, error) {
	out, err := exec.Command("git", "config", "--global", "user.signingkey").Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}
