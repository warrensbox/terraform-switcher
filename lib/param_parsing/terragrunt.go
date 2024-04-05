package param_parsing

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
)

const terraGruntFileName = "terragrunt.hcl"

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

func GetVersionFromTerragrunt(params Params) Params {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	if lib.CheckFileExist(filePath) {
		logger.Infof("Reading configuration from %s", filePath)
		parser := hclparse.NewParser()
		hclFile, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			logger.Fatalf("Unable to parse HCL file %s", filePath)
			os.Exit(1)
		}
		var versionFromTerragrunt terragruntVersionConstraints
		diagnostics = gohcl.DecodeBody(hclFile.Body, nil, &versionFromTerragrunt)
		if diagnostics.HasErrors() {
			logger.Fatal("Could not decode body of HCL file.")
			os.Exit(1)
		}
		version, err := lib.GetSemver(versionFromTerragrunt.TerraformVersionConstraint, params.MirrorURL)
		if err != nil {
			logger.Fatal("Could not determine semantic version")
			os.Exit(1)
		}
		params.Version = version
	}
	return params
}

func terraGruntFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	return lib.CheckFileExist(filePath)
}
