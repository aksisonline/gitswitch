package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Shell int

const (
	ShellUnknown Shell = iota
	ShellZsh
	ShellBash
	ShellFish
)

type Framework int

const (
	FrameworkRaw Framework = iota
	FrameworkOMZ
	FrameworkP10k
	FrameworkStarship
)

const marker = "# gitswitch shell integration"

// DetectShell infers the current shell from $SHELL.
func DetectShell() Shell {
	s := filepath.Base(os.Getenv("SHELL"))
	switch s {
	case "zsh":
		return ShellZsh
	case "bash":
		return ShellBash
	case "fish":
		return ShellFish
	default:
		return ShellUnknown
	}
}

// DetectFramework inspects environment variables to identify the prompt framework.
func DetectFramework() Framework {
	if os.Getenv("STARSHIP_SHELL") != "" {
		return FrameworkStarship
	}
	if os.Getenv("POWERLEVEL9K_MODE") != "" || strings.HasPrefix(os.Getenv("POWERLEVEL9K_LEFT_PROMPT_ELEMENTS"), "P") {
		return FrameworkP10k
	}
	// P10k also sets POWERLEVEL9K_ prefixed vars; check theme name as fallback
	if strings.Contains(os.Getenv("ZSH_THEME"), "powerlevel") {
		return FrameworkP10k
	}
	if os.Getenv("ZSH") != "" || os.Getenv("ZSH_THEME") != "" {
		return FrameworkOMZ
	}
	return FrameworkRaw
}

// RCFile returns the shell rc file path for the given shell.
func RCFile(sh Shell) string {
	home, _ := os.UserHomeDir()
	switch sh {
	case ShellZsh:
		return filepath.Join(home, ".zshrc")
	case ShellBash:
		rc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(rc); err == nil {
			return rc
		}
		return filepath.Join(home, ".bash_profile")
	case ShellFish:
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return filepath.Join(home, ".bashrc")
	}
}

// IsInstalled checks whether the gitswitch marker exists in the given file.
func IsInstalled(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), marker)
}

// nudgeSnippetZsh returns the zsh nudge + prompt + completion snippet.
func nudgeSnippetZsh() string {
	return `
` + marker + ` begin
__gitswitch_prompt() {
  git rev-parse --git-dir > /dev/null 2>&1 || return
  local info nick color
  info=$(gitswitch current --prompt 2>/dev/null)
  [[ -z "$info" ]] && return
  nick=$(echo "$info" | cut -f1)
  color=$(echo "$info" | cut -f2)
  [[ -z "$color" ]] && color=141
  echo "%F{$color}[${nick}]%f"
}

__gitswitch_nudge() {
  local root
  root=$(git rev-parse --show-toplevel 2>/dev/null) || { unset __GITSWITCH_LAST_REPO; return; }
  [[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return
  export __GITSWITCH_LAST_REPO="$root"
  gitswitch record --path "$root" 2>/dev/null
  local result
  result=$(gitswitch recommend 2>/dev/null) || return
  local nickname name email
  IFS=$'\t' read -r nickname name email <<< "$result"
  printf "gitswitch: this repo usually uses %s <%s> — switch? [y/N] " "$name" "$email"
  read -k 1 -s reply; echo
  [[ "$reply" =~ ^[Yy]$ ]] && gitswitch switch "$nickname"
}
autoload -Uz add-zsh-hook
add-zsh-hook chpwd __gitswitch_nudge
__gitswitch_nudge
PROMPT='$(__gitswitch_prompt)'"$PROMPT"
autoload -U compinit; compinit
source <(gitswitch completion zsh)
` + marker + ` end
`
}

// nudgeSnippetBash returns the bash nudge + prompt + completion snippet.
func nudgeSnippetBash() string {
	return `
` + marker + ` begin
__gitswitch_prompt() {
  git rev-parse --git-dir > /dev/null 2>&1 || return
  local info nick color
  info=$(gitswitch current --prompt 2>/dev/null)
  [[ -z "$info" ]] && return
  nick=$(echo "$info" | cut -f1)
  color=$(echo "$info" | cut -f2)
  [[ -z "$color" ]] && color=141
  printf '\[\e[38;5;%sm\][%s]\[\e[0m\] ' "$color" "$nick"
}

__gitswitch_nudge() {
  # short-circuit when $PWD hasn't changed to avoid git rev-parse on every prompt
  [[ "$PWD" == "$__GITSWITCH_LAST_PWD" ]] && return
  export __GITSWITCH_LAST_PWD="$PWD"
  local root
  root=$(git rev-parse --show-toplevel 2>/dev/null) || { unset __GITSWITCH_LAST_REPO; return; }
  [[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return
  export __GITSWITCH_LAST_REPO="$root"
  gitswitch record --path "$root" 2>/dev/null
  local result
  result=$(gitswitch recommend 2>/dev/null) || return
  local nickname name email
  IFS=$'\t' read -r nickname name email <<< "$result"
  printf "gitswitch: this repo usually uses %s <%s> — switch? [y/N] " "$name" "$email"
  read -n 1 -s reply; echo
  [[ "$reply" =~ ^[Yy]$ ]] && gitswitch switch "$nickname"
}
PS1='$(__gitswitch_prompt)'"$PS1"
PROMPT_COMMAND="__gitswitch_nudge${PROMPT_COMMAND:+; $PROMPT_COMMAND}"
__gitswitch_nudge
source <(gitswitch completion bash)
` + marker + ` end
`
}

// nudgeSnippetFish returns the fish nudge + cd hook + completion snippet.
func nudgeSnippetFish() string {
	return `
` + marker + ` begin
function __gitswitch_prompt
  git rev-parse --git-dir > /dev/null 2>&1; or return
  set info (gitswitch current --prompt 2>/dev/null)
  test -z "$info"; and return
  set parts (string split \t $info)
  set nick $parts[1]
  set color $parts[2]
  test -z "$color"; and set color 141
  printf '\e[38;5;%sm[%s]\e[0m' $color $nick
end

function __gitswitch_nudge
  set root (git rev-parse --show-toplevel 2>/dev/null); or begin; set -e __GITSWITCH_LAST_REPO; return; end
  if test "$root" = "$__GITSWITCH_LAST_REPO"
    return
  end
  set -gx __GITSWITCH_LAST_REPO $root
  gitswitch record --path $root 2>/dev/null
  set result (gitswitch recommend 2>/dev/null); or return
  set parts (string split \t $result)
  set nickname $parts[1]
  set name $parts[2]
  set email $parts[3]
  printf "gitswitch: this repo usually uses %s <%s> — switch? [y/N] " $name $email
  read -n 1 reply
  if string match -qi 'y' $reply
    gitswitch switch $nickname
  end
end

function __gitswitch_cd_hook --on-variable PWD
  __gitswitch_nudge
end
__gitswitch_nudge
function fish_right_prompt
  __gitswitch_prompt
end
gitswitch completion fish | source
` + marker + ` end
`
}

// starshipSnippet returns the starship.toml custom module block.
func starshipSnippet() string {
	return `
[custom.gitswitch]
command = "gitswitch current --short"
when = "git rev-parse --git-dir > /dev/null 2>&1"
symbol = " "
style = "bold cyan"
format = "[$symbol($output)]($style) "
`
}

// omzPluginContent returns the oh-my-zsh plugin file content.
func omzPluginContent() string {
	return `# gitswitch oh-my-zsh plugin
__gitswitch_prompt() {
  git rev-parse --git-dir > /dev/null 2>&1 || return
  local info nick color
  info=$(gitswitch current --prompt 2>/dev/null)
  [[ -z "$info" ]] && return
  nick=$(echo "$info" | cut -f1)
  color=$(echo "$info" | cut -f2)
  [[ -z "$color" ]] && color=141
  echo "%F{$color}[${nick}]%f"
}

__gitswitch_nudge() {
  local root
  root=$(git rev-parse --show-toplevel 2>/dev/null) || { unset __GITSWITCH_LAST_REPO; return; }
  [[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return
  export __GITSWITCH_LAST_REPO="$root"
  gitswitch record --path "$root" 2>/dev/null
  local result
  result=$(gitswitch recommend 2>/dev/null) || return
  local nickname name email
  IFS=$'\t' read -r nickname name email <<< "$result"
  printf "gitswitch: this repo usually uses %s <%s> — switch? [y/N] " "$name" "$email"
  read -k 1 -s reply; echo
  [[ "$reply" =~ ^[Yy]$ ]] && gitswitch switch "$nickname"
}
autoload -Uz add-zsh-hook
add-zsh-hook chpwd __gitswitch_nudge
__gitswitch_nudge
PROMPT='$(__gitswitch_prompt)'"$PROMPT"
autoload -U compinit; compinit
source <(gitswitch completion zsh)
`
}

// p10kSnippet returns the P10k segment function for manual insertion.
func p10kSnippet() string {
	return `
` + marker + ` begin
function prompt_gitswitch() {
  git rev-parse --git-dir > /dev/null 2>&1 || return
  local info nick color
  info=$(gitswitch current --prompt 2>/dev/null)
  [[ -z "$info" ]] && return
  nick=$(echo "$info" | cut -f1)
  color=$(echo "$info" | cut -f2)
  [[ -z "$color" ]] && color=141
  p10k segment -f "$color" -t "[$nick]"
}

__gitswitch_nudge() {
  local root
  root=$(git rev-parse --show-toplevel 2>/dev/null) || { unset __GITSWITCH_LAST_REPO; return; }
  [[ "$root" == "$__GITSWITCH_LAST_REPO" ]] && return
  export __GITSWITCH_LAST_REPO="$root"
  gitswitch record --path "$root" 2>/dev/null
  local result
  result=$(gitswitch recommend 2>/dev/null) || return
  local nickname name email
  IFS=$'\t' read -r nickname name email <<< "$result"
  printf "gitswitch: this repo usually uses %s <%s> — switch? [y/N] " "$name" "$email"
  read -k 1 -s reply; echo
  [[ "$reply" =~ ^[Yy]$ ]] && gitswitch switch "$nickname"
}
autoload -Uz add-zsh-hook
add-zsh-hook chpwd __gitswitch_nudge
__gitswitch_nudge
autoload -U compinit; compinit
source <(gitswitch completion zsh)
` + marker + ` end
`
}

// Install writes the appropriate integration for the detected framework.
// Returns a human-readable description of what was done.
func Install(sh Shell, fw Framework) (string, error) {
	home, _ := os.UserHomeDir()

	switch fw {
	case FrameworkStarship:
		return installStarship(home)
	case FrameworkOMZ:
		return installOMZ(home)
	case FrameworkP10k:
		return installP10k(sh, home)
	default:
		return installRaw(sh)
	}
}

func installStarship(home string) (string, error) {
	tomlPath := filepath.Join(home, ".config", "starship.toml")
	if IsInstalled(tomlPath) {
		return fmt.Sprintf("already installed in %s", tomlPath), nil
	}
	// Ensure file exists
	if err := os.MkdirAll(filepath.Dir(tomlPath), 0755); err != nil {
		return "", err
	}
	f, err := os.OpenFile(tomlPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	snippet := marker + " begin\n" + starshipSnippet() + "\n" + marker + " end\n"
	if _, err := f.WriteString(snippet); err != nil {
		return "", err
	}
	return fmt.Sprintf("wrote [custom.gitswitch] block to %s", tomlPath), nil
}

func installOMZ(home string) (string, error) {
	pluginDir := filepath.Join(home, ".oh-my-zsh", "custom", "plugins", "gitswitch")
	pluginFile := filepath.Join(pluginDir, "gitswitch.plugin.zsh")
	if _, err := os.Stat(pluginFile); err == nil {
		return fmt.Sprintf("plugin already exists at %s", pluginFile), nil
	}
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(pluginFile, []byte(omzPluginContent()), 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("wrote plugin to %s\n  → add 'gitswitch' to the plugins array in ~/.zshrc", pluginFile), nil
}

func installP10k(sh Shell, home string) (string, error) {
	rcFile := RCFile(sh)
	if IsInstalled(rcFile) {
		return fmt.Sprintf("nudge hook already installed in %s", rcFile), nil
	}
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(p10kSnippet()); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"wrote nudge hook to %s\n  → for the prompt segment, add 'gitswitch' to POWERLEVEL9K_LEFT_PROMPT_ELEMENTS in ~/.p10k.zsh",
		rcFile,
	), nil
}

// Uninstall removes the gitswitch marker block from the rc file for the
// provided shell, removes it from starship.toml if present, and removes the
// OMZ plugin file if present.
// Returns a human-readable description of what was removed.
func Uninstall(sh Shell, _ Framework) (string, error) {
	home, _ := os.UserHomeDir()
	var removed []string

	// rc file (raw, p10k, bash)
	rcFile := RCFile(sh)
	if IsInstalled(rcFile) {
		if err := removeMarkerBlock(rcFile); err != nil {
			return "", fmt.Errorf("could not clean %s: %w", rcFile, err)
		}
		removed = append(removed, rcFile)
	}

	// starship.toml
	tomlPath := filepath.Join(home, ".config", "starship.toml")
	if IsInstalled(tomlPath) {
		if err := removeMarkerBlock(tomlPath); err != nil {
			return "", fmt.Errorf("could not clean %s: %w", tomlPath, err)
		}
		removed = append(removed, tomlPath)
	}

	// OMZ plugin file
	omzPlugin := filepath.Join(home, ".oh-my-zsh", "custom", "plugins", "gitswitch", "gitswitch.plugin.zsh")
	if _, err := os.Stat(omzPlugin); err == nil {
		if err := os.Remove(omzPlugin); err != nil {
			return "", fmt.Errorf("could not remove %s: %w", omzPlugin, err)
		}
		removed = append(removed, omzPlugin)
	}

	if len(removed) == 0 {
		return "nothing to remove — gitswitch shell integration was not installed", nil
	}
	return fmt.Sprintf("removed gitswitch integration from: %s", strings.Join(removed, ", ")), nil
}

// removeMarkerBlock strips the lines between (and including) the begin/end
// marker lines from path, writing the result atomically via a temp file + rename
// with the original file mode preserved.
func removeMarkerBlock(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	inside := false
	found := false
	for _, line := range lines {
		if strings.Contains(line, marker+" begin") {
			inside = true
			found = true
			continue
		}
		if strings.Contains(line, marker+" end") {
			if !inside {
				return fmt.Errorf("unbalanced marker in %s: found end without begin", path)
			}
			inside = false
			continue
		}
		if !inside {
			out = append(out, line)
		}
	}
	if inside {
		return fmt.Errorf("unbalanced marker in %s: begin without end", path)
	}
	if !found {
		return nil
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".gitswitch-uninstall-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.WriteString(strings.Join(out, "\n")); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Chmod(tmpName, info.Mode()); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

func installRaw(sh Shell) (string, error) {
	rcFile := RCFile(sh)
	if IsInstalled(rcFile) {
		return fmt.Sprintf("already installed in %s", rcFile), nil
	}

	var snippet string
	switch sh {
	case ShellZsh:
		snippet = nudgeSnippetZsh()
	case ShellBash:
		snippet = nudgeSnippetBash()
	case ShellFish:
		snippet = nudgeSnippetFish()
	default:
		snippet = nudgeSnippetBash()
	}

	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(snippet); err != nil {
		return "", err
	}
	return fmt.Sprintf("wrote gitswitch integration to %s", rcFile), nil
}
