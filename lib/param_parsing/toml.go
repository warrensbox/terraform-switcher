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

		if viperParser.Get("arch") != nil {
			params.Arch = os.ExpandEnv(viperParser.GetString("arch"))
			logger.Debugf("Using \"arch\" from %q: %q", tomlPath, params.Arch)
		}
		if viperParser.Get("bin") != nil {
			params.CustomBinaryPath = os.ExpandEnv(viperParser.GetString("bin"))
			logger.Debugf("Using \"bin\" from %q: %q", tomlPath, params.CustomBinaryPath)
		}
		if viperParser.Get("log-level") != nil {
			params.LogLevel = viperParser.GetString("log-level")
			logger.Debugf("Using \"log-level\" from %q: %q", tomlPath, params.LogLevel)
		}
		if viperParser.Get("version") != nil {
			params.Version = viperParser.GetString("version")
			logger.Debugf("Using \"version\" from %q: %q", tomlPath, params.Version)
		}
		if viperParser.Get("default-version") != nil {
			params.DefaultVersion = viperParser.GetString("default-version")
			logger.Debugf("Using \"default-version\" from %q: %q", tomlPath, params.DefaultVersion)
		}
		if configKey := "product"; viperParser.Get(configKey) != nil {
			params.Product = viperParser.GetString(configKey)
			logger.Debugf("Using %q from %q: %q", configKey, tomlPath, params.Product)
		}
	}
	return params, nil
}

func tomlFileExists(params Params) bool {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
