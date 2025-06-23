//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
)

var ansiCodesRegex = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

func TestGetParameters_arch_from_args(t *testing.T) {
	expected := "arch_from_args"
	os.Args = []string{"cmd", "--arch=" + expected}
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

	os.Setenv("BIN_DIR_FROM_TOML", "/usr/local/bin")
	os.Setenv("INSTALL_DIR_FROM_TOML", "/tmp")

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

	expected = "/tmp"
	actual = params.InstallPath
	if actual != expected {
		t.Error("InstallPath Param was not as expected. Actual: " + actual + ", Expected: " + expected)
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

	expected = "1.5.4"
	actual = params.DefaultVersion
	if actual != expected {
		t.Error("DefaultVersion Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "opentofu"
	actual = params.Product
	if actual != expected {
		t.Error("Product Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	expected = "NOTICE"
	actual = params.LogLevel
	if actual != expected {
		t.Error("LogLevel Param was not as expected. Actual: " + actual + ", Expected: " + expected)
	}

	os.Unsetenv("BIN_DIR_FROM_TOML")
	os.Unsetenv("INSTALL_DIR_FROM_TOML")

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
}

func TestGetParameters_toml_params_are_overridden_by_cli(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	expected := "../../test-data/integration-tests/test_tfswitchtoml"
	os.Args = []string{"cmd", "--chdir=" + expected, "--bin=/usr/test/bin", "--product=terraform", "--arch=arch_from_args", "1.6.0"}
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

	expected = "arch_from_args"
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
	err := lib.InstallProductVersion(product, params.DryRun, params.Version, params.CustomBinaryPath, params.InstallPath, params.MirrorURL, params.Arch)
	if err != nil || lib.FileExistsAndIsNotDir(installFileVersionPath) {
		t.Error("Dry run should NOT install any files.")
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

//nolint:revive // FIXME: the 3rd argument is not used %-/ 10-Mar-2025
func checkExpectedPrecedenceVersion(t *testing.T, expectedVersion string, _ string) {
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

func TestVersionFlagOutput(t *testing.T) {
	flagName := "--version"
	expectedOutput := "Version: "
	goCommandArgs := []string{"run", "../../main.go", flagName}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	if !strings.HasPrefix(string(out), expectedOutput) {
		t.Fatalf("Expected %q, got: %q", expectedOutput, string(out))
	}

	t.Logf("Success: %q", string(out))
}

func TestHelpFlagOutput(t *testing.T) {
	flagName := "--help"
	expectedOutput := "Usage: "
	goCommandArgs := []string{"run", "../../main.go", flagName}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	if !strings.HasPrefix(string(out), expectedOutput) {
		t.Fatalf("Expected %q, got: %q", expectedOutput, string(out))
	}

	t.Logf("Success: %q", string(out))
}

func TestDryRunFlagOutput(t *testing.T) {
	flagName := "--dry-run"
	testVersion := "1.10.5"
	expectedOutput := fmt.Sprintf(" INFO [DRY-RUN] Would have attempted to install version %q  \n", testVersion)
	goCommandArgs := []string{"run", "../../main.go", flagName, testVersion}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	re := regexp.MustCompile(ansiCodesRegex)
	outNoANSI := func(str string) string {
		return re.ReplaceAllString(str, "")
	}(string(out))

	if !strings.HasSuffix(outNoANSI, expectedOutput) {
		t.Fatalf("Expected %q, got: %q", expectedOutput, outNoANSI)
	}

	t.Logf("Success: %q", outNoANSI)
}

func TestNoColorFlagOutput(t *testing.T) {
	flagName := "--no-color"
	goCommandArgs := []string{"run", "../../main.go", flagName, "--dry-run", "1.10.5"}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	matched, err := regexp.MatchString(ansiCodesRegex, string(out))
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	if matched {
		t.Fatalf("Expected no ANSI color codes in output, but found some: %q", string(out))
	} else {
		t.Log("Success: no ANSI color codes in output")
	}
}

func TestForceColorFlagOutput(t *testing.T) {
	flagName := "--force-color"
	goCommandArgs := []string{"run", "../../main.go", flagName, "--dry-run", "1.10.5"}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	matched, err := regexp.MatchString(ansiCodesRegex, string(out))
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	if !matched {
		t.Fatalf("Expected ANSI color codes in output, but found none: %q", string(out))
	} else {
		t.Log("Success: found ANSI color codes in output")
	}
}

func TestNoAndForceColorFlagsOutput(t *testing.T) {
	flagNameForceColor := "--force-color"
	flagNameNoColor := "--no-color"

	expectedOutput := "FATAL (init) Cannot force color and disable color at the same time. Please choose either of them."

	goCommandArgs := []string{"run", "../../main.go", "--dry-run", flagNameForceColor, flagNameNoColor, "1.10.5"}

	t.Logf("Testing %q and %q flags both present", flagNameForceColor, flagNameNoColor)

	out, _ := exec.Command("go", goCommandArgs...).CombinedOutput()

	if !strings.Contains(string(out), expectedOutput) {
		t.Errorf("Expected %q, got: %q", expectedOutput, out)
	} else {
		t.Logf("Success: %q", expectedOutput)
	}
}
