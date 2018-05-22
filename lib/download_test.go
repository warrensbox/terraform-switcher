package lib_test

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	lib "github.com/warren-veerasingam/terraform-switcher/lib"
)

// const (
// 	hashiURL       = "https://releases.hashicorp.com/terraform/"
// 	installFile    = "terraform"
// 	installVersion = "terraform_"
// 	binLocation    = "/usr/local/bin/terraform"
// 	installPath    = "/.terraform.versions_test/"
// 	macOS          = "_darwin_amd64.zip"
// )

func TestDownloadURL(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"
	installVersion := "terraform_"
	installPath := "/.terraform.versions_test/"
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
		log.Printf("Creating directory for teraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for teraform: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest terraform version */
	lowestVersion := "0.0.1"

	url := hashiURL + lowestVersion + "/" + installVersion + lowestVersion + macOS
	expectedFile := usr.HomeDir + installPath + installVersion + lowestVersion + macOS
	installedFile, _ := lib.DownloadFromURL(installLocation, url)

	if installedFile == expectedFile {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Download file mismatches expected file")
	}

	/* test download latest terraform version */
	latestVersion := "0.11.7"

	url = hashiURL + latestVersion + "/" + installVersion + latestVersion + macOS
	expectedFile = usr.HomeDir + installPath + installVersion + latestVersion + macOS
	installedFile, _ = lib.DownloadFromURL(installLocation, url)

	if installedFile == expectedFile {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Log("Download file name matches expected file")
	} else {
		t.Logf("Expected file name %v", expectedFile)
		t.Logf("Downloaded file name %v", installedFile)
		t.Error("Downoad file name mismatches expected file")
	}

	cleanUp(installLocation)
}

func TestDownloadedFileExist(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"
	installFile := "terraform"
	installVersion := "terraform_"
	installPath := "/.terraform.versions_test/"
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
		log.Printf("Creating directory for teraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for teraform: %v", installLocation)
			panic(err)
		}
	}

	/* test download lowest terraform version */
	lowestVersion := "0.0.1"

	url := hashiURL + lowestVersion + "/" + installVersion + lowestVersion + macOS
	expectedFile := usr.HomeDir + installPath + installVersion + lowestVersion + macOS
	installedFile, _ := lib.DownloadFromURL(installLocation, url)

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Downoad file mismatches expected file")
	}

	/* test download latest terraform version */
	latestVersion := "0.11.7"

	url = hashiURL + latestVersion + "/" + installVersion + latestVersion + macOS
	expectedFile = usr.HomeDir + installPath + installVersion + latestVersion + macOS
	installFile, _ = lib.DownloadFromURL(installLocation, url)

	if checkFileExist(expectedFile) {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installFile)
		t.Error("Downoad file mismatches expected file")
	}

	cleanUp(installLocation)
}

func TestURLValid(t *testing.T) {

	hashiURL := "https://releases.hashicorp.com/terraform/"

	url, err := url.ParseRequestURI(hashiURL)
	if err != nil {
		t.Errorf("Valid URL provided:  %v", err)
		t.Errorf("Invalid URL %v", err)
	} else {
		t.Logf("Valid URL from %v", url)
	}
}

func cleanUp(path string) {
	removeContents(path)
	removeFiles(path)
}

func removeFiles(src string) {
	files, err := filepath.Glob(src)
	if err != nil {

		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
