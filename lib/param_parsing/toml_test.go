package param_parsing

import (
	"testing"
)

func prepare() Params {
	var params Params
	params.ChDirPath = "../../test-data/test_tfswitchtoml"
	return params
}

func TestGetParamsTOML_BinaryPath(t *testing.T) {
	expected := "/usr/local/bin/terraform_from_toml"
	params := prepare()
	params = getParamsTOML(params)
	if params.CustomBinaryPath != expected {
		t.Log("Actual:", params.CustomBinaryPath)
		t.Log("Expected:", expected)
		t.Error("BinaryPath not matching")
	}
}

func TestGetParamsTOML_Version(t *testing.T) {
	expected := "0.11.3_toml"
	params := prepare()
	params = getParamsTOML(params)
	if params.Version != expected {
		t.Log("Actual:", params.Version)
		t.Log("Expected:", expected)
		t.Error("Version not matching")
	}
}
