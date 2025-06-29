//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

func prepare(tomlPath string) (Params, error) {
	if tomlPath == "" {
		tomlPath = "../../test-data/integration-tests/test_tfswitchtoml"
	}

	var params Params
	params.TomlDir = tomlPath
	logger = lib.InitLogger("DEBUG")
	return getParamsTOML(params)
}

func TestGetParamsTOML_BinaryPath(t *testing.T) {
	expected := "/usr/local/bin/terraform_from_toml"
	os.Setenv("BIN_DIR_FROM_TOML", "/usr/local/bin") // TOML value utilizes env var expansion
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.CustomBinaryPath
	if actual != expected {
		t.Errorf("%s not matching. Got %q, expected %q", "CustomBinaryPath", actual, expected)
	}
	os.Unsetenv("BIN_DIR_FROM_TOML")
}

func TestGetParamsTOML_InstallPath(t *testing.T) {
	expected := "/tmp"
	os.Setenv("INSTALL_DIR_FROM_TOML", "/tmp") // TOML value utilizes env var expansion
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.InstallPath
	if actual != expected {
		t.Errorf("%s not matching. Got %q, expected %q", "InstallPath", actual, expected)
	}
	os.Unsetenv("INSTALL_DIR_FROM_TOML")
}

func TestGetParamsTOML_Version(t *testing.T) {
	expected := "1.6.2"
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.Version
	if actual != expected {
		t.Errorf("%s not matching. Got %q, expected %q", "Version", actual, expected)
	}
}

func TestGetParamsTOML_Default_Version(t *testing.T) {
	expected := "1.5.4"
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.DefaultVersion
	if actual != expected {
		t.Errorf("%s not matching. Got %q, expected %q", "DefaultVersion", actual, expected)
	}
}

func TestGetParamsTOML_log_level(t *testing.T) {
	expected := "NOTICE"
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.LogLevel
	if actual != expected {
		t.Errorf("%s not matching. Got %q, expected %q", "LogLevel", actual, expected)
	}
}

func TestGetParamsTOML_no_color(t *testing.T) {
	expected := false
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.NoColor
	if actual != expected {
		t.Errorf("%s not matching. Got %v, expected %v", "NoColor", actual, expected)
	}
}

func TestGetParamsTOML_force_color(t *testing.T) {
	expected := false
	params, err := prepare("../../test-data/integration-tests/test_tfswitchtoml")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	actual := params.ForceColor
	if actual != expected {
		t.Errorf("%s not matching. Got %v, expected %v", "ForceColor", actual, expected)
	}
}

func TestGetParamsTOML_no_file(t *testing.T) {
	params, err := prepare("../../test-data/skip-integration-tests/test_no_file")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %q", params.Version)
	}
}

func TestGetParamsTOML_error_in_file(t *testing.T) {
	params, err := prepare("../../test-data/skip-integration-tests/test_tfswitchtoml_error")
	if err == nil {
		t.Fatalf("Expected error for reading erroneous TOML file. Got nil")
	}
	if params.Version != "" {
		t.Errorf("Expected empty version string. Got: %q", params.Version)
	}
}
