//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/warrensbox/terraform-switcher/lib"
)

// Test is expected to pick up `terragrunt.hcl` file
func TestGetVersionFromTerragrunt(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_terragrunt_hcl"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	v1, v1Err := version.NewVersion("0.13")
	if v1Err != nil {
		t.Fatalf("Error parsing v1 version: %v", v1Err)
	}
	v2, v2Err := version.NewVersion("0.14")
	if v2Err != nil {
		t.Fatalf("Error parsing v2 version: %v", v2Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Fatalf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	t.Logf("Testing whether %q is >= %q and < %q", actualVersion, v1, v2)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Fatalf("Determined version is not between %q and %q", v1, v2)
	}
	t.Logf("The %q is >= %q and < %q (expected)", actualVersion, v1, v2)
}

// Test is expected to pick up `root.hcl` file
func TestGetVersionFromTerragrunt_root_hcl(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_terragrunt_root_hcl"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	v1, v1Err := version.NewVersion("0.13")
	if v1Err != nil {
		t.Fatalf("Error parsing v1 version: %v", v1Err)
	}
	v2, v2Err := version.NewVersion("0.14")
	if v2Err != nil {
		t.Fatalf("Error parsing v2 version: %v", v2Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Fatalf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	t.Logf("Testing whether %q is >= %q and < %q", actualVersion, v1, v2)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Fatalf("Determined version is not between %q and %q", v1, v2)
	}
	t.Logf("The %q is >= %q and < %q (expected)", actualVersion, v1, v2)
}

// Test is expected to pick up custom file, specified via env var
func TestGetVersionFromTerragrunt_with_env_var(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	os.Setenv("TF_TERRAGRUNT_CONFIG_FILE_NAME", "custom.hcl")
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_terragrunt_hcl"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	v1, v1Err := version.NewVersion("1.12")
	if v1Err != nil {
		t.Fatalf("Error parsing v1 version: %v", v1Err)
	}
	v2, v2Err := version.NewVersion("1.13")
	if v2Err != nil {
		t.Fatalf("Error parsing v2 version: %v", v2Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Fatalf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	t.Logf("Testing whether %q is >= %q and < %q", actualVersion, v1, v2)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Fatalf("Determined version is not between %q and %q", v1, v2)
	}
	t.Logf("The %q is >= %q and < %q (expected)", actualVersion, v1, v2)
	os.Unsetenv("TF_TERRAGRUNT_CONFIG_FILE_NAME")
}

// Test is expected to pick up `terragrunt.hcl` file because
// the one specified via env var does not exist
func TestGetVersionFromTerragrunt_with_env_var_non_existing_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("TRACE")
	os.Setenv("TF_TERRAGRUNT_CONFIG_FILE_NAME", "nonexisting_file.hcl")
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_terragrunt_hcl"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	v1, v1Err := version.NewVersion("0.13")
	if v1Err != nil {
		t.Fatalf("Error parsing v1 version: %v", v1Err)
	}
	v2, v2Err := version.NewVersion("0.14")
	if v2Err != nil {
		t.Fatalf("Error parsing v2 version: %v", v2Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Fatalf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	t.Logf("Testing whether %q is >= %q and < %q", actualVersion, v1, v2)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Fatalf("Determined version is not between %q and %q", v1, v2)
	}
	t.Logf("The %q is >= %q and < %q (expected)", actualVersion, v1, v2)
	os.Unsetenv("TF_TERRAGRUNT_CONFIG_FILE_NAME")
}

func TestGetVersionTerragrunt_with_no_terragrunt_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if params.Version != "" {
		t.Fatalf("Version should be empty")
	}
}

func TestGetVersionTerragrunt_with_no_version(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_terragrunt_no_version"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if params.Version != "" {
		t.Fatalf("Version should be empty")
	}
}

func TestGetVersionFromTerragrunt_erroneous_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_terragrunt_error_hcl"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := ""
	if params.Version != expected {
		t.Fatalf("Expected version %q, got %q", expected, params.Version)
	}
}
