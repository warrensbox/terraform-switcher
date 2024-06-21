package lib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CreateSymlink : create symlink or copy file to bin directory if windows
func CreateSymlink(cwd string, dir string) error {
	// If we are on windows the symlink is not working correctly.
	// Copy the desired terraform binary to the path environment.
	if runtime.GOOS == "windows" {
		r, err := os.Open(cwd)
		if err != nil {
			logger.Debugf("Unable to open source binary: %s", cwd)
			return errors.New(fmt.Sprintf("Unable to open source binary: %s", cwd))
		}
		defer r.Close()

		w, err := os.Create(dir + ".exe")
		if err != nil {
			logger.Debugf("Could not create target binary: %s.exe", dir)
			return errors.New(fmt.Sprintf("Could not create target binary: %s.exe", dir))
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
			logger.Debugf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, dir, dir, err)
			return errors.New(fmt.Sprintf("Unable to create new symlink %s : %s", dir, err))
		}
	}
	return nil
}

// RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) error {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		logger.Debugf(`
		Unable to stat symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, symlinkPath, symlinkPath, err)
		return errors.New(fmt.Sprintf("Unable to stat symlink %s : %s", symlinkPath, err))
	} else {
		errRemove := os.Remove(symlinkPath)

		if errRemove != nil {
			logger.Debugf(`
			Unable to remove symlink.
			Maybe symlink already exist. Try removing existing symlink manually.
			Try running "unlink %s" to remove existing symlink.
			If error persist, you may not have the permission to create a symlink at %s.
			Error: %s
			`, symlinkPath, symlinkPath, errRemove)
			return errors.New(fmt.Sprintf("Unable to remove symlink %s : %s", symlinkPath, err))
		}
	}
	return nil
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

// ChangeSymlink : move symlink to existing binary for Terraform
//
// Deprecated: This function has been deprecated in favor of ChangeProductSymlink and will be removed in v2.0.0
func ChangeSymlink(binVersionPath string, binPath string) {
	product := getLegacyProduct()
	err := ChangeProductSymlink(product, binVersionPath, binPath)
	if err != nil {
		logger.Fatal(err)
	}
}

// ChangeProductSymlink : move symlink for product to existing binary
func ChangeProductSymlink(product Product, binVersionPath string, userBinPath string) error {

	homedir := GetHomeDirectory() //get user's home directory
	homeBinPath := filepath.Join(homedir, "bin", product.GetExecutableName())
	possibleInstallLocations := []string{userBinPath, homeBinPath}
	var err error

	for _, location := range possibleInstallLocations {
		if CheckDirExist(Path(location)) {
			/* remove current symlink if exist*/
			symlinkExist := CheckSymlink(location)
			if symlinkExist {
				_ = RemoveSymlink(location)
			}

			/* set symlink to desired version */
			err = CreateSymlink(binVersionPath, location)
			if err == nil {
				logger.Debugf("Symlink created at %s", location)
				return nil
			}
		}
	}

	if err == nil {
		msg := fmt.Sprintf("Unable to find existing directory in %s. %s",
			strings.Join(possibleInstallLocations, " or "),
			"Manually create one of them and try again.")
		err = errors.New(msg)
	}

	return err
}
