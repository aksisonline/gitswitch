package history

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RepoHistory struct {
	Pinned     string         `json:"pinned,omitempty"` // manually pinned nickname, always wins
	Identities map[string]int `json:"identities"`       // nickname -> count (auto-learned)
	LastUsed   string         `json:"last_used"`
}

type History struct {
	Repos map[string]RepoHistory `json:"repos"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gitswitch"), nil
}

func historyPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

// Load reads history.json; returns an empty History if the file does not exist.
// If the file is corrupted, backs it up as history.json.bak before starting fresh.
func Load() (*History, error) {
	path, err := historyPath()
	if err != nil {
		return &History{Repos: make(map[string]RepoHistory)}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Repos: make(map[string]RepoHistory)}, nil
		}
		return nil, err
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		// back up corrupted file so the user can inspect it, then start fresh
		_ = os.Rename(path, path+".bak")
		return &History{Repos: make(map[string]RepoHistory)}, nil
	}
	if h.Repos == nil {
		h.Repos = make(map[string]RepoHistory)
	}
	return &h, nil
}

// Save writes history to disk atomically (temp file + rename) to prevent corruption
// from concurrent writers.
func Save(h *History) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".history-*.json")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, filepath.Join(dir, "history.json"))
}

// recordInHistory increments the usage count for nickname in h without touching the file.
// Exported for use in tests.
func recordInHistory(h *History, repoKey, nickname string) {
	rh, ok := h.Repos[repoKey]
	if !ok {
		rh = RepoHistory{Identities: make(map[string]int)}
	}
	if rh.Identities == nil {
		rh.Identities = make(map[string]int)
	}
	rh.Identities[nickname]++
	rh.LastUsed = nickname
	h.Repos[repoKey] = rh
}

// Record increments the usage count for nickname under repoKey.
func Record(repoKey, nickname string) error {
	if repoKey == "" || nickname == "" {
		return nil
	}
	h, err := Load()
	if err != nil {
		return err
	}
	recordInHistory(h, repoKey, nickname)
	return Save(h)
}

// Recommend returns the suggested nickname for repoKey.
// Pinned identity always wins over auto-learned counts.
// Auto-learned threshold: ≥3 uses AND ≥60% share AND differs from currentNickname.
func Recommend(repoKey, currentNickname string) (nickname string, ok bool) {
	if repoKey == "" {
		return "", false
	}
	h, err := Load()
	if err != nil {
		return "", false
	}
	return recommendFromHistory(h, repoKey, currentNickname)
}

func recommendFromHistory(h *History, repoKey, currentNickname string) (string, bool) {
	rh, exists := h.Repos[repoKey]
	if !exists {
		return "", false
	}
	if rh.Pinned != "" {
		if rh.Pinned == currentNickname {
			return "", false
		}
		return rh.Pinned, true
	}
	var topNick string
	var topCount, total int
	for nick, count := range rh.Identities {
		total += count
		if count > topCount {
			topCount = count
			topNick = nick
		}
	}
	if total == 0 || topCount < 3 {
		return "", false
	}
	if float64(topCount)/float64(total) < 0.60 {
		return "", false
	}
	if topNick == currentNickname {
		return "", false
	}
	return topNick, true
}

// Pin permanently sets the recommended identity for repoKey.
func Pin(repoKey, nickname string) error {
	if repoKey == "" || nickname == "" {
		return nil
	}
	h, err := Load()
	if err != nil {
		return err
	}
	rh := h.Repos[repoKey]
	if rh.Identities == nil {
		rh.Identities = make(map[string]int)
	}
	rh.Pinned = nickname
	h.Repos[repoKey] = rh
	return Save(h)
}

// Unpin clears the pinned identity for repoKey, falling back to auto-learned counts.
func Unpin(repoKey string) error {
	if repoKey == "" {
		return nil
	}
	h, err := Load()
	if err != nil {
		return err
	}
	rh, exists := h.Repos[repoKey]
	if !exists {
		return nil
	}
	rh.Pinned = ""
	h.Repos[repoKey] = rh
	return Save(h)
}

// GetRepoKey resolves the repo key for the current working directory.
// Tries git remote URL first, falls back to absolute repo root path.
func GetRepoKey() string {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err == nil {
		key := strings.TrimSpace(string(out))
		if key != "" {
			return key
		}
	}
	out, err = exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		key := strings.TrimSpace(string(out))
		if key != "" {
			return key
		}
	}
	return ""
}

// GetRepoKeyForPath resolves the repo key for a given directory path.
func GetRepoKeyForPath(dir string) string {
	out, err := exec.Command("git", "-C", dir, "remote", "get-url", "origin").Output()
	if err == nil {
		key := strings.TrimSpace(string(out))
		if key != "" {
			return key
		}
	}
	out, err = exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
	if err == nil {
		key := strings.TrimSpace(string(out))
		if key != "" {
			return key
		}
	}
	return ""
}

// marshalHistory encodes a History to JSON bytes (used by tests).
func marshalHistory(h *History) ([]byte, error) {
	return json.MarshalIndent(h, "", "  ")
}

// loadFromPath reads a History from an explicit file path (used by tests).
func loadFromPath(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	if h.Repos == nil {
		h.Repos = make(map[string]RepoHistory)
	}
	return &h, nil
}
