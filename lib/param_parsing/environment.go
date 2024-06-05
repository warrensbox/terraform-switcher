package param_parsing

import "os"

func GetParamsFromEnvironment(params Params) Params {
	if envVersion := os.Getenv("TF_VERSION"); envVersion != "" {
		params.Version = envVersion
	}
	if envProduct := os.Getenv("TF_PRODUCT"); envProduct != "" {
		params.Product = envProduct
	}
	return params
}
