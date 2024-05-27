package lib

import (
	"runtime"
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
	InstallDir                = ".terraform.versions"
	recentFile                = "RECENT"
	tfDarwinArm64StartVersion = "1.0.2"
	DefaultProductId          = "terraform"
)
