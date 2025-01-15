package param_parsing

import "os"

func GetParamsFromEnvironment(params Params) Params {
	if envArch := os.Getenv("TF_ARCH"); envArch != "" {
		params.Arch = envArch
		logger.Debugf("Using architecture from environment variable \"TF_ARCH\": %q", envArch)
	}
	if envVersion := os.Getenv("TF_VERSION"); envVersion != "" {
		params.Version = envVersion
		logger.Debugf("Using version from environment variable \"TF_VERSION\": %q", envVersion)
	}
	if envDefaultVersion := os.Getenv("TF_DEFAULT_VERSION"); envDefaultVersion != "" {
		params.DefaultVersion = envDefaultVersion
		logger.Debugf("Using default-version from environment variable \"TF_DEFAULT_VERSION\": %q", envDefaultVersion)
	}
	if envProduct := os.Getenv("TF_PRODUCT"); envProduct != "" {
		params.Product = envProduct
		logger.Debugf("Using product from environment variable \"TF_PRODUCT\": %q", envProduct)
	}
	return params
}
