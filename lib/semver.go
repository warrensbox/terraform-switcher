//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"fmt"
	"sort"

	semver "github.com/hashicorp/go-version"
)

// GetSemver : returns version that will be installed based on server constraint provided
func GetSemver(tfconstraint string, mirrorURL string) (string, error) {
	listAll := true
	tflist, errTFList := getTFList(mirrorURL, listAll) // get list of versions
	if errTFList != nil {
		return "", fmt.Errorf("Error getting list of versions from %q: %v", mirrorURL, errTFList)
	}
	logger.Infof("Reading required version from constraint: %q", tfconstraint)
	tfversion, err := SemVerParser(&tfconstraint, tflist)
	return tfversion, err
}

// SemVerParser  : Goes through the list of versions, returns a valid version for contraint provided
func SemVerParser(tfconstraint *string, tflist []string) (string, error) {
	tfversion := ""
	constraints, err := semver.NewConstraint(*tfconstraint) // NewConstraint returns a Constraints instance that a Version instance can be checked against
	if err != nil {
		return "", fmt.Errorf("Error parsing constraint: %s", err)
	}
	versions := make([]*semver.Version, len(tflist))
	// put tfversion into semver object
	for i, tfvals := range tflist {
		version, err := semver.NewVersion(tfvals) // NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
		if err != nil {
			return "", fmt.Errorf("Error parsing constraint: %s", err)
		}
		versions[i] = version
	}

	sort.Sort(sort.Reverse(semver.Collection(versions)))

	for _, element := range versions {
		if constraints.Check(element) { // Validate a version against a constraint
			tfversion = element.String()
			if validVersionFormat(tfversion) { // check if version format is correct
				logger.Infof("Matched version: %q", tfversion)
				return tfversion, nil
			}
			PrintInvalidTFVersion()
		}
	}

	return "", fmt.Errorf("Did not find version matching constraint: %s", *tfconstraint)
}

// PrintInvalidTFVersion Print invalid TF version
func PrintInvalidTFVersion() {
	logger.Error("Version does not exist or invalid terraform version format.\n\tFormat should be #.#.# or #.#.#-@# where # are numbers and @ are word characters.\n\tFor example, 1.11.7 and 0.11.9-beta1 are valid versions")
}

// PrintInvalidMinorTFVersion Print invalid minor TF version
func PrintInvalidMinorTFVersion() {
	logger.Error("Invalid minor terraform version format.\n\tFormat should be #.# where # are numbers. For example, 1.11 is valid version")
}
