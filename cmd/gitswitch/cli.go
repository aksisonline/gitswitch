package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aksisonline/gitswitch/internal/storage"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// This file is the command line's "frame": help output, the version flag,
// argument and flag validation, shell completion, and a spinner for the
// commands that wait on the network. The commands themselves live in main.go,
// channels.go, guard.go and xcode.go.

const docsURL = "https://gitswitch.dev/docs"

// usageError marks a mistake in how a command was called (wrong arguments, an
// unknown flag, a bad flag value) as opposed to the command itself failing.
// main prints a pointer to that command's --help and exits 2, the conventional
// status for usage errors.
type usageError struct {
	err error
	cmd *cobra.Command
}

func (e *usageError) Error() string { return e.err.Error() }
func (e *usageError) Unwrap() error { return e.err }

func usageErrorf(cmd *cobra.Command, format string, a ...any) error {
	return &usageError{err: fmt.Errorf(format, a...), cmd: cmd}
}

// argRange validates the positional argument count and, when it is wrong,
// says what the command expects in the words of its own usage line instead of
// cobra's "accepts 3 arg(s), received 1".
func argRange(minArgs, maxArgs int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= minArgs && len(args) <= maxArgs {
			return nil
		}
		expected := strings.TrimSpace(strings.TrimPrefix(cmd.Use, cmd.Name()))
		switch {
		case maxArgs == 0:
			return usageErrorf(cmd, "'%s' takes no arguments (got %d)", cmd.CommandPath(), len(args))
		case len(args) < minArgs:
			return usageErrorf(cmd, "missing arguments — expected: %s %s", cmd.CommandPath(), expected)
		default:
			return usageErrorf(cmd, "too many arguments — expected: %s %s", cmd.CommandPath(), expected)
		}
	}
}

func exactArgs(n int) cobra.PositionalArgs { return argRange(n, n) }

var noArgs = argRange(0, 0)

// oneOf validates an enumerated flag value.
func oneOf(cmd *cobra.Command, flag, value string, allowed ...string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	var named []string
	for _, a := range allowed {
		if a != "" {
			named = append(named, a)
		}
	}
	return usageErrorf(cmd, "invalid --%s %q — must be one of: %s", flag, value, strings.Join(named, ", "))
}

// completeNicknames offers stored profile nicknames for shell completion.
func completeNicknames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 || store == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return nicknameCompletions(), cobra.ShellCompDirectiveNoFileComp
}

func nicknameCompletions() []string {
	profiles, err := store.Load()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, p.Nickname+"\t"+p.Email)
	}
	return out
}

func fixedCompletions(values ...string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return values, cobra.ShellCompDirectiveNoFileComp
	}
}

// unknownTargetError explains a bare `gitswitch <word>` that is neither a
// command nor a profile, with the closest matches of each.
func unknownTargetError(root *cobra.Command, word string, profiles []storage.Profile) error {
	var hints []string
	hints = append(hints, root.SuggestionsFor(word)...)
	for _, p := range profiles {
		if strings.HasPrefix(strings.ToLower(p.Nickname), strings.ToLower(word)) || levenshtein(p.Nickname, word) <= 2 {
			hints = append(hints, p.Nickname)
		}
	}
	msg := fmt.Sprintf("no profile or command named %q", word)
	if len(hints) > 0 {
		msg += " — did you mean: " + strings.Join(hints, ", ") + "?"
	}
	return usageErrorf(root, "%s", msg)
}

// levenshtein is the edit distance between a and b (case-insensitive).
func levenshtein(a, b string) int {
	ra, rb := []rune(strings.ToLower(a)), []rune(strings.ToLower(b))
	prev := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur := make([]int, len(rb)+1)
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(rb)]
}

// ── help ─────────────────────────────────────────────────────────────────────

var (
	helpHeading = lipgloss.NewStyle().Bold(true).Foreground(cPrimary)
	helpCommand = lipgloss.NewStyle().Foreground(cAccent)
	helpBrand   = lipgloss.NewStyle().Bold(true).Foreground(cPrimary)
)

// helpRows renders name/description pairs as an aligned two-column block.
func helpRows(rows [][2]string) string {
	width := 0
	for _, r := range rows {
		width = max(width, lipgloss.Width(r[0]))
	}
	var b strings.Builder
	for _, r := range rows {
		pad := strings.Repeat(" ", width-lipgloss.Width(r[0]))
		fmt.Fprintf(&b, "  %s%s   %s\n", helpCommand.Render(r[0]), pad, r[1])
	}
	return b.String()
}

func helpSection(b *strings.Builder, title, body string) {
	if strings.TrimSpace(body) == "" {
		return
	}
	fmt.Fprintf(b, "%s\n%s\n", helpHeading.Render(title), strings.TrimRight(body, "\n")+"\n")
}

// flagRows lists a flag set as name/description pairs, hidden flags excluded.
func flagRows(fs *pflag.FlagSet) [][2]string {
	var rows [][2]string
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		name := "    --" + f.Name
		if f.Shorthand != "" {
			name = "-" + f.Shorthand + ", --" + f.Name
		}
		if f.Value.Type() != "bool" {
			name += " <" + f.Value.Type() + ">"
		}
		usage := f.Usage
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "[]" && !strings.Contains(usage, "default") {
			usage += fmt.Sprintf(" (default %s)", f.DefValue)
		}
		rows = append(rows, [2]string{name, usage})
	})
	return rows
}

// visibleCommands returns cmd's runnable, non-hidden subcommands.
func visibleCommands(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, c := range cmd.Commands() {
		if c.IsAvailableCommand() || c.Name() == "help" {
			out = append(out, c)
		}
	}
	return out
}

// rootHelp is what `gitswitch --help` prints: the three ways to call it, the
// commands by group, the global flags, and a few things to try.
func rootHelp(cmd *cobra.Command) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n%s.\n\n", helpBrand.Render("gitswitch"), styleDim(version), cmd.Short)

	helpSection(&b, "Usage", helpRows([][2]string{
		{"gitswitch", "open the profile picker"},
		{"gitswitch <nickname>", "switch to a profile"},
		{"gitswitch <command> [flags]", "run a command"},
	}))

	byGroup := map[string][][2]string{}
	for _, c := range visibleCommands(cmd) {
		byGroup[c.GroupID] = append(byGroup[c.GroupID], [2]string{c.Name(), c.Short})
	}
	for _, g := range cmd.Groups() {
		helpSection(&b, strings.TrimSuffix(g.Title, ":"), helpRows(byGroup[g.ID]))
	}
	helpSection(&b, "Other", helpRows(byGroup[""]))

	helpSection(&b, "Flags", helpRows([][2]string{
		{"-h, --help", "Show help for gitswitch or any command"},
		{"-v, --version", "Show the installed version"},
	}))

	helpSection(&b, "Examples", helpRows([][2]string{
		{"gitswitch login", "connect a GitHub account"},
		{"gitswitch work", "switch to the profile nicknamed \"work\""},
		{"gitswitch pin work", "always use \"work\" in the current repo"},
		{"gitswitch shell", "set up prompt, auto-switch and the commit guard"},
		{"gitswitch doctor", "check that everything is wired up"},
	}))

	fmt.Fprintf(&b, "%s\n", styleDim("Run 'gitswitch <command> --help' for a command's flags and examples."))
	fmt.Fprintf(&b, "%s\n", styleDim("Docs: "+docsURL))
	return b.String()
}

// commandHelp is `gitswitch <command> --help`.
func commandHelp(cmd *cobra.Command) string {
	var b strings.Builder
	desc := cmd.Long
	if desc == "" {
		desc = cmd.Short
	}
	fmt.Fprintf(&b, "%s\n\n", strings.TrimSpace(desc))

	usage := "  " + cmd.UseLine()
	if cmd.HasAvailableSubCommands() {
		usage += "\n  " + cmd.CommandPath() + " <command> [flags]"
	}
	helpSection(&b, "Usage", usage)

	if len(cmd.Aliases) > 0 {
		helpSection(&b, "Aliases", "  "+cmd.Name()+", "+strings.Join(cmd.Aliases, ", "))
	}

	var subs [][2]string
	for _, c := range visibleCommands(cmd) {
		if c.Name() != "help" {
			subs = append(subs, [2]string{c.Name(), c.Short})
		}
	}
	helpSection(&b, "Commands", helpRows(subs))

	flags := flagRows(cmd.LocalFlags())
	sort.SliceStable(flags, func(i, j int) bool { // help last
		return !strings.Contains(flags[i][0], "--help") && strings.Contains(flags[j][0], "--help")
	})
	helpSection(&b, "Flags", helpRows(flags))

	if cmd.Example != "" {
		helpSection(&b, "Examples", cmd.Example)
	}
	fmt.Fprintf(&b, "%s\n", styleDim("Docs: "+docsURL+"/cli/commands"))
	return b.String()
}

// installHelp wires the custom help and usage output into cobra.
func installHelp(root *cobra.Command) {
	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		// Make sure -h/--help shows up in the command's own flag list.
		cmd.InitDefaultHelpFlag()
		if f := cmd.Flags().Lookup("help"); f != nil {
			f.Usage = "Show help for this command"
		}
		if cmd == root {
			fmt.Fprint(cmd.OutOrStdout(), rootHelp(cmd))
			return
		}
		fmt.Fprint(cmd.OutOrStdout(), commandHelp(cmd))
	})
	// Usage is only shown for usage errors, and main prints a one-line pointer
	// for those — never the whole help text in the middle of an error.
	root.SetUsageFunc(func(cmd *cobra.Command) error { return nil })
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return &usageError{err: err, cmd: cmd}
	})
}

// normalizeArgs maps the version spellings people try — `--v`, `-V` — onto
// `--version`, so all of them work as the first argument.
func normalizeArgs(args []string) {
	if len(args) == 2 && (args[1] == "--v" || args[1] == "-V") {
		args[1] = "--version"
	}
}

// reportError prints err the way every gitswitch command does and returns the
// process exit status: 2 for a usage mistake, 1 for anything else.
func reportError(err error) int {
	if err == nil {
		return 0
	}
	if msg := err.Error(); msg != "" {
		fmt.Fprintln(os.Stderr, styleErr("✗ "+msg))
	}
	var ue *usageError
	if errors.As(err, &ue) {
		path := "gitswitch"
		if ue.cmd != nil {
			path = ue.cmd.CommandPath()
		}
		fmt.Fprintln(os.Stderr, styleDim(fmt.Sprintf("  Run '%s --help' for usage.", path)))
		return 2
	}
	// cobra's own "unknown command" / "unknown flag" errors are usage errors too.
	if strings.HasPrefix(err.Error(), "unknown ") {
		fmt.Fprintln(os.Stderr, styleDim("  Run 'gitswitch --help' for usage."))
		return 2
	}
	return 1
}

// ── installs on PATH ─────────────────────────────────────────────────────────

// installsOnPath returns every distinct gitswitch binary reachable through
// PATH, in PATH order (the first is the one that runs). More than one means
// two installs are shadowing each other — typically a script-installed copy in
// /usr/local/bin sitting next to a Homebrew one, which then appears to
// "never update" because the stale copy keeps winning.
func installsOnPath() []string {
	name := "gitswitch"
	if os.PathSeparator == '\\' {
		name += ".exe"
	}
	var out []string
	seen := map[string]bool{}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == "" {
			continue
		}
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil || info.IsDir() || info.Mode()&0o111 == 0 {
			continue
		}
		real := path
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			real = resolved
		}
		if !seen[real] {
			seen[real] = true
			out = append(out, real)
		}
	}
	return out
}

// ── spinner ──────────────────────────────────────────────────────────────────

// withSpinner runs fn while showing a spinner and label on stderr. When stderr
// is not a terminal (a pipe, CI, an agent) it prints nothing at all, so
// captured output stays clean.
func withSpinner(label string, fn func() error) error {
	stat, err := os.Stderr.Stat()
	if err != nil || stat.Mode()&os.ModeCharDevice == 0 {
		return fn()
	}
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		tick := time.NewTicker(80 * time.Millisecond)
		defer tick.Stop()
		for i := 0; ; i++ {
			fmt.Fprintf(os.Stderr, "\r%s %s", lipgloss.NewStyle().Foreground(cPrimary).Render(frames[i%len(frames)]), label)
			select {
			case <-done:
				fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", lipgloss.Width(label)+2))
				return
			case <-tick.C:
			}
		}
	}()
	err = fn()
	close(done)
	<-finished
	return err
}
