package lib

import (
	"os"
)

func CheckDirWritable(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Path doesn't exist: %q", path)
		return false
	}

	err = nil
	if !info.IsDir() {
		logger.Errorf("Path isn't a directory: %q", path)
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		logger.Errorf("Path is not writable by the user: %q", path)
		return false
	}

	return true
}
