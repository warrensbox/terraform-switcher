package lib

import (
	"runtime"
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
		home := GetHomeDirectory()
		defaultBin = home + "/bin/terraform.exe"
	}
	return defaultBin
}

const (
	DefaultMirror             = "https://releases.hashicorp.com/terraform"
	DefaultLatest             = ""
	installFile               = "terraform"
	InstallDir                = ".terraform.versions"
	recentFile                = "RECENT"
	tfDarwinArm64StartVersion = "1.0.2"
	VersionPrefix             = "terraform_"
)
