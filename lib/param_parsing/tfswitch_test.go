package param_parsing

import (
	"testing"
)

func TestGetParamsFromTfSwitch(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/test_tfswitchrc"
	params = GetParamsFromTfSwitch(params)
	expected := "0.10.5_tfswitch"
	if params.Version != expected {
		t.Error("Version from tfswitchrc not read correctly. Actual: " + params.Version + ", Expected: " + expected)
	}
}
