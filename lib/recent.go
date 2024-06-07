package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type RecentFile struct {
	Terraform []string `json:"terraform"`
	OpenTofu  []string `json:"opentofu"`
}

func addRecent(requestedVersion string, installPath string, product Product) {
	if !validVersionFormat(requestedVersion) {
		logger.Errorf("The version %q is not a valid version string and won't be stored", requestedVersion)
		return
	}
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFile
	if CheckFileExist(recentFilePath) {
		unmarshalRecentFileData(recentFilePath, &recentFileData)
	}
	prependRecentVersionToList(requestedVersion, product, &recentFileData)
	saveRecentFile(recentFileData, recentFilePath)
}

func prependRecentVersionToList(version string, product Product, r *RecentFile) {
	sliceToCheck := product.GetRecentVersionProduct(r)
	for versionIndex, versionValue := range sliceToCheck {
		if versionValue == version {
			sliceToCheck = append(sliceToCheck[:versionIndex], sliceToCheck[versionIndex+1:]...)
		}
	}
	sliceToCheck = append([]string{version}, sliceToCheck...)

	product.SetRecentVersionProduct(r, sliceToCheck)
}

func getRecentVersions(installPath string, product Product) ([]string, error) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFile
	unmarshalRecentFileData(recentFilePath, &recentFileData)
	listOfRecentVersions := product.GetRecentVersionProduct(&recentFileData)
	var maxCount int
	if len(listOfRecentVersions) >= 5 {
		maxCount = 5
	} else {
		maxCount = len(listOfRecentVersions)
	}
	var returnedRecentVersions []string
	for i := 0; i < maxCount; i++ {
		returnedRecentVersions = append(returnedRecentVersions, listOfRecentVersions[i])
	}
	return returnedRecentVersions, nil
}

func unmarshalRecentFileData(recentFilePath string, recentFileData *RecentFile) {
	recentFileContent, err := os.ReadFile(recentFilePath)
	if err != nil {
		logger.Errorf("Could not open recent versions file %q", recentFilePath)
	}
	if len(string(recentFileContent)) >= 1 && string(recentFileContent[0:1]) != "{" {
		convertOldRecentFile(recentFileContent, recentFileData)
	} else {
		err = json.Unmarshal(recentFileContent, &recentFileData)
		if err != nil {
			logger.Errorf("Could not unmarshal recent versions content from %q file", recentFilePath)
		}
	}
}

func convertOldRecentFile(content []byte, recentFileData *RecentFile) {
	lines := strings.Split(string(content), "\n")
	for _, s := range lines {
		if s != "" {
			recentFileData.Terraform = append(recentFileData.Terraform, s)
		}
	}
}

func saveRecentFile(data RecentFile, path string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Could not marshal data to JSON: %v", err)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		logger.Errorf("Could not save file %q: %v", path, err)
	}
}
