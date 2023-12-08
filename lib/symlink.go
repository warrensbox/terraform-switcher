package lib

import (
	"fmt"
	"io"
	"log"
	"os"
)

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) {

	err := os.Symlink(cwd, dir)
	if err != nil {
		log.Println("Unable to create symlink. Trying to copy file (os without symlink permissions)")

		err := CopyFile(cwd, dir)
		if err != nil {
			log.Fatalf(`
			Unable to create new symlink or copy file.
			Maybe symlink or file already exist. Try removing existing symlink or file manually.
			Try running "unlink %s" to remove existing symlink or "rm %s" to remove existing file.
			If error persist, you may not have the permission to create a symlink or file at %s.
			Error: %s
			`, dir, dir, dir, err)
			os.Exit(1)
		}
	}
}

// RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatalf(`
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
			log.Fatalf(`
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
		if _, err := os.Stat(symlinkPath); err == nil {
			return true
		} else {
			return false
		}
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return true
	}

	return false
}

// ChangeSymlink : move symlink to existing binary
func ChangeSymlink(binVersionPath string, binPath string) {
	fmt.Println("ca passe la - ChangeSymlink")
	//installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file
	binPath = InstallableBinLocation(binPath)

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)
	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(binVersionPath, binPath)

}
