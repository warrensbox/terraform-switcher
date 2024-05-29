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
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_versiontf"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, _ = GetVersionFromVersionsTF(params)
	v1, _ := version.NewVersion("1.0.5")
	actualVersion, _ := version.NewVersion(params.Version)
	if !actualVersion.Equal(v1) {
		t.Errorf("Determined version is not 1.0.5, but %s", params.Version)
	}
}

func TestGetVersionFromVersionsTF_impossible_constraints(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_matching_constraints"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
	expectedError := "did not find version matching constraint: ~> 1.0.0, =1.0.5, <= 1.0.4"
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
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_existent"
	params.MirrorURL = lib.GetProductById("terraform").GetDefaultMirrorUrl()
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error got nil")
	} else {
		expected := "did not find version matching constraint: > 99999.0.0"
		if fmt.Sprint(err) != expected {
			t.Errorf("Expected error %q, got %q", expected, err)
		}
	}
}
