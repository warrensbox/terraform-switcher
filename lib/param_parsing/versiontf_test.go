package param_parsing

import (
	"github.com/hashicorp/go-version"
	"testing"
)

func TestGetVersionFromVersionsTF(t *testing.T) {
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/test_versiontf"
	params = GetVersionFromVersionsTF(params)
	v1, _ := version.NewVersion("1.0.0")
	v2, _ := version.NewVersion("2.0.0")
	actualVersion, _ := version.NewVersion(params.Version)
	if !actualVersion.GreaterThanOrEqual(v1) || !actualVersion.LessThan(v2) {
		t.Error("Determined version is not between 1.0.0 and 2.0.0")
	}
}
