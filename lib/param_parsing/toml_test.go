package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"testing"
)

func prepare() Params {
	var params Params
	params.ChDirPath = "../../test-data/test_tfswitchtoml"
	logger = lib.InitLogger("DEBUG")
	return params
}

func TestGetParamsTOML_BinaryPath(t *testing.T) {
	expected := "/usr/local/bin/terraform_from_toml"
	params := prepare()
	params = getParamsTOML(params)
	if params.CustomBinaryPath != expected {
		t.Errorf("BinaryPath not matching. Got %v, expected %v", params.CustomBinaryPath, expected)
	}
}

func TestGetParamsTOML_Version(t *testing.T) {
	expected := "0.11.4"
	params := prepare()
	params = getParamsTOML(params)
	if params.Version != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.Version, expected)
	}
}

func TestGetParamsTOML_log_level(t *testing.T) {
	expected := "NOTICE"
	params := prepare()
	params = getParamsTOML(params)
	if params.LogLevel != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.LogLevel, expected)
	}
}
