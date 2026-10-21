package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/git"
	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/storage"
)

func setupGuardTestEnv(t *testing.T) (*storage.Store, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	cfg := filepath.Join(home, "gitconfig")
	if err := os.WriteFile(cfg, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", cfg)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	st, err := storage.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	if err := exec.Command("git", "-C", repo, "init").Run(); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "-C", repo, "remote", "add", "origin", "git@github.com:acme/repo.git").Run(); err != nil {
		t.Fatal(err)
	}
	return st, repo
}

func TestCheckGuard_ZeroProfilesAllowsCommit(t *testing.T) {
	st, repo := setupGuardTestEnv(t)
	if err := checkGuard(st, repo); err != nil {
		t.Errorf("expected nil with zero profiles, got: %v", err)
	}
}

func TestCheckGuard_NoEmailConfiguredBlocks(t *testing.T) {
	st, repo := setupGuardTestEnv(t)
	if err := st.Add("work", "Alice", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}

	err := checkGuard(st, repo)
	if err == nil || !strings.Contains(err.Error(), "no git user.email is configured") {
		t.Errorf("expected no-email block, got: %v", err)
	}
}

func TestCheckGuard_UnknownEmailBlocks(t *testing.T) {
	st, repo := setupGuardTestEnv(t)
	if err := st.Add("work", "Alice", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "config", "--global", "user.email", "stranger@other.com").Run(); err != nil {
		t.Fatal(err)
	}

	err := checkGuard(st, repo)
	if err == nil || !strings.Contains(err.Error(), "does not match any gitswitch profile") {
		t.Errorf("expected unknown-email block, got: %v", err)
	}
}

func TestCheckGuard_PinnedRepoMismatchBlocks(t *testing.T) {
	st, repo := setupGuardTestEnv(t)
	if err := st.Add("personal", "Alice", "alice@personal.dev", "", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := st.Add("work", "Alice Corp", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}

	// Global is personal, but repo is pinned in history.json to work without local git config matching.
	if err := exec.Command("git", "config", "--global", "user.email", "alice@personal.dev").Run(); err != nil {
		t.Fatal(err)
	}
	repoKey := history.GetRepoKeyForPath(repo)
	if err := history.Pin(repoKey, "work"); err != nil {
		t.Fatal(err)
	}

	err := checkGuard(st, repo)
	if err == nil || !strings.Contains(err.Error(), "pinned to 'work'") {
		t.Errorf("expected pinned mismatch block, got: %v", err)
	}

	// Once local config is set to work's email (as `gitswitch pin work` writes), checkGuard passes.
	if err := exec.Command("git", "-C", repo, "config", "--local", "user.email", "alice@work.com").Run(); err != nil {
		t.Fatal(err)
	}
	if err := checkGuard(st, repo); err != nil {
		t.Errorf("expected pinned repo with matching local email to pass, got: %v", err)
	}

	// An explicit GIT_AUTHOR_EMAIL override to personal in a repo pinned to work must still be blocked.
	t.Setenv("GIT_AUTHOR_EMAIL", "alice@personal.dev")
	err = checkGuard(st, repo)
	if err == nil || !strings.Contains(err.Error(), "pinned to 'work'") {
		t.Errorf("expected env override against pin to be blocked, got: %v", err)
	}
}

func TestCheckGuard_LearnedHistoryMismatchBlocks(t *testing.T) {
	st, repo := setupGuardTestEnv(t)
	if err := st.Add("personal", "Alice", "alice@personal.dev", "", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := st.Add("work", "Alice Corp", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}

	repoKey := history.GetRepoKeyForPath(repo)
	for i := 0; i < history.AutoPinThreshold; i++ {
		if _, err := history.Record(repoKey, "work"); err != nil {
			t.Fatal(err)
		}
	}

	// Active global is personal → should block because repo usually commits as work.
	if err := exec.Command("git", "config", "--global", "user.email", "alice@personal.dev").Run(); err != nil {
		t.Fatal(err)
	}
	err := checkGuard(st, repo)
	if err == nil || !strings.Contains(err.Error(), "usually commits as 'work'") {
		t.Errorf("expected learned history mismatch block, got: %v", err)
	}

	// Switching global to work → passes.
	if err := exec.Command("git", "config", "--global", "user.email", "alice@work.com").Run(); err != nil {
		t.Fatal(err)
	}
	if err := checkGuard(st, repo); err != nil {
		t.Errorf("expected matching learned profile to pass, got: %v", err)
	}
}

// TestGuardHookEndToEndGitCommit builds the real gitswitch binary, enables the
// global pre-commit guard hook, and runs a real `git commit` to verify that a
// wrong-author commit is refused before entering git log and a matching commit
// succeeds.
func TestGuardHookEndToEndGitCommit(t *testing.T) {
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, "gitswitch")
	if out, err := exec.Command("go", "build", "-o", binPath, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build gitswitch: %v\n%s", err, out)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	st, repo := setupGuardTestEnv(t)
	if err := st.Add("personal", "Alice", "alice@personal.dev", "", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := st.Add("work", "Alice Corp", "alice@work.com", "", "", ""); err != nil {
		t.Fatal(err)
	}
	if err := st.SetActive("personal"); err != nil {
		t.Fatal(err)
	}

	if err := git.InstallGuardHook(); err != nil {
		t.Fatalf("InstallGuardHook: %v", err)
	}

	// Pin the repo to 'work' in history.json, but set git config to 'personal'.
	repoKey := history.GetRepoKeyForPath(repo)
	if err := history.Pin(repoKey, "work"); err != nil {
		t.Fatal(err)
	}
	_ = exec.Command("git", "config", "--global", "user.name", "Alice").Run()
	_ = exec.Command("git", "config", "--global", "user.email", "alice@personal.dev").Run()

	readme := filepath.Join(repo, "README.md")
	if err := os.WriteFile(readme, []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "-C", repo, "add", "README.md").Run(); err != nil {
		t.Fatal(err)
	}

	// 1. Commit with mismatched author must be refused by the pre-commit hook!
	out, err := exec.Command("git", "-C", repo, "commit", "-m", "wrong account commit").CombinedOutput()
	if err == nil {
		t.Fatalf("expected git commit to be refused by pre-commit guard, but it succeeded:\n%s", out)
	}
	if !strings.Contains(string(out), "pinned to 'work'") {
		t.Errorf("expected guard error output to mention 'pinned to work', got:\n%s", out)
	}

	// Verify nothing entered git history.
	if err := exec.Command("git", "-C", repo, "rev-parse", "--verify", "HEAD").Run(); err == nil {
		t.Fatal("commit entered git history despite guard refusal")
	}

	// 2. Fix identity (`git config --local user.email alice@work.com`) and commit again -> succeeds!
	_ = exec.Command("git", "-C", repo, "config", "--local", "user.name", "Alice Corp").Run()
	_ = exec.Command("git", "-C", repo, "config", "--local", "user.email", "alice@work.com").Run()
	if out, err := exec.Command("git", "-C", repo, "commit", "-m", "right account commit").CombinedOutput(); err != nil {
		t.Fatalf("expected git commit with matching pinned identity to succeed, got: %v\n%s", err, out)
	}
}
