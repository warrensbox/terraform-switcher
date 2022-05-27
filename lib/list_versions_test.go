package lib_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	semver "github.com/hashicorp/go-version"

	"github.com/warrensbox/terraform-switcher/lib"
)

const (
	hashiURL = "https://api.releases.hashicorp.com/v1/releases/terraform"
)

// Used for constructing dummy release during tests
func NewRelease(s string) *lib.Release {
	v, err := semver.NewVersion(s)
	if err != nil {
		fmt.Println("Got here errorR")
		log.Fatalln(err)
	}
	return &lib.Release{Version: v}
}

// TestGetTFList : Get list from hashicorp

func TestGetTFReleases(t *testing.T) {

	listAll := true
	list, err := lib.GetTFReleases(hashiURL, listAll)
	if err != nil {
		log.Fatalln(err)
	}

	// Release metadata from https://releases.hashicorp.com/terraform/0.1.0
	// All fields besides the 3 used in lib.Release have been stripped out
	jSON := []byte(`{
	  "builds":[{"arch":"amd64","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_darwin_amd64.zip"},{"arch":"386","os":"linux","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_386.zip"},{"arch":"amd64","os":"linux","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_amd64.zip"},{"arch":"386","os":"windows","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_windows_386.zip"}],
	  "timestamp_created": "2017-07-12T06:41:24.000Z",
	  "version": "0.1.0"
	}`)
	var val lib.Release
	if err := json.Unmarshal(jSON, &val); err != nil {
		log.Fatalf("%s: %s", err, jSON)
	}
	var exists bool

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(list)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(&val, s.Index(i).Interface()) == true {
				exists = true
			}
		}
	}
	if !exists {
		log.Fatalf("Not able to find Release version: %s\n", val.Version)
	} else {
		t.Log("Write versions exist (expected)")
	}

}

//TestRemoveDuplicateVersions :  test to removed duplicate
func TestRemoveDuplicateVersions(t *testing.T) {

	jSON := []byte(`{
	  "builds":[{"arch":"amd64","os":"darwin","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_darwin_amd64.zip"},{"arch":"386","os":"linux","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_386.zip"},{"arch":"amd64","os":"linux","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_amd64.zip"},{"arch":"386","os":"windows","url":"https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_windows_386.zip"}],
	  "timestamp_created": "2017-07-12T06:41:24.000Z",
	  "version": "0.1.0"
	}`)
	var val lib.Release
	if err := json.Unmarshal(jSON, &val); err != nil {
		log.Fatalf("%s: %s", err, jSON)
	}

	var test_array = []*lib.Release{
		NewRelease("0.0.1"),
		NewRelease("0.0.2"),
		NewRelease("0.0.3"),
		NewRelease("0.0.1"),
		NewRelease("0.1.0"),
		NewRelease("0.12.0-beta1"),
		NewRelease("0.12.0-beta1"),
	}

	list := lib.RemoveDuplicateVersions(test_array)

	if len(list) == len(test_array) {
		fmt.Println(test_array[0].Version.Equal(test_array[3].Version))
		log.Fatalf("Not able to remove duplicate: %v\n", test_array)
	} else {
		t.Log("Write versions exist (expected)")
	}
}

//TestValidVersionFormat : test if func returns valid version format
// more regex testing at https://rubular.com/r/UvWXui7EU2icSb
func TestValidVersionFormat(t *testing.T) {

	var version string
	version = "0.11.8"

	valid := lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.a"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "22323"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "@^&*!)!"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9-beta1"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "0.12.0-rc2"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-boom"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-1"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

}
