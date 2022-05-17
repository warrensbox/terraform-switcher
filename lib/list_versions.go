package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	releases, err := GetTFReleases(mirrorURL, preRelease)
	if err != nil {
		return nil, err
	}
	if preRelease {
		semver := fmt.Sprintf(`%s{1}\.\d+\-[a-zA-z]+\d*`, version)
		r, err := regexp.Compile(semver)
		if err != nil {
			return nil, err
		}
		for _, release := range releases {
			if r.MatchString(release.Version) {
				fmt.Printf("Matched version: %s\n", release.Version)
				return release, nil
			}
		}
		return nil, fmt.Errorf("Error: no match for requested version: %s", version)
	} else {
		version = fmt.Sprintf("~> %v", version)
		semv, err := SemVerParser(&version, releases)
		return semv, err
	}
}

// httpGet : generic http get client for the given url and query parameters.
func httpGet(url *url.URL, values url.Values) (*http.Response, error) {
	url.RawQuery = values.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("[Error] : Retrieving contents from url %s\n: %q", url, err)

	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Error] : non-200 response code during request: %d: Http status: %s", res.StatusCode, res.Status)
	}
	return res, nil
}

// getReleases : subfunc for GetTFReleases, used in a loop to get all terraform releases given the hashicorp url
func getReleases(url *url.URL, values url.Values) ([]*Release, error) {
	var releases []*Release
	resp, errURL := httpGet(url, values)
	if errURL != nil {
		return nil, fmt.Errorf("[Error] : Getting url: %q", errURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
	}

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body.Bytes(), &releases); err != nil {
		return nil, fmt.Errorf("%q: %s", err, body.String())
	}
	return releases, nil
}

//GetTFReleases :  Get all terraform releases given the hashicorp url
func GetTFReleases(mirrorURL string, preRelease bool) ([]*Release, error) {
	limit := 20
	u, err := url.Parse(mirrorURL)
	if err != nil {
		return nil, fmt.Errorf("[Error] : parsing url: %q", err)
	}
	values := u.Query()
	values.Set("limit", strconv.Itoa(limit))
	releaseSet, err := getReleases(u, values)
	if err != nil {
		return nil, err
	}
	var releases []*Release
	releases = append(releases, releaseSet...)
	for len(releaseSet) == limit {
		values.Set("after", releaseSet[len(releaseSet)-1].TimestampCreated.String())
		releaseSet, err = getReleases(u, values)
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
	url, err := url.Parse(mirrorURL + "/" + requestedVersion)
	if err != nil {
		return nil, fmt.Errorf("[Error] : parsing URL: %q", err)
	}
	resp, errURL := httpGet(url, nil)
	if errURL != nil {
		return nil, fmt.Errorf("[Error] : Getting url: %q", errURL)
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		return nil, fmt.Errorf("[Error]: parsing http response body: %q", err)
	}
	var release *Release
	if err := json.Unmarshal(body.Bytes(), &release); err != nil {
		return nil, fmt.Errorf("%q: %s", err, body)
	}
	return release, nil

}

//removePreReleases : Removes any prerelease versions from a given slice of Release.
func removePreReleases(releases []*Release) []*Release {
	realReleases := []*Release{}
	for i, r := range releases {
		if !r.IsPrerelease {
			realReleases = append(realReleases, releases[i])
		}
	}
	return realReleases
}

//VersionExist : check if requested version exist
func VersionExist(rel *Release, releases []*Release) bool {
	for _, r := range releases {
		if rel == r {
			return true
		}
	}
	return false
}

//RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []*Release) []*Release {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []*Release{}

	for _, val := range elements {
		versionOnly := strings.TrimSuffix(val.Version, " *recent")
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
