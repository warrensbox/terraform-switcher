package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type tfVersionList struct {
	tflist []string
}

type Release struct {
	IsPrerelease     bool            `json:"is_prerelease"`
	TimestampCreated strfmt.DateTime `json:"timestamp_created"`
	Version          string          `json:"version"`
}

//GetTFLatest :  Get the latest terraform version given the hashicorp url
func GetTFLatest(mirrorURL string, preRelease bool) (*Release, error) {

	result, error := GetTFReleases(mirrorURL, preRelease)
	if error != nil {
		return nil, error
	}
	for i := range result {
		if !result[i].IsPrerelease {
			return result[i], nil
		}
	}

	return nil, nil
}

//GetTFLatestImplicit :  Get the latest implicit terraform version given the hashicorp url
func GetTFLatestImplicit(mirrorURL string, preRelease bool, version string) (string, error) {

	if preRelease == true {
		//TODO: use GetTFList() instead of GetTFURLBody
		releases, error := GetTFReleases(mirrorURL, preRelease)
		if error != nil {
			return "", error
		}
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver := fmt.Sprintf(`\/(%s{1}\.\d+\-[a-zA-z]+\d*)\/?"`, version)
		r, err := regexp.Compile(semver)
		if err != nil {
			return "", err
		}
		for i := range releases {
			if r.MatchString(releases[i].Version) {
				str := r.FindString(releases[i].Version)
				return str, nil
			}
		}
	} else if preRelease == false {
		tflist, err := GetTFReleases(mirrorURL, preRelease)
		version = fmt.Sprintf("~> %v", version)
		semv, err := SemVerParser(&version, tflist)
		if err != nil {
			return "", err
		}
		return semv, nil
	}
	return "", nil
}

func httpGet(url, limit string, timestamp strfmt.DateTime) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("limit", limit)
	if timestamp.String() != time.RFC3339 {
		q.Add("after", timestamp.String())
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		// should clean this up
		return nil, errors.New("issue during request")
	}
	return res, nil
}

func getReleases(mirrorURL, limit string, timestamp strfmt.DateTime) ([]*Release, error) {
	var releases []*Release
	resp, errURL := httpGet(mirrorURL, limit, timestamp)
	if errURL != nil {
		log.Printf("[Error] : Getting url: %v", errURL)
		os.Exit(1)
		return nil, errURL
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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

func GetTFReleases(mirrorURL string, preRelease bool) ([]*Release, error) {
	// temp testing
	limit := "20"
	t, _ := time.Parse(time.RFC3339, time.RFC3339)
	tmpTime := strfmt.DateTime(t)

	rel, _ := getReleases(mirrorURL, limit, tmpTime)

	// this is ugly
	lastTimestamp := tmpTime
	currentTimestamp := strfmt.NewDateTime()

	var releases []*Release

	for lastTimestamp != currentTimestamp {
		rel, _ = getReleases(mirrorURL, limit, rel[len(rel)-1].TimestampCreated)
		currentTimestamp = rel[len(rel)-1].TimestampCreated
		if len(releases) > 0 {
			lastTimestamp = releases[len(releases)-1].TimestampCreated
		}
		releases = append(releases, rel...)
	}

	if !preRelease {
		for i, r := range releases {
			if r.IsPrerelease {
				releases = removePreReleases(releases, i)
			}
		}
	}
	return releases, nil
}

func removePreReleases(slice []*Release, s int) []*Release {
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
