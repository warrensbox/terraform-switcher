package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
	"path/filepath"
	"strings"
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
	}
	return params, nil
}

func terraformVersionFileExists(params Params) bool {
	filePath := filepath.Join(params.ChDirPath, terraformVersionFileName)
	return lib.CheckFileExist(filePath)
}
