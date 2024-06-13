package param_parsing

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/warrensbox/terraform-switcher/lib"
)

const tfSwitchTOMLFileName = ".tfswitch.toml"

// getParamsTOML parses everything in the toml file, return required version and bin path
func getParamsTOML(params Params) (Params, error) {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	if tomlFileExists(params) {
		logger.Infof("Reading configuration from %q", tomlPath)
		configfileName := lib.GetFileName(tfSwitchTOMLFileName)
		viperParser := viper.New()
		viperParser.SetConfigType("toml")
		viperParser.SetConfigName(configfileName)
		viperParser.AddConfigPath(params.TomlDir)

		errs := viperParser.ReadInConfig() // Find and read the config file
		if errs != nil {
			logger.Errorf("Could not to read %q: %v", tomlPath, errs)
			return params, errs
		}

		if viperParser.Get("bin") != nil {
			params.CustomBinaryPath = os.ExpandEnv(viperParser.GetString("bin"))
		}
		if viperParser.Get("log-level") != nil {
			params.LogLevel = viperParser.GetString("log-level")
		}
		if viperParser.Get("version") != nil {
			params.Version = viperParser.GetString("version")
		}
		if configKey := "product"; viperParser.Get(configKey) != nil {
			params.Product = viperParser.GetString(configKey)
		}
	}
	return params, nil
}

func tomlFileExists(params Params) bool {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
