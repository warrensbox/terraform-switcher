package param_parsing

import "os"

func GetParamsFromEnvironment(params Params) Params {
	params.Version = os.Getenv("TF_VERSION")
	return params
}
