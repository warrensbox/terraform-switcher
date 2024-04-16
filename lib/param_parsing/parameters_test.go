package param_parsing

import (
	"github.com/pborman/getopt"
	"os"
	"testing"
)

func TestGetParameters_version_from_args(t *testing.T) {
	expected := "0.13args"
	os.Args = []string{"cmd", expected}
	params := GetParameters()
	actual := params.Version
	if actual != expected {
		t.Error("Version Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func TestGetParameters_params_are_overridden_by_toml_file(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	expected := "../../test-data/integration-tests/test_tfswitchtoml"
	os.Args = []string{"cmd", "--chdir=" + expected}
	params := GetParameters()
	actual := params.ChDirPath
	if actual != expected {
		t.Error("ChDir Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "/usr/local/bin/terraform_from_toml"
	actual = params.CustomBinaryPath
	if actual != expected {
		t.Error("CustomBinaryPath Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}
	expected = "0.11.4"
	actual = params.Version
	if actual != expected {
		t.Error("Version Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}
}
func TestGetParameters_toml_params_are_overridden_by_cli(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	expected := "../../test-data/integration-tests/test_tfswitchtoml"
	os.Args = []string{"cmd", "--chdir=" + expected, "--bin=/usr/test/bin"}
	params := GetParameters()
	actual := params.ChDirPath
	if actual != expected {
		t.Error("ChDir Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "/usr/test/bin"
	actual = params.CustomBinaryPath
	if actual != expected {
		t.Error("CustomBinaryPath Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}
	expected = "0.11.4"
	actual = params.Version
	if actual != expected {
		t.Error("Version Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}
}

func TestGetParameters_check_config_precedence(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	os.Args = []string{"cmd", "--chdir=../../test-data/integration-tests/test_precedence"}
	parameters := prepare()
	parameters = GetParameters()
	expected := "0.11.3"
	if parameters.Version != expected {
		t.Error("Version Param was not as expected. Actual: " + parameters.Version + ", Expected: " + expected)
	}
}
