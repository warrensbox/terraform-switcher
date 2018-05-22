package main

/*
* Version 0.0.1
* Compatible with Mac OS X ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terraform directory if does not exist
* 2- Download zip file from url to /usr/local/terraform
* 3- Unzip the file to /usr/local/terraform
* 4- Rename the file from `terraform` to `terraform_version`
* 5- Remove the downloaded zip file
* 6- Read the existing symlink for terraform (Check if it's a homebrew symlink)
* 7- Remove that symlink (Check if it's a homebrew symlink)
* 8- Create new symlink to binary  `terraform_version`
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	lib "github.com/warren-veerasingam/terraform-switcher/lib"
)

type tfVersionList struct {
	tflist []string
}

const (
	hashiURL       = "https://releases.hashicorp.com/terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	binLocation    = "/usr/local/bin/terraform"
	installPath    = "/.terraform.versions/"
	macOS          = "_darwin_amd64.zip"
)

func main() {

	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	fmt.Printf("Current user: %v \n", usr.HomeDir)

	/* set installation location */
	installLocation := usr.HomeDir + installPath

	/* set default binary path for terraform */
	installedBinPath := binLocation

	/* find terraform binary location if terraform is already installed*/
	cmd := lib.NewCommand("terraform")
	next := cmd.Find()
	//existed := false

	/* overrride installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		fmt.Printf("Found installation path: %v \n", path)
		installedBinPath = path
	}

	fmt.Printf("Terraform binary path: %v", installedBinPath)

	/* Create local installation directory if it does not exist */
	lib.CreateDirIfNotExist(installLocation)

	/* Get list of terraform versions from hashicorp releases */
	resp, errURL := http.Get(hashiURL)
	if errURL != nil {
		log.Printf("Error getting url: %v", errURL)
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Printf("Error reading body: %v", errBody)
		return
	}

	bodyString := string(body)
	result := strings.Split(bodyString, "\n")

	var tfVersionList tfVersionList

	for i := range result {
		//getting versions from body; should return match /X.X.X/
		r, _ := regexp.Compile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)

		if r.MatchString(result[i]) {
			str := r.FindString(result[i])
			trimstr := strings.Trim(str, "/") //remove "/" from /X.X.X/
			tfVersionList.tflist = append(tfVersionList.tflist, trimstr)
		}
	}

	/* prompt user to select version of terraform */
	prompt := promptui.Select{
		Label: "Select Terraform version",
		Items: tfVersionList.tflist,
	}

	_, version, errPrompt := prompt.Run()

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	fmt.Printf("Terraform version %q selected\n", version)

	/* check if selected version already downloaded */
	fileExist := lib.CheckFileExist(installLocation + installVersion + version)

	/* if selected version already exist, */
	if fileExist {
		/* remove current symlink and set new symlink to desired version */
		lib.RemoveSymlink(installedBinPath)

		/* set symlink to desired version */
		lib.CreateSymlink(installLocation+installVersion+version, installedBinPath)
		fmt.Printf("Swicthed terraform to version %q \n", version)
		os.Exit(0)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := hashiURL + version + "/" + installVersion + version + macOS
	zipFile, _ := lib.DownloadFromURL(installLocation, url)

	fmt.Printf("Downloaded zipFile: %v \n", zipFile)

	/* unzip the downloaded zipfile */
	files, errUnzip := lib.Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	fmt.Println("Unzipped: " + strings.Join(files, "\n"))

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	lib.RenameFile(installLocation+installFile, installLocation+installVersion+version)

	/* remove zipped file to clear clutter */
	lib.RemoveFiles(installLocation + installVersion + version + macOS)

	/* remove current symlink and set new symlink to desired version  */
	lib.RemoveSymlink(installedBinPath)

	/* set symlink to desired version */
	lib.CreateSymlink(installLocation+installVersion+version, installedBinPath)
	fmt.Printf("Swicthed terraform to version %q \n", version)
	os.Exit(0)
}
