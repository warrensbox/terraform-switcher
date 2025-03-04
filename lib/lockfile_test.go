package lib

import (
	"path/filepath"
	"testing"
	"time"
)

// Test Locking
func TestLocking(t *testing.T) {
	var lockFile string = ".tfswitch.lock"
	lockFilePath := filepath.Join(t.TempDir(), lockFile)

	t.Logf("Testing lock acquirement: %s", lockFilePath)

	// Acquire lock
	if lockedFile, err := acquireLock(lockFilePath, 1, 1*time.Second); err == nil {
		t.Logf("Lock acquired successfully: %s", lockFilePath)

		// Concurrent lock
		t.Logf("Testing concurrent lock acquirement: %s", lockFilePath)
		if _, err := acquireLock(lockFilePath, 1, 1*time.Second); err == nil {
			t.Errorf("Concurrent lock acquired successfully: %s. This is NOT expected!", lockFilePath)
		} else {
			t.Logf("Concurrent lock failed: %s. This is expected.", lockFilePath)
		}

		// Release lock
		releaseLock(lockFilePath, lockedFile)
		if CheckFileExist(lockFilePath) {
			t.Errorf("Lock %s still exists. This is NOT expected!", lockFilePath)
		} else {
			t.Logf("Lock released successfully: %s", lockFilePath)
		}
	} else {
		t.Errorf("Failed to acquire lock: %s", lockFilePath)
	}
}
