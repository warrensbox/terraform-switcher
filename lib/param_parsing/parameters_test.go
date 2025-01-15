package param_parsing

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
)

func TestGetParameters_arch_from_args(t *testing.T) {
	expected := "amd64args"
	os.Args = []string{"cmd", expected}
	params := GetParameters()
	actual := params.Arch
	if actual != expected {
		t.Error("Arch Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

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
	params := Params{}
	params = initParams(params)
	params.TomlDir = "../../test-data/integration-tests/test_tfswitchtoml"
	params = populateParams(params)

	actual := params.ChDirPath
	if actual != expected {
		t.Error("ChDir Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "/usr/local/bin/terraform_from_toml"
	actual = params.CustomBinaryPath
	if actual != expected {
		t.Error("CustomBinaryPath Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "amd64"
	actual = params.Arch
	if actual != expected {
		t.Error("Arch Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "1.6.2"
	actual = params.Version
	if actual != expected {
		t.Error("Version Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "opentofu"
	actual = params.Product
	if actual != expected {
		t.Error("Product Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func TestGetParameters_toml_params_are_overridden_by_cli(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	expected := "../../test-data/integration-tests/test_tfswitchtoml"
	os.Args = []string{"cmd", "--chdir=" + expected, "--bin=/usr/test/bin", "--product=terraform", "--arch=amd64", "1.6.0"}
	params := Params{}
	params = initParams(params)
	params.TomlDir = expected
	params = populateParams(params)

	actual := params.ChDirPath
	if actual != expected {
		t.Error("ChDir Param was not parsed correctly. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "/usr/test/bin"
	actual = params.CustomBinaryPath
	if actual != expected {
		t.Error("CustomBinaryPath Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "1.6.0"
	actual = params.Version
	if actual != expected {
		t.Error("Version Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "terraform"
	actual = params.Product
	if actual != expected {
		t.Error("Product Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "amd64"
	actual = params.Arch
	if actual != expected {
		t.Error("Arch Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
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
	params := Params{}
	params = initParams(params)
	params.TomlDir = expected
	params = populateParams(params)

	installLocation := lib.GetInstallLocation(params.InstallPath)
	product := lib.GetProductById(params.Product)
	if product == nil {
		t.Error("Nil product returned")
	}
	installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, product.GetVersionPrefix()+params.Version))
	// Make sure the file tfswitch WOULD download is absent
	_ = os.Remove(installFileVersionPath)
	lib.InstallProductVersion(product, params.DryRun, params.Version, params.CustomBinaryPath, params.InstallPath, params.MirrorURL, params.Arch)
	if lib.FileExistsAndIsNotDir(installFileVersionPath) {
		t.Error("Dry run should NOT download any files.")
	}
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func writeTestFile(t *testing.T, basePath string, fileName string, fileContent string) {
	fullPath := filepath.Join(basePath, fileName)
	if err := os.WriteFile(fullPath, []byte(fileContent), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(fullPath)
	})
}

func checkExpectedPrecedenceVersion(t *testing.T, expectedVersion string, expectedDefaultVersion string) {
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--chdir=../../test-data/skip-integration-tests/test_precedence"}
	parameters := Params{}
	parameters = initParams(parameters)
	parameters.TomlDir = "../../test-data/skip-integration-tests/test_precedence"
	parameters = populateParams(parameters)
	if parameters.Version != expectedVersion {
		t.Error("Version Param was not as expected. Actual: " + parameters.Version + ", Expected: " + expectedVersion)
	}
}

func TestGetParameters_check_config_precedence(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	chDir := "../../test-data/skip-integration-tests/test_precedence"
	checkExpectedPrecedenceVersion(t, "", "")

	// Create TfSwitch TOML
	tfSwitchTOMLContent := `
bin = "/usr/local/bin/terraform_from_toml"
version = "0.11.3"
default-version = "0.12.1"
`
	writeTestFile(t, chDir, ".tfswitch.toml", tfSwitchTOMLContent)
	checkExpectedPrecedenceVersion(t, "0.11.3", "0.12.1")

	// Create tfswitchrc file
	writeTestFile(t, chDir, ".tfswitchrc", "0.10.5")
	checkExpectedPrecedenceVersion(t, "0.10.5", "0.12.1")

	// Create terraform-version file
	writeTestFile(t, chDir, ".terraform-version", "0.11.0")
	checkExpectedPrecedenceVersion(t, "0.11.0", "0.12.1")

	// Create terraform file
	terraformFileContent := `
terraform {
	required_version = "0.14.1"
}
`
	writeTestFile(t, chDir, "main.tf", terraformFileContent)
	checkExpectedPrecedenceVersion(t, "0.14.1", "0.12.1")

	// Create terraform file
	terragruntContent := `
terraform_version_constraint = ">= 0.13, < 0.14"
`
	writeTestFile(t, chDir, "terragrunt.hcl", terragruntContent)
	checkExpectedPrecedenceVersion(t, "0.13.7", "0.12.1")

	// Test with environment variable
	os.Setenv("TF_VERSION", "0.11.31.env")
	checkExpectedPrecedenceVersion(t, "0.11.31.env", "0.12.1")
	os.Setenv("TF_DEFAULT_VERSION", "0.12.2.env")
	checkExpectedPrecedenceVersion(t, "0.11.31.env", "0.12.2.env")

	// Test passing command line argument to override all
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--chdir=../../test-data/skip-integration-tests/test_precedence", "--default=0.12.3", "1.4.5"}
	parameters := GetParameters()
	if expectedVersion := "1.4.5"; parameters.Version != expectedVersion {
		t.Error("Version Param was not as expected. Actual: " + parameters.Version + ", Expected: " + expectedVersion)
	}
	if expectedVersion := "0.12.3"; parameters.DefaultVersion != expectedVersion {
		t.Error("DefaultVersion Param was not as expected. Actual: " + parameters.DefaultVersion + ", Expected: " + expectedVersion)
	}

	os.Unsetenv("TF_VERSION")
	os.Unsetenv("TF_DEFAULT_VERSION")
}

func checkExpectedPrecedenceProduct(t *testing.T, baseDir string, expectedProduct lib.Product) {
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", fmt.Sprintf("--chdir=%s", baseDir)}
	parameters := Params{}
	parameters = initParams(parameters)
	parameters.TomlDir = baseDir
	parameters = populateParams(parameters)
	if parameters.Product != expectedProduct.GetId() {
		t.Error("Product Param was not as expected. Actual: " + parameters.Product + ", Expected: " + expectedProduct.GetId())
	}

	// Check ProductEntity
	if parameters.ProductEntity.GetId() != expectedProduct.GetId() {
		t.Error("ProductEntity Param was not as expected. Actual: " + parameters.Product + ", Expected: " + expectedProduct.GetId())
	}
}

func TestGetParameters_check_product_precedence(t *testing.T) {
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	terraformProduct := lib.GetProductById("terraform")
	openTofuProduct := lib.GetProductById("opentofu")

	t.Log("Testing without configuration")
	checkExpectedPrecedenceProduct(t, tempDir, terraformProduct)

	// Create TfSwitch TOML
	tfSwitchTOMLContent := `
bin = "/usr/local/bin/terraform_from_toml"
product = "opentofu"
`
	writeTestFile(t, tempDir, ".tfswitch.toml", tfSwitchTOMLContent)
	t.Log("Testing with TOML configuration")
	checkExpectedPrecedenceProduct(t, tempDir, openTofuProduct)

	// Test passing command line argument to override all
	os.Setenv("TF_PRODUCT", "terraform")
	t.Log("Testing with environment variable")
	checkExpectedPrecedenceProduct(t, tempDir, terraformProduct)

	os.Unsetenv("TF_VERSION")
}
