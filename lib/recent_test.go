package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_convertData(t *testing.T) {
	recentFileContent, err := os.ReadFile("../test-data/recent/RECENT_OLD_FORMAT")
	if err != nil {
		logger.Error("Could not open file ../test-data/recent/RECENT_OLD_FORMAT")
	}

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
	assert.Equal(t, "{\"terraform\":[\"1.2.3\",\"4.5.6\"],\"openTofu\":[\"6.6.6\"]}", string(content))
}

func Test_getRecentVersionsForTerraform(t *testing.T) {
	logger = InitLogger("DEBUG")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json/", distributionTerraform)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"1.2.3", "4.5.6"}, strings)
}

func Test_getRecentVersionsForOpenTofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json", distributionOpenTofu)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"6.6.6"}, strings)
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
	bytes, err := os.ReadFile(filepath.Join(temp, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\",\"3.7.1\",\"3.7.2\"],\"openTofu\":null}", string(bytes))
	addRecent("3.7.2", temp, distributionTerraform)
	assert.Equal(t, "{\"terraform\":[\"3.7.2\",\"3.7.0\",\"3.7.1\"],\"openTofu\":null}", string(bytes))

	addRecent("1.1.1", temp, distributionOpenTofu)
	bytes, err = os.ReadFile(filepath.Join(temp, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.2\",\"3.7.0\",\"3.7.1\"],\"openTofu\":[\"1.1.1\"]}", string(bytes))
}

func Test_prependExistingVersionIsMovingToTop(t *testing.T) {
	var recentFileData = RecentFiles{
		Terraform: []string{"1.2.3", "4.5.6", "7.7.7"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", "", distributionTerraform, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])

	prependRecentVersionToList("1.2.3", "", distributionTerraform, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "1.2.3", recentFileData.Terraform[0])
	assert.Equal(t, "7.7.7", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}

func Test_prependNewVersion(t *testing.T) {
	var recentFileData = RecentFiles{
		Terraform: []string{"1.2.3", "4.5.6"},
		OpenTofu:  []string{"6.6.6"},
	}
	prependRecentVersionToList("7.7.7", "", distributionTerraform, &recentFileData)
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, "7.7.7", recentFileData.Terraform[0])
	assert.Equal(t, "1.2.3", recentFileData.Terraform[1])
	assert.Equal(t, "4.5.6", recentFileData.Terraform[2])
}
