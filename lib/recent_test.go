package lib

import (
	"os"
	"path"
	"testing"
)

// Test_appendRecentVersionToList : Test appendRecentVersionToList method
func Test_appendRecentVersionToList(t *testing.T) {
	// Empty slice, adding first entry
	var expected []string = []string{"1.5.3"}
	actual := appendRecentVersionToList([]string{}, "1.5.3")
	if err := compareLists(actual, expected); err != nil {
		t.Error(err)
	}

	// Adding second entry
	expected = []string{"1.5.3", "1.5.1"}
	actual = appendRecentVersionToList([]string{"1.5.1"}, "1.5.3")
	if err := compareLists(actual, expected); err != nil {
		t.Error(err)
	}

	// Adding fourth, popping off last entry
	expected = []string{"1.0.0", "1.5.3", "1.5.1"}
	actual = appendRecentVersionToList([]string{"1.5.3", "1.5.1", "2.0.0"}, "1.0.0")
	if err := compareLists(actual, expected); err != nil {
		t.Error(err)
	}

	// Adding duplicate, ensure it's moved
	expected = []string{"1.5.3", "1.5.1", "2.0.0"}
	actual = appendRecentVersionToList([]string{"1.5.1", "1.5.3", "2.0.0"}, "1.5.3")
	if err := compareLists(actual, expected); err != nil {
		t.Error(err)
	}

	// Adding same version and ensure nothing changes
	expected = []string{"1.5.3", "1.5.1", "2.0.0"}
	actual = appendRecentVersionToList([]string{"1.5.3", "1.5.1", "2.0.0"}, "1.5.3")
	if err := compareLists(actual, expected); err != nil {
		t.Error(err)
	}
}

// Test_addRecentVersion_no_version : Test addRecentVersion with no recent version file
func Test_addRecentVersion_no_file(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	product := GetProductById("terraform")

	addRecentVersion(product, "1.2.3", tempDir)

	// Ensure recent versions file matches expected
	recentData, err := os.ReadFile(path.Join(tempDir, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Fatal(err)
	}
	expectedRecentData := "{\"Terraform\":[\"1.2.3\"],\"OpenTofu\":null}"
	if string(recentData) != expectedRecentData {
		t.Errorf("Recent file data does not match expected. Expected: %q, actual: %q", expectedRecentData, string(recentData))
	}
}

// Test_addRecentVersion_no_version : Test addRecentVersion with pre-existing configuration
func Test_addRecentVersion_pre_existing(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tfVersionsDir := path.Join(tempDir, ".terraform.versions")
	os.Mkdir(tfVersionsDir, 0750)
	recentFile := path.Join(tfVersionsDir, "RECENT")
	os.WriteFile(recentFile, []byte("{\"Terraform\":[\"1.2.0\"],\"OpenTofu\":null}"), 0640)

	product := GetProductById("terraform")

	addRecentVersion(product, "1.2.1", tempDir)

	// Ensure recent versions file matches expected
	recentData, err := os.ReadFile(path.Join(tempDir, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Fatal(err)
	}
	expectedRecentData := "{\"Terraform\":[\"1.2.1\",\"1.2.0\"],\"OpenTofu\":null}"
	if string(recentData) != expectedRecentData {
		t.Errorf("Recent file data does not match expected. Expected: %q, actual: %q", expectedRecentData, string(recentData))
	}
}

// Test_convertLegacyRecentData_empty_string : Test convertLegacyRecentData with empty string
func Test_convertLegacyRecentData_empty_string(t *testing.T) {
	var expected = []string{}
	var recentFile RecentFile
	convertLegacyRecentData([]byte(""), &recentFile)
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
}

// Test_convertLegacyRecentData_only_affects_terraform : Test convertLegacyRecentData ensuring that it only modifies Terraform
func Test_convertLegacyRecentData_only_affects_terraform(t *testing.T) {
	var expected = []string{"2.0.0", "3.0.0"}
	var recentFile RecentFile = RecentFile{
		Terraform: []string{},
		OpenTofu:  []string{"1.2.3", "1.2.4"},
	}
	convertLegacyRecentData([]byte("2.0.0\n3.0.0"), &recentFile)
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
}

// Test_getRecentFileData : Test getRecentFileData with valid file
func Test_getRecentFileData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tfVersionsDir := path.Join(tempDir, ".terraform.versions")
	os.Mkdir(tfVersionsDir, 0750)
	recentFilePath := path.Join(tfVersionsDir, "RECENT")
	os.WriteFile(recentFilePath, []byte("{\"Terraform\":[\"1.2.0\",\"1.4.3\"],\"OpenTofu\":[\"2.0.0\", \"2.1.0\"]}"), 0640)

	recentFile := getRecentFileData(tempDir)

	var expected []string = []string{"1.2.0", "1.4.3"}
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
	expected = []string{"2.0.0", "2.1.0"}
	if err := compareLists(recentFile.OpenTofu, expected); err != nil {
		t.Error(err)
	}
}

// Test_getRecentFileData_legacy_data : Test getRecentFileData with file with legacy format
func Test_getRecentFileData_legacy_data(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tfVersionsDir := path.Join(tempDir, ".terraform.versions")
	os.Mkdir(tfVersionsDir, 0750)
	recentFilePath := path.Join(tfVersionsDir, "RECENT")
	os.WriteFile(recentFilePath, []byte("1.5.2\n1.6.1\n1.7.0"), 0640)

	recentFile := getRecentFileData(tempDir)

	var expected []string = []string{"1.5.2", "1.6.1", "1.7.0"}
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
	expected = []string{}
	if err := compareLists(recentFile.OpenTofu, expected); err != nil {
		t.Error(err)
	}
}

// Test_getRecentFileData_no_file : Test getRecentFileData with no file
func Test_getRecentFileData_no_file(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	recentFile := getRecentFileData(tempDir)

	var expected []string = []string{}
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
	if err := compareLists(recentFile.OpenTofu, expected); err != nil {
		t.Error(err)
	}
}

// Test_getRecentFileData_invalid_file_data : Test getRecentFileData with invalid file data
func Test_getRecentFileData_invalid_file_data(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tfVersionsDir := path.Join(tempDir, ".terraform.versions")
	os.Mkdir(tfVersionsDir, 0750)
	recentFilePath := path.Join(tfVersionsDir, "RECENT")
	os.WriteFile(recentFilePath, []byte("{NOTVALIDJSON}"), 0640)

	recentFile := getRecentFileData(tempDir)

	var expected []string = []string{}
	if err := compareLists(recentFile.Terraform, expected); err != nil {
		t.Error(err)
	}
	if err := compareLists(recentFile.OpenTofu, expected); err != nil {
		t.Error(err)
	}
}
