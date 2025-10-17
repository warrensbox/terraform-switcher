//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import "github.com/warrensbox/terraform-switcher/lib"

// MatchVersionRequirement : Check if a given version meets a version requirement mandated by the configuration.
// - Return `0` if the version matches the requirement, otherwise return `2`.
// - Return `1` for any kind of general error (like version format parsing error).
func MatchVersionRequirement(parameters Params) int {
	// Sanity check MatchVersionRequirement parameter value and return `1` on error
	if !lib.IsValidVersionFormat(parameters.MatchVersionRequirement) {
		logger.Errorf("Failed to validate version format: %q", parameters.MatchVersionRequirement)
		lib.PrintInvalidTFVersion()
		return 1
	}

	// Fall back version requirement to a version from cmdline or to a default version arg (if either is provided)
	// Version from cmdline has precedence over default version arg
	if parameters.VersionRequirement == "" {
		if parameters.Version != "" {
			parameters.VersionRequirement = parameters.Version
		} else {
			parameters.VersionRequirement = parameters.DefaultVersion
		}
	}

	// If version requirement is still undefined, treat any version as acceptable
	if parameters.VersionRequirement == "" {
		logger.Warnf("No version requirement found to match against (version %q is acceptable)", parameters.MatchVersionRequirement)
		return 0
	}

	// Return `0` if the version to match meets the version requirement
	_, err := lib.SemVerParser(&parameters.VersionRequirement, []string{parameters.MatchVersionRequirement})
	if err == nil {
		logger.Infof("Version %q matches requirement %q", parameters.MatchVersionRequirement, parameters.VersionRequirement)
		return 0
	}
	// Otherwise return `2`
	logger.Errorf("Version %q mismatches requirement %q", parameters.MatchVersionRequirement, parameters.VersionRequirement)
	return 2
}
