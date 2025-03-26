package param_parsing

import (
	"os"
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

func prepare() Params {
	var params Params
	params.TomlDir = "../../test-data/integration-tests/test_tfswitchtoml"
	logger = lib.InitLogger("DEBUG")
	return params
}

func TestGetParamsTOML_BinaryPath(t *testing.T) {
	expected := "/usr/local/bin/terraform_from_toml"
	os.Setenv("BIN_DIR_FROM_TOML", "/usr/local/bin")
	params := prepare()
	params, err := getParamsTOML(params)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.CustomBinaryPath != expected {
		t.Errorf("BinaryPath not matching. Got %v, expected %v", params.CustomBinaryPath, expected)
	}
	os.Unsetenv("BIN_DIR_FROM_TOML")
}

func TestGetParamsTOML_InstallPath(t *testing.T) {
	expected := "/tmp"
	os.Setenv("INSTALL_DIR_FROM_TOML", "/tmp")
	params := prepare()
	params, err := getParamsTOML(params)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if params.InstallPath != expected {
		t.Errorf("InstallPath not matching. Got %q, expected %q", params.InstallPath, expected)
	}
	os.Unsetenv("INSTALL_DIR_FROM_TOML")
}

func TestGetParamsTOML_Version(t *testing.T) {
	expected := "1.6.2"
	params := prepare()
	params, err := getParamsTOML(params)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.Version != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.Version, expected)
	}
}

func TestGetParamsTOML_Default_Version(t *testing.T) {
	expected := "1.5.4"
	params := prepare()
	params, err := getParamsTOML(params)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.DefaultVersion != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.DefaultVersion, expected)
	}
}

func TestGetParamsTOML_log_level(t *testing.T) {
	expected := "NOTICE"
	params := prepare()
	params, err := getParamsTOML(params)
	if err != nil {
		t.Fatalf("Got error '%s'", err)
	}
	if params.LogLevel != expected {
		t.Errorf("Version not matching. Got %v, expected %v", params.LogLevel, expected)
	}
}

func TestGetParamsTOML_no_file(t *testing.T) {
	var params Params
	params.TomlDir = "../../test-data/skip-integration-tests/test_no_file"
	params, err := getParamsTOML(params)
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
	params.TomlDir = "../../test-data/skip-integration-tests/test_tfswitchtoml_error"
	params, err := getParamsTOML(params)
	if err == nil {
		t.Errorf("Expected error for reading erroneous toml file. Got nil")
	}
	if params.Version != "" {
		t.Errorf("Version should be empty")
	}
}
