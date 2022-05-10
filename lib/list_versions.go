package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
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
type Releases struct {
	Builds []struct {
		Arch string `json:"arch"`
		Os   string `json:"os"`
		Url  string `json:"url"`
	} `json:"builds"`
	IsPrerelease bool   `json:"is_prerelease"`
	LicenseClass string `json:"license_class"`
	Name         string `json:"name"`
	Status       struct {
		State            string    `json:"state"`
		TimestampUpdated time.Time `json:"timestamp_updated"`
	} `json:"status"`
	TimestampCreated     time.Time `json:"timestamp_created"`
	TimestampUpdated     time.Time `json:"timestamp_updated"`
	UrlShasums           string    `json:"url_shasums"`
	UrlShasumsSignatures []string  `json:"url_shasums_signatures"`
	UrlSourceRepository  string    `json:"url_source_repository"`
	Version              string    `json:"version"`
}

//GetTFList :  Get the list of available terraform version given the hashicorp url
func GetTFList(mirrorURL string, preRelease bool) ([]string, error) {

	result, error := GetTFURLBody(mirrorURL)
	if error != nil {
		return nil, error
	}

	var tfVersionList tfVersionList
	var semver string
	if preRelease == true {
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver = `(\d+\.\d+\.\d+)(-[a-zA-z]+\d*)?`
	} else if preRelease == false {
		// Getting versions from body; should return match /X.X.X/ where X is a number
		// without the ending '"' pre-release folders would be tried and break.
		semver = `(\d+\.\d+\.\d+)`
	}
	var versions []string
	for i := range result {
		versions = append(versions, result[i].Version)
	}

	r, _ := regexp.Compile(semver)
	for i := range result {
		if r.MatchString(versions[i]) {
			str := r.FindString(versions[i])
			trimstr := strings.Trim(str, "/\"") //remove "/" from /X.X.X/
			tfVersionList.tflist = append(tfVersionList.tflist, trimstr)
		}
	}

	if len(tfVersionList.tflist) == 0 {
		fmt.Printf("Cannot get list from mirror: %s\n", mirrorURL)
	}

	return tfVersionList.tflist, nil

}

//GetTFLatest :  Get the latest terraform version given the hashicorp url
func GetTFLatest(mirrorURL string) (string, error) {

	result, error := GetTFURLBody(mirrorURL)
	if error != nil {
		return "", error
	}
	// Getting versions from body; should return match /X.X.X/ where X is a number
	semver := `\/(\d+\.\d+\.\d+)\/?"`
	r, _ := regexp.Compile(semver)

	var versions []string
	for i := range result {
		versions = append(versions, result[i].Version)
	}
	for i := range result {
		if r.MatchString(versions[i]) {
			str := r.FindString(versions[i])
			trimstr := strings.Trim(str, "/") //remove "/" from /X.X.X/
			return trimstr, nil
		}
	}

	return "", nil
}

//GetTFLatestImplicit :  Get the latest implicit terraform version given the hashicorp url
func GetTFLatestImplicit(mirrorURL string, preRelease bool, version string) (string, error) {

	if preRelease == true {
		//TODO: use GetTFList() instead of GetTFURLBody
		result, error := GetTFURLBody(mirrorURL)
		if error != nil {
			return "", error
		}
		// Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
		semver := fmt.Sprintf(`\/(%s{1}\.\d+\-[a-zA-z]+\d*)\/?"`, version)
		r, err := regexp.Compile(semver)
		if err != nil {
			return "", err
		}
		var versions []string
		for i := range result {
			versions = append(versions, result[i].Version)
		}

		for i := range versions {
			if r.MatchString(versions[i]) {
				str := r.FindString(versions[i])
				trimstr := strings.Trim(str, "/\"") //remove '/' or '"' from /X.X.X/" or /X.X.X"
				return trimstr, nil
			}
		}
	} else if preRelease == false {
		listAll := false
		tflist, _ := GetTFList(mirrorURL, listAll) //get list of versions
		version = fmt.Sprintf("~> %v", version)
		semv, err := SemVerParser(&version, tflist)
		if err != nil {
			return "", err
		}
		return semv, nil
	}
	return "", nil
}

//GetTFURLBody : Get list of terraform versions from hashicorp releases
func GetTFURLBody(mirrorURL string) ([]Releases, error) {

	var releases []Releases

	resp, errURL := http.Get(mirrorURL)
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
