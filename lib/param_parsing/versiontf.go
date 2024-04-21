package param_parsing

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
)

func GetVersionFromVersionsTF(params Params) (Params, error) {
	logger.Infof("Reading version from terraform module at %q", params.ChDirPath)
	module, err := tfconfig.LoadModule(params.ChDirPath)
	if err != nil {
		logger.Errorf("Could not load terraform module at %q", params.ChDirPath)
		return params, err.Err()
	}
	tfConstraint := module.RequiredCore[0]
	version, err2 := lib.GetSemver(tfConstraint, params.MirrorURL)
	if err2 != nil {
		logger.Errorf("No version found matching %q", tfConstraint)
		return params, err2
	}
	params.Version = version
	return params, nil
}

func isTerraformModule(params Params) bool {
	module, err := tfconfig.LoadModule(params.ChDirPath)
	return err == nil && len(module.RequiredCore) > 0
}
