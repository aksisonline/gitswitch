package credential

import (
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHost string
		wantProt string
		wantPath string
		wantUser string
		wantNoPort string
	}{
		{
			name:     "standard four keys",
			input:    "protocol=https\nhost=github.com\npath=org/repo.git\nusername=alice\n\n",
			wantProt: "https", wantHost: "github.com", wantPath: "org/repo.git", wantUser: "alice",
			wantNoPort: "github.com",
		},
		{
			name:     "blank line terminates before EOF junk",
			input:    "protocol=https\nhost=github.com\n\nprotocol=ignored\n",
			wantProt: "https", wantHost: "github.com",
			wantNoPort: "github.com",
		},
		{
			name:     "EOF without blank line",
			input:    "protocol=https\nhost=github.com\n",
			wantProt: "https", wantHost: "github.com",
			wantNoPort: "github.com",
		},
		{
			name:     "host with port stripped",
			input:    "protocol=https\nhost=github.com:443\n\n",
			wantProt: "https", wantHost: "github.com:443",
			wantNoPort: "github.com",
		},
		{
			name:     "unknown keys ignored",
			input:    "protocol=https\nhost=github.com\nwwwauth[]=Basic\nquit=0\n\n",
			wantProt: "https", wantHost: "github.com",
			wantNoPort: "github.com",
		},
		{
			name:     "value containing equals sign",
			input:    "password=ab=cd\nhost=github.com\n\n",
			wantHost: "github.com",
			wantNoPort: "github.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := ParseRequest(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ParseRequest error: %v", err)
			}
			if req.Protocol != tt.wantProt {
				t.Errorf("Protocol = %q, want %q", req.Protocol, tt.wantProt)
			}
			if req.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", req.Host, tt.wantHost)
			}
			if req.Path != tt.wantPath {
				t.Errorf("Path = %q, want %q", req.Path, tt.wantPath)
			}
			if req.Username != tt.wantUser {
				t.Errorf("Username = %q, want %q", req.Username, tt.wantUser)
			}
			if got := req.HostNoPort(); got != tt.wantNoPort {
				t.Errorf("HostNoPort() = %q, want %q", got, tt.wantNoPort)
			}
		})
	}
}

func TestHostNoPort(t *testing.T) {
	tests := []struct{ in, want string }{
		{"github.com", "github.com"},
		{"github.com:443", "github.com"},
		{"github.corp.com:8443", "github.corp.com"},
		{"[::1]", "[::1]"},
		{"[::1]:443", "[::1]"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := (Request{Host: tt.in}).HostNoPort(); got != tt.want {
			t.Errorf("HostNoPort(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestWriteResponse(t *testing.T) {
	var b strings.Builder
	req := Request{Protocol: "https", Host: "github.com:443"}
	writeResponse(&b, req, "alice", "ghp_secret")
	want := "protocol=https\nhost=github.com:443\nusername=alice\npassword=ghp_secret\n\n"
	if b.String() != want {
		t.Errorf("writeResponse =\n%q\nwant\n%q", b.String(), want)
	}
}

func TestWriteResponseOmitsEmptyProtocolHost(t *testing.T) {
	var b strings.Builder
	writeResponse(&b, Request{}, "alice", "tok")
	want := "username=alice\npassword=tok\n\n"
	if b.String() != want {
		t.Errorf("writeResponse =\n%q\nwant\n%q", b.String(), want)
	}
}

// withTokenFetcher swaps the package-level tokenFetcher for the duration of a
// test and restores it afterward.
func withTokenFetcher(t *testing.T, fn func(host, ghUser string) (string, error)) {
	t.Helper()
	orig := tokenFetcher
	tokenFetcher = fn
	t.Cleanup(func() { tokenFetcher = orig })
}
