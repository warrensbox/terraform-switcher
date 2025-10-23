//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
)

const paramTypeVersionTF = "Terraform module"

func getConstraintFromVersionsTF(params Params) (Params, error) {
	if !isTerraformModule(params) {
		return params, nil
	}

	var tfConstraints []string

	relPath, err := lib.GetRelativePath(params.ChDirPath)
	if err != nil {
		return params, err
	}

	logger.Infof("Reading version constraint from %s at %q", paramTypeVersionTF, relPath)
	module, _ := tfconfig.LoadModule(relPath) // nolint:errcheck // covered by conditional below
	if module.Diagnostics.HasErrors() {
		return params, fmt.Errorf("Could not load %s at %q", paramTypeVersionTF, relPath)
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

	params.VersionRequirement = strings.Join(tfConstraints, ", ")
	logger.Debugf("Using version constraint from %s at %q: %q", paramTypeVersionTF, relPath, params.VersionRequirement)
	return params, nil
}

func GetVersionFromVersionsTF(params Params) (Params, error) {
	params, err := getConstraintFromVersionsTF(params)
	if err != nil {
		return params, err
	}

	if params.MatchVersionRequirement == "" {
		version, err2 := lib.GetSemver(params.VersionRequirement, params.MirrorURL)
		if err2 != nil {
			logger.Errorf("No version found matching %q", params.VersionRequirement)
			return params, err2
		}
		params.Version = version
		logger.Debugf("Using version from %s: %q", paramTypeVersionTF, params.Version)
	}
	return params, nil
}

func isTerraformModule(params Params) bool {
	relPath, errRelPath := lib.GetRelativePath(params.ChDirPath)
	if errRelPath != nil {
		logger.Warn(errRelPath)
		return false
	}

	module, err := tfconfig.LoadModule(relPath)
	if err != nil {
		logger.Warnf("Error parsing %s: %v", paramTypeVersionTF, err)
		return false
	}
	if len(module.RequiredCore) == 0 {
		logger.Debugf("No required version constraints defined by %s at %q", paramTypeVersionTF, relPath)
	}
	return len(module.RequiredCore) > 0
}
