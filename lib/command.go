package lib

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// string `windows` has 12 occurrences, make it a constant (goconst)
const windows = "windows"

// Command : type string
type Command struct {
	name string
}

// NewCommand : get command
func NewCommand(name string) *Command {
	return &Command{name: name}
}

// PathList : get bin path list
func (cmd *Command) PathList() []string {
	path := os.Getenv("PATH")
	return strings.Split(path, string(os.PathListSeparator))
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// IsRegularFile : check if the given path points to a regular file
func IsRegularFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.Mode().IsRegular()
}

func isExecutable(path string) bool {
	if isDir(path) {
		return false
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if runtime.GOOS == windows {
		return true
	}

	if fileInfo.Mode()&0o111 != 0 {
		return true
	}

	return false
}

// Find : find all bin path
func (cmd *Command) Find() func() string {
	pathChan := make(chan string)
	go func() {
		for _, p := range cmd.PathList() {
			if !isDir(p) {
				continue
			}
			fileList, err := os.ReadDir(p)
			if err != nil {
				continue
			}

			for _, f := range fileList {
				path := filepath.Join(p, f.Name())
				if isExecutable(path) && f.Name() == cmd.name {
					pathChan <- path
				}
			}
		}
		pathChan <- ""
	}()

	return func() string {
		return <-pathChan
	}
}
