package lib

import (
	"reflect"
	"testing"
)

const (
	hashiURL = "https://releases.hashicorp.com/terraform/"
)

// TestGetTFList : Get list from hashicorp
func TestGetTFList(t *testing.T) {

	list, _ := getTFList(hashiURL, true)

	val := "0.1.0"
	var exists bool

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(list)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
			}
		}
	}

	if !exists {
		t.Errorf("Not able to find version: %s", val)
	} else {
		t.Log("Write versions exist (expected)")
	}

}

// TestRemoveDuplicateVersions :  test to removed duplicate
func TestRemoveDuplicateVersions(t *testing.T) {

	testArray := []string{"0.0.1", "0.0.2", "0.0.3", "0.0.1", "0.12.0-beta1", "0.12.0-beta1"}

	list := removeDuplicateVersions(testArray)

	if len(list) == len(testArray) {
		t.Errorf("Not able to remove duplicate: %s\n", testArray)
	} else {
		t.Log("Write versions exist (expected)")
	}
}

// TestValidVersionFormat : test if func returns valid version format
// more regex testing at https://rubular.com/r/UvWXui7EU2icSb
func TestValidVersionFormat(t *testing.T) {

	var version string
	version = "0.11.8"

	valid := validVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9"

	valid = validVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.a"

	valid = validVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "22323"

	valid = validVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "@^&*!)!"

	valid = validVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9-beta1"

	valid = validVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "0.12.0-rc2"

	valid = validVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-boom"

	valid = validVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.4-1"

	valid = validVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		t.Errorf("Failed to verify version format: %s\n", version)
	}

}
