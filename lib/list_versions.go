package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		Os   string `json:"os"`
		Url  string `json:"url"`
	} `json:"builds"`
	IsPrerelease     bool            `json:"is_prerelease"`
	TimestampCreated strfmt.DateTime `json:"timestamp_created"`
	Version          string          `json:"version"`
}

//GetTFLatest :  Get the latest terraform version given the hashicorp url
func GetTFLatest(mirrorURL string, preRelease bool) (Release, error) {
	releases, error := GetTFReleases(mirrorURL, preRelease)
	if error != nil {
		return Release{}, error
	}
	for i := range releases {
		if !releases[i].IsPrerelease {
			return releases[i], nil
		}
	}
	return Release{}, nil
}

//GetTFLatestImplicit :  Get the latest implicit terraform version given the hashicorp url
func GetTFLatestImplicit(mirrorURL string, preRelease bool, version string) (Release, error) {
	if preRelease {
		version = fmt.Sprintf(`%s{1}\.\d+\-[a-zA-z]+\d*`, version)
	} else if !preRelease {
		version = fmt.Sprintf("~> %v", version)
	}
	releases, err := GetTFReleases(mirrorURL, preRelease)
	//version := fmt.Sprintf(`%s{1}\.\d+\-[a-zA-z]+\d*`, version)
	semv, err := SemVerParser(&version, releases)
	if err != nil {
		return Release{}, err
	}
	return semv, nil
}

func httpGet(url string, queryParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("issue during request (%d: %q)", res.StatusCode, res.Status)
		return nil, errors.New(err)
	}
	return res, nil
}

func getReleases(mirrorURL string, queryParams map[string]string) ([]Release, error) {
	var releases []Release
	resp, errURL := httpGet(mirrorURL, queryParams)
	if errURL != nil {
		log.Printf("[Error] : Getting url: %v", errURL)
		os.Exit(1)
		return nil, errURL
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Error] : Retrieving contents from url: %s", mirrorURL)
		os.Exit(1)
	}

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(body.Bytes(), &releases); err != nil {
		log.Fatalf("%s: %s", err, body.String())
	}
	return releases, nil
}

func GetTFReleases(mirrorURL string, preRelease bool) ([]Release, error) {
	limit := 20
	queryParams := map[string]string{"limit": strconv.Itoa(limit)}
	releaseSet, _ := getReleases(mirrorURL, queryParams)

	var releases []Release
	releaseSet, _ = getReleases(mirrorURL, queryParams)
	releases = append(releases, releaseSet...)
	for len(releaseSet) == limit {
		queryParams["after"] = releaseSet[len(releaseSet)-1].TimestampCreated.String()
		releaseSet, _ = getReleases(mirrorURL, queryParams)
		releases = append(releases, releaseSet...)
	}

	if !preRelease {
		for i, r := range releases {
			if r.IsPrerelease {
				releases = removePreReleases(releases, i)
			}
		}
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Version > releases[j].Version
	})
	return releases, nil
}

func GetTFRelease(mirrorURL, requestedVersion string) (Release, error) {
	resp, errURL := httpGet(mirrorURL+"/"+requestedVersion, nil)
	if errURL != nil {
		log.Printf("[Error] : Getting url: %v", errURL)
		os.Exit(1)
		return Release{}, errURL
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Error] : Retrieving contents from url: %s", mirrorURL)
		os.Exit(1)
	}

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, resp.Body); err != nil {
		log.Fatal(err)
	}
	var release Release
	if err := json.Unmarshal(body.Bytes(), &release); err != nil {
		log.Fatalf("%s: %s", err, body.String())
	}
	return release, nil

}

func removePreReleases(slice []Release, s int) []Release {
	return append(slice[:s], slice[s+1:]...)
}

//VersionExist : check if requested version exist
func VersionExist(val interface{}, array interface{}) (exists bool) {

	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
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
		if encountered[versionOnly] == true {
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
