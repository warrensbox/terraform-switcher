package param_parsing

import (
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/warrensbox/terraform-switcher/lib"
)

const tfSwitchTOMLFileName = ".tfswitch.toml"

// Return base directory of TOML file
func getTomlFileDirectory() string {
	return lib.GetHomeDirectory()
}

// getParamsTOML parses everything in the toml file, return required version and bin path
func getParamsTOML(params Params) (Params, error) {
	baseDir := getTomlFileDirectory()
	tomlPath := filepath.Join(baseDir, tfSwitchTOMLFileName)
	if tomlFileExists() {
		logger.Infof("Reading configuration from %q", tomlPath)
		configfileName := lib.GetFileName(tfSwitchTOMLFileName)
		viperParser := viper.New()
		viperParser.SetConfigType("toml")
		viperParser.SetConfigName(configfileName)
		viperParser.AddConfigPath(baseDir)

		errs := viperParser.ReadInConfig() // Find and read the config file
		if errs != nil {
			logger.Errorf("Could not to read %q: %v", tomlPath, errs)
			return params, errs
		}

		if viperParser.Get("bin") != nil {
			params.CustomBinaryPath = viperParser.GetString("bin")
		}
		if viperParser.Get("log-level") != nil {
			params.LogLevel = viperParser.GetString("log-level")
		}
		if viperParser.Get("version") != nil {
			params.Version = viperParser.GetString("version")
		}
	}
	return params, nil
}

func tomlFileExists() bool {
	baseDir := getTomlFileDirectory()
	tomlPath := filepath.Join(baseDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
