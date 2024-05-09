package param_parsing

import (
	"testing"
)

func TestGetParamsFromTerraformVersion(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/integration-tests/test_terraform-version"
	params, err := GetParamsFromTerraformVersion(params)
	expected := "0.11.0"
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.Version != expected {
		t.Errorf("Version from .terraform-version not read correctly. Got: %v, Expect: %v", params.Version, expected)
	}
}

func TestGetParamsFromTerraformVersion_no_file(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, err := GetParamsFromTerraformVersion(params)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %v", params.Version)
	}
}
