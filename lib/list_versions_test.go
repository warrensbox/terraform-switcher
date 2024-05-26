package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const (
	hashiURL = "https://releases.hashicorp.com/terraform/"

	hashicorpBody = `
	<li>
	<a href="/terraform/0.12.3-beta1/">terraform_0.12.3-beta1</a>
	</li>
	<li>
	<a href="/terraform/0.12.2/">terraform_0.12.2</a>
	</li>
	<li>
	<a href="/terraform/0.12.1/">terraform_0.12.1</a>
	</li>
	<li>
	<a href="/terraform/0.12.0/">terraform_0.12.0</a>
	</li>
	<li>
	<a href="/terraform/0.12.0-rc1/">terraform_0.12.0-rc1</a>
	</li>
	<li>
	<a href="/terraform/0.12.0-beta2/">terraform_0.12.0-beta2</a>
	</li>
	<li>
	<a href="/terraform/0.11.13/">terraform_0.11.13</a>
	</li>
`

	openTofuBody = `
<!DOCTYPE html>
<html>
<head>
	<title>OpenTofu releases</title>
</head>
<body>
<ul><li><a href="/tofu/1.7.1-beta1/">tofu_1.7.1-beta1</a></li><li><a href="/tofu/1.7.0/">tofu_1.7.0</a></li><li><a href="/tofu/1.7.0-rc1/">tofu_1.7.0-rc1</a></li><li><a href="/tofu/1.7.0-beta1/">tofu_1.7.0-beta1</a></li><li><a href="/tofu/1.7.0-alpha1/">tofu_1.7.0-alpha1</a></li><li><a href="/tofu/1.6.2/">tofu_1.6.2</a></li><li><a href="/tofu/1.6.0-alpha1/">tofu_1.6.0-alpha1</a></li></ul>
</body>
</html>
`
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

func compareLists(actual []string, expected []string) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("Slices are not equal length: Expected: %v, actual: %v", expected, actual)
	}

	for i, v := range expected {
		if v != actual[i] {
			return fmt.Errorf("Elements are not the same. Expected: " + v + ", actual: " + actual[i])
		}
	}
	return nil
}

func getMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/hashicorp/":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(hashicorpBody))
		case "/opentofu/":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(openTofuBody))
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}

// TestGetVersionsFromBodyHashicorp :  test hashicorp release body
func TestGetVersionsFromBodyHashicorp(t *testing.T) {
	var testTfVersionList tfVersionList
	getVersionsFromBody(hashicorpBody, false, &testTfVersionList)
	expectedVersion := []string{"0.12.2", "0.12.1", "0.12.0", "0.11.13"}
	if err := compareLists(testTfVersionList.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	getVersionsFromBody(hashicorpBody, true, &testTfVersionListPre)
	expectedVersion = []string{"0.12.3-beta1", "0.12.2", "0.12.1", "0.12.0", "0.12.0-rc1", "0.12.0-beta2", "0.11.13"}
	if err := compareLists(testTfVersionListPre.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}
}

// TestGetVersionsFromBodyOpenTofu :  test OpenTofu release body
func TestGetVersionsFromBodyOpenTofu(t *testing.T) {
	var testTfVersionList tfVersionList
	getVersionsFromBody(openTofuBody, false, &testTfVersionList)
	expectedVersion := []string{"1.7.0", "1.6.2"}
	if err := compareLists(testTfVersionList.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}

	// Test pre-release
	var testTfVersionListPre tfVersionList
	getVersionsFromBody(openTofuBody, true, &testTfVersionListPre)
	expectedVersion = []string{"1.7.1-beta1", "1.7.0", "1.7.0-rc1", "1.7.0-beta1", "1.7.0-alpha1", "1.6.2", "1.6.0-alpha1"}
	if err := compareLists(testTfVersionListPre.tflist, expectedVersion); err != nil {
		t.Errorf("Parsed version does not match expected versions: %v", err)
	}
}

// TestGetTFLatest : Test getTFLatest
func TestGetTFLatest(t *testing.T) {
	server := getMockServer()
	defer server.Close()

	version, err := getTFLatest(fmt.Sprintf("%s/%s", server.URL, "hashicorp"))
	if err != nil {
		t.Error(err)
	}
	expectedVersion := "0.12.2"
	if version != expectedVersion {
		t.Errorf("Expected latest version does not match. Expected: %s, actual: %s", expectedVersion, version)
	}
}

// TestGetTFLatestImplicit : Test getTFLatestImplicit
func TestGetTFLatestImplicit(t *testing.T) {
	tName := "version=%s_preRelease=%v"
	t.Run(fmt.Sprintf(tName, "0.11.0", false), func(t *testing.T) { testGetTFLatestImplicit(t, "0.11.0", false, "0.11.13") })
	t.Run(fmt.Sprintf(tName, "0.11.0", true), func(t *testing.T) { testGetTFLatestImplicit(t, "0.11.0", true, "0.11.13") })
	t.Run(fmt.Sprintf(tName, "0.11", false), func(t *testing.T) { testGetTFLatestImplicit(t, "0.11", false, "0.12.2") })
	t.Run(fmt.Sprintf(tName, "0.11", true), func(t *testing.T) { testGetTFLatestImplicit(t, "0.11", true, "0.12.3-beta1") })
}

func testGetTFLatestImplicit(t *testing.T, version string, preRelease bool, expectedVersion string) {
	server := getMockServer()
	defer server.Close()

	version, err := getTFLatestImplicit(fmt.Sprintf("%s/%s", server.URL, "hashicorp"), preRelease, version)
	if err != nil {
		t.Error(err)
	}
	if version != expectedVersion {
		t.Errorf("Expected latest version does not match. Expected: %s, actual: %s", expectedVersion, version)
	}
}

// TestGetTFURLBody :  Test getTFURLBody method
func TestGetTFURLBody(t *testing.T) {
	server := getMockServer()
	defer server.Close()

	body, err := getTFURLBody(fmt.Sprintf("%s/%s", server.URL, "hashicorp"))
	if err != nil {
		t.Error(err)
	}
	if body != hashicorpBody {
		t.Errorf("Body not returned correctly. Expected: %s, actual: %s", hashicorpBody, body)
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
