package lib

import (
	"runtime"
)

var (
	PubKeyId     = "72D7468F"
	PubKeyPrefix = "hashicorp_"
	PubKeyUri    = "https://www.hashicorp.com/.well-known/pgp-key.txt"
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
	distTerraform             = "terraform"
	distTofu                  = "tofu"
	installFile               = "terraform"
	InstallDir                = ".terraform.versions"
	pubKeySuffix              = ".asc"
	recentFile                = "RECENT"
	tfDarwinArm64StartVersion = "1.0.2"
	VersionPrefix             = "terraform_"
)
