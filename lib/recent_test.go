package lib

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertData(t *testing.T) {
	recentFileContent := []byte("1.5.6\n0.13.0-rc1\n1.0.11\n")

	var recentFileData RecentFile
	convertOldRecentFile(recentFileContent, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, 0, len(recentFileData.OpenTofu))
	assert.Equal(t, "1.5.6", recentFileData.Terraform[0])
	assert.Equal(t, "0.13.0-rc1", recentFileData.Terraform[1])
	assert.Equal(t, "1.0.11", recentFileData.Terraform[2])

	// Test with empty data
	recentFileContent = []byte("")
	recentFileData = RecentFile{}
	convertOldRecentFile(recentFileContent, &recentFileData)
	assert.Equal(t, 0, len(recentFileData.Terraform))
	assert.Equal(t, 0, len(recentFileData.OpenTofu))
}

// Test_unmarshalRecentFileData_conversion : Test unmarshalRecentFileData with old version format
func Test_unmarshalRecentFileData_conversion(t *testing.T) {
	var expectedRecentFileData RecentFile

	t.Log("Test empty file")
	expectedRecentFileData = RecentFile{}
	performUnmarshalRecentFileDataTest(t, "", &expectedRecentFileData)

	t.Log("Test one version")
	expectedRecentFileData = RecentFile{
		Terraform: []string{"1.3.2"},
	}
	performUnmarshalRecentFileDataTest(t, "1.3.2", &expectedRecentFileData)

	t.Log("Test trailing new line")
	expectedRecentFileData = RecentFile{
		Terraform: []string{"1.3.2"},
	}
	performUnmarshalRecentFileDataTest(t, "1.3.2\n", &expectedRecentFileData)

	t.Log("Test multiple versions")
	expectedRecentFileData = RecentFile{
		Terraform: []string{"1.3.2", "1.2.3"},
	}
	performUnmarshalRecentFileDataTest(t, "1.3.2\n1.2.3\n", &expectedRecentFileData)
}

// Test_unmarshalRecentFileData : Test unmarshalRecentFileData
func Test_unmarshalRecentFileData(t *testing.T) {
	var expectedRecentFileData RecentFile

	t.Log("Test only Terraform")
	expectedRecentFileData = RecentFile{
		Terraform: []string{"1.5.0", "1.6.0"},
	}
	performUnmarshalRecentFileDataTest(t, `{"terraform": ["1.5.0", "1.6.0"]}`, &expectedRecentFileData)

	t.Log("Test only OpenTofu")
	expectedRecentFileData = RecentFile{
		OpenTofu: []string{"1.5.0", "1.6.0"},
	}
	performUnmarshalRecentFileDataTest(t, `{"opentofu": ["1.5.0", "1.6.0"]}`, &expectedRecentFileData)

	t.Log("Test both")
	expectedRecentFileData = RecentFile{
		Terraform: []string{"1.2.3", "1.3.2"},
		OpenTofu:  []string{"1.5.0", "1.6.0"},
	}
	performUnmarshalRecentFileDataTest(t, `{"terraform": ["1.2.3", "1.3.2"], "opentofu": ["1.5.0", "1.6.0"]}`, &expectedRecentFileData)
}

func performUnmarshalRecentFileDataTest(t *testing.T, recentFileContent string, expectedRecentFileData *RecentFile) {
	temp, err := os.MkdirTemp("", "recent-test")
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	defer os.RemoveAll(temp)
	pathToTempFile := filepath.Join(temp, "recent.json")

	os.WriteFile(pathToTempFile, []byte(recentFileContent), 0600)

	var recentFileData = RecentFile{}
	unmarshalRecentFileData(pathToTempFile, &recentFileData)
	t.Log("Comparing Terraform versions")
	assert.Equal(t, expectedRecentFileData.Terraform, recentFileData.Terraform)
	t.Log("Comparing OpenTofu versions")
	assert.Equal(t, expectedRecentFileData.OpenTofu, recentFileData.OpenTofu)
}

func Test_saveFile(t *testing.T) {
	var recentFileData = RecentFile{
		Terraform: []string{"1.2.3", "4.5.6"},
		OpenTofu:  []string{"6.6.6"},
	}
	temp, err := os.MkdirTemp("", "recent-test")
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	defer func(path string) {
		_ = os.RemoveAll(temp)
	}(temp)
	pathToTempFile := filepath.Join(temp, "recent.json")
	saveRecentFile(recentFileData, pathToTempFile)

	content, err := os.ReadFile(pathToTempFile)
	if err != nil {
		t.Errorf("Could not read converted file %v", pathToTempFile)
	}
	assert.Equal(t, "{\"terraform\":[\"1.2.3\",\"4.5.6\"],\"opentofu\":[\"6.6.6\"]}", string(content))
}

func Test_getRecentVersionsForTerraform(t *testing.T) {
	logger = InitLogger("DEBUG")
	product := GetProductById("terraform")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json/", product)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, 5, len(strings))
	assert.Equal(t, []string{"1.2.3", "4.5.6", "4.5.7", "4.5.8", "4.5.9"}, strings)
}

func Test_getRecentVersionsForOpenTofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	product := GetProductById("opentofu")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json", product)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"6.6.6"}, strings)
}

func Test_addRecent(t *testing.T) {
	logger = InitLogger("DEBUG")
	terraform := GetProductById("terraform")
	opentofu := GetProductById("opentofu")
	temp, err := os.MkdirTemp("", "recent-test")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(temp)
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	addRecent("3.7.0", temp, terraform)
	addRecent("3.7.1", temp, terraform)
	addRecent("3.7.2", temp, terraform)
	filePath := filepath.Join(temp, ".terraform.versions", "RECENT")
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Could not open file %v", filePath)
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.2\",\"3.7.1\",\"3.7.0\"],\"opentofu\":null}", string(bytes))
	addRecent("3.7.0", temp, terraform)
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Could not open file %v", filePath)
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\",\"3.7.2\",\"3.7.1\"],\"opentofu\":null}", string(bytes))

	addRecent("1.1.1", temp, opentofu)
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\",\"3.7.2\",\"3.7.1\"],\"opentofu\":[\"1.1.1\"]}", string(bytes))
}

func Test_prependExistingVersionIsMovingToTop(t *testing.T) {
	product := GetProductById("terraform")
	var recentFileData = RecentFile{
		Terraform: []string{"1.2.3", "4.5.6", "7.7.7"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", product, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])

	prependRecentVersionToList("1.2.3", product, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "1.2.3", recentFileData.Terraform[0])
	assert.Equal(t, "7.7.7", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}

func Test_prependNewVersion(t *testing.T) {
	product := GetProductById("terraform")
	var recentFileData = RecentFile{
		Terraform: []string{"1.2.3", "4.5.6", "4.5.7", "4.5.8", "4.5.9"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", product, &recentFileData)
	assert.Equal(t, 6, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}
