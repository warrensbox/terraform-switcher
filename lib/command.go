package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Command struct {
	name string
}

func NewCommand(name string) *Command {
	return &Command{name: name}
}

func (this *Command) PathList() []string {
	path := os.Getenv("PATH")
	return strings.Split(path, string(os.PathListSeparator))
}

func isDir(path string) bool {
	file_info, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return file_info.IsDir()
}

func isExecutable(path string) bool {
	if isDir(path) {
		return false
	}

	file_info, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	if file_info.Mode()&0111 != 0 {
		return true
	}

	return false
}

func (this *Command) Find() func() string {
	path_chan := make(chan string)
	go func() {
		for _, p := range this.PathList() {
			if !isDir(p) {
				continue
			}
			file_list, err := ioutil.ReadDir(p)
			if err != nil {
				continue
			}

			for _, f := range file_list {
				path := filepath.Join(p, f.Name())
				if isExecutable(path) && f.Name() == this.name {
					path_chan <- path
				}
			}
		}
		path_chan <- ""
	}()

	return func() string {
		return <-path_chan
	}
}
