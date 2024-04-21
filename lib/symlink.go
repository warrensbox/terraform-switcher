package lib

import (
	"io"
	"os"
	"runtime"
)

// CreateSymlink : create symlink or copy file to bin directory if windows
func CreateSymlink(cwd string, dir string) {
	// If we are on windows the symlink is not working correctly.
	// Copy the desired terraform binary to the path environment.
	if runtime.GOOS == "windows" {
		r, err := os.Open(cwd)
		if err != nil {
			logger.Fatalf("Unable to open source binary: %s", cwd)
			os.Exit(1)
		}
		defer r.Close()

		w, err := os.Create(dir + ".exe")
		if err != nil {
			logger.Fatalf("Could not create target binary: %s", dir+".exe")
			os.Exit(1)
		}
		defer func() {
			if c := w.Close(); err == nil {
				err = c
			}
		}()
		_, err = io.Copy(w, r)
	} else {
		err := os.Symlink(cwd, dir)
		if err != nil {
			logger.Fatalf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, dir, dir, err)
			os.Exit(1)
		}
	}
}

// RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		logger.Fatalf(`
		Unable to stat symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, symlinkPath, symlinkPath, err)
		os.Exit(1)
	} else {
		errRemove := os.Remove(symlinkPath)

		if errRemove != nil {
			logger.Fatalf(`
			Unable to remove symlink.
			Maybe symlink already exist. Try removing existing symlink manually.
			Try running "unlink %s" to remove existing symlink.
			If error persist, you may not have the permission to create a symlink at %s.
			Error: %s
			`, symlinkPath, symlinkPath, errRemove)
			os.Exit(1)
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

// ChangeSymlink : move symlink to existing binary
func ChangeSymlink(binVersionPath string, binPath string) {

	//installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file
	binPath = installableBinLocation(binPath)

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)
	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(binVersionPath, binPath)

}
