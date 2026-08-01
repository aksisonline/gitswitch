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

// AutoPinThreshold is how many times the same nickname must be used in a repo
// before gitswitch trusts the pattern enough to act on it on its own —
// auto-pinning a not-yet-pinned repo, or recommending a switch away from the
// currently active profile.
const AutoPinThreshold = 3

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
		return nil, err
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

// saveToPath writes h to an explicit file path atomically (temp file + rename).
func saveToPath(h *History, path string) error {
	dir := filepath.Dir(path)
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
	return os.Rename(tmpName, path)
}

// Save writes history to the default path atomically (temp file + rename) to
// prevent corruption from concurrent writers.
func Save(h *History) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	return saveToPath(h, filepath.Join(dir, "history.json"))
}

// recordInHistory increments the usage count for nickname in h without touching the file.
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

// recordAt loads history from histPath, increments the count, and saves back.
// Returns nickname's updated usage count for repoKey.
// It is used by Record (with the default path) and by tests (with a temp path).
func recordAt(histPath, repoKey, nickname string) (count int, err error) {
	h, err := loadFromPath(histPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
		h = &History{Repos: make(map[string]RepoHistory)}
	}
	recordInHistory(h, repoKey, nickname)
	if err := saveToPath(h, histPath); err != nil {
		return 0, err
	}
	return h.Repos[repoKey].Identities[nickname], nil
}

// Record increments the usage count for nickname under repoKey and returns its
// new total — callers compare this against AutoPinThreshold to decide whether
// to auto-pin, without a second read of history.json.
// It holds an exclusive advisory lock for the duration of the read-modify-write
// to prevent lost updates when multiple shells call "gitswitch record" concurrently.
func Record(repoKey, nickname string) (count int, err error) {
	if repoKey == "" || nickname == "" {
		return 0, nil
	}
	dir, err := configDir()
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, err
	}
	err = withLock(dir, func() error {
		var lockErr error
		count, lockErr = recordAt(filepath.Join(dir, "history.json"), repoKey, nickname)
		return lockErr
	})
	return count, err
}

// Recommend returns the suggested nickname for repoKey.
// Pinned identity always wins over auto-learned counts.
// Auto-learned threshold: ≥AutoPinThreshold uses AND ≥60% share AND differs from currentNickname.
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
	if total == 0 || topCount < AutoPinThreshold {
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
// Locked the same way as Record, so a pin can never race a concurrent
// "gitswitch record" (or another Pin) into a lost update.
func Pin(repoKey, nickname string) error {
	if repoKey == "" || nickname == "" {
		return nil
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return withLock(dir, func() error {
		path := filepath.Join(dir, "history.json")
		h, err := loadFromPath(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			h = &History{Repos: make(map[string]RepoHistory)}
		}
		rh := h.Repos[repoKey]
		if rh.Identities == nil {
			rh.Identities = make(map[string]int)
		}
		rh.Pinned = nickname
		h.Repos[repoKey] = rh
		return saveToPath(h, path)
	})
}

// Unpin clears the pinned identity for repoKey, falling back to auto-learned counts.
func Unpin(repoKey string) error {
	if repoKey == "" {
		return nil
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return withLock(dir, func() error {
		path := filepath.Join(dir, "history.json")
		h, err := loadFromPath(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		rh, exists := h.Repos[repoKey]
		if !exists {
			return nil
		}
		rh.Pinned = ""
		h.Repos[repoKey] = rh
		return saveToPath(h, path)
	})
}

// GetPinned returns the pinned nickname for repoKey, or "" if none.
func GetPinned(repoKey string) string {
	h, err := Load()
	if err != nil {
		return ""
	}
	return h.Repos[repoKey].Pinned
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
