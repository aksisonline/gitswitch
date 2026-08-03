package main

import "testing"

func TestCommandAllowsMissingGit(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{"bare invocation", []string{"gitswitch"}, true},
		{"doctor", []string{"gitswitch", "doctor"}, true},
		{"add", []string{"gitswitch", "add"}, false},
	}
	for _, c := range cases {
		if got := commandAllowsMissingGit(c.args); got != c.want {
			t.Errorf("%s: commandAllowsMissingGit(%v) = %v, want %v", c.name, c.args, got, c.want)
		}
	}
}
