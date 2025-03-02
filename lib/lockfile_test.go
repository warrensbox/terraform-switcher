package lib

import (
	"testing"
)

// Test Locking
func TestLocking(t *testing.T) {
	var lockFile string = "lockfile.lock"

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
