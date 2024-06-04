package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type RecentFiles struct {
	Terraform []string `json:"terraform"`
	OpenTofu  []string `json:"opentofu"`
}

func addRecent(requestedVersion string, installPath string, distribution string) {
	if !validVersionFormat(requestedVersion) {
		logger.Errorf("The version %q is not a valid version string and won't be stored", requestedVersion)
		return
	}
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	if CheckFileExist(recentFilePath) {
		unmarshalRecentFileData(recentFilePath, &recentFileData)
	}
	prependRecentVersionToList(requestedVersion, distribution, &recentFileData)
	saveRecentFile(recentFileData, recentFilePath)
}

func prependRecentVersionToList(version, distribution string, r *RecentFiles) {
	var sliceToCheck []string
	switch distribution {
	case distributionTerraform:
		sliceToCheck = r.Terraform
	case distributionOpenTofu:
		sliceToCheck = r.OpenTofu
	}
	for versionIndex, versionValue := range sliceToCheck {
		if versionValue == version {
			sliceToCheck = append(sliceToCheck[:versionIndex], sliceToCheck[versionIndex+1:]...)
		}
	}
	sliceToCheck = append([]string{version}, sliceToCheck...)

	switch distribution {
	case distributionTerraform:
		r.Terraform = sliceToCheck
	case distributionOpenTofu:
		r.OpenTofu = sliceToCheck
	}
}

func getRecentVersions(installPath string, dist string) ([]string, error) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	unmarshalRecentFileData(recentFilePath, &recentFileData)
	var listOfRecentVersions []string
	switch dist {
	case distributionTerraform:
		listOfRecentVersions = recentFileData.Terraform
	case distributionOpenTofu:
		listOfRecentVersions = recentFileData.OpenTofu
	}
	var maxCount int
	if len(listOfRecentVersions) >= 3 {
		maxCount = 3
	} else {
		maxCount = len(listOfRecentVersions)
	}
	for i := 0; i < maxCount; i++ {
		listOfRecentVersions[i] = listOfRecentVersions[i] + " *recent"
	}
	return listOfRecentVersions, nil
}

func unmarshalRecentFileData(recentFilePath string, recentFileData *RecentFiles) {
	recentFileContent, err := os.ReadFile(recentFilePath)
	if err != nil {
		logger.Errorf("Could not open recent versions file %q", recentFilePath)
	}
	if string(recentFileContent[0:1]) != "{" {
		convertOldRecentFile(recentFileContent, recentFileData)
	} else {
		err = json.Unmarshal(recentFileContent, &recentFileData)
		if err != nil {
			logger.Errorf("Could not unmarshal recent versions content from %q file", recentFilePath)
		}
	}
}

func convertOldRecentFile(content []byte, recentFileData *RecentFiles) {
	lines := strings.Split(string(content), "\n")
	for _, s := range lines {
		if s != "" {
			recentFileData.Terraform = append(recentFileData.Terraform, s)
		}
	}
}

func saveRecentFile(data RecentFiles, path string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Could not marshal data to JSON: %v", err)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		logger.Errorf("Could not save file %q: %v", path, err)
	}
}
