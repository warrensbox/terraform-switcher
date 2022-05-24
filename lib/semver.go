package lib

import (
	"fmt"
	"sort"

	semver "github.com/hashicorp/go-version"
)

// GetSemver : returns version that will be installed based on server constaint provided
func GetSemver(tfconstraint *string, mirrorURL *string) (string, error) {

	//listAll := true
	//tflist, _ := GetTFList(*mirrorURL, listAll) //get list of versions
	tflist, _ := GetTFReleases(*mirrorURL, true)
	fmt.Printf("Reading required version from constraint: %s\n", *tfconstraint)
	tfversion, err := SemVerParser(tfconstraint, tflist)
	return tfversion, err
}

// ValidateSemVer : Goes through the list of terraform version, return a valid tf version for contraint provided
func SemVerParser(tfconstraint *string, tflist []Release) (string, error) {
	tfversion := ""
	constraints, err := semver.NewConstraint(*tfconstraint) //NewConstraint returns a Constraints instance that a Version instance can be checked against
	if err != nil {
		return "", fmt.Errorf("error parsing constraint: %s", err)
	}
	versions := make([]*semver.Version, len(tflist))
	//put tfversion into semver object
	for i, tfvals := range tflist {
		version, err := semver.NewVersion(tfvals.Version) //NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
		if err != nil {
			return "", fmt.Errorf("error parsing constraint: %s", err)
		}
		versions[i] = version
	}

	sort.Sort(sort.Reverse(semver.Collection(versions)))

	for _, element := range versions {
		if constraints.Check(element) { // Validate a version against a constraint
			tfversion = element.String()
			fmt.Printf("Matched version: %s\n", tfversion)
			if ValidVersionFormat(tfversion) { //check if version format is correct
				return tfversion, nil
			}
		}
	}

	PrintInvalidTFVersion()
	return "", fmt.Errorf("error parsing constraint: %s", *tfconstraint)
}

// Print invalid TF version
func PrintInvalidTFVersion() {
	fmt.Println("Version does not exist or invalid terraform version format.\n Format should be #.#.# or #.#.#-@# where # are numbers and @ are word characters.\n For example, 0.11.7 and 0.11.9-beta1 are valid versions")
}

// Print invalid TF version
func PrintInvalidMinorTFVersion() {
	fmt.Println("Invalid minor terraform version format. Format should be #.# where # are numbers. For example, 0.11 is valid version")
}
