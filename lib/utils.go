package lib

import (
	"bufio"
	"fmt"
	"github.com/pborman/getopt"
	"os"
	"path/filepath"
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
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}

func CurrentActiveVersion(installPath string) {
	installLocation = getInstallLocation(installPath)
	currentFile := filepath.Join(installLocation, currentFileName)
	file, err := os.Open(currentFile)
	defer file.Close()
	if err != nil {
		logger.Fatalf("Could not open file '%q'", err)
	}
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fmt.Println("Terraform-Version:", scanner.Text())
	}
}
