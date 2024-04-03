package lib

import (
	"log"
	"os"
)

// FileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func closeFileHandlers(handlers []*os.File) {
	for _, handler := range handlers {
		log.Println("Closing file handler", handler.Name())
		_ = handler.Close()
	}
}
