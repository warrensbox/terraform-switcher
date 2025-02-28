package lib

import (
	"fmt"
	"os"
	"time"
)

// Create exclusive lock
func acquireLock(lockFile string, LockWaitMaxAttempts int, LockWaitInterval time.Duration) error {
	var waitForLock bool = true
	var LockWaitCount int = 0

	logger.Debugf("Attempting to acquire lock %q", lockFile)
	for waitForLock {
		LockWaitCount++
		if _, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0o644); err != nil {
			if LockWaitCount > LockWaitMaxAttempts {
				return fmt.Errorf("Unable to acquire lock %q", lockFile)
			}

			logger.Infof("Waiting for lock %q to be released (attempt %d out of %d)", lockFile, LockWaitCount, LockWaitMaxAttempts)

			if lockFileInfo, err := os.Stat(lockFile); err == nil {
				logger.Debugf("Lock %q last modification time: %s", lockFile, lockFileInfo.ModTime())
			} else {
				logger.Warnf("Unable to get lock %q last modification time: %v", lockFile, err)
			}

			time.Sleep(LockWaitInterval)
		} else {
			waitForLock = false
			logger.Debugf("Acquired lock %q", lockFile)
			break
		}
	}

	return nil
}

// Release lock file
func releaseLock(lockFile string) {
	if exist := CheckFileExist(lockFile); exist {
		if err := os.Remove(lockFile); err != nil {
			logger.Debugf("Error releasing lock %q: %v", lockFile, err)
		} else {
			logger.Debugf("Releasing lock %q", lockFile)
		}
	} else {
		logger.Debugf("Lock %q not found", lockFile)
	}
}
