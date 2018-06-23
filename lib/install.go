package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

const (
	hashiURL       = "https://releases.hashicorp.com/terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	binLocation    = "/usr/local/bin/terraform"
	installPath    = "/.terraform.versions/"
	macOS          = "_darwin_amd64.zip"
	linux          = "_darwin_amd64.zip"
	recentFile     = "RECENT"
)

var (
	installLocation  = "/tmp"
	installedBinPath = "/tmp"
)

func init() {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	/* set installation location */
	installLocation = usr.HomeDir + installPath

	/* set default binary path for terraform */
	installedBinPath = binLocation

	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand("terraform")
	next := cmd.Find()
	//existed := false

	/* overrride installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		fmt.Printf("Found installation path: %v \n", path)
		installedBinPath = path
	}

	fmt.Printf("Terraform binary path: %v \n", installedBinPath)

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

}

//Install : Install the provided version in the argument
func Install(tfversion string) {

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	/* check if selected version already downloaded */
	fileExist := CheckFileExist(installLocation + installVersion + tfversion)

	/* if selected version already exist, */
	if fileExist {
		/* remove current symlink if exist*/
		exist := CheckFileExist(installedBinPath)

		if !exist {
			fmt.Println("Symlink does not exist")
		} else {
			RemoveSymlink(installedBinPath)
		}

		/* set symlink to desired version */
		CreateSymlink(installLocation+installVersion+tfversion, installedBinPath)
		fmt.Printf("Swicthed terraform to version %q \n", tfversion)
		os.Exit(0)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := hashiURL + tfversion + "/" + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip"
	zipFile, _ := DownloadFromURL(installLocation, url)

	fmt.Printf("Downloaded zipFile: %v \n", zipFile)

	/* unzip the downloaded zipfile */
	files, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	fmt.Println("Unzipped: " + strings.Join(files, "\n"))

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	RenameFile(installLocation+installFile, installLocation+installVersion+tfversion)

	/* remove zipped file to clear clutter */
	RemoveFiles(installLocation + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip")

	/* remove current symlink if exist*/
	exist := CheckFileExist(installedBinPath)

	if !exist {
		fmt.Println("Symlink does not exist")
	} else {
		fmt.Println("Symlink exist")
		RemoveSymlink(installedBinPath)
	}

	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+tfversion, installedBinPath)
	fmt.Printf("Swicthed terraform to version %q \n", tfversion)
	os.Exit(0)
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string) {

	semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				fmt.Println("file corrupted")
				RemoveFiles(installLocation + recentFile)
				CreateRecentFile(requestedVersion)
				return
			}
		}

		versionExist := VersionExist(requestedVersion, lines)

		if !versionExist {
			if len(lines) >= 3 {
				_, lines = lines[len(lines)-1], lines[:len(lines)-1]

				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			} else {
				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			}
		}

	} else {
		CreateRecentFile(requestedVersion)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions() ([]string, error) {

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				RemoveFiles(installLocation + recentFile)
				return nil, errRead
			}
		}

		return lines, nil
	}

	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {
	WriteLines([]string{requestedVersion}, installLocation+recentFile)
}
