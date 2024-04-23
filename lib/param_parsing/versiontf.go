package param_parsing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
)

func GetVersionFromVersionsTF(params Params) (Params, error) {
	var tfConstraints []string
	var exactConstraints []string

	curDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Could not get current working directory: %v", err)
	}

	absPath := params.ChDirPath
	if !filepath.IsAbs(params.ChDirPath) {
		absPath, err = filepath.Abs(params.ChDirPath)
		if err != nil {
			logger.Fatalf("Could not derive absolute path to %q: %v", params.ChDirPath, err)
		}
	}

	relPath, err := filepath.Rel(curDir, absPath)
	if err != nil {
		logger.Fatalf("Could not derive relative path to %q: %v", params.ChDirPath, err)
	}

	logger.Infof("Reading version from Terraform module at %q", relPath)
	module, _ := tfconfig.LoadModule(params.ChDirPath)
	if module.Diagnostics.HasErrors() {
		logger.Fatalf("Could not load Terraform module at %q", params.ChDirPath)
	}

	requiredVersions := module.RequiredCore

	for key := range requiredVersions {
		tfConstraint := requiredVersions[key]
		tfConstraintParts := strings.Fields(tfConstraint)

		if len(tfConstraintParts) > 2 {
			logger.Fatalf("Invalid version constraint found: %q", tfConstraint)
		} else if len(tfConstraintParts) == 1 {
			exactConstraints = append(exactConstraints, tfConstraint)
			tfConstraint = "= " + tfConstraintParts[0]
		}

		if tfConstraintParts[0] == "=" {
			exactConstraints = append(exactConstraints, tfConstraint)
		}

		tfConstraints = append(tfConstraints, tfConstraint)
	}

	if len(exactConstraints) > 0 && len(tfConstraints) > 1 {
		return params, fmt.Errorf("exact constraint (%q) cannot be combined with other conditions", strings.Join(exactConstraints, ", "))
	}

	tfConstraint := strings.Join(tfConstraints, ", ")

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
