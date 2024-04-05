package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
	"strings"
)

const tfSwitchFileName = ".tfswitchrc"

func GetParamsFromTfSwitch(params Params) Params {
	filePath := params.ChDirPath + "/" + tfSwitchFileName
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %s", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Fatalf("Could not read file content %s: %v", filePath, err)
			os.Exit(1)
		}
		params.Version = strings.TrimSpace(string(content))
	}
	return params
}

func tfSwitchFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + tfSwitchFileName
	return lib.CheckFileExist(filePath)
}
