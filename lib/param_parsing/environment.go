//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
	"reflect"
	"strconv"

	"github.com/warrensbox/terraform-switcher/lib"
)

func GetParamsFromEnvironment(params Params) Params {
	reflectedParams := reflect.ValueOf(&params)
	for _, envVar := range paramMappings {
		description := envVar.description
		env := envVar.env
		param := envVar.param
		ptype := envVar.ptype
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
		if !paramKey.CanSet() {
			logger.Warnf("Parameter %q cannot be set, skipping assignment from environment variable %q", param, env)
			continue
		}

		envVarValue := os.Getenv(env)
		if envVarValue == "" {
			logger.Tracef("Environment variable %q value is empty for %q parameter, skipping assignment", env, toml)
			continue
		}

		logger.Debugf("%s (%q) from environment variable %q: %q", description, toml, env, envVarValue)

		switch ptype {
		case reflect.Bool:
			// Inherit `gookit/color` lib's behavior: whatever the value is, set it to true
			// E.g. NO_COLOR: https://github.com/gookit/color/blob/master/color.go#L49
			paramKey.SetBool(true)
		case reflect.Int64:
			// Convert string to int64 (base 10, 64-bit size)
			parsedInt, err := strconv.ParseInt(envVarValue, 10, 64)
			if err != nil {
				logger.Fatalf(
					"Failed to parse environment variable %q value %q, skipping assignment",
					env, envVarValue,
				)
				continue
			}
			if env == "TF_HTTP_TIMEOUT" {
				if !lib.IsValidHTTPTimeoutValue(parsedInt) {
					logger.Fatalf(
						"Environment variable %q value %q must be a non-negative valid integer, skipping assignment",
						env, envVarValue,
					)
					continue
				}
			}
			paramKey.SetInt(parsedInt)
		case reflect.String:
			paramKey.SetString(envVarValue)
		default:
			logger.Errorf(
				"Internal error: unhandled switch case for \"%T\" type of %q parameter (env var %q)",
				ptype, param, toml,
			)
			continue
		}
	}
	return params
}
