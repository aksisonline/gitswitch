package credential

import (
	"strings"
	"testing"

	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/storage"
)

// newIsolatedStore points HOME at a temp dir so storage and history read/write
// an empty, isolated ~/.config/gitswitch.
func newIsolatedStore(t *testing.T) *storage.Store {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	st, err := storage.New()
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	return st
}

func TestGet_ActiveProfile(t *testing.T) {
	st := newIsolatedStore(t)
	if err := st.Add("personal", "Alice", "alice@personal.dev", "", "", "alice-gh"); err != nil {
		t.Fatal(err)
	}
	if err := st.SetActive("personal"); err != nil {
		t.Fatal(err)
	}

	var gotHost, gotUser string
	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		gotHost, gotUser = host, ghUser
		return "ghp_token123", nil
	})

	var b strings.Builder
	req := Request{Protocol: "https", Host: "github.com"}
	if err := Get(req, st, "", &b); err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if gotUser != "alice-gh" {
		t.Errorf("fetcher called with ghUser %q, want alice-gh", gotUser)
	}
	if gotHost != "github.com" {
		t.Errorf("fetcher called with host %q, want github.com", gotHost)
	}
	want := "protocol=https\nhost=github.com\nusername=alice-gh\npassword=ghp_token123\n\n"
	if b.String() != want {
		t.Errorf("output =\n%q\nwant\n%q", b.String(), want)
	}
}

func TestGet_PinnedRepoOverridesActive(t *testing.T) {
	st := newIsolatedStore(t)
	_ = st.Add("personal", "Alice", "alice@personal.dev", "", "", "alice-gh")
	_ = st.Add("work", "Alice W", "alice@work.com", "", "", "work-gh")
	_ = st.SetActive("personal")

	repoKey := "git@github.com:acme/widget.git"
	if err := history.Pin(repoKey, "work"); err != nil {
		t.Fatal(err)
	}

	var gotUser string
	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		gotUser = ghUser
		return "ghp_work", nil
	})

	var b strings.Builder
	req := Request{Protocol: "https", Host: "github.com"}
	if err := Get(req, st, repoKey, &b); err != nil {
		t.Fatal(err)
	}

	if gotUser != "work-gh" {
		t.Errorf("fetcher called with ghUser %q, want work-gh (pinned profile)", gotUser)
	}
	if !strings.Contains(b.String(), "username=work-gh") || !strings.Contains(b.String(), "password=ghp_work") {
		t.Errorf("output did not use pinned profile: %q", b.String())
	}
}

func TestGet_NoTokenWritesNothing(t *testing.T) {
	st := newIsolatedStore(t)
	_ = st.Add("personal", "Alice", "alice@personal.dev", "", "", "alice-gh")
	_ = st.SetActive("personal")

	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		return "", nil // account not authenticated for this host
	})

	var b strings.Builder
	req := Request{Protocol: "https", Host: "github.com"}
	if err := Get(req, st, "", &b); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 0 {
		t.Errorf("expected empty output, got %q", b.String())
	}
}

func TestGet_NoActiveProfileWritesNothing(t *testing.T) {
	st := newIsolatedStore(t) // empty store, nothing active

	called := false
	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		called = true
		return "tok", nil
	})

	var b strings.Builder
	req := Request{Protocol: "https", Host: "github.com"}
	if err := Get(req, st, "", &b); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 0 {
		t.Errorf("expected empty output, got %q", b.String())
	}
	if called {
		t.Error("tokenFetcher should not be called when no profile resolves")
	}
}

func TestGet_ActiveProfileEmptyGHUserWritesNothing(t *testing.T) {
	st := newIsolatedStore(t)
	_ = st.Add("personal", "Alice", "alice@personal.dev", "", "", "") // no gh_user
	_ = st.SetActive("personal")

	called := false
	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		called = true
		return "tok", nil
	})

	var b strings.Builder
	if err := Get(Request{Protocol: "https", Host: "github.com"}, st, "", &b); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 0 || called {
		t.Errorf("expected silent passthrough; output=%q called=%v", b.String(), called)
	}
}

func TestResolveGHUser_HostPortStrippedForFetch(t *testing.T) {
	st := newIsolatedStore(t)
	_ = st.Add("personal", "Alice", "alice@personal.dev", "", "", "alice-gh")
	_ = st.SetActive("personal")

	var gotHost string
	withTokenFetcher(t, func(host, ghUser string) (string, error) {
		gotHost = host
		return "tok", nil
	})

	var b strings.Builder
	if err := Get(Request{Protocol: "https", Host: "github.com:443"}, st, "", &b); err != nil {
		t.Fatal(err)
	}
	if gotHost != "github.com" {
		t.Errorf("fetcher host = %q, want github.com (port stripped)", gotHost)
	}
}
