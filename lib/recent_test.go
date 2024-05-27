package lib

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_convertData(t *testing.T) {
	recentFileContent, err := os.ReadFile("../test-data/RECENT")
	if err != nil {
		logger.Error("Could not open file ../test-data/RECENT")
	}

	var recentFileData RecentFiles
	convertData(recentFileContent, &recentFileData)
	assert.Equal(t, "1.5.6", recentFileData.Terraform[0])
	assert.Equal(t, "0.13.0-rc1", recentFileData.Terraform[1])
	assert.Equal(t, "1.0.11", recentFileData.Terraform[2])
	assert.Equal(t, 3, len(recentFileData.Terraform))
	assert.Equal(t, 0, len(recentFileData.Tofu))
}
