package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_convertData(t *testing.T) {
	recentFileContent := []byte("1.5.6\n0.13.0-rc1\n1.0.11\n")

	var recentFileData RecentFiles
	convertOldRecentFile(recentFileContent, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, 0, len(recentFileData.OpenTofu))
	assert.Equal(t, "1.5.6", recentFileData.Terraform[0])
	assert.Equal(t, "0.13.0-rc1", recentFileData.Terraform[1])
	assert.Equal(t, "1.0.11", recentFileData.Terraform[2])
}

func Test_saveFile(t *testing.T) {
	var recentFileData = RecentFiles{
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
	strings, err := getRecentVersions("../test-data/recent/recent_as_json/", distributionTerraform)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, 5, len(strings))
	assert.Equal(t, []string{"1.2.3 *recent", "4.5.6 *recent", "4.5.7 *recent", "4.5.8 *recent", "4.5.9 *recent"}, strings)
}

func Test_getRecentVersionsForOpenTofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json", distributionOpenTofu)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"6.6.6 *recent"}, strings)
}

func Test_addRecent(t *testing.T) {
	logger = InitLogger("DEBUG")
	temp, err := os.MkdirTemp("", "recent-test")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(temp)
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	addRecent("3.7.0", temp, distributionTerraform)
	addRecent("3.7.1", temp, distributionTerraform)
	addRecent("3.7.2", temp, distributionTerraform)
	filePath := filepath.Join(temp, ".terraform.versions", "RECENT")
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Could not open file %v", filePath)
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.2\",\"3.7.1\",\"3.7.0\"],\"opentofu\":null}", string(bytes))
	addRecent("3.7.0", temp, distributionTerraform)
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Could not open file %v", filePath)
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\",\"3.7.2\",\"3.7.1\"],\"opentofu\":null}", string(bytes))

	addRecent("1.1.1", temp, distributionOpenTofu)
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\",\"3.7.2\",\"3.7.1\"],\"opentofu\":[\"1.1.1\"]}", string(bytes))
}

func Test_prependExistingVersionIsMovingToTop(t *testing.T) {
	var recentFileData = RecentFiles{
		Terraform: []string{"1.2.3", "4.5.6", "7.7.7"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", distributionTerraform, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])

	prependRecentVersionToList("1.2.3", distributionTerraform, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "1.2.3", recentFileData.Terraform[0])
	assert.Equal(t, "7.7.7", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}

func Test_prependNewVersion(t *testing.T) {
	var recentFileData = RecentFiles{
		Terraform: []string{"1.2.3", "4.5.6", "4.5.7", "4.5.8", "4.5.9"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", distributionTerraform, &recentFileData)
	assert.Equal(t, 6, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}
