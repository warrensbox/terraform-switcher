package param_parsing

import (
	"os"
	"reflect"
)

func GetParamsFromEnvironment(params Params) Params {
	type envVar struct {
		name        string
		param       string
		description string
	}
	envVars := []envVar{
		{name: "TF_ARCH", param: "Arch", description: "CPU architecture"},
		{name: "TF_BINARY_PATH", param: "CustomBinaryPath", description: "custom binary path"},
		{name: "TF_DEFAULT_VERSION", param: "DefaultVersion", description: "default version"},
		{name: "TF_INSTALL_PATH", param: "InstallPath", description: "custom install path"},
		{name: "TF_PRODUCT", param: "Product", description: "product"},
		{name: "TF_VERSION", param: "Version", description: "version"},
	}

	reflectedParams := reflect.ValueOf(&params)
	for _, envVar := range envVars {
		name := envVar.name
		param := envVar.param
		description := envVar.description

		if len(name) == 0 {
			logger.Errorf("Internal error: environment variable name is empty for parameter %q mapping, skipping assignment", param)
			continue
		}
		if len(param) == 0 {
			logger.Errorf("Internal error: parameter name is empty for environment variable %q mapping, skipping assignment", name)
			continue
		}
		if len(description) == 0 {
			description = param
		}

		f := reflect.Indirect(reflectedParams).FieldByName(param)
		if f.Kind() != reflect.String {
			logger.Warnf("Parameter %q is not a string, skipping assignment from environment variable %q", param, name)
			continue
		}

		if envVarValue := os.Getenv(name); envVarValue != "" {
			logger.Debugf("Using %s from environment variable %q: %q", description, name, envVarValue)
			if !f.CanSet() {
				logger.Warnf("Parameter %q cannot be set, skipping assignment from environment variable %q", param, name)
			}
			f.SetString(envVarValue)
		}
	}
	return params
}
