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
	convertData(recentFileContent, &recentFileData)
	assert.Equal(t, "1.5.6", recentFileData.Terraform[0])
	assert.Equal(t, "0.13.0-rc1", recentFileData.Terraform[1])
	assert.Equal(t, "1.0.11", recentFileData.Terraform[2])
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, 0, len(recentFileData.Tofu))
}

func Test_saveFile(t *testing.T) {
	var recentFileData = RecentFiles{
		Terraform: []string{"1.2.3", "4.5.6"},
		Tofu:      []string{"6.6.6"},
	}
	temp, err := os.MkdirTemp("", "recent-test")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(temp)
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	pathToTempFile := filepath.Join(temp, "recent.json")
	saveFile(recentFileData, pathToTempFile)

	content, err := os.ReadFile(pathToTempFile)
	if err != nil {
		t.Errorf("Could not read converted file %v", pathToTempFile)
	}
	assert.Equal(t, "{\"terraform\":[\"1.2.3\",\"4.5.6\"],\"tofu\":[\"6.6.6\"]}", string(content))
}

func Test_getRecentVersionsGenericForTerraform(t *testing.T) {
	logger = InitLogger("DEBUG")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json/", distTerraform)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"1.2.3", "4.5.6"}, strings)
}

func Test_getRecentVersionsGenericForTofu(t *testing.T) {
	logger = InitLogger("DEBUG")
	strings, err := getRecentVersions("../test-data/recent/recent_as_json", distTofu)
	if err != nil {
		t.Error("Unable to get versions from recent file")
	}
	assert.Equal(t, []string{"6.6.6"}, strings)
}

func Test_addRecentGenericForTerraform(t *testing.T) {
	logger = InitLogger("DEBUG")
	temp, err := os.MkdirTemp("", "recent-test")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(temp)
	if err != nil {
		t.Errorf("Could not create temporary directory")
	}
	addRecent("3.7.0", temp, distTerraform)
	bytes, err := os.ReadFile(filepath.Join(temp, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\"],\"tofu\":null}", string(bytes))

	addRecent("1.1.1", temp, distTofu)
	bytes, err = os.ReadFile(filepath.Join(temp, ".terraform.versions", "RECENT"))
	if err != nil {
		t.Error("Could not open file")
		t.Error(err)
	}
	assert.Equal(t, "{\"terraform\":[\"3.7.0\"],\"tofu\":[\"1.1.1\"]}", string(bytes))
}
