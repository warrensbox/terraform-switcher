package lib

import (
	"log"
	"os"
)

//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) {

	err := os.Symlink(cwd, dir)
	if err != nil {
		log.Fatalf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try running removing existing symlink.
		Try running "unlink %s" to remove existing symlink.
		Maybe you do not have privilege to create symlink at %s.
		Error: %s
		`, dir, dir, err)
		panic(err)
	}
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatalf(`
		Unable to remove symlink.
		Try running removing existing symlink.
		Try running "unlink %s" to remove existing symlink.
		Maybe you do not have privilege to remove symlink at %s.
		Error: %s
		`, symlinkPath, symlinkPath, err)
		panic(err)
	} else {
		errRemove := os.Remove(symlinkPath)
		if errRemove != nil {
			log.Fatalf(`
			Unable to remove symlink.
			Try running removing existing symlink.
			Try running "unlink %s" to remove existing symlink.
			Maybe you do not have privilege to remove symlink at %s.
			Error: %s
			`, symlinkPath, symlinkPath, errRemove)
			panic(errRemove)
		}
	}
}

// CheckSymlink : check file is symlink
func CheckSymlink(symlinkPath string) bool {

	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return true
	}

	return false
}
