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

	"github.com/manifoldco/promptui"
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

	args := os.Args

	if len(os.Args) > 1 {
		switch os := args[1]; os {
		case "--version":
			fmt.Println(version)
		case "version":
			fmt.Println(version)
		case "-v":
			fmt.Println(version)
		}
	} else {

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

	}
}
