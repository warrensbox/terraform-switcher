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

func addRecent(requestedVersion string, installPath string) {
	addRecentGeneric(requestedVersion, installPath, "terraform")
}

func addRecentGeneric(requestedVersion string, installPath string, dist string) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	if CheckFileExist(recentFilePath) {
		unmarshal(recentFilePath, recentFileData)
	}
	var sliceToCheck []string
	if dist == "terraform" {
		sliceToCheck = recentFileData.Terraform
	} else if dist == "tofu" {
		sliceToCheck = recentFileData.Tofu
	}
	for _, v := range sliceToCheck {
		//TODO Check for valid version format
		if v == requestedVersion {
			// entry already exists. Nothing to do
			return
		}
	}
	if dist == "terraform" {
		recentFileData.Terraform = append(recentFileData.Terraform, requestedVersion)
	} else if dist == "tofu" {
		recentFileData.Tofu = append(recentFileData.Tofu, requestedVersion)
	}
	saveFile(recentFileData, recentFilePath)
}

func getRecentVersions(installPath string) ([]string, error) {
	return getRecentVersionsGeneric(installPath, "terraform")
}

func getRecentVersionsGeneric(installPath string, dist string) ([]string, error) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	unmarshal(recentFilePath, recentFileData)
	if dist == "terraform" {
		return recentFileData.Terraform, nil
	} else if dist == "tofu" {
		return recentFileData.Tofu, nil
	}
	return nil, nil
}

func unmarshal(recentFilePath string, recentFileData RecentFiles) {
	recentFileContent, err := os.ReadFile(recentFilePath)
	if err != nil {
		logger.Errorf("Could not open file %v", recentFilePath)
	}
	if !strings.HasPrefix(string(recentFileContent), "{") {
		convertData(recentFileContent, &recentFileData)
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
