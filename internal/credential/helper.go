// Package credential implements a git credential helper for gitswitch.
//
// gitswitch never stores tokens. It acts purely as a router: given a git
// credential request (protocol/host), it resolves which gitswitch profile
// applies to the current repo, then delegates to the gh CLI to fetch that
// account's token (`gh auth token --hostname <host> --user <ghUser>`). The
// token lives in gh's own secure storage; gitswitch only forwards it to git.
//
// For any case it cannot serve (gh missing, no profile, account not
// authenticated for the host), the helper writes nothing and exits 0 so git
// falls through to the next helper in the chain (osxkeychain, gh, prompt).
package credential

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/aksisonline/gitswitch/internal/history"
	"github.com/aksisonline/gitswitch/internal/storage"
)

// Request is the parsed git credential protocol input.
type Request struct {
	Protocol string
	Host     string // may include :port; use HostNoPort() for matching/queries
	Path     string
	Username string
}

// HostNoPort returns Host with any trailing :port stripped
// (github.com:443 -> github.com). IPv6 hosts are left untouched.
func (r Request) HostNoPort() string {
	h := r.Host
	if strings.HasPrefix(h, "[") {
		// bracketed IPv6 literal, possibly with :port after the bracket
		if i := strings.LastIndex(h, "]"); i >= 0 {
			return h[:i+1]
		}
		return h
	}
	if i := strings.LastIndex(h, ":"); i >= 0 {
		return h[:i]
	}
	return h
}

// ParseRequest reads key=value lines from r until a blank line or EOF.
// Unknown keys are ignored. Lines without '=' are skipped.
func ParseRequest(r io.Reader) (Request, error) {
	var req Request
	sc := bufio.NewScanner(r)
	// git limits lines to 65535 bytes; allow generous buffer headroom.
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			break // blank line terminates the request
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		key := line[:eq]
		val := line[eq+1:]
		switch key {
		case "protocol":
			req.Protocol = val
		case "host":
			req.Host = val
		case "path":
			req.Path = val
		case "username":
			req.Username = val
		}
	}
	if err := sc.Err(); err != nil {
		return req, err
	}
	return req, nil
}

// tokenFetcher is the seam tests stub. Default shells out to `gh auth token`.
var tokenFetcher = ghAuthToken

// ghAuthToken runs: gh auth token --hostname <host> --user <ghUser>
// Returns ("", nil) on any failure (gh missing, account not authenticated for
// the host, etc.). It never returns a non-nil error — callers treat an empty
// token as "cannot serve this request" and stay silent.
func ghAuthToken(host, ghUser string) (string, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return "", nil
	}
	args := []string{"auth", "token", "--hostname", host}
	if ghUser != "" {
		args = append(args, "--user", ghUser)
	}
	out, err := exec.Command("gh", args...).Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}

// resolveGHUser picks which profile's gh_user to use for this request.
//
//  1. Determine the globally active profile.
//  2. Ask history for a per-repo recommendation (pinned/learned identity for
//     repoKey). history.Recommend returns ("",false) when the recommendation
//     equals the active nickname, so the active profile is the natural
//     fallback and the two branches converge on the same result.
//
// Returns "" if no profile applies (the caller then stays silent).
func resolveGHUser(req Request, st *storage.Store, repoKey string) string {
	active, _ := st.GetActive()
	currentNick := ""
	if active != nil {
		currentNick = active.Nickname
	}

	if nick, ok := history.Recommend(repoKey, currentNick); ok {
		if p, err := st.Get(nick); err == nil && p != nil {
			return p.GHUser
		}
	}

	if active != nil {
		return active.GHUser
	}
	return ""
}

// Get handles the `get`/`fill` operation: resolve the applicable profile,
// fetch its token via gh, and write a credential response to w. On any miss
// (no profile, gh unavailable, account not authed for host) it writes nothing
// so git falls through to the next helper. Always returns nil (exit 0).
func Get(req Request, st *storage.Store, repoKey string, w io.Writer) error {
	ghUser := resolveGHUser(req, st, repoKey)
	if ghUser == "" {
		return nil
	}

	token, _ := tokenFetcher(req.HostNoPort(), ghUser)
	if token == "" {
		return nil
	}

	writeResponse(w, req, ghUser, token)
	return nil
}

// writeResponse emits the git credential reply: protocol/host (echoed back if
// present), username, password, then the terminating blank line.
func writeResponse(w io.Writer, req Request, username, token string) {
	var b strings.Builder
	if req.Protocol != "" {
		fmt.Fprintf(&b, "protocol=%s\n", req.Protocol)
	}
	if req.Host != "" {
		fmt.Fprintf(&b, "host=%s\n", req.Host)
	}
	fmt.Fprintf(&b, "username=%s\n", username)
	fmt.Fprintf(&b, "password=%s\n", token)
	b.WriteByte('\n')
	io.WriteString(w, b.String())
}
