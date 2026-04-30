package history

import (
	"os"
	"path/filepath"
	"testing"
)

const testRepo = "git@github.com:org/repo.git"

// newHistory builds a History value inline for test readability.
func newHistory(repos map[string]RepoHistory) *History {
	return &History{Repos: repos}
}

func rh(pinned string, counts map[string]int) RepoHistory {
	if counts == nil {
		counts = make(map[string]int)
	}
	return RepoHistory{Pinned: pinned, Identities: counts}
}

// ── Recommend: auto-learned ──────────────────────────────────────────────────

func TestRecommend_NoHistory(t *testing.T) {
	h := newHistory(nil)
	_, ok := recommendFromHistory(h, testRepo, "")
	if ok {
		t.Error("expected no recommendation for empty history")
	}
}

func TestRecommend_RepoNotInHistory(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		"other-repo": rh("", map[string]int{"work": 5}),
	})
	_, ok := recommendFromHistory(h, testRepo, "")
	if ok {
		t.Error("expected no recommendation for unknown repo")
	}
}

func TestRecommend_BelowMinCount(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 2}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "")
	if ok {
		t.Errorf("count=2 below threshold, should not recommend, got %q", nick)
	}
}

func TestRecommend_BelowShareThreshold(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 3, "personal": 3}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "")
	if ok {
		t.Errorf("50%% share below 60%% threshold, should not recommend, got %q", nick)
	}
}

func TestRecommend_ExactlyAt60Percent(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 3, "personal": 2}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "work" {
		t.Errorf("3/5 = 60%% should pass threshold, got %q ok=%v", nick, ok)
	}
}

func TestRecommend_MeetsThreshold(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 9, "personal": 1}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "work" {
		t.Errorf("expected 'work', got %q ok=%v", nick, ok)
	}
}

func TestRecommend_AlreadyOnTopIdentity(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 9, "personal": 1}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "work")
	if ok {
		t.Errorf("already on recommended identity, should not nudge, got %q", nick)
	}
}

// ── Recommend: pinned ────────────────────────────────────────────────────────

func TestRecommend_PinnedWinsOverCounts(t *testing.T) {
	// counts would recommend "work" (90%), but pinned says "personal"
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("personal", map[string]int{"work": 9, "personal": 1}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "personal" {
		t.Errorf("pinned should override counts, got %q ok=%v", nick, ok)
	}
}

func TestRecommend_PinnedWinsWithNoCountHistory(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("work", nil),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "work" {
		t.Errorf("pinned with no counts should still recommend, got %q ok=%v", nick, ok)
	}
}

func TestRecommend_PinnedAlreadyActive(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("work", nil),
	})
	nick, ok := recommendFromHistory(h, testRepo, "work")
	if ok {
		t.Errorf("already on pinned identity, should not nudge, got %q", nick)
	}
}

func TestRecommend_UnpinFallsBackToCounts(t *testing.T) {
	// After unpin, counts-based logic should take over
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 9, "personal": 1}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "work" {
		t.Errorf("after unpin, counts should drive recommendation, got %q ok=%v", nick, ok)
	}
}

func TestRecommend_UnpinWithCountsBelowThreshold(t *testing.T) {
	// Unpinned and counts too low → silence
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 1}),
	})
	nick, ok := recommendFromHistory(h, testRepo, "other")
	if ok {
		t.Errorf("unpinned with insufficient counts should not recommend, got %q", nick)
	}
}

// ── Pin / Unpin round-trip ───────────────────────────────────────────────────

func TestPinUnpinRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "history.json")

	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 9, "personal": 1}),
	})
	data, err := marshalHistory(h)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	// pin
	h.Repos[testRepo] = func() RepoHistory {
		r := h.Repos[testRepo]; r.Pinned = "personal"; return r
	}()
	nick, ok := recommendFromHistory(h, testRepo, "work")
	if !ok || nick != "personal" {
		t.Errorf("after pin: expected 'personal', got %q ok=%v", nick, ok)
	}

	// unpin
	h.Repos[testRepo] = func() RepoHistory {
		r := h.Repos[testRepo]; r.Pinned = ""; return r
	}()
	nick, ok = recommendFromHistory(h, testRepo, "other")
	if !ok || nick != "work" {
		t.Errorf("after unpin: expected 'work' from counts, got %q ok=%v", nick, ok)
	}
}

// ── Record ───────────────────────────────────────────────────────────────────

func TestRecord_IncrementsCount(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("", map[string]int{"work": 2}),
	})
	rh := h.Repos[testRepo]
	rh.Identities["work"]++
	rh.LastUsed = "work"
	h.Repos[testRepo] = rh

	if h.Repos[testRepo].Identities["work"] != 3 {
		t.Errorf("expected count 3, got %d", h.Repos[testRepo].Identities["work"])
	}
	if h.Repos[testRepo].LastUsed != "work" {
		t.Errorf("expected last_used 'work', got %q", h.Repos[testRepo].LastUsed)
	}
}

func TestRecord_DoesNotTouchPinned(t *testing.T) {
	h := newHistory(map[string]RepoHistory{
		testRepo: rh("personal", map[string]int{"work": 2}),
	})
	// simulate record call (increments count, never touches Pinned)
	rh := h.Repos[testRepo]
	rh.Identities["work"]++
	rh.LastUsed = "work"
	h.Repos[testRepo] = rh

	if h.Repos[testRepo].Pinned != "personal" {
		t.Errorf("record should not touch pinned field, got %q", h.Repos[testRepo].Pinned)
	}
}

// ── JSON round-trip ──────────────────────────────────────────────────────────

func TestLoadSaveRoundtrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "history.json")

	h := newHistory(map[string]RepoHistory{
		testRepo: rh("work", map[string]int{"work": 5, "aks": 1}),
	})
	data, err := marshalHistory(h)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	h2, err := loadFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	repo := h2.Repos[testRepo]
	if repo.Pinned != "work" {
		t.Errorf("pinned: expected 'work', got %q", repo.Pinned)
	}
	if repo.Identities["work"] != 5 {
		t.Errorf("count: expected 5, got %d", repo.Identities["work"])
	}
	if repo.Identities["aks"] != 1 {
		t.Errorf("count: expected 1 for aks, got %d", repo.Identities["aks"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	// loadFromPath on missing file returns error, but Load() should return empty History
	_, err := loadFromPath(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Error("expected error reading nonexistent file")
	}
}
