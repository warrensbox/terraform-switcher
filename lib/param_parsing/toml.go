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

		if configKey := "arch"; viperParser.Get(configKey) != nil {
			params.Arch = os.ExpandEnv(viperParser.GetString(configKey))
			logger.Debugf("CPU architecture (%q) from %q: %q", configKey, tomlPath, params.Arch)
		}
		if configKey := "bin"; viperParser.Get(configKey) != nil {
			params.CustomBinaryPath = os.ExpandEnv(viperParser.GetString(configKey))
			logger.Debugf("Custom binary path (%q) from %q: %q", configKey, tomlPath, params.CustomBinaryPath)
		}
		if configKey := "install"; viperParser.Get(configKey) != nil {
			params.InstallPath = viperParser.GetString(configKey)
			logger.Debugf("Custom install path (%q) from %q: %q", configKey, tomlPath, params.InstallPath)
		}
		if configKey := "log-level"; viperParser.Get(configKey) != nil {
			params.LogLevel = viperParser.GetString(configKey)
			logger.Debugf("Logging level (%q) from %q: %q", configKey, tomlPath, params.LogLevel)
		}
		if configKey := "version"; viperParser.Get(configKey) != nil {
			params.Version = viperParser.GetString(configKey)
			logger.Debugf("Installation version (%q) from %q: %q", configKey, tomlPath, params.Version)
		}
		if configKey := "default-version"; viperParser.Get(configKey) != nil {
			params.DefaultVersion = viperParser.GetString(configKey)
			logger.Debugf("Fallback version (%q) from %q: %q", configKey, tomlPath, params.DefaultVersion)
		}
		if configKey := "product"; viperParser.Get(configKey) != nil {
			params.Product = viperParser.GetString(configKey)
			logger.Debugf("Product name (%q) from %q: %q", configKey, tomlPath, params.Product)
		}
	}
	return params, nil
}

func tomlFileExists(params Params) bool {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
