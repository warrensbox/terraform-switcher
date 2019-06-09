package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
)

const (
	hashiURL       = "https://releases.hashicorp.com/terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	//binLocation    = "/usr/local/bin/terraform"
	installPath = "/.terraform.versions/"
	recentFile  = "RECENT"
)

var (
	installLocation = "/tmp"
	//installedBinPath = "/tmp"
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
	installedBinPath := "/usr/local/bin/terraform"

	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand("terraform")
	next := cmd.Find()

	/* overrride installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		installedBinPath = path
	}

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(installedBinPath)

	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

}

//Install : Install the provided version in the argument
func Install(tfversion string, binPath string) {

	if !ValidVersionFormat(tfversion) {
		fmt.Printf("The provided terraform version format does not exist - %s. Try `tfswitch -l` to see all available versions.\n", tfversion)
		os.Exit(1)
	}

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	/* check if selected version already downloaded */
	fileExist := CheckFileExist(installLocation + installVersion + tfversion)

	/* if selected version already exist, */
	if fileExist {

		/* remove current symlink if exist*/
		symlinkExist := CheckSymlink(binPath)

		if symlinkExist {
			RemoveSymlink(binPath)
		}

		/* set symlink to desired version */
		CreateSymlink(installLocation+installVersion+tfversion, binPath)
		fmt.Printf("Switched terraform to version %q \n", tfversion)
		AddRecent(tfversion) //add to recent file for faster lookup
		os.Exit(0)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := hashiURL + tfversion + "/" + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip"
	zipFile, errDownload := DownloadFromURL(installLocation, url)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		fmt.Println(errDownload)
		os.Exit(1)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	RenameFile(installLocation+installFile, installLocation+installVersion+tfversion)

	/* remove zipped file to clear clutter */
	RemoveFiles(installLocation + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip")

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)

	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+tfversion, binPath)
	fmt.Printf("Switched terraform to version %q \n", tfversion)
	AddRecent(tfversion) //add to recent file for faster lookup
	os.Exit(0)
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string) {

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !ValidVersionFormat(line) {
				fmt.Println("File dirty. Recreating cache file.")
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

		lines, errRead := ReadLines(installLocation + recentFile)
		outputRecent := []string{}

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			/* 	checks if versions in the recent file are valid.
			If any version is invalid, it will be consider dirty
			and the recent file will be removed
			*/
			if !ValidVersionFormat(line) {
				RemoveFiles(installLocation + recentFile)
				return nil, errRead
			}

			/* 	output can be confusing since it displays the 3 most recent used terraform version
			append the string *recent to the output to make it more user friendly
			*/
			outputRecent = append(outputRecent, fmt.Sprintf("%s *recent", line))
		}

		return outputRecent, nil
	}

	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {
	WriteLines([]string{requestedVersion}, installLocation+recentFile)
}
