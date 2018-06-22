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
	"log"
	"os"

	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/terraform-switcher/lib"
)

const (
	hashiURL       = "https://releases.hashicorp.com/terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	binLocation    = "/usr/local/bin/terraform"
	installPath    = "/.terraform.versions/"
	macOS          = "_darwin_amd64.zip"
	linux          = "_darwin_amd64.zip"
)

var version = "0.0.1\n"

// var (
// 	installLocation  = "/tmp"
// 	installedBinPath = "/tmp"
// )

func main() {
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tfswitch", "something")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message", "something")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	if *versionFlag {
		fmt.Println(version)
	} else if *helpFlag {
		UsageMessage()
	} else {

		if len(args) == 1 {
			semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(args[0]) {
				requestedVersion := args[0]

				//check if version exist before downloading it
				tflist, _ := lib.GetTFList(hashiURL)
				exist := lib.VersionExist(requestedVersion, tflist)

				if exist {
					lib.Install(requestedVersion)
				} else {
					fmt.Println("Not a valid terraform version")
				}

			} else {
				fmt.Println("Not a valid terraform version")
				fmt.Println("Args must be a valid terraform version")
				UsageMessage()
			}

		} else if len(args) == 0 {

			// os.Exit(-1)
			tflist, _ := lib.GetTFList(hashiURL)

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

			fmt.Printf("Terraform version %q selected\n", tfversion)

			lib.Install(tfversion)
		} else {
			UsageMessage()
		}
	}
}

func UsageMessage() {
	fmt.Println("\n\nInvalid Selection")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}
