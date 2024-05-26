package param_parsing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
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
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func TestGetParameters_toml_params_are_overridden_by_cli(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
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
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func TestGetParameters_dry_run_wont_download_anything(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	expected := "../../test-data/integration-tests/test_versiontf"
	os.Args = []string{"cmd", "--chdir=" + expected, "--bin=/tmp", "--dry-run"}
	params := GetParameters()
	installLocation := lib.GetInstallLocation(params.InstallPath)
	installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, lib.VersionPrefix+params.Version))
	// Make sure the file tfswitch WOULD download is absent
	_ = os.Remove(installFileVersionPath)
	lib.InstallVersion(params.DryRun, params.Version, params.CustomBinaryPath, params.InstallPath, params.MirrorURL)
	if lib.FileExistsAndIsNotDir(installFileVersionPath) {
		t.Error("Dry run should NOT download any files.")
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func checkExpectedPrecedenceVersion(t *testing.T, expectedVersion string) {
	os.Args = []string{"cmd", "--chdir=../../test-data/integration-tests/test_precedence"}
	parameters := GetParameters()
	expected := "0.13.7"
	if parameters.Version != expected {
		t.Error("Version Param was not as expected. Actual: " + parameters.Version + ", Expected: " + expected)
	}
}

func TestGetParameters_check_config_precedence(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	chDir := "../../test-data/integration-tests/test_precedence"
	checkExpectedPrecedenceVersion(t, "")

	// Create TfSwitch TOML
	tfSwitchTOMLContent := `
bin = "/usr/local/bin/terraform_from_toml"
version = "0.11.3"
`
	if err := os.WriteFile(filepath.Join(chDir, ".tfswitch.toml"), []byte(tfSwitchTOMLContent), 0666); err != nil {
		t.Error(err)
	}
	checkExpectedPrecedenceVersion(t, "0.11.3")

	// Create tfswitchrc file
	if err := os.WriteFile(filepath.Join(chDir, ".tfswitchrc"), []byte("0.10.5"), 0666); err != nil {
		t.Error(err)
	}
	checkExpectedPrecedenceVersion(t, "0.10.5")

	// Create terraform-version file
	if err := os.WriteFile(filepath.Join(chDir, ".terraform-version"), []byte("0.11.0"), 0666); err != nil {
		t.Error(err)
	}
	checkExpectedPrecedenceVersion(t, "0.11.0")

	// Create terraform file
	terraformFileContent := `
terraform {
	required_version = "0.14.1"
}
`
	if err := os.WriteFile(filepath.Join(chDir, "main.tf"), []byte(terraformFileContent), 0666); err != nil {
		t.Error(err)
	}
	checkExpectedPrecedenceVersion(t, "0.14.1")

	// Create terraform file
	terragruntContent := `
terraform_version_constraint = ">= 0.13, < 0.14"
`
	if err := os.WriteFile(filepath.Join(chDir, "terragrunt.hcl"), []byte(terragruntContent), 0666); err != nil {
		t.Error(err)
	}
	checkExpectedPrecedenceVersion(t, "0.13.9")

	// Test with environment variable
	os.Setenv("TERRAFORM_VERSION", "0.11.31.env")
	checkExpectedPrecedenceVersion(t, "0.11.31.env")
	os.Unsetenv("TERRAFORM_VERSION")
}
