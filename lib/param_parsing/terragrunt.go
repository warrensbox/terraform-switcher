//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/warrensbox/terraform-switcher/lib"
)

const terraGruntFileName = "terragrunt.hcl"

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

func GetVersionFromTerragrunt(params Params) (Params, error) {
	filePath := filepath.Join(params.ChDirPath, terraGruntFileName)
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %q", filePath)
		parser := hclparse.NewParser()
		hclFile, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			return params, fmt.Errorf("unable to parse HCL file %q", filePath)
		}
		var versionFromTerragrunt terragruntVersionConstraints
		diagnostics = gohcl.DecodeBody(hclFile.Body, nil, &versionFromTerragrunt)
		// do not fail on failure to decode the body, as it may f.e. miss a required block,
		// though we don't want to fail execution because of that
		if diagnostics.HasErrors() {
			logger.Errorf(diagnostics.Error())
		}
		if versionFromTerragrunt.TerraformVersionConstraint == "" {
			logger.Infof("No terraform version constraint in %q", filePath)
			return params, nil
		}
		version, err := lib.GetSemver(versionFromTerragrunt.TerraformVersionConstraint, params.MirrorURL)
		if err != nil {
			return params, fmt.Errorf("no version found matching %q", versionFromTerragrunt.TerraformVersionConstraint)
		}
		params.Version = version
		logger.Debugf("Using version from %q: %q", filePath, params.Version)
	}
	return params, nil
}

func terraGruntFileExists(params Params) bool {
	filePath := filepath.Join(params.ChDirPath, terraGruntFileName)
	return lib.CheckFileExist(filePath)
}
