package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
)

type Release struct {
	Builds []struct {
		Arch string `json:"arch"`
		OS   string `json:"os"`
		URL  string `json:"url"`
	} `json:"builds"`
	IsPrerelease     bool            `json:"is_prerelease"`
	TimestampCreated strfmt.DateTime `json:"timestamp_created"`
	Version          string          `json:"version"`
}

//GetTFLatest :  Get the latest terraform version given the hashicorp url
func GetTFLatest(mirrorURL string, preRelease bool) (*Release, error) {
	releases, error := GetTFReleases(mirrorURL, preRelease)
	if error != nil {
		return nil, error
	}
	for i := range releases {
		if !releases[i].IsPrerelease {
			return releases[i], nil
		}
	}
	return nil, nil
}

//GetTFLatestImplicit :  Get the latest implicit terraform version given the hashicorp url
func GetTFLatestImplicit(mirrorURL string, preRelease bool, version string) (*Release, error) {
	if preRelease {
		version = fmt.Sprintf(`%s{1}\.\d+\-[a-zA-z]+\d*`, version)
	} else if !preRelease {
		version = fmt.Sprintf("~> %v", version)
	}
	releases, err := GetTFReleases(mirrorURL, preRelease)
	if err != nil {
		return nil, err
	}
	semv, err := SemVerParser(&version, releases)
	if err != nil {
		return nil, err
	}
	return semv, nil
}

// httpGet : generic http get client for the given url and query parameters.
func httpGet(url string, queryParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("issue during request (%d: %q)", res.StatusCode, res.Status)
	}
	return res, nil
}

// getReleases : subfunc for GetTFReleases, used in a loop to get all terraform releases given the hashicorp url
func getReleases(mirrorURL string, queryParams map[string]string) ([]*Release, error) {
	var releases []*Release
	resp, errURL := httpGet(mirrorURL, queryParams)
	if errURL != nil {
		return nil, fmt.Errorf("[Error] : Getting url: %v", errURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Error] : Retrieving contents from url: %s", mirrorURL)
	}

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body.Bytes(), &releases); err != nil {
		return nil, fmt.Errorf("%s: %s", err, body.String())
	}
	return releases, nil
}

//GetTFReleases :  Get all terraform releases given the hashicorp url
func GetTFReleases(mirrorURL string, preRelease bool) ([]*Release, error) {
	limit := 20
	queryParams := map[string]string{"limit": strconv.Itoa(limit)}
	releaseSet, err := getReleases(mirrorURL, queryParams)
	if err != nil {
		return nil, err
	}
	var releases []*Release
	releases = append(releases, releaseSet...)
	for len(releaseSet) == limit {
		queryParams["after"] = releaseSet[len(releaseSet)-1].TimestampCreated.String()
		releaseSet, err = getReleases(mirrorURL, queryParams)
		if err != nil {
			return nil, err
		}
		releases = append(releases, releaseSet...)
	}

	if !preRelease {
		releases = removePreReleases(releases)
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Version > releases[j].Version
	})
	return releases, nil
}

//GetTFRelease :  Get the requested terraform release given the hashicorp url
func GetTFRelease(mirrorURL, requestedVersion string) (*Release, error) {
	resp, errURL := httpGet(mirrorURL+"/"+requestedVersion, nil)
	if errURL != nil {
		return nil, fmt.Errorf("[Error] : Getting url: %v", errURL)
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		return nil, err
	}
	var release *Release
	if err := json.Unmarshal(body.Bytes(), &release); err != nil {
		return nil, fmt.Errorf("%s: %s", err, body.String())
	}
	return release, nil

}

//removePreReleases : Removes any prerelease versions from a given slice of Release.
func removePreReleases(releases []*Release) []*Release {
	for i, r := range releases {
		if r.IsPrerelease {
			releases = append(releases[:i], releases[i+1:]...)
		}
	}
	return releases
}

//VersionExist : check if requested version exist
func VersionExist(val interface{}, array interface{}) (exists bool) {

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
	}

	return exists
}

//RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, val := range elements {
		versionOnly := strings.Trim(val, " *recent")
		if encountered[versionOnly] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[versionOnly] = true
			// Append to result slice.
			result = append(result, val)
		}
	}
	// Return the new slice.
	return result
}

// ValidVersionFormat : returns valid version format
/* For example: 0.1.2 = valid
// For example: 0.1.2-beta1 = valid
// For example: 0.1.2-alpha = valid
// For example: a.1.2 = invalid
// For example: 0.1. 2 = invalid
*/
func ValidVersionFormat(version string) bool {

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
func ValidMinorVersionFormat(version string) bool {

	// Getting versions from body; should return match /X.X./ where X is a number
	semverRegex := regexp.MustCompile(`^(\d+\.\d+)$`)

	return semverRegex.MatchString(version)
}
