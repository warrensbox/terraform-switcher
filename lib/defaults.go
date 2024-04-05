package lib

import (
	"github.com/mitchellh/go-homedir"
	"runtime"
)

// GetDefaultBin Get default binary path
func GetDefaultBin() string {
	var defaultBin = "/usr/local/bin/terraform"
	if runtime.GOOS == "windows" {
		home, err := homedir.Dir()
		if err != nil {
			logger.Fatal("Could not detect home directory.")
		}
		defaultBin = home + "/bin/terraform.exe"
	}
	return defaultBin
}
