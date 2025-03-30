package lib

const (
	DefaultMirror             = "https://releases.hashicorp.com/terraform"
	DefaultLatest             = ""
	InstallDir                = ".terraform.versions"
	pubKeySuffix              = ".asc"
	recentFile                = "RECENT"
	tfDarwinArm64StartVersion = "1.0.2"
	DefaultProductId          = "terraform" // nolint:revive // FIXME: var-naming: const DefaultProductId should be DefaultProductID (revive)
)
