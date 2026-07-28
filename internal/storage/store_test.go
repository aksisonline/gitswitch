package storage

import "testing"

func TestUpdatePreservesTokenRefWhenOmitted(t *testing.T) {
	s := NewAt(t.TempDir())
	if err := s.Add("work", "Alice", "a@work.com", "SIG", "~/.ssh/id_work", "alice-work"); err != nil {
		t.Fatal(err)
	}
	if err := s.Update("work", Profile{
		Nickname: "work",
		UserName: "Alice",
		Email:    "a@work.com",
		SignKey:  "SIG",
		SSHKey:   "~/.ssh/id_work",
		GHUser:   "alice-work",
		TokenRef: "gitswitch:work:github.com",
	}); err != nil {
		t.Fatal(err)
	}

	// TUI edit omits TokenRef — must not wipe the OAuth keychain pointer.
	if err := s.Update("work", Profile{
		Nickname: "work",
		UserName: "Alice A",
		Email:    "a@work.com",
		SignKey:  "SIG",
		SSHKey:   "~/.ssh/id_work",
		GHUser:   "alice-work",
	}); err != nil {
		t.Fatal(err)
	}
	p, err := s.Get("work")
	if err != nil {
		t.Fatal(err)
	}
	if p.TokenRef != "gitswitch:work:github.com" {
		t.Fatalf("TokenRef not preserved: got %q", p.TokenRef)
	}
	if p.UserName != "Alice A" {
		t.Fatalf("UserName not updated: got %q", p.UserName)
	}
}

func TestUpdateCanReplaceTokenRefExplicitly(t *testing.T) {
	s := NewAt(t.TempDir())
	_ = s.Add("work", "Alice", "a@work.com", "", "", "alice")
	_ = s.Update("work", Profile{
		Nickname: "work", UserName: "Alice", Email: "a@work.com",
		GHUser: "alice", TokenRef: "old-ref",
	})
	_ = s.Update("work", Profile{
		Nickname: "work", UserName: "Alice", Email: "a@work.com",
		GHUser: "alice", TokenRef: "new-ref",
	})
	p, _ := s.Get("work")
	if p.TokenRef != "new-ref" {
		t.Fatalf("expected new-ref, got %q", p.TokenRef)
	}
}
