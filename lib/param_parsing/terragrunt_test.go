package param_parsing

import (
	"github.com/hashicorp/go-version"
	"testing"
)

func TestGetVersionFromTerragrunt(t *testing.T) {
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/test_terragrunt_hcl"
	params = GetVersionFromTerragrunt(params)
	v1, _ := version.NewVersion("0.13")
	v2, _ := version.NewVersion("0.14")
	actualVersion, _ := version.NewVersion(params.Version)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Error("Determined version is not between 0.13 and 0.14")
	}
}
