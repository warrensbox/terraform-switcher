package param_parsing

import (
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
)

const versionTfFileName = "version.tf"

func GetVersionFromVersionsTF(params Params) Params {
	filePath := params.ChDirPath + "/" + versionTfFileName
	if lib.CheckFileExist(filePath) {
		fmt.Printf("Reading version from %s\n", filePath)
		module, err := tfconfig.LoadModule(params.ChDirPath)
		if err != nil {
			log.Fatal("Could not load terraform module")
		}
		tfconstraint := module.RequiredCore[0]
		version, err2 := lib.GetSemver(tfconstraint, params.MirrorURL)
		if err2 != nil {
			log.Fatal("Could not determine semantic version")
		}
		params.Version = version
	}
	return params
}

func versionTFFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + versionTfFileName
	return lib.CheckFileExist(filePath)
}
