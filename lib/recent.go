package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type RecentFiles struct {
	Terraform []string `json:"terraform"`
	Tofu      []string `json:"tofu"`
}

func addRecent(requestedVersion string, installPath string, dist string) {
	if !validVersionFormat(requestedVersion) {
		logger.Errorf("The version %s is not a valid version string and won't be stored", requestedVersion)
		return
	}
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	if CheckFileExist(recentFilePath) {
		unmarshal(recentFilePath, &recentFileData)
	}
	var sliceToCheck []string
	if dist == distTerraform {
		sliceToCheck = recentFileData.Terraform
	} else if dist == distTofu {
		sliceToCheck = recentFileData.Tofu
	}
	for _, v := range sliceToCheck {
		if v == requestedVersion {
			// entry already exists. Nothing to do
			return
		}
	}
	if dist == distTerraform {
		recentFileData.Terraform = append(recentFileData.Terraform, requestedVersion)
	} else if dist == distTofu {
		recentFileData.Tofu = append(recentFileData.Tofu, requestedVersion)
	}
	saveFile(recentFileData, recentFilePath)
}

func getRecentVersions(installPath string, dist string) ([]string, error) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	unmarshal(recentFilePath, &recentFileData)
	if dist == distTerraform {
		return recentFileData.Terraform, nil
	} else if dist == distTofu {
		return recentFileData.Tofu, nil
	}
	return nil, nil
}

func unmarshal(recentFilePath string, recentFileData *RecentFiles) {
	recentFileContent, err := os.ReadFile(recentFilePath)
	if err != nil {
		logger.Errorf("Could not open file %v", recentFilePath)
	}
	if string(recentFileContent[0:1]) != "{" {
		convertData(recentFileContent, recentFileData)
	} else {
		err = json.Unmarshal(recentFileContent, &recentFileData)
		if err != nil {
			logger.Errorf("Could not unmarshal content of %v", recentFilePath)
		}
	}
}

func convertData(content []byte, recentFileData *RecentFiles) {
	lines := strings.Split(string(content), "\n")
	for _, s := range lines {
		if s != "" {
			recentFileData.Terraform = append(recentFileData.Terraform, s)
		}
	}
}

func saveFile(data RecentFiles, path string) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Could not Marshal data to json. %q", err)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		logger.Errorf("Could not save file %v", path)
	}
}
