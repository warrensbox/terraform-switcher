package lib

import (
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

// GetDefaultBin Get default binary path
func GetDefaultBin() string {
	var defaultBin = "/usr/local/bin/terraform"
	if runtime.GOOS == "windows" {
		home, err := homedir.Dir()
		if err != nil {
			logger.Fatal("Could not detect home directory.")
			os.Exit(1)
		}
		defaultBin = home + "/bin/terraform.exe"
	}
	return defaultBin
}
