package param_parsing

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/warrensbox/terraform-switcher/lib"
)

const terraGruntFileName = "terragrunt.hcl"

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

func GetVersionFromTerragrunt(params Params) (Params, error) {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %q", filePath)
		parser := hclparse.NewParser()
		hclFile, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			return params, fmt.Errorf("unable to parse HCL file %q", filePath)
			//logger.Fatalf("Unable to parse HCL file %q", filePath)
		}
		var versionFromTerragrunt terragruntVersionConstraints
		diagnostics = gohcl.DecodeBody(hclFile.Body, nil, &versionFromTerragrunt)
		if diagnostics.HasErrors() {
			return params, fmt.Errorf("could not decode body of HCL file %q", filePath)
		}
		version, err := lib.GetSemver(versionFromTerragrunt.TerraformVersionConstraint, params.MirrorURL)
		if err != nil {
			return params, fmt.Errorf("no version found matching %q", versionFromTerragrunt.TerraformVersionConstraint)
		}
		params.Version = version
	}
	return params, nil
}

func terraGruntFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	return lib.CheckFileExist(filePath)
}
