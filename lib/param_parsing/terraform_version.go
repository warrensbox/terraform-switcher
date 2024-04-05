package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
	"strings"
)

const terraformVersionFileName = ".terraform-version"

func GetParamsFromTerraformVersion(params Params) Params {
	filePath := params.ChDirPath + "/" + terraformVersionFileName
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %s", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Errorf("Could not read file content at %s:%v", filePath, err)
			os.Exit(1)
		}
		params.Version = strings.TrimSpace(string(content))
	}
	return params
}

func terraformVersionFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + terraformVersionFileName
	return lib.CheckFileExist(filePath)
}
