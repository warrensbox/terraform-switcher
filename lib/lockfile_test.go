package lib

import (
	"path/filepath"
	"testing"
)

// Test Locking
func TestLocking(t *testing.T) {
	var lockFile string = ".tfswitch.lock"
	tmpDir := t.TempDir()
	lockFilePath := filepath.Join(tmpDir, lockFile)

	// Acquire lock
	if lockedFile, err := acquireLock(lockFilePath, 1, 1); err == nil {
		t.Logf("Lock acquired successfully: %s", err)

		// Release lock
		releaseLock(lockFilePath, lockedFile)
		if CheckFileExist(lockFilePath) {
			t.Errorf("Failed to release lock: %s", lockFilePath)
		} else {
			t.Logf("Lock released successfully: %s", lockFilePath)
		}
	} else {
		t.Errorf("Failed to acquire lock: %s", lockFilePath)
	}
}
