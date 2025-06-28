//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
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

		// Find and read the config file
		if err := viperParser.ReadInConfig(); err != nil {
			logger.Errorf("Could not to read %q: %v", tomlPath, err)
			return params, err
		}

		reflectedParams := reflect.ValueOf(&params)
		for _, configKey := range paramMappings {
			description := configKey.description
			param := configKey.param
			ptype := configKey.ptype
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

			paramKey := reflect.Indirect(reflectedParams).FieldByName(param)
			if !paramKey.CanSet() {
				logger.Errorf("Internal error: parameter %q cannot be set, skipping assignment from TOML key %q", param, toml)
				continue
			}

			if viperParser.Get(toml) != nil {
				configKeyValue := viperParser.Get(toml)

				if reflect.TypeOf(configKeyValue).Kind() != ptype {
					logger.Warnf(
						"TOML key %q is not a %s but a %s, skipping assignment of %q parameter from TOML",
						toml, ptype.String(), reflect.TypeOf(configKeyValue).Kind(), param,
					)
					continue
				}

				switch toml {
				case "bin", "install":
					envExpandedConfigKeyValue := os.ExpandEnv(configKeyValue.(string))
					logger.Debugf(
						"Expanded environment variables in %q TOML key value (if any): %q -> %q",
						toml, configKeyValue, envExpandedConfigKeyValue,
					)
					configKeyValue = envExpandedConfigKeyValue
				}

				logger.Debugf("%s (%q) from %q: %v", description, toml, tomlPath, configKeyValue)

				switch ptype {
				case reflect.Bool:
					paramKey.SetBool(configKeyValue.(bool))
				case reflect.String:
					paramKey.SetString(configKeyValue.(string))
				default:
					logger.Errorf(
						"Internal error: unhandled switch case for \"%T\" type of %q parameter (TOML key %q)",
						ptype, param, toml,
					)
					continue
				}
			}
		}
	}
	return params, nil
}

func tomlFileExists(params Params) bool {
	tomlPath := filepath.Join(params.TomlDir, tfSwitchTOMLFileName)
	return lib.CheckFileExist(tomlPath)
}
