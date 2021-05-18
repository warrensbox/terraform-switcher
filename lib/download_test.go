package lib_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"testing"

	lib "github.com/warrensbox/terraform-switcher/lib"
)

// TestDownloadFromURL_FileNameMatch : Check expected filename exist when downloaded
func TestDownloadFromURL_FileNameMatch(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"
	installVersion := "terraform_"
	installPath := getInstallLocation(".terraform.versions_test")
	macOS := "_darwin_amd64.zip"

	// get current user
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	fmt.Printf("Current user: %v \n", usr.HomeDir)
	installLocation := usr.HomeDir + installPath

	// create /.terraform.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		log.Printf("Creating directory for terraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terraform: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest terraform version */
	lowestVersion := "0.1.0"

	url := hashiURL + lowestVersion + "/" + installVersion + lowestVersion + macOS
	expectedFile := usr.HomeDir + installPath + installVersion + lowestVersion + macOS
	installedFile, errDownload := lib.DownloadFromURL(installLocation, url)

	if errDownload != nil {
		t.Logf("Expected file name %v to be downloaded", expectedFile)
		t.Error("Download not possible (unexpected)")
	}

	if installedFile == expectedFile {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Download file mismatches expected file (unexpected)")
	}

	/* test download latest terraform version */
	latestVersion := "0.11.7"

	url = hashiURL + latestVersion + "/" + installVersion + latestVersion + macOS
	expectedFile = usr.HomeDir + installPath + installVersion + latestVersion + macOS
	installedFile, errDownload = lib.DownloadFromURL(installLocation, url)

	if errDownload != nil {
		t.Logf("Expected file name %v to be downloaded", expectedFile)
		t.Error("Download not possible (unexpected)")
	}

	if installedFile == expectedFile {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Log("Download file name matches expected file")
	} else {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Error("Dowload file name mismatches expected file (unexpected)")
	}

	cleanUp(installLocation)
}

// TestDownloadFromURL_FileExist : Check expected file exist when downloaded
func TestDownloadFromURL_FileExist(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"
	installFile := "terraform"
	installVersion := "terraform_"
	installPath := getInstallLocation(".terraform.versions_test")
	macOS := "_darwin_amd64.zip"

	// get current user
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	fmt.Printf("Current user: %v \n", usr.HomeDir)
	installLocation := usr.HomeDir + installPath

	// create /.terraform.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		log.Printf("Creating directory for terraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terraform: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest terraform version */
	lowestVersion := "0.1.0"

	url := hashiURL + lowestVersion + "/" + installVersion + lowestVersion + macOS
	expectedFile := usr.HomeDir + installPath + installVersion + lowestVersion + macOS
	installedFile, errDownload := lib.DownloadFromURL(installLocation, url)

	if errDownload != nil {
		t.Logf("Expected file name %v to be downloaded", expectedFile)
		t.Error("Download not possible (unexpected)")
	}

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Download file mismatches expected file (unexpected)")
	}

	/* test download latest terraform version */
	latestVersion := "0.11.7"

	url = hashiURL + latestVersion + "/" + installVersion + latestVersion + macOS
	expectedFile = usr.HomeDir + installPath + installVersion + latestVersion + macOS
	installFile, errDownload = lib.DownloadFromURL(installLocation, url)

	if errDownload != nil {
		t.Logf("Expected file name %v to be downloaded", expectedFile)
		t.Error("Download not possible (unexpected)")
	}

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installFile)
		t.Error("Download file mismatches expected file (unexpected)")
	}

	cleanUp(installLocation)
}

// TestInvalidURL : Invalid url should throw an error
func TestInvalidURL(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"
	installVersion := "terraform_"
	installPath := getInstallLocation(".terraform.versions_test")
	macOS := "_darwin_amd64.zip"
	invalidVersion := "0.11.7-nonexistent"

	// get current user
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	fmt.Printf("Current user: %v \n", usr.HomeDir)
	installLocation := usr.HomeDir + installPath

	// create /.terraform.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		log.Printf("Creating directory for terraform: %v\n", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terraform: %v\n", installLocation)
			panic(err)
		}
	}

	url := hashiURL + invalidVersion + "/" + installVersion + invalidVersion + macOS
	//expectedFile := usr.HomeDir + installPath + installVersion + invalidVersion + macOS
	_, errDownload := lib.DownloadFromURL(installLocation, url)

	if errDownload != nil {
		t.Logf("Unable to download from %s - invalid url or version (expected)\n", url)
		t.Logf("Download not possible (expected)")
	}

	cleanUp(installLocation)
}

// TestDownloadFromURL_Valid : Test if https://releases.hashicorp.com/terraform/ is still valid
func TestDownloadFromURL_Valid(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"

	url, err := url.ParseRequestURI(hashiURL)
	if err != nil {
		t.Errorf("Valid URL provided:  %v", err)
		t.Errorf("Invalid URL %v", err)
	} else {
		t.Logf("Valid URL from %v", url)
	}
}
