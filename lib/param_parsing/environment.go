package param_parsing

import (
	"os"
	"reflect"
)

func GetParamsFromEnvironment(params Params) Params {
	reflectedParams := reflect.ValueOf(&params)
	for _, envVar := range paramMappings {
		description := envVar.description
		env := envVar.env
		param := envVar.param
		toml := envVar.toml

		if len(env) == 0 {
			logger.Errorf("Internal error: environment variable name is empty for parameter %q mapping, skipping assignment", param)
			continue
		}
		if len(param) == 0 {
			logger.Errorf("Internal error: parameter name is empty for environment variable %q mapping, skipping assignment", env)
			continue
		}
		if len(description) == 0 {
			description = param
		}

		paramKey := reflect.Indirect(reflectedParams).FieldByName(param)
		if paramKey.Kind() != reflect.String {
			logger.Warnf("Parameter %q is not a string, skipping assignment from environment variable %q", param, env)
			continue
		}

		if envVarValue := os.Getenv(env); envVarValue != "" {
			logger.Debugf("%s (%q) from environment variable %q: %q", description, toml, env, envVarValue)
			if !paramKey.CanSet() {
				logger.Warnf("Parameter %q cannot be set, skipping assignment from environment variable %q", param, env)
			}
			paramKey.SetString(envVarValue)
		}
	}
	return params
}
