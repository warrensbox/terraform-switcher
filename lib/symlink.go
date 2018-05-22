package lib

import (
	"os"
)

//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) error {

	err := os.Symlink(cwd, dir)
	if err != nil {
		return err
	}
	return nil
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) error {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		return err
	}
	os.Remove(symlinkPath)
	return nil
}
