package lib

import (
	"fmt"
	"os"
	"time"
)

// Create exclusive lock
func acquireLock(lockFile string, lockWaitMaxAttempts int, lockWaitInterval time.Duration) error {
	logger.Debugf("Attempting to acquire lock %q", lockFile)

	for lockAttempt := range lockWaitMaxAttempts {
		lockAttemptCounter := lockAttempt + 1

		if _, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0o644); err == nil {
			logger.Debugf("Acquired lock %q", lockFile)
			return nil
		}

		logger.Infof("Waiting for lock %q to be released (attempt %d out of %d)", lockFile, lockAttemptCounter, lockWaitMaxAttempts)

		if lockFileInfo, err := os.Stat(lockFile); err == nil {
			logger.Debugf("Lock %q last modification time: %s", lockFile, lockFileInfo.ModTime())
		} else {
			logger.Warnf("Unable to get lock %q last modification time: %v", lockFile, err)
		}

		if lockAttemptCounter < lockWaitMaxAttempts {
			time.Sleep(lockWaitInterval)
		}
	}

	return fmt.Errorf("Failed to acquire lock %q", lockFile)
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
