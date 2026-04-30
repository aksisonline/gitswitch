//go:build !windows

package history

import (
	"os"
	"path/filepath"
	"syscall"
)

// withLock acquires an exclusive advisory flock on a lock file inside dir,
// calls fn, then releases the lock.  If the lock file cannot be created or
// flock fails we fall back to calling fn without a lock so that the shell
// hook still works in degraded environments (e.g. network filesystems that
// do not support flock).
func withLock(dir string, fn func() error) error {
	lockPath := filepath.Join(dir, ".history.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		// Cannot create lock file — proceed without lock (best effort).
		return fn()
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		// flock not supported (e.g. some NFS mounts) — proceed without lock.
		return fn()
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN) //nolint:errcheck
	return fn()
}
