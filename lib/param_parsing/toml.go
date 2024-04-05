package param_parsing

import (
	"github.com/spf13/viper"
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
)

const tfSwitchTOMLFileName = ".tfswitch.toml"

// getParamsTOML parses everything in the toml file, return required version and bin path
func getParamsTOML(params Params) Params {
	tomlPath := params.ChDirPath + "/" + tfSwitchTOMLFileName
	if tomlFileExists(params) {
		logger.Infof("Reading configuration from %q", tomlPath)
		configfileName := lib.GetFileName(tfSwitchTOMLFileName)
		viperParser := viper.New()
		viperParser.SetConfigType("toml")
		viperParser.SetConfigName(configfileName)
		viperParser.AddConfigPath(params.ChDirPath)

		errs := viperParser.ReadInConfig() // Find and read the config file
		if errs != nil {
			logger.Fatalf("Could not to read %q: %v", tomlPath, errs)
			os.Exit(1)
		}

		params.Version = viperParser.GetString("version") // Attempt to get the version if it's provided in the toml
		params.CustomBinaryPath = viperParser.GetString("bin")
	}
	return params
}

func tomlFileExists(params Params) bool {
	tomlPath := params.ChDirPath + "/" + tfSwitchTOMLFileName
	return lib.CheckFileExist(tomlPath)
}
