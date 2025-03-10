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
	if runtime.GOOS == windows {
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
	}

	if errRemove := os.Remove(symlinkPath); errRemove != nil {
		return fmt.Errorf(`
		Unable to remove symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %q" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %q.
		Error: %v
		`, symlinkPath, symlinkPath, errRemove)
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
//
//nolint:gocyclo
func ChangeProductSymlink(product Product, binVersionPath string, userBinPath string) error {
	homedir := GetHomeDirectory() // get user's home directory
	homeBinPath := filepath.Join(homedir, "bin", product.GetExecutableName())

	var err error
	var locationsFmt string

	// Possible install locations with boolean property as to whether to attempt to create
	type installLocations struct {
		path   string
		create bool
	}
	possibleInstallLocations := []installLocations{
		{path: userBinPath, create: false},
		{path: homeBinPath, create: true},
	}

	for idx, location := range possibleInstallLocations {
		isFallback := false
		if idx > 0 {
			isFallback = true
		}
		convertedPath := ConvertExecutableExt(location.path)
		possibleInstallLocations[idx].path = convertedPath
		locationsFmt += fmt.Sprintf("\n\t• №%d: %q (create: %-5t, isFallack: %t)", idx+1, convertedPath, location.create, isFallback)
	}
	logger.Noticef("Possible install locations:%s", locationsFmt)

	for idx, location := range possibleInstallLocations {
		dirPath := Path(location.path)
		attempt := idx + 1

		if attempt > 1 {
			logger.Warnf("Falling back to install to %q directory", dirPath)
		}

		logger.Noticef("Attempting to install to %q directory (possible install location №%d)", dirPath, attempt)

		// If directory does not exist, check if we should create it, otherwise skip
		if !CheckDirExist(dirPath) {
			logger.Warnf("Installation directory %q doesn't exist!", dirPath)
			if location.create {
				logger.Infof("Creating %q directory", dirPath)
				err = os.MkdirAll(dirPath, 0o755)
				if err != nil {
					logger.Errorf("Unable to create %q directory: %v", dirPath, err)
					continue
				}
			} else {
				continue
			}
		} else if !CheckIsDir(dirPath) {
			logger.Warnf("The %q is not a directory!", dirPath)
			continue
		}
		logger.Noticef("Installation location: %q", location.path)

		/* remove current symlink if exist*/
		if CheckSymlink(location.path) {
			_ = RemoveSymlink(location.path)
		}

		/* set symlink to desired version */
		err = CreateSymlink(binVersionPath, location.path)
		if err == nil {
			logger.Noticef("Symlink created at %q", location.path)

			// Print helper message to export PATH if the directory is not in PATH only for non-Windows systems,
			// as it's all complicated on Windows. See https://github.com/warrensbox/terraform-switcher/issues/558
			if runtime.GOOS != windows {
				isDirInPath := false

				for _, envPathElement := range strings.Split(os.Getenv("PATH"), ":") {
					expandedEnvPathElement := strings.TrimRight(strings.Replace(envPathElement, "~", homedir, 1), "/")

					if expandedEnvPathElement == strings.TrimRight(dirPath, "/") {
						isDirInPath = true
						break
					}
				}

				if !isDirInPath {
					logger.Warnf("Run `export PATH=\"$PATH:%s\"` to append %q to $PATH", dirPath, location.path)
				}
			}

			return nil
		}
	}

	if err == nil {
		return fmt.Errorf("None of the installation directories exist:%s\n\t%s", locationsFmt,
			"Manually create one of them and try again")
	}

	return err
}
