//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
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

func getVersionsFromJson(product Product, body string, preRelease bool, tfVersionList *tfVersionList) error {
	versionList, err := product.GetVersionsFromJson([]byte(body))
	if err != nil {
		return err
	}
	var semver string
	if preRelease {
		semver = regexSemVer.Full.String()
	} else {
		semver = regexSemVer.Patch.String()
	}
	semver = "^" + semver + "$"
	r, err := regexp.Compile(semver)
	if err != nil {
		logger.Fatalf("Error compiling %q regex: %v", semver, err)
		return err
	}
	for _, versionItx := range versionList {
		if r.MatchString(versionItx) {
			logger.Warnf("Adding version: %s", versionItx)
			tfVersionList.tflist = append(tfVersionList.tflist, versionItx)
		}
	}

	// Sort versions
	sort.Slice(tfVersionList.tflist, func(i, j int) bool {
		iVersion, err := version.NewSemver(tfVersionList.tflist[i])
		if err != nil {
			logger.Warn("Failed to parse version: %s", tfVersionList.tflist[i])
			return true
		}
		jVersion, err := version.NewSemver(tfVersionList.tflist[j])
		if err != nil {
			logger.Warn("Failed to parse version: %s", tfVersionList.tflist[j])
			return false
		}
		return iVersion.GreaterThan(jVersion)
	})

	return nil
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
func getTFList(product Product, mirrorURL string, preRelease bool) ([]string, error) {
	logger.Debug("Getting list of versions")
	body, err := getTFURLBody(mirrorURL)
	if err != nil {
		return nil, err
	}

	var tfVerList tfVersionList
	err = getVersionsFromJson(product, body, preRelease, &tfVerList)
	if err != nil {
		getVersionsFromBody(body, preRelease, &tfVerList)
	}

	if len(tfVerList.tflist) == 0 {
		logger.Errorf("Cannot get version list from mirror: %s", mirrorURL)
	}
	return tfVerList.tflist, nil
}

// getTFLatest : Get the latest version given the mirror URL
func getTFLatest(product Product, mirrorURL string) (string, error) {
	versions, err := getTFList(product, mirrorURL, false)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("No Versions available")
	}
	return versions[len(versions)-1], nil
}

// getTFLatestImplicit : Get the latest implicit version given the mirror URL
func getTFLatestImplicit(product Product, mirrorURL string, preRelease bool, version string) (string, error) {
	tflist, errTFList := getTFList(product, mirrorURL, preRelease) // get list of versions
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

// getTFURLBody : Get list of versions from the mirror URL
func getTFURLBody(mirrorURL string) (string, error) {
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	isJson := strings.HasSuffix(mirrorURL, ".json")
	if !hasSlash && !isJson {
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

// IsValidVersionFormat : Public wrapper for validVersionFormat
func IsValidVersionFormat(version string, validation ...*regexp.Regexp) bool {
	return validVersionFormat(version, validation...)
}

// ShowLatestVersion : Show latest stable version given the mirror URL
func ShowLatestVersion(product Product, mirrorURL string) {
	tfversion, err := getTFLatest(product, mirrorURL)
	if err != nil {
		logger.Fatalf("Error getting latest version from %q: %v", mirrorURL, err)
	}

	fmt.Printf("%s\n", tfversion)
}

// ShowLatestImplicitVersion : show latest implicit version given the mirror URL
func ShowLatestImplicitVersion(product Product, requestedVersion, mirrorURL string, preRelease bool) {
	if validVersionFormat(requestedVersion, regexSemVer.Minor) || (validVersionFormat(requestedVersion, regexSemVer.Patch) && !preRelease) {
		tfversion, err := getTFLatestImplicit(product, mirrorURL, preRelease, requestedVersion)
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
