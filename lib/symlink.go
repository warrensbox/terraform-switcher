package lib

import (
	"os"
)

//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) error {

	if err := os.Symlink(cwd, dir); err != nil {
		return err
	}
	return nil
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) error {

	if _, err := os.Lstat(symlinkPath); err != nil {
		return err
	}
	os.Remove(symlinkPath)
	return nil
}
