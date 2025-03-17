package param_parsing

import (
	"path/filepath"
	"reflect"

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

		reflectedParams := reflect.ValueOf(&params)
		for _, configKey := range paramMappings {
			description := configKey.description
			param := configKey.param
			toml := configKey.toml

			if len(toml) == 0 {
				logger.Errorf("Internal error: TOML key name is empty for parameter %q mapping, skipping assignment", param)
				continue
			}
			if len(param) == 0 {
				logger.Errorf("Internal error: parameter name is empty for TOML key %q mapping, skipping assignment", toml)
				continue
			}
			if len(description) == 0 {
				description = param
			}

			f := reflect.Indirect(reflectedParams).FieldByName(param)
			if f.Kind() != reflect.String {
				logger.Warnf("Parameter %q is not a string, skipping assignment from TOML key %q", param, toml)
				continue
			}
			if viperParser.Get(toml) != nil {
				configKeyValue := viperParser.GetString(toml)
				logger.Debugf("%s (%q) from %q: %q", description, toml, tomlPath, configKeyValue)
				if !f.CanSet() {
					logger.Warnf("Parameter %q cannot be set, skipping assignment from TOML key %q", param, toml)
				}
				f.SetString(configKeyValue)
			}
		}
	}
	return params, nil
}

func tomlFileExists(params Params) bool {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
