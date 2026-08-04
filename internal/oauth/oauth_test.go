package oauth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// withTestServer points httpClient at srv for the duration of the test,
// mirroring credential.tokenFetcher's swap-seam idiom.
func withTestServer(t *testing.T, srv *httptest.Server) {
	t.Helper()
	orig := httpClient
	httpClient = srv.Client()
	t.Cleanup(func() { httpClient = orig })
}

func TestFetchVerifiedEmailsUserEmailsHappyPath(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// apiBaseURL treats any non-"github.com" host as GHES and appends
		// /api/v3 — exercise that same path shape here.
		if r.URL.Path != "/api/v3/user/emails" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		fmt.Fprint(w, `[
			{"email": "unverified@example.com", "verified": false},
			{"email": "primary@example.com", "verified": true},
			{"email": "secondary@example.com", "verified": true}
		]`)
	}))
	defer srv.Close()
	withTestServer(t, srv)

	got := FetchVerifiedEmails("token", srv.Listener.Addr().String())
	want := map[string]bool{"primary@example.com": true, "secondary@example.com": true}
	if len(got) != len(want) {
		t.Fatalf("got %v, want 2 verified emails", got)
	}
	for _, e := range got {
		if !want[e] {
			t.Errorf("unexpected email %q in result %v", e, got)
		}
	}
}

func TestFetchVerifiedEmailsFallsBackToPublicUserEmail(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v3/user/emails":
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"message": "Requires user:email scope"}`)
		case "/api/v3/user":
			fmt.Fprint(w, `{"login": "alice", "email": "public@example.com"}`)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()
	withTestServer(t, srv)

	got := FetchVerifiedEmails("token", srv.Listener.Addr().String())
	if len(got) != 1 || got[0] != "public@example.com" {
		t.Fatalf("got %v, want [public@example.com] via /user fallback", got)
	}
}

func TestFetchVerifiedEmailsReturnsNilOnTotalFailure(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"message": "insufficient scope"}`)
	}))
	defer srv.Close()
	withTestServer(t, srv)

	if got := FetchVerifiedEmails("token", srv.Listener.Addr().String()); got != nil {
		t.Fatalf("expected nil on total failure, got %v", got)
	}
}
