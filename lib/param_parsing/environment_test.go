package param_parsing

import (
	"os"
	"testing"
)

func TestGetParamsFromEnvironment_version_from_env(t *testing.T) {
	var params Params
	expected := "1.0.0_from_env"
	_ = os.Setenv("TF_VERSION", expected)
	params = initParams(params)
	params = GetParamsFromEnvironment(params)
	if params.Version != expected {
		t.Error("Determined version is not matchching. Got " + params.Version + ", expected " + expected)
	}
}
