//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/pborman/getopt/v2"
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

	expectedBool := false
	actualBool := params.NoColor
	if actualBool != expectedBool {
		t.Errorf("NoColor Param was not as expected. Actual: %v, Expected: %v", actualBool, expectedBool)
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
	err := lib.InstallProductVersion(product, params.DryRun, params.ShowRequiredFlag, params.Version, params.CustomBinaryPath, params.InstallPath, params.MirrorURL, params.MirrorDownloadURL, params.Arch)
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
		t.Errorf("Expected %q, got: %q", expectedOutput, string(out))
	} else {
		t.Logf("Success: %q", string(out))
	}
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
		t.Errorf("Expected %q, got: %q", expectedOutput, string(out))
	} else {
		t.Logf("Success: %q", string(out))
	}
}

func TestDryRunFlagOutput(t *testing.T) {
	flagName := "--dry-run"
	testVersion := "1.10.5"
	expectedOutput := " INFO [DRY-RUN] No changes will be made  \n"
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

	if !strings.Contains(outNoANSI, expectedOutput) {
		t.Errorf("Expected %q, got: %q", expectedOutput, outNoANSI)
	} else {
		t.Logf("Success: %q", outNoANSI)
	}
}

func TestShowRequiredFlagOutput(t *testing.T) {
	flagName := "--show-required"
	testVersion := "1.10.5"
	expectedOutput := fmt.Sprintf("%s\n", testVersion)
	goCommandArgs := []string{"run", "../../main.go", flagName, testVersion}

	t.Logf("Testing %q flag output", flagName)

	out, err := exec.Command("go", goCommandArgs...).Output()
	if err != nil {
		t.Fatalf("Unexpected failure: \"%v\", output: %q", err, string(out))
	}

	re := regexp.MustCompile(ansiCodesRegex)
	outNoANSI := func(str string) string {
		return re.ReplaceAllString(str, "")
	}(string(out))

	if !strings.HasSuffix(outNoANSI, expectedOutput) {
		t.Errorf("Expected %q, got: %q", expectedOutput, outNoANSI)
	} else {
		t.Logf("Success: %q", outNoANSI)
	}
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
		t.Errorf("Expected no ANSI color codes in output, but found some: %q", string(out))
	} else {
		t.Log("Success: no ANSI color codes in output")
	}
}

func TestForceColorFlagOutput(t *testing.T) {
	flagName := "--force-color"
	if color.SupportColor() {
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
			t.Errorf("Expected ANSI color codes in output, but found none: %q", string(out))
		} else {
			t.Log("Success: found ANSI color codes in output")
		}
	} else {
		t.Skipf("Skipping test for %q flag as terminal doesn't support colors", flagName)
	}
}

func TestNoAndForceColorFlagsOutput(t *testing.T) {
	flagNameForceColor := "--force-color"
	flagNameNoColor := "--no-color"

	expectedOutput := "FATAL (init) Cannot force color and disable color at the same time. Please choose either of them."

	goCommandArgs := []string{"run", "../../main.go", "--dry-run", flagNameForceColor, flagNameNoColor, "1.10.5"}

	t.Logf("Testing %q and %q flags both present", flagNameForceColor, flagNameNoColor)

	out, err := exec.Command("go", goCommandArgs...).CombinedOutput() // nolint:errcheck // We want to test the output even if it fails
	if err == nil {
		t.Fatalf("Expected an error, but got none. Output: %q", string(out))
	}

	if !strings.Contains(string(out), expectedOutput) {
		t.Errorf("Expected %q, got: %q", expectedOutput, out)
	} else {
		t.Logf("Success: %q", expectedOutput)
	}
}

func checkExpectedHTTPTimeoutPrecedence(t *testing.T, chDir string, expected int64) {
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--chdir=" + chDir}
	parameters := Params{}
	parameters = initParams(parameters)
	parameters.TomlDir = chDir
	parameters = populateParams(parameters)
	if parameters.HTTPTimeout != expected {
		t.Errorf("HTTPTimeout Param was not as expected. Actual: %d, Expected: %d", parameters.HTTPTimeout, expected)
	}
}

// TestGetParameters_check_http_timeout_precedence verifies the full HTTP
// timeout precedence chain: CLI flag > environment variable > TOML > default.
func TestGetParameters_check_http_timeout_precedence(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	chDir := "../../test-data/skip-integration-tests/test_precedence"

	// Default value when no TOML, environment variable or CLI flag is set
	checkExpectedHTTPTimeoutPrecedence(t, chDir, 60)

	// TOML overrides default
	writeTestFile(t, chDir, ".tfswitch.toml", "http-timeout = 90\n")
	checkExpectedHTTPTimeoutPrecedence(t, chDir, 90)

	// Environment variable overrides TOML
	os.Setenv("TF_HTTP_TIMEOUT", "150")
	checkExpectedHTTPTimeoutPrecedence(t, chDir, 150)

	// CLI flag overrides environment variable and TOML
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--chdir=" + chDir, "--http-timeout=30"}
	parameters := Params{}
	parameters = initParams(parameters)
	parameters.TomlDir = chDir
	parameters = populateParams(parameters)
	if expected := int64(30); parameters.HTTPTimeout != expected {
		t.Errorf("HTTPTimeout Param was not as expected. Actual: %d, Expected: %d", parameters.HTTPTimeout, expected)
	}

	os.Unsetenv("TF_HTTP_TIMEOUT")
}

// TestGetParameters_http_timeout_zero_disables_timeout is a regression test for
// the documented (but easy to overlook) behavior that `0` is accepted as a
// valid HTTP timeout value, and disables the timeout entirely rather than
// making requests fail instantly.
func TestGetParameters_http_timeout_zero_disables_timeout(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})
	os.Args = []string{"cmd", "--http-timeout=0"}
	parameters := GetParameters()
	if expected := int64(0); parameters.HTTPTimeout != expected {
		t.Errorf("HTTPTimeout Param was not as expected. Actual: %d, Expected: %d", parameters.HTTPTimeout, expected)
	}
}

// TestGetParameters_http_timeout_is_enforced_on_stalled_requests is an end-to-end
// regression test proving that a configured HTTP timeout (here, via the CLI flag,
// exactly as a user would set it) is actually wired into `http.DefaultClient` and
// terminates a request to a server that never responds.
func TestGetParameters_http_timeout_is_enforced_on_stalled_requests(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	originalDefaultClient := http.DefaultClient
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
		http.DefaultClient = originalDefaultClient
	})

	// A server that stalls indefinitely, simulating a hung/unresponsive remote endpoint.
	stall := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		<-stall
	}))
	t.Cleanup(server.Close)

	// Configure the shortest meaningful HTTP timeout via the CLI flag.
	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--http-timeout=1"}
	parameters := GetParameters()
	if expected := int64(1); parameters.HTTPTimeout != expected {
		t.Fatalf("HTTPTimeout Param was not as expected. Actual: %d, Expected: %d", parameters.HTTPTimeout, expected)
	}

	start := time.Now()
	_, err := http.DefaultClient.Get(server.URL) //nolint:noctx,bodyclose // request is expected to error out before a body exists
	elapsed := time.Since(start)

	// Unblock the handler unconditionally (before any assertion that could fail/return
	// early) so `server.Close()` registered above doesn't hang waiting for it to finish.
	close(stall)

	if err == nil {
		t.Fatal("Expected a timeout error for a stalled request, got none")
	}

	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Errorf("Expected a timeout error, got: %v", err)
	}

	if elapsed >= 5*time.Second {
		t.Errorf("Request took too long to time out: %v (expected close to the configured 1s timeout)", elapsed)
	}
}

// TestGetParameters_http_timeout_does_not_abort_requests_within_the_limit complements
// the above test: it proves a configured timeout does not prematurely abort requests
// that complete well within it.
func TestGetParameters_http_timeout_does_not_abort_requests_within_the_limit(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	originalDefaultClient := http.DefaultClient
	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
		http.DefaultClient = originalDefaultClient
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	getopt.CommandLine = getopt.New()
	os.Args = []string{"cmd", "--http-timeout=5"}
	parameters := GetParameters()
	if expected := int64(5); parameters.HTTPTimeout != expected {
		t.Fatalf("HTTPTimeout Param was not as expected. Actual: %d, Expected: %d", parameters.HTTPTimeout, expected)
	}

	resp, err := http.DefaultClient.Get(server.URL) //nolint:noctx // exercising the configured client as-is
	if err != nil {
		t.Fatalf("Unexpected error for a fast-responding server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code. Got: %d, expected: %d", resp.StatusCode, http.StatusOK)
	}
}
