package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type RecentFiles struct {
	Terraform []string `json:"terraform"`
	OpenTofu  []string `json:"openTofu"`
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
	prependRecentVersionToList(requestedVersion, installPath, distribution, &recentFileData)
	saveRecentFile(recentFileData, recentFilePath)
}

func prependRecentVersionToList(version, installPath, distribution string, r *RecentFiles) {
	var sliceToCheck []string
	if distribution == distributionTerraform {
		sliceToCheck = r.Terraform
	} else if distribution == distributionOpenTofu {
		sliceToCheck = r.OpenTofu
	}
	for versionIndex, versionValue := range sliceToCheck {
		if versionValue == version {
			sliceToCheck = append(sliceToCheck[:versionIndex], sliceToCheck[versionIndex+1:]...)
		}
	}
	sliceToCheck = append([]string{version}, sliceToCheck...)

	//TODO delete files that are falling of the first three slice elements
	//if len(sliceToCheck) > 3 {
	//	deleteDownloadedBinaries(installPath, distribution, sliceToCheck[3:])
	//	sliceToCheck = sliceToCheck[0:2]
	//}

	if distribution == distributionTerraform {
		r.Terraform = sliceToCheck
	} else if distribution == distributionOpenTofu {
		r.OpenTofu = sliceToCheck
	}
}

func deleteDownloadedBinaries(installPath, distribution string, versions []string) {
	installLocation := GetInstallLocation(installPath)
	for _, versionToDelete := range versions {
		var fileToDelete string
		if distribution == distributionTerraform {
			fileToDelete = ConvertExecutableExt(TerraformPrefix + versionToDelete)
		}
		filePathToDelete := filepath.Join(installLocation, fileToDelete)
		logger.Debugf("Deleting obsolete binary %v", filePathToDelete)
		_ = os.Remove(filePathToDelete)
	}
}

func getRecentVersions(installPath string, dist string) ([]string, error) {
	installLocation := GetInstallLocation(installPath)
	recentFilePath := filepath.Join(installLocation, recentFile)
	var recentFileData RecentFiles
	unmarshalRecentFileData(recentFilePath, &recentFileData)
	var listOfRecentVersions []string
	if dist == distributionTerraform {
		listOfRecentVersions = recentFileData.Terraform
	} else if dist == distributionOpenTofu {
		listOfRecentVersions = recentFileData.OpenTofu
	}
	for index, versionString := range listOfRecentVersions {
		listOfRecentVersions[index] = versionString + " *recent"
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
