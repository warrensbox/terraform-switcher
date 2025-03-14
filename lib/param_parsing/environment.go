package param_parsing

import "os"

func GetParamsFromEnvironment(params Params) Params {
	if envVar := os.Getenv("TF_ARCH"); envVar != "" {
		params.Arch = envVar
		logger.Debugf("Using architecture from environment variable \"TF_ARCH\": %q", envVar)
	}
	if envVar := os.Getenv("TF_VERSION"); envVar != "" {
		params.Version = envVar
		logger.Debugf("Using version from environment variable \"TF_VERSION\": %q", envVar)
	}
	if envVar := os.Getenv("TF_DEFAULT_VERSION"); envVar != "" {
		params.DefaultVersion = envVar
		logger.Debugf("Using default version from environment variable \"TF_DEFAULT_VERSION\": %q", envVar)
	}
	if envVar := os.Getenv("TF_PRODUCT"); envVar != "" {
		params.Product = envVar
		logger.Debugf("Using product from environment variable \"TF_PRODUCT\": %q", envVar)
	}
	if envVar := os.Getenv("TF_BINARY_PATH"); envVar != "" {
		params.CustomBinaryPath = envVar
		logger.Debugf("Using custom binary path from environment variable \"TF_BINARY_PATH\": %q", envVar)
	}
	if envVar := os.Getenv("TF_INSTALL_PATH"); envVar != "" {
		params.InstallPath = envVar
		logger.Debugf("Using custom install path from environment variable \"TF_INSTALL_PATH\": %q", envVar)
	}
	return params
}
