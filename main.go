package main

/*
* Version 0.5.0
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
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/terraform-switcher/lib"
)

const (
	hashiURL = "https://releases.hashicorp.com/terraform/"
)

var version = "0.5.0\n"

func main() {
	listAllFlag := getopt.BoolLong("list-all", 'l', "list all versions of terraform - including beta and rc")
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tfswitch")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}
	rcfile := dir + "/.tfswitchrc"

	if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
	} else if *helpFlag {
		usageMessage()
	} else if *listAllFlag {
		listAll := true //set list all true - all versions including beta and rc will be displayed
		installOption(listAll)
	} else {

		if len(args) == 1 { //if tf version is provided in command line

			if lib.ValidVersionFormat(args[0]) {
				requestedVersion := args[0]
				listAll := true //set list all true - all versions including beta and rc will be displayed

				tflist, _ := lib.GetTFList(hashiURL, listAll)       //get list of versions
				exist := lib.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

				if exist {
					lib.AddRecent(requestedVersion) //add to recent file for faster lookup
					lib.Install(requestedVersion)
				} else {
					log.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				}

			} else {
				log.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				fmt.Println("Args must be a valid terraform version")
				usageMessage()
			}

		} else if _, err := os.Stat(rcfile); err == nil && len(args) == 0 { //if there is a .tfswitchrc file, and no commmand line arguments
			fmt.Println("Reading required terraform version ...")

			fileContents, err := ioutil.ReadFile(rcfile)
			if err != nil {
				log.Println("Failed to read .tfswitchrc file. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md")
				log.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			tfversion := strings.TrimSuffix(string(fileContents), "\n")

			if lib.ValidVersionFormat(tfversion) { //check if version is correct
				lib.AddRecent(string(tfversion)) //add to RECENT file for faster lookup
				lib.Install(string(tfversion))
			} else {
				log.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				os.Exit(1)
			}
		} else if len(args) == 0 { //if there are no commmand line arguments
			listAll := false
			tflist, _ := lib.GetTFList(hashiURL, listAll)
			recentVersions, _ := lib.GetRecentVersions() //get recent versions from RECENT file
			tflist = append(recentVersions, tflist...)   //append recent versions to the top of the list
			tflist = lib.RemoveDuplicateVersions(tflist) //remove duplicate version

			/* prompt user to select version of terraform */
			prompt := promptui.Select{
				Label: "Select Terraform version",
				Items: tflist,
			}

			_, tfversion, errPrompt := prompt.Run()

			if errPrompt != nil {
				log.Printf("Prompt failed %v\n", errPrompt)
				os.Exit(1)
			}

			lib.AddRecent(tfversion) //add to recent file for faster lookup
			lib.Install(tfversion)

		} else {
			usageMessage()
		}
	}
}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}

/* installOption : displays & installs tf version */
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func installOption(listAll bool) {

	tflist, _ := lib.GetTFList(hashiURL, listAll) //get list of versions
	recentVersions, _ := lib.GetRecentVersions()  //get recent versions from RECENT file
	tflist = append(recentVersions, tflist...)    //append recent versions to the top of the list
	tflist = lib.RemoveDuplicateVersions(tflist)  //remove duplicate version

	/* prompt user to select version of terraform */
	prompt := promptui.Select{
		Label: "Select Terraform version",
		Items: tflist,
	}

	_, tfversion, errPrompt := prompt.Run()

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	lib.AddRecent(tfversion) //add to recent file for faster lookup
	lib.Install(tfversion)
}
