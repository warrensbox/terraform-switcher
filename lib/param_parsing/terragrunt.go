package param_parsing

import (
	"fmt"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
)

const terraGruntFileName = "terragrunt.hcl"

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

func GetVersionFromTerragrunt(params Params) Params {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	if lib.CheckFileExist(filePath) {
		fmt.Printf("Reading configuration from %s\n", filePath)
		parser := hclparse.NewParser()
		hclFile, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			log.Fatal("Unable to parse HCL file", filePath)
		}
		var versionFromTerragrunt terragruntVersionConstraints
		diagnostics = gohcl.DecodeBody(hclFile.Body, nil, &versionFromTerragrunt)
		if diagnostics.HasErrors() {
			log.Fatal("Could not decode body of HCL file.")
		}
		version, err := lib.GetSemver(versionFromTerragrunt.TerraformVersionConstraint, params.MirrorURL)
		if err != nil {
			log.Fatal("Could not determine semantic version")
		}
		params.Version = version
	}
	return params
}

func terraGruntFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + terraGruntFileName
	return lib.CheckFileExist(filePath)
}
