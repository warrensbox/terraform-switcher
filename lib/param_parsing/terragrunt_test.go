//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/warrensbox/terraform-switcher/lib"
)

func TestGetVersionFromTerragrunt(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_terragrunt_hcl"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	v1, v1Err := version.NewVersion("0.13")
	if v1Err != nil {
		t.Errorf("Error parsing v1 version: %v", v1Err)
	}
	v2, v2Err := version.NewVersion("0.14")
	if v2Err != nil {
		t.Errorf("Error parsing v2 version: %v", v2Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Errorf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Error("Determined version is not between 0.13 and 0.14")
	}
}

func TestGetVersionTerragrunt_with_no_terragrunt_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_no_file"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if params.Version != "" {
		t.Error("Version should be empty")
	}
}

func TestGetVersionTerragrunt_with_no_version(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_terragrunt_no_version"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if params.Version != "" {
		t.Error("Version should be empty")
	}
}

func TestGetVersionFromTerragrunt_erroneous_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_terragrunt_error_hcl"
	params, err := GetVersionFromTerragrunt(params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := ""
	if params.Version != expected {
		t.Errorf("Expected version %q, got %q", expected, params.Version)
	}
}
