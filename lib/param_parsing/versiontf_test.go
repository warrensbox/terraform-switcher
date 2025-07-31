//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/warrensbox/terraform-switcher/lib"
)

func TestGetVersionFromVersionsTF_matches_version(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	var getVerErr error
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_versiontf"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, getVerErr = GetVersionFromVersionsTF(params)
	if getVerErr != nil {
		t.Errorf("Error getting version from Terraform module: %v", getVerErr)
	}
	v1, v1Err := version.NewVersion("1.0.5")
	if v1Err != nil {
		t.Errorf("Error parsing v1 version: %v", v1Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Errorf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	if !actualVersion.Equal(v1) {
		t.Errorf("Determined version is not 1.0.5, but %s", params.Version)
	}
}

// This test ensures we do fall back to FOSS Terraform when asked for versions
// prior to 1.6.0, and that the version is parsed correctly when doing so.
func TestGetVersionFromVersionsTF_fossFallback_matches_product(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	var getVerErr error
	params = initParams(params)
	params.Product = "opentofu"
	params.FossFallback = true
	params.ProductEntity = lib.GetProductById(params.Product)
	params.CustomBinaryPath = "/tmp/tofu"
	params.ChDirPath = "../../test-data/integration-tests/test_versiontf"
	params.MirrorURL = params.ProductEntity.GetDefaultMirrorUrl()
	params, getVerErr = GetVersionFromVersionsTF(params)
	if getVerErr != nil {
		t.Errorf("Error getting version from Terraform module: %v", getVerErr)
	}
	expected := "terraform"
	if params.Product != expected {
		t.Errorf("Error falling back to FOSS Terraform, product is not terraform, but %s", params.Product)
	}
	v1, v1Err := version.NewVersion("1.0.5")
	if v1Err != nil {
		t.Errorf("Error parsing v1 version: %v", v1Err)
	}
	actualVersion, actualVersionErr := version.NewVersion(params.Version)
	if actualVersionErr != nil {
		t.Errorf("Error parsing actualVersion version: %v", actualVersionErr)
	}
	if !actualVersion.Equal(v1) {
		t.Errorf("Determined version is not 1.0.5, but %s", params.Version)
	}
}

// This test ensures we don't fall back to Terraform for versions on or after 1.6.0
func TestGetVersionFromVersionsTF_fossFallback_only_for_foss_versions(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.Product = "opentofu"
	params.FossFallback = true
	params.ProductEntity = lib.GetProductById(params.Product)
	params.CustomBinaryPath = "/tmp/tofu"
	params.ChDirPath = "../../test-data/integration-tests/test_versiontf_foss"
	params.MirrorURL = params.ProductEntity.GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
  // expectedError := "Error getting version from OpenTofu module: Did not find version matching constraint: =1.6.6"
	expectedError := "Did not find version matching constraint: =1.6.6"
	if err == nil {
		t.Errorf("Expected error '%s', got nil", expectedError)
	} else {
		if err.Error() == expectedError {
			t.Logf("Got expected error '%s'", err)
		} else {
			t.Errorf("Got unexpected error '%s'", err)
		}
	}
}

func TestGetVersionFromVersionsTF_impossible_constraints(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ProductEntity = lib.GetProductById("terraform")
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_matching_constraints"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
	expectedError := "Did not find version matching constraint: ~> 1.0.0, =1.0.5, <= 1.0.4"
	if err == nil {
		t.Errorf("Expected error '%s', got nil", expectedError)
	} else {
		if err.Error() == expectedError {
			t.Logf("Got expected error '%s'", err)
		} else {
			t.Errorf("Got unexpected error '%s'", err)
		}
	}
}

func TestGetVersionFromVersionsTF_erroneous_file(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_error"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error got nil")
	} else {
		expected := "Malformed constraint: ~527> 1.0.0"
		if fmt.Sprint(err) != expected {
			t.Errorf("Expected error %q, got %q", expected, err)
		}
	}
}

func TestGetVersionFromVersionsTF_non_existent_constraint(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ProductEntity = lib.GetProductById("terraform")
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_existent"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error got nil")
	} else {
		expected := "Did not find version matching constraint: > 99999.0.0"
		if fmt.Sprint(err) != expected {
			t.Errorf("Expected error %q, got %q", expected, err)
		}
	}
}
