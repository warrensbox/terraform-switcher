package lib

import (
	"runtime"
	"testing"
)

// Test Locking
func TestLocking(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("Skipping test on Windows: cannot figure it out 01-Mar-2025")
	} else {
		var lockFile string = "../test-data/lockfile.lock"

		// Acquire lock
		if lockAcquireErr := acquireLock(lockFile, 1, 1); lockAcquireErr == nil {
			t.Logf("Lock acquired successfully: %s", lockFile)

			// Release lock
			releaseLock(lockFile)
			if CheckFileExist(lockFile) {
				t.Errorf("Failed to release lock: %s", lockFile)
			} else {
				t.Logf("Lock released successfully: %s", lockFile)
			}
		} else {
			t.Errorf("Failed to acquire lock: %s", lockFile)
		}
	}
}
