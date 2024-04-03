package param_parsing

import (
	"testing"
)

func TestGetParamsFromTerraformVersion(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/test_terraform-version"
	params = GetParamsFromTerraformVersion(params)
	expected := "0.11.0"
	if params.Version != expected {
		t.Errorf("Version from .terraform-version not read correctly. Got: %v, Expect: %v", params.Version, expected)
	}
}
