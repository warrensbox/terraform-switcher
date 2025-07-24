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

// Semantic version regexes without `^` and `$` anchors
// Follows https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
var regexSemVer = struct {
	Full             *regexp.Regexp
	Minor            *regexp.Regexp
	Patch            *regexp.Regexp
	PreReleaseSuffix *regexp.Regexp
}{
	Full:             regexp.MustCompile(`(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`),
	Minor:            regexp.MustCompile(`(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)`),
	Patch:            regexp.MustCompile(`(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)`),
	PreReleaseSuffix: regexp.MustCompile(`\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`),
}

func getVersionsFromBody(body string, preRelease bool, tfVersionList *tfVersionList) {
	var semver string
	// Without the ending '"' pre-release folders would be tried and break.
	if preRelease {
		semver = `\/?` + regexSemVer.Full.String() + `/?"`
	} else if !preRelease {
		semver = `\/?` + regexSemVer.Patch.String() + `\/?"`
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

// getTFList : Get the list of available versions given the mirror URL
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

// getTFLatest : Get the latest version given the mirror URL
func getTFLatest(mirrorURL string) (string, error) {
	result, err := getTFURLBody(mirrorURL)
	if err != nil {
		return "", err
	}
	// Getting versions from body; should return match /X.X.X/ where X is a number
	semver := `\/?` + regexSemVer.Patch.String() + `\/?"`
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

// getTFLatestImplicit : Get the latest implicit version given the mirror URL
func getTFLatestImplicit(mirrorURL string, preRelease bool, version string) (string, error) {
	if preRelease {
		// TODO: use getTFList() instead of getTFURLBody
		body, err := getTFURLBody(mirrorURL)
		if err != nil {
			return "", err
		}
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver := `\/?` + version + regexSemVer.PreReleaseSuffix.String() + `\/?"`
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

// getTFURLBody : Get list of versions from the mirror URL
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

// versionExist : Check if requested version exists
func versionExist(val, array any) (exists bool) {
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
		logger.Fatalf("Internal error: expected \"slice\", got %q", reflect.TypeOf(array).Kind())
	}
	return exists
}

// removeDuplicateVersions : Remove duplicate versions from a slice of strings
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

// validVersionFormat : Return true if valid semantic version provided based on the type of version requested
// Caveat: Passing no validation argument validates the full semantic version format to provide backward compatibility
func validVersionFormat(version string, validation ...*regexp.Regexp) bool {
	var semverRegex *regexp.Regexp

	switch len(validation) {
	case 0:
		semverRegex = regexSemVer.Full
	case 1:
		semverRegex = validation[0]
		// regexSemVer.PreReleaseSuffix is a special use case, hence do not accept it as valid validation argument
		if semverRegex == regexSemVer.PreReleaseSuffix {
			logger.Fatalf("Internal error: invalid \"validation\" argument value")
		}
	default:
		logger.Fatalf("Internal error: invalid number of arguments (must be 1 or 2, got %d)", 1+len(validation))
	}

	semverRegex = regexp.MustCompile(`^` + semverRegex.String() + `$`)

	return semverRegex.MatchString(version)
}

// ShowLatestVersion : Show latest stable version given the mirror URL
func ShowLatestVersion(mirrorURL string) {
	tfversion, err := getTFLatest(mirrorURL)
	if err != nil {
		logger.Fatalf("Error getting latest version from %q: %v", mirrorURL, err)
	}

	fmt.Printf("%s\n", tfversion)
}

// ShowLatestImplicitVersion : show latest implicit version given the mirror URL
func ShowLatestImplicitVersion(requestedVersion, mirrorURL string, preRelease bool) {
	if validVersionFormat(requestedVersion, regexSemVer.Minor) || (validVersionFormat(requestedVersion, regexSemVer.Patch) && !preRelease) {
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
		if preRelease {
			PrintInvalidMinorTFVersion()
		} else {
			printInvalidVersionFormat()
		}
	}
}
