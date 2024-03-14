package lib

import (
	"github.com/mitchellh/go-homedir"
	"log"
	"runtime"
)

// GetDefaultBin Get default binary path
func GetDefaultBin() string {
	var defaultBin = "/usr/local/bin/terraform"
	if runtime.GOOS == "windows" {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf(`Could not detect home directory.`)
		}
		defaultBin = home + `/bin/terraform.exe`
	}
	return defaultBin
}
