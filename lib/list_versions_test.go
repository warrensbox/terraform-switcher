package lib_test

import (
	"github.com/go-openapi/strfmt"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
	"reflect"
	"testing"
)

const (
	hashiURL = "https://api.releases.hashicorp.com/v1/releases/terraform"
)

// TestGetTFList : Get list from hashicorp

func TestGetTFReleases(t *testing.T) {

	listAll := true
	list, _ := lib.GetTFReleases(hashiURL, listAll)

	//val := lib.Release{Version: "0.1.0"}
	time, parseErr := strfmt.ParseDateTime("2017-07-12T06:41:24.000Z")
	if parseErr != nil {
		log.Fatalf(parseErr.Error())
	}
	val := lib.Release{
		Builds: []lib.Builds{
			{
				Arch: "amd64",
				Os:   "darwin",
				Url:  "https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_darwin_amd64.zip",
			},
			{
				Arch: "386",
				Os:   "linux",
				Url:  "https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_386.zip",
			},
			{
				Arch: "amd64",
				Os:   "linux",
				Url:  "https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_linux_amd64.zip",
			},
			{
				Arch: "386",
				Os:   "windows",
				Url:  "https://releases.hashicorp.com/terraform/0.1.0/terraform_0.1.0_windows_386.zip",
			},
		},
		TimestampCreated: time,
		Version:          "0.1.0",
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
		log.Fatalf("Not able to find version: %s\n", val.Version)
	} else {
		t.Log("Write versions exist (expected)")
	}

}

//TestRemoveDuplicateVersions :  test to removed duplicate
func TestRemoveDuplicateVersions(t *testing.T) {

	test_array := []string{"0.0.1", "0.0.2", "0.0.3", "0.0.1", "0.12.0-beta1", "0.12.0-beta1"}

	list := lib.RemoveDuplicateVersions(test_array)

	if len(list) == len(test_array) {
		log.Fatalf("Not able to remove duplicate: %s\n", test_array)
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
