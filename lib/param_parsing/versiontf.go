package param_parsing

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
)

const versionTfFileName = "version.tf"

func GetVersionFromVersionsTF(params Params) Params {
	filePath := params.ChDirPath + "/" + versionTfFileName
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading version from %q", filePath)
		module, err := tfconfig.LoadModule(params.ChDirPath)
		if err != nil {
			logger.Fatal("Could not load terraform module")
			os.Exit(1)
		}
		tfconstraint := module.RequiredCore[0]
		version, err2 := lib.GetSemver(tfconstraint, params.MirrorURL)
		if err2 != nil {
			logger.Fatalf("No version found matching %q", tfconstraint)
			os.Exit(1)
		}
		params.Version = version
	}
	return params
}

func versionTFFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + versionTfFileName
	return lib.CheckFileExist(filePath)
}
