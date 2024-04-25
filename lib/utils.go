package lib

import (
	"fmt"
	"github.com/pborman/getopt"
	"os"
	"os/exec"
	"strings"
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
	out, err := exec.Command(installPath, "--version").Output()
	if err != nil {
		logger.Fatal(err)
	}
	outputString := string(out)
	result := strings.FieldsFunc(outputString, func(c rune) bool { return c == '\n' || c == '\r' })
	fmt.Print(result[0])
}
