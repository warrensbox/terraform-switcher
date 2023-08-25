package lib_test

import (
	"os"
	"os/user"
	"runtime"
	"testing"
)

// TestAddRecent : Create a file, check filename exist,
// rename file, check new filename exit
func GetInstallLocation(installPath string) string {
	return string(os.PathSeparator) + installPath + string(os.PathSeparator)
}

func getInstallFile(installFile string) string {
	if runtime.GOOS == "windows" {
		return installFile + ".exe"
	}

	return installFile
}

func TestInstall(t *testing.T) {

	t.Run("User should exist",
		func(t *testing.T) {
			_, errCurr := user.Current()
			if errCurr != nil {
				t.Errorf("Unable to get user %v [unexpected]", errCurr)
			}

			_, errCurr = os.UserHomeDir()
			if errCurr != nil {
				t.Errorf("Unable to get user home directory: %v [unexpected]", errCurr)
			}
		},
	)
}
