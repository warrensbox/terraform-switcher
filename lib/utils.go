package lib

import (
	"os"
)

// FileExistsAndIsNotDir checks if a file exists and is not a directory before we try using it to prevent further errors
func FileExistsAndIsNotDir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func closeFileHandlers(handlers []*os.File) {
	for _, handler := range handlers {
		logger.Debugf("Closing file handler %q", handler.Name())
		_ = handler.Close()
	}
}
