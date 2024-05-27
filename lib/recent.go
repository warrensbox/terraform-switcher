package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RecentFile struct {
	Terraform []string
	OpenTofu  []string
}

func appendRecentVersionToList(versions []string, requestedVersion string) []string {
	// Check for requestedVersion in versions list
	for versionIndex, versionVal := range versions {
		if versionVal == requestedVersion {
			versions = append(versions[:versionIndex], versions[versionIndex+1:]...)
		}
	}

	// Add new version to start of slice
	versions = append([]string{requestedVersion}, versions...)
	if len(versions) > 3 {
		versions = versions[0:2]
	}

	return versions
}

// addRecent : add to recent file
func addRecentVersion(product Product, requestedVersion string, installPath string) {
	installLocation = GetInstallLocation(installPath) //get installation location -  this is where we will put our terraform binary file
	recentFilePath := filepath.Join(installLocation, recentFile)

	// Obtain pre-existing latest version
	recentData := getRecentFileData(installPath)

	if product.GetId() == "terraform" {
		recentData.Terraform = appendRecentVersionToList(recentData.Terraform, requestedVersion)
	} else if product.GetId() == "opentofu" {
		recentData.OpenTofu = appendRecentVersionToList(recentData.OpenTofu, requestedVersion)
	}

	// Write new versions back to recent files
	recentVersionFh, err := os.OpenFile(recentFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("Error to open %q file for writing: %v", recentFilePath, err)
		return
	}
	defer recentVersionFh.Close()

	// Marhsall data and write to file
	jsonData, err := json.Marshal(recentData)
	if err != nil {
		logger.Warnf("Error during marshalling recent versions data from %s file: %v. Ignoring", recentFilePath, err)
	}

	_, err = recentVersionFh.Write(jsonData)
	if err != nil {
		logger.Warnf("Error writing recent versions file (%q): %v. Ignoring", recentFilePath, err)
	}
}

func convertLegacyRecentData(content []byte, recentFileContent *RecentFile) {
	lines := strings.Split(string(content), "\n")
	for _, s := range lines {
		if s != "" {
			recentFileContent.Terraform = append(recentFileContent.Terraform, s)
		}
	}
}

func getRecentFileData(installPath string) RecentFile {
	installLocation = GetInstallLocation(installPath) //get installation location -  this is where we will put our terraform binary file
	recentFilePath := filepath.Join(installLocation, recentFile)
	var outputRecent RecentFile

	fileExist := CheckFileExist(recentFilePath)
	if fileExist {
		content, err := ioutil.ReadFile(recentFilePath)
		if err != nil {
			logger.Warnf("Error opening recent versions file (%q): %v. Ignoring", recentFilePath, err)
			return outputRecent
		}

		if !strings.HasPrefix(string(content), "{") {
			convertLegacyRecentData(content, &outputRecent)
		}

		err = json.Unmarshal(content, &outputRecent)
		if err != nil {
			logger.Warnf("Error during unmarshalling recent versions data from %s file: %v. Ignoring", recentFilePath, err)
		}
	}
	return outputRecent
}

// getRecentVersions : get recent version from file
func getRecentVersions(product Product, installPath string) []string {
	outputRecent := getRecentFileData(installPath)

	if product.GetId() == "terraform" {
		return outputRecent.Terraform
	} else if product.GetId() == "opentofu" {
		return outputRecent.OpenTofu
	}

	// Catch-all for unmatched product
	return []string{}
}
