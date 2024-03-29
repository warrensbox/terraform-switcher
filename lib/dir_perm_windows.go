package lib

import (
	"os"
)

func CheckDirWritable(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		logger.Info("Path doesn't exist")
		return false
	}

	err = nil
	if !info.IsDir() {
		logger.Info("Path isn't a directory")
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		logger.Info("Write permission bit is not set on this file for user")
		return false
	}

	return true
}
