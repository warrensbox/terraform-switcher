//go:build !windows
// +build !windows

package lib

import "golang.org/x/sys/unix"

//Check if user has permission to directory :
//dir=path to file
//return bool
func CheckDirWritable(dir string) bool {
	return unix.Access(dir, unix.W_OK) == nil
}
