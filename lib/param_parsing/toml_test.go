package param_parsing

import (
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

func prepare() Params {
	var params Params
	params.ChDirPath = "../../test-data/integration-tests/test_tfswitchtoml"
	logger = lib.InitLogger("DEBUG")
	return params
}

func TestGetParamsTOML_BinaryPath(t *testing.T) {
	expected := "/usr/local/bin/terraform_from_toml"
	params := prepare()
	params, err := getParamsTOML(params, params.ChDirPath)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.CustomBinaryPath != expected {
		t.Errorf("BinaryPath not matching. Got %v, expected %v", params.CustomBinaryPath, expected)
	}
}

func TestGetParamsTOML_Version(t *testing.T) {
	expected := "0.11.4"
	params := prepare()
	params, err := getParamsTOML(params, params.ChDirPath)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.Version != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.Version, expected)
	}
}

func TestGetParamsTOML_log_level(t *testing.T) {
	expected := "NOTICE"
	params := prepare()
	params, err := getParamsTOML(params, params.ChDirPath)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.LogLevel != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.LogLevel, expected)
	}
}

func TestGetParamsTOML_no_file(t *testing.T) {
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, err := getParamsTOML(params, params.ChDirPath)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %v", params.Version)
	}
}

func TestGetParamsTOML_error_in_file(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params.ChDirPath = "../../test-data/skip-integration-tests/test_tfswitchtoml_error"
	params, err := getParamsTOML(params, params.ChDirPath)
	if err == nil {
		t.Errorf("Expected error for reading erroneous toml file. Got nil")
	}
	if params.Version != "" {
		t.Errorf("Version should be empty")
	}
}
