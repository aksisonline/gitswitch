//go:build windows

package history

// withLock is a no-op on Windows.  Shell hooks (and therefore concurrent
// record calls) are not supported on Windows, so no advisory lock is needed.
func withLock(_ string, fn func() error) error {
	return fn()
}
