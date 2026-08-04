package secrets

import (
	"errors"
	"os"

	"github.com/zalando/go-keyring"
)

const service = "gitswitch"

// Store is the interface for token storage backends.
type Store interface {
	Set(ref, token string) error
	Available() bool
}

// Default returns the best available Store for the current machine.
// Prefers the OS keychain; falls back to nothing if unavailable.
// Callers can force a specific backend via GITSWITCH_SECRETS_BACKEND env var.
func Default() Store {
	if os.Getenv("GITSWITCH_SECRETS_BACKEND") == "none" {
		return &noopStore{}
	}
	k := &KeychainStore{}
	if k.Available() {
		return k
	}
	return &noopStore{}
}

// KeychainStore uses the OS keychain via go-keyring.
// On macOS: system Keychain (not iCloud — go-keyring never sets kSecAttrSynchronizable).
// On Linux: libsecret / GNOME Keyring / KeePassXC via D-Bus.
// On Windows: Windows Credential Manager.
type KeychainStore struct{}

func (k *KeychainStore) Available() bool {
	// Probe with a harmless write+delete.
	probe := "__gitswitch_probe__"
	if err := keyring.Set(service, probe, probe); err != nil {
		return false
	}
	_ = keyring.Delete(service, probe)
	return true
}

func (k *KeychainStore) Set(ref, token string) error {
	return keyring.Set(service, ref, token)
}

// noopStore is a no-op fallback used when no backend is available.
type noopStore struct{}

func (n *noopStore) Available() bool       { return false }
func (n *noopStore) Set(_, _ string) error { return errors.New("no secret store available") }
