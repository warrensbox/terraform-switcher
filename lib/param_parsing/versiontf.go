package param_parsing

import (
	semver "github.com/hashicorp/go-version"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
)

func GetVersionFromVersionsTF(params Params) (Params, error) {
	var tfConstraints []string

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
		// Check if the version contraint is valid
		constraint, constraintErr := semver.NewConstraint(requiredVersions[key])
		if constraintErr != nil {
			logger.Errorf("Invalid version constraint found: %q", requiredVersions[key])
			return params, constraintErr
		}
		// It's valid. Add to list
		tfConstraints = append(tfConstraints, constraint.String())
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
