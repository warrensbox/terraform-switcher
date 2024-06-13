package param_parsing

import (
	"os"
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

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

func TestGetParamsFromEnvironment_product_from_env(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	expected := "opentofu"
	_ = os.Setenv("TF_PRODUCT", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	_ = os.Unsetenv("TF_PRODUCT")
	if params.Product != expected {
		t.Error("Determined version is not matching. Got " + params.Product + ", expected " + expected)
	}
}
