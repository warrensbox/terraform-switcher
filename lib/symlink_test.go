package lib

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mitchellh/go-homedir"
)

// TestCreateSymlink : check if symlink exist-remove if exist,
// create symlink, check if symlink exist, remove symlink
func TestCreateSymlink(t *testing.T) {
	testSymlinkDest := "/test-tfswitcher-dest"
	testSymlinkSrc := "/test-tfswitcher-src"
	if runtime.GOOS == windows {
		testSymlinkSrc += ".exe"
	}

	home, err := homedir.Dir()
	if err != nil {
		t.Errorf("Could not detect home directory.")
	}
	symlinkPathSrc := filepath.Join(home, testSymlinkSrc)
	symlinkPathDest := filepath.Join(home, testSymlinkDest)

	// Create file for test as windows does not like no source
	create, err := os.Create(symlinkPathDest)
	if err != nil {
		t.Errorf("Could not create test dest file for symlink at %v", symlinkPathDest)
	}
	defer create.Close()

	if runtime.GOOS != windows {
		ln, _ := os.Readlink(symlinkPathSrc) // nolint:errcheck // covered by conditional below

		if ln != symlinkPathDest {
			t.Logf("Symlink does not exist %v [expected]", ln)
		} else {
			t.Logf("Symlink exist %v [expected]", ln)
			_ = os.Remove(symlinkPathSrc)
			t.Logf("Removed existing symlink for testing purposes")
		}
	}

	lnCreateErr := CreateSymlink(symlinkPathDest, symlinkPathSrc)
	if lnCreateErr != nil {
		t.Errorf("Could not create symlink at %q to %q: %v", symlinkPathSrc, symlinkPathDest, lnCreateErr)
	}

	if runtime.GOOS == windows {
		_, err := os.Stat(symlinkPathSrc)
		if err != nil {
			t.Logf("Could not stat file copy at %v. [unexpected]", symlinkPathSrc)
			t.Error("File copy was not created.")
		} else {
			t.Logf("File copy exists at %v [expected]", symlinkPathSrc)
		}
	} else {
		lnCheck, _ := os.Readlink(symlinkPathSrc) // nolint:errcheck // covered by conditional below
		if lnCheck == symlinkPathDest {
			t.Logf("Symlink exist %v [expected]", lnCheck)
		} else {
			t.Logf("Symlink does not exist %v [unexpected]", lnCheck)
			t.Error("Symlink was not created")
		}
	}

	symlinkPathSrcErr := os.Remove(symlinkPathSrc)
	if symlinkPathSrcErr != nil {
		t.Logf("Could not remove %q: %v [internal failure]", symlinkPathSrc, symlinkPathSrcErr)
	}
	symlinkPathDestErr := os.Remove(symlinkPathDest)
	if symlinkPathDestErr != nil {
		t.Logf("Could not remove %q: %v [internal failure]", symlinkPathDest, symlinkPathDestErr)
	}
}

// TestRemoveSymlink : check if symlink exist-create if does not exist,
// remove symlink, check if symlink exist
func TestRemoveSymlink(t *testing.T) {
	testSymlinkSrc := "/test-tfswitcher-src"
	testSymlinkDest := "/test-tfswitcher-dest"

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	symlinkPathSrc := filepath.Join(homedir, testSymlinkSrc)
	symlinkPathDest := filepath.Join(homedir, testSymlinkDest)

	ln, _ := os.Readlink(symlinkPathSrc) // nolint:errcheck // covered by conditional below

	if ln != symlinkPathDest {
		t.Logf("Symlink does exist %v [expected]", ln)
		t.Log("Creating symlink")
		if err := os.Symlink(symlinkPathDest, symlinkPathSrc); err != nil {
			t.Error(err)
		}
	}

	RemoveSymlink(symlinkPathSrc)             // nolint:errcheck // covered by conditional below
	lnCheck, _ := os.Readlink(symlinkPathSrc) // nolint:errcheck // covered by conditional below
	if lnCheck == symlinkPathDest {
		t.Logf("Symlink should not exist %v [unexpected]", lnCheck)
		t.Error("Symlink was not removed")
	} else {
		t.Logf("Symlink was removed  %v [expected]", lnCheck)
	}
}

// TestCheckSymlink : Create symlink, test if file is symlink
func TestCheckSymlink(t *testing.T) {
	testSymlinkSrc := "/test-tfswitcher-src"
	testSymlinkDest := "/test-tfswitcher-dest"

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	symlinkPathSrc := filepath.Join(homedir, testSymlinkSrc)
	symlinkPathDest := filepath.Join(homedir, testSymlinkDest)

	ln, _ := os.Readlink(symlinkPathSrc) // nolint:errcheck // it is okay to ignore error here

	if ln != symlinkPathDest {
		t.Log("Creating symlink")
		if err := os.Symlink(symlinkPathDest, symlinkPathSrc); err != nil {
			t.Error(err)
		}
	}

	symlinkExist := CheckSymlink(symlinkPathSrc)

	if symlinkExist {
		t.Logf("Symlink does exist %v [expected]", ln)
	} else {
		t.Logf("Symlink does not exist %v [unexpected]", ln)
	}

	symlinkPathSrcErr := os.Remove(symlinkPathSrc)
	if symlinkPathSrcErr != nil {
		t.Logf("Could not remove %q: %v [internal failure]", symlinkPathSrc, symlinkPathSrcErr)
	}
}
