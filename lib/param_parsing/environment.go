package param_parsing

import "os"

func GetParamsFromEnvironment(params Params) Params {
	if envVersion := os.Getenv("TF_VERSION"); envVersion != "" {
		params.Version = envVersion
		logger.Debugf("Found environment variable TF_VERSION: %q", envVersion)
	}
	if envProduct := os.Getenv("TF_PRODUCT"); envProduct != "" {
		params.Product = envProduct
		logger.Debugf("Found environment variable TF_PRODUCT: %q", envProduct)
	}
	return params
}
