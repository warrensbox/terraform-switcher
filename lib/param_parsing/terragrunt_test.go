package param_parsing

import (
	"github.com/hashicorp/go-version"
	"github.com/warrensbox/terraform-switcher/lib"
	"testing"
)

func TestGetVersionFromTerragrunt(t *testing.T) {
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/test_terragrunt_hcl"
	params, _ = GetVersionFromTerragrunt(params)
	v1, _ := version.NewVersion("0.13")
	v2, _ := version.NewVersion("0.14")
	actualVersion, _ := version.NewVersion(params.Version)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Error("Determined version is not between 0.13 and 0.14")
	}
}

func TestGetVersionTerragrunt_with_no_terragrunt_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/test_no_file"
	params, _ = GetVersionFromTerragrunt(params)
	if params.Version != "" {
		t.Error("Version should be empty")
	}
}

func TestGetVersionFromTerragrunt_erroneous_file(t *testing.T) {
	var params Params
	logger = lib.InitLogger("DEBUG")
	params = initParams(params)
	params.ChDirPath = "../../test-data/test_terragrunt_error_hcl"
	params, err := GetVersionFromTerragrunt(params)
	if err == nil {
		t.Error("Expected error but got none.")
	} else {
		expectedError := "could not decode body of HCL file \"../../test-data/test_terragrunt_error_hcl/terragrunt.hcl\""
		if err.Error() != expectedError {
			t.Errorf("Expected error to be '%q', got '%q'", expectedError, err)
		}
	}
}
