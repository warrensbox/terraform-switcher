//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pborman/getopt"
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

func UsageMessage() {
	getopt.PrintUsage(os.Stderr)
}

// GetRelativePath : get relative path from absolute path
func GetRelativePath(absPath string) (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Could not get current working directory: %v", err)
	}

	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return "", fmt.Errorf("Could not derive absolute path to %q: %v", absPath, err)
		}
	}

	relPath, err := filepath.Rel(curDir, absPath)
	if err != nil {
		return "", fmt.Errorf("Could not derive relative path to %q: %v", absPath, err)
	}

	return relPath, nil
}

// RemoveDuplicateStrings : deduplicate slice of strings
func RemoveDuplicateStrings(slice []string) []string {
	seen := map[string]bool{}
	res := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			res = append(res, item)
		}
	}
	return res
}
