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
}

func TestGetParameters_set_product_entity(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	os.Args = []string{"cmd", "--product=opentofu"}
	params := GetParameters()

	if expected := "opentofu"; params.ProductEntity.GetId() != expected {
		t.Errorf("Incorrect product entity set on params. Expected: %q, Actual: %q", expected, params.ProductEntity.GetId())
	}

	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--product=terraform"}
	params = GetParameters()

	if expected := "terraform"; params.ProductEntity.GetId() != expected {
		t.Errorf("Incorrect product entity set on params. Expected: %q, Actual: %q", expected, params.ProductEntity.GetId())
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

func TestGetParameters_check_config_precedence(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	os.Args = []string{"cmd", "--chdir=../../test-data/integration-tests/test_precedence"}
	parameters := GetParameters()
	expected := "0.11.3"
	if parameters.Version != expected {
		t.Error("Version Param was not as expected. Actual: " + parameters.Version + ", Expected: " + expected)
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}
