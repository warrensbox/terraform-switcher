package param_parsing

import (
	"testing"
)

func TestGetParamsFromTerraformVersion(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/integration-tests/test_terraform-version"
	params, _ = GetParamsFromTerraformVersion(params)
	expected := "0.11.0"
	if params.Version != expected {
		t.Errorf("Version from .terraform-version not read correctly. Got: %v, Expect: %v", params.Version, expected)
	}
}

func TestGetParamsFromTerraformVersion_no_file(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, _ = GetParamsFromTerraformVersion(params)
	if params.Version != "" {
		t.Errorf("Expected emtpy version string. Got: %v", params.Version)
	}
}
