package lib

import (
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
			return fmt.Errorf("Unable to open source binary: %q", cwd)
		}
		defer r.Close()

		w, err := os.Create(dir)
		if err != nil {
			return fmt.Errorf("Could not create target binary: %q", dir)
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
			return fmt.Errorf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %q" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %q.
		Error: %v
		`, dir, dir, err)
		}
	}
	return nil
}

// RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) error {
	_, err := os.Lstat(symlinkPath)
	if err != nil {
		return fmt.Errorf(`
		Unable to stat symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %q" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %q.
		Error: %v
		`, symlinkPath, symlinkPath, err)
	} else {
		errRemove := os.Remove(symlinkPath)

		if errRemove != nil {
			return fmt.Errorf(`
			Unable to remove symlink.
			Maybe symlink already exist. Try removing existing symlink manually.
			Try running "unlink %q" to remove existing symlink.
			If error persist, you may not have the permission to create a symlink at %q.
			Error: %v
			`, symlinkPath, symlinkPath, errRemove)
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
	homedir := GetHomeDirectory() // get user's home directory
	homeBinPath := filepath.Join(homedir, "bin", product.GetExecutableName())
	possibleInstallLocations := []string{Path(userBinPath), Path(homeBinPath)}
	var err error

	for _, location := range possibleInstallLocations {
		if CheckDirExist(location) {
			/* remove current symlink if exist*/
			symlinkExist := CheckSymlink(location)
			if symlinkExist {
				_ = RemoveSymlink(location)
			}

			/* set symlink to desired version */
			err = CreateSymlink(binVersionPath, location)
			if err == nil {
				logger.Debugf("Symlink created at %q", location)
				return nil
			}
		}
	}

	if err == nil {
		return fmt.Printf("None of the installation directories exist: \"%s\". %s\n",
			strings.Join(possibleInstallLocations, `", "`),
			"Manually create one of them and try again.")
	}

	return err
}
