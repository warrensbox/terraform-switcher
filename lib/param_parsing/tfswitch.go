package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
	"path/filepath"
	"strings"
)

const tfSwitchFileName = ".tfswitchrc"

func GetParamsFromTfSwitch(params Params) (Params, error) {
	filePath := filepath.Join(params.ChDirPath, tfSwitchFileName)
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %q", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Errorf("Could not read file content from %q: %v", filePath, err)
			return params, err
		}
		params.Version = strings.TrimSpace(string(content))
	}
	return params, nil
}

func tfSwitchFileExists(params Params) bool {
	filePath := filepath.Join(params.ChDirPath, tfSwitchFileName)
	return lib.CheckFileExist(filePath)
}
