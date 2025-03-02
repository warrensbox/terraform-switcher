package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// Acquire exclusive lock
func acquireLock(lockFile string, lockWaitMaxAttempts int, lockWaitInterval time.Duration) (*lockedfile.File, error) {
	logger.Debugf("Attempting to acquire lock %q", lockFile)

	for lockAttempt := range lockWaitMaxAttempts {
		lockAttemptCounter := lockAttempt + 1

		if file, err := lockedfile.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0o644); err == nil {
			logger.Debugf("Acquired lock %q", lockFile)
			return file, nil
		}

		logger.Infof("Waiting for lock %q to be released (attempt %d out of %d)", lockFile, lockAttemptCounter, lockWaitMaxAttempts)

		if lockFileInfo, err := os.Stat(lockFile); err == nil {
			logger.Debugf("Lock %q last modification time: %s", lockFile, lockFileInfo.ModTime())
		} else {
			logger.Warnf("Unable to get lock %q last modification time: %w", lockFile, err)
		}

		if lockAttemptCounter < lockWaitMaxAttempts {
			time.Sleep(lockWaitInterval)
		}
	}

	return nil, fmt.Errorf("Failed to acquire lock %q", lockFile)
}

// Release lock file
func releaseLock(lockFile string, lockedFile *lockedfile.File) {
	logger.Debugf("Releasing lock %q", lockFile)

	if lockedFile == nil {
		logger.Warnf("Lock is `nil` on %q", lockFile)
		removeLock(lockFile)
		return
	}

	if err := lockedFile.Close(); err != nil {
		logger.Warnf("Failed to release lock %q: %w", lockFile, err)
	} else {
		logger.Debugf("Released lock %q", lockFile)
	}

	removeLock(lockFile)
}

// Remove lock file
func removeLock(lockFile string) {
	logger.Debugf("Removing lock %q", lockFile)

	if exist := CheckFileExist(lockFile); exist {
		if err := os.Remove(lockFile); err != nil {
			logger.Warnf("Failed to remove lock %q: %w", lockFile, err)
		} else {
			logger.Debugf("Removed lock %q", lockFile)
		}
	} else {
		logger.Warnf("Lock %q doesn't exist", lockFile)
	}
}
