package param_parsing

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/warrensbox/terraform-switcher/lib"
)

const terraformVersionFileName = ".terraform-version"

func GetParamsFromTerraformVersion(params Params) (Params, error) {
	filePath := filepath.Join(params.ChDirPath, terraformVersionFileName)
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %q", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Errorf("Could not read file content at %q: %v", filePath, err)
			return params, err
		}
		params.Version = strings.TrimSpace(string(content))
		logger.Debugf("Using version from %q: %q", filePath, params.Version)
	}
	return params, nil
}

func terraformVersionFileExists(params Params) bool {
	filePath := filepath.Join(params.ChDirPath, terraformVersionFileName)
	return lib.CheckFileExist(filePath)
}
