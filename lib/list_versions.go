//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type tfVersionList struct {
	tflist []string
}

func getVersionsFromBody(body string, preRelease bool, tfVersionList *tfVersionList) {
	var semver string
	if preRelease {
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver = `\/?(\d+\.\d+\.\d+)(-[a-zA-z]+\d*)?/?"`
	} else if !preRelease {
		// Getting versions from body; should return match /X.X.X/ where X is a number
		// without the ending '"' pre-release folders would be tried and break.
		semver = `\/?(\d+\.\d+\.\d+)\/?"`
	}
	r, err := regexp.Compile(semver)
	if err != nil {
		logger.Fatalf("Error compiling %q regex: %v", semver, err)
	}

	matches := r.FindAllString(body, -1)
	if matches == nil {
		return
	}
	for _, match := range matches {
		trimstr := strings.Trim(match, "/\"") // remove '/' or '"' from /X.X.X/" or /X.X.X"
		tfVersionList.tflist = append(tfVersionList.tflist, trimstr)
	}
}

// getTFList :  Get the list of available versions given the mirror URL
func getTFList(mirrorURL string, preRelease bool) ([]string, error) {
	logger.Debug("Getting list of versions")
	result, err := getTFURLBody(mirrorURL)
	if err != nil {
		return nil, err
	}

	var tfVerList tfVersionList
	getVersionsFromBody(result, preRelease, &tfVerList)

	if len(tfVerList.tflist) == 0 {
		logger.Errorf("Cannot get version list from mirror: %s", mirrorURL)
	}
	return tfVerList.tflist, nil
}

// getTFLatest :  Get the latest terraform version given the hashicorp url
func getTFLatest(mirrorURL string) (string, error) {
	result, err := getTFURLBody(mirrorURL)
	if err != nil {
		return "", err
	}
	// Getting versions from body; should return match /X.X.X/ where X is a number
	semver := `\/?(\d+\.\d+\.\d+)\/?"`
	r, errSemVer := regexp.Compile(semver)
	if errSemVer != nil {
		return "", fmt.Errorf("Error compiling %q regex: %v", semver, errSemVer)
	}
	bodyLines := strings.Split(result, "\n")
	for i := range result {
		if r.MatchString(bodyLines[i]) {
			str := r.FindString(bodyLines[i])
			trimstr := strings.Trim(str, "/\"") // remove '/' or '"' from /X.X.X/" or /X.X.X"
			return trimstr, nil
		}
	}
	return "", nil
}

// getTFLatestImplicit :  Get the latest implicit terraform version given the hashicorp url
func getTFLatestImplicit(mirrorURL string, preRelease bool, version string) (string, error) {
	if preRelease {
		// TODO: use getTFList() instead of getTFURLBody
		body, err := getTFURLBody(mirrorURL)
		if err != nil {
			return "", err
		}
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver := fmt.Sprintf(`\/?(%s{1}\.\d+\-[a-zA-z]+\d*)\/?"`, version)
		r, errReSemVer := regexp.Compile(semver)
		if errReSemVer != nil {
			return "", errReSemVer
		}
		versions := strings.Split(body, "\n")
		for i := range versions {
			if r.MatchString(versions[i]) {
				str := r.FindString(versions[i])
				trimstr := strings.Trim(str, "/\"") // remove '/' or '"' from /X.X.X/" or /X.X.X"
				return trimstr, nil
			}
		}
	} else if !preRelease {
		listAll := false
		tflist, errTFList := getTFList(mirrorURL, listAll) // get list of versions
		if errTFList != nil {
			return "", fmt.Errorf("Error getting list of versions from %q: %v", mirrorURL, errTFList)
		}

		version = fmt.Sprintf("~> %v", version)
		semv, err := SemVerParser(&version, tflist)
		if err != nil {
			return "", err
		}
		return semv, nil
	}
	return "", nil
}

// getTFURLBody : Get list of terraform versions from hashicorp releases
func getTFURLBody(mirrorURL string) (string, error) {
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	if !hasSlash {
		// if it does not have slash - append slash
		mirrorURL = fmt.Sprintf("%s/", mirrorURL)
	}
	resp, errURL := http.Get(mirrorURL) // nolint:gosec // `mirrorURL' is expected to be variable
	if errURL != nil {
		logger.Fatalf("Error getting url: %v", errURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Fatalf("Error retrieving contents from url: %s", mirrorURL)
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		logger.Fatalf("Error reading body: %v", errBody)
	}

	bodyString := string(body)

	return bodyString, nil
}

// versionExist : check if requested version exist
func versionExist(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
				return exists
			}
		}
	default:
		panic("unhandled default case")
	}
	return exists
}

// removeDuplicateVersions : remove duplicate version
func removeDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, val := range elements {
		versionOnly := strings.TrimSuffix(val, " *recent")
		if !encountered[versionOnly] {
			// Record this element as an encountered element.
			encountered[versionOnly] = true
			// Append to result slice.
			result = append(result, val)
		}
	}
	// Return the new slice.
	return result
}

// validVersionFormat : returns valid version format
/* For example: 0.1.2 = valid
// For example: 0.1.2-beta1 = valid
// For example: 0.1.2-alpha = valid
// For example: a.1.2 = invalid
// For example: 0.1. 2 = invalid
*/
func validVersionFormat(version string) bool {
	// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
	// Follow https://semver.org/spec/v1.0.0-beta.html
	// Check regular expression at https://rubular.com/r/ju3PxbaSBALpJB
	semverRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)(-[a-zA-z]+\d*)?$`)
	return semverRegex.MatchString(version)
}

// ValidMinorVersionFormat : returns valid MINOR version format
/* For example: 0.1 = valid
// For example: a.1.2 = invalid
// For example: 0.1.2 = invalid
*/
func validMinorVersionFormat(version string) bool {
	// Getting versions from body; should return match /X.X./ where X is a number
	semverRegex := regexp.MustCompile(`^(\d+\.\d+)$`)

	return semverRegex.MatchString(version)
}

// ShowLatestVersion show install latest stable tf version
func ShowLatestVersion(mirrorURL string) {
	tfversion, err := getTFLatest(mirrorURL)
	if err != nil {
		logger.Fatalf("Error getting latest version from %q: %v", mirrorURL, err)
	}

	fmt.Printf("%s\n", tfversion)
}

// ShowLatestImplicitVersion show latest - argument (version) must be provided
func ShowLatestImplicitVersion(requestedVersion, mirrorURL string, preRelease bool) {
	if validMinorVersionFormat(requestedVersion) {
		tfversion, err := getTFLatestImplicit(mirrorURL, preRelease, requestedVersion)
		if err != nil {
			logger.Fatalf("Error getting latest implicit version %q from %q: %v", requestedVersion, mirrorURL, err)
		}

		if len(tfversion) > 0 {
			fmt.Printf("%s\n", tfversion)
		} else {
			logger.Fatalf("Requested version does not exist: %q.\n\tTry `tfswitch -l` to see all available versions", requestedVersion)
		}
	} else {
		PrintInvalidMinorTFVersion()
	}
}
