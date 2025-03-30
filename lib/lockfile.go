package lib

import (
	"fmt"
	"os"
	"time"
)

// Acquire exclusive lock
func acquireLock(lockFile string, lockWaitMaxAttempts int, lockWaitInterval time.Duration) (*os.File, error) {
	logger.Debugf("Attempting to acquire lock %q", lockFile)

	for lockAttempt := 1; lockAttempt <= lockWaitMaxAttempts; lockAttempt++ {
		if file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0o644); err == nil {
			logger.Debugf("Acquired lock %q", lockFile)
			return file, nil
		}

		logger.Infof("Waiting for lock %q to be released (attempt %d out of %d)", lockFile, lockAttempt, lockWaitMaxAttempts)

		if lockFileInfo, err := os.Stat(lockFile); err == nil {
			logger.Debugf("Lock %q last modification time: %s", lockFile, lockFileInfo.ModTime())
		} else {
			logger.Warnf("Unable to get lock %q last modification time: %v", lockFile, err)
		}

		if lockAttempt < lockWaitMaxAttempts {
			time.Sleep(lockWaitInterval)
		}
	}

	return nil, fmt.Errorf("Failed to acquire lock %q", lockFile)
}

// Release and remove lock
func releaseLock(lockFile string, lockedFH *os.File) {
	logger.Debugf("Releasing lock %q", lockFile)

	if lockedFH == nil {
		logger.Warnf("Lock is `nil` on %q", lockFile)
		if CheckFileExist(lockFile) {
			logger.Warnf("Lock %q exists. This is NOT expected!", lockFile)
		}
		return
	}

	if err := lockedFH.Close(); err != nil {
		logger.Warnf("Failed to release lock %q: %v", lockFile, err)
	} else {
		logger.Debugf("Released lock %q", lockFile)
	}

	logger.Debugf("Removing lock %q", lockFile)

	if CheckFileExist(lockFile) {
		if err := os.Remove(lockFile); err != nil {
			logger.Warnf("Failed to remove lock %q: %v", lockFile, err)
		} else {
			logger.Debugf("Removed lock %q", lockFile)
		}
	} else {
		logger.Warnf("Lock %q doesn't exist", lockFile)
	}
}
