//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/warrensbox/terraform-switcher/lib"
)

// Order of precedence for Terragrunt file names: first has highest precedence
var terragruntFileNames = []string{"terragrunt.hcl", "root.hcl"}

const (
	paramTypeTerragrunt        = "Terragrunt"
	terragruntConfigEnvVarName = "TF_TERRAGRUNT_CONFIG_FILE_NAME"
)

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

func terragruntFileNamesNew() []string {
	terragruntFileNamesNew := terragruntFileNames

	// Allow custom Terragrunt file name via env var
	if terragruntFileName := os.Getenv(terragruntConfigEnvVarName); terragruntFileName != "" {
		logger.Infof("Found %q env var: %q", terragruntConfigEnvVarName, terragruntFileName)

		// Take only the base name of the value to avoid path injection
		terragruntFileNameBase := filepath.Base(terragruntFileName)
		if terragruntFileNameBase != terragruntFileName {
			logger.Warnf("Stripping path from %q -> %q", terragruntFileName, terragruntFileNameBase)
			terragruntFileName = terragruntFileNameBase
		}

		// Prepend to the list to make custom file have highest precedence if it exists
		logger.Debugf("Prepending %q to the list of legit %s configuration files", terragruntFileName, paramTypeTerragrunt)
		terragruntFileNamesNew = append([]string{terragruntFileName}, terragruntFileNamesNew...)
	}

	// Deduplicate while preserving order, `lib.RemoveDuplicateStrings` keeps the first occurrence and drops later ones
	terragruntFileNamesNew = lib.RemoveDuplicateStrings(terragruntFileNamesNew)

	return terragruntFileNamesNew
}

func GetVersionFromTerragrunt(params Params) (Params, error) {
	relPath, errRelPath := lib.GetRelativePath(params.ChDirPath)
	if errRelPath != nil {
		return params, errRelPath
	}

	var versionFromTerragrunt terragruntVersionConstraints

	// Iterate over possible Terragrunt files and break on first found version constraint
	for _, terragruntFileName := range terragruntFileNamesNew() {
		filePath := filepath.Join(relPath, terragruntFileName)
		if !lib.IsRegularFile(filePath) {
			if lib.CheckFileExist(filePath) {
				logger.Warnf("Skipping non-regular %s configuration file %q", paramTypeTerragrunt, filePath)
			} else {
				logger.Tracef("Skipping non-existing %s configuration file %q", paramTypeTerragrunt, filePath)
			}
			continue
		}

		logger.Infof("Reading %s configuration from %q", paramTypeTerragrunt, filePath)
		parser := hclparse.NewParser()
		hclFile, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			logger.Errorf("Unable to parse %s HCL file %q", paramTypeTerragrunt, filePath)
			continue
		}
		diagnostics = gohcl.DecodeBody(hclFile.Body, nil, &versionFromTerragrunt)
		if diagnostics.HasErrors() {
			logger.Errorf(diagnostics.Error())
			continue
		}

		if versionFromTerragrunt.TerraformVersionConstraint != "" {
			params.VersionRequirement = versionFromTerragrunt.TerraformVersionConstraint
			logger.Debugf("Version requirement from %s configuration at %q: %q", paramTypeTerragrunt, filePath, params.VersionRequirement)
			break
		}

		logger.Debugf("No terraform version constraint found in %s configuration at %q", paramTypeTerragrunt, filePath)
	}

	if versionFromTerragrunt.TerraformVersionConstraint == "" {
		return params, nil
	}

	if params.MatchVersionRequirement == "" {
		version, err := lib.GetSemver(params.VersionRequirement, params.MirrorURL)
		if err != nil {
			return params, fmt.Errorf("no version found matching %q", params.VersionRequirement)
		}
		params.Version = version
		logger.Debugf("Using version from %s configuration: %q", paramTypeTerragrunt, params.Version)
	}

	return params, nil
}
