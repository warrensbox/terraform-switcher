package param_parsing

import (
	"testing"
)

func TestGetParamsFromTfSwitch(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/test_tfswitchrc"
	params, _ = GetParamsFromTfSwitch(params)
	expected := "0.10.5"
	if params.Version != expected {
		t.Error("Version from tfswitchrc not read correctly. Actual: " + params.Version + ", Expected: " + expected)
	}
}

func TestGetParamsFromTfSwitch_no_file(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/test_no_file"
	params, _ = GetParamsFromTfSwitch(params)
	if params.Version != "" {
		t.Errorf("Expected emtpy version string. Got: %v", params.Version)
	}
}
