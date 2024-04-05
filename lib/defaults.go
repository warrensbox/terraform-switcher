package lib

import (
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
)

var (
	PubKeyId     = "72D7468F"
	PubKeyPrefix = "hashicorp_"
	PubKeyUri    = "https://www.hashicorp.com/.well-known/pgp-key.txt"
)

const (
	pubKeySuffix = ".asc"
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

const (
	DefaultMirror = "https://releases.hashicorp.com/terraform"
	DefaultLatest = ""
)
