//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"testing"
)

func TestGetParamsFromTfSwitch(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/integration-tests/test_tfswitchrc"
	params, err := GetParamsFromTfSwitch(params)
	expected := "0.10.5"
	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}
	if params.Version != expected {
		t.Error("Version from tfswitchrc not read correctly. Actual: " + params.Version + ", Expected: " + expected)
	}
}

func TestGetParamsFromTfSwitch_no_file(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, err := GetParamsFromTfSwitch(params)
	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %v", params.Version)
	}
}

func TestGetParamsFromTfSwitch_no_version(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_version"
	params, err := GetParamsFromTfSwitch(params)
	if err != nil {
		t.Errorf("Expected no error. Got: %v", err)
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %v", params.Version)
	}
}
