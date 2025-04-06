//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

func TestGetParamsFromEnvironment_arch_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "amd64_from_env"
	_ = os.Setenv("TF_ARCH", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_ARCH")
	if params.Arch != expected {
		t.Error("Determined arch is not matching. Got " + params.Arch + ", expected " + expected)
	}
}

func TestGetParamsFromEnvironment_version_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "1.0.0_from_env"
	_ = os.Setenv("TF_VERSION", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_VERSION")
	if params.Version != expected {
		t.Error("Determined version is not matching. Got " + params.Version + ", expected " + expected)
	}
}

func TestGetParamsFromEnvironment_default_version_from_env(t *testing.T) {
	var params Params
	expected := "1.0.0_from_env"
	_ = os.Setenv("TF_DEFAULT_VERSION", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_DEFAULT_VERSION")
	if params.DefaultVersion != expected {
		t.Error("Determined default version is not matching. Got " + params.DefaultVersion + ", expected " + expected)
	}
}

func TestGetParamsFromEnvironment_product_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "opentofu"
	_ = os.Setenv("TF_PRODUCT", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_PRODUCT")
	if params.Product != expected {
		t.Error("Determined product is not matching. Got " + params.Product + ", expected " + expected)
	}
}

func TestGetParamsFromEnvironment_bin_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "custom_binary_path_from_env"
	_ = os.Setenv("TF_BINARY_PATH", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_BINARY_PATH")
	if params.CustomBinaryPath != expected {
		t.Errorf("Determined custom binary path is not matching. Got %q, expected %q", params.CustomBinaryPath, expected)
	}
}

func TestGetParamsFromEnvironment_install_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "/custom_install_path_from_env"
	_ = os.Setenv("TF_INSTALL_PATH", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_INSTALL_PATH")
	if params.InstallPath != expected {
		t.Errorf("Determined custom install path is not matching. Got %q, expected %q", params.InstallPath, expected)
	}
}

func TestGetParamsFromEnvironment_log_level_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "DEBUG"
	_ = os.Setenv("TF_LOG_LEVEL", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_LOG_LEVEL")
	if params.LogLevel != expected {
		t.Errorf("Determined log level is not matching. Got %q, expected %q", params.LogLevel, expected)
	}
}
