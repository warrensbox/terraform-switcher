package lib

import (
	"fmt"
	"os"
)

func CheckDirWritable(path string) bool {

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("Path doesn't exist")
		return false
	}

	err = nil
	if !info.IsDir() {
		fmt.Println("Path isn't a directory")
		return false
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		fmt.Println("Write permission bit is not set on this file for user")
		return false
	}

	return true
}
