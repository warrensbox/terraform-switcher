package main

/*
* Version 0.6.0
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
	"github.com/spf13/viper"
	lib "github.com/warrensbox/terraform-switcher/lib"
)

const (
	hashiURL   = "https://releases.hashicorp.com/terraform/"
	defaultBin = "/usr/local/bin/terraform"
)

var version = "0.7.0\n"

func main() {
	//The default binary path is /usr/local/bin/terraform
	custBinPath := getopt.StringLong("bin", 'b', defaultBin, "Custom binary path. For example: /Users/username/bin/terraform")
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	versionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()
	pathDir := lib.Path(*custBinPath)
	binDirExist := lib.CheckDirExist(pathDir)

	if !binDirExist {
		fmt.Printf("Binary path does not exist: %s\n", pathDir)
		fmt.Printf("Please create binary path: %s for terraform installation\n", pathDir)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	rcfile := dir + "/.tfswitchrc" //settings for .tfswitchrc file
	configfile := dir + "/.tfswitch.toml"

	if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
	} else if *helpFlag {
		usageMessage()
	} else {

		if _, err := os.Stat(configfile); err == nil {
			fmt.Println("Reading required terraform version from .tfswitch.toml")
			tfversion := ""
			binPath := *custBinPath
			configfileName := lib.GetFileName(".tfswitch.toml")
			viper.SetConfigType("toml")
			viper.SetConfigName(configfileName)
			viper.AddConfigPath(dir)

			errs := viper.ReadInConfig() // Find and read the config file
			if errs != nil {
				fmt.Println("Unable to read .tfswitch.toml provided") // Handle errors reading the config file
				fmt.Println(err)
				os.Exit(1)
			}

			checkDefault := strings.Compare(binPath, defaultBin)

			bin := viper.Get("bin")
			if checkDefault != -1 && bin != nil {
				binPath = bin.(string)
			}
			version := viper.Get("version")
			if version != nil {
				tfversion = version.(string)
			}

			if len(args) == 1 {
				fmt.Println("ARGs and conf file - delete")
				requestedVersion := args[0]
				listAll := true                                     //set list all true - all versions including beta and rc will be displayed
				tflist, _ := lib.GetTFList(hashiURL, listAll)       //get list of versions
				exist := lib.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

				if exist {
					tfversion = requestedVersion
				}
			}

			pathDir := lib.Path(binPath)
			binDirExist := lib.CheckDirExist(pathDir)

			if !binDirExist {
				fmt.Printf("Binary path does not exist: %s\n", pathDir)
				fmt.Printf("Please create binary path: %s for terraform installation\n", pathDir)
				os.Exit(1)
			} else if *listAllFlag {
				listAll := true //set list all true - all versions including beta and rc will be displayed
				installOption(listAll, &binPath)
			} else if tfversion == "" {
				// if *listAllFlag {
				// 	listAll := true //set list all true - all versions including beta and rc will be displayed
				// 	installOption(listAll, &binPath)
				// } else {
				listAll := false //set list all false - only official release will be displayed
				installOption(listAll, &binPath)
				//}
				os.Exit(0)
			} else {
				if lib.ValidVersionFormat(tfversion) { //check if version is correct
					lib.Install(tfversion, binPath)
				} else {
					fmt.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
					os.Exit(1)
				}
			}

		} else if _, err := os.Stat(rcfile); err == nil && len(args) == 0 { //if there is a .tfswitchrc file, and no commmand line arguments
			fmt.Println("NO ARGs but rc file - delete")
			fmt.Println("Reading required terraform version .tfswitchrc ")

			fileContents, err := ioutil.ReadFile(rcfile)
			if err != nil {
				fmt.Println("Failed to read .tfswitchrc file. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md")
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			tfversion := strings.TrimSuffix(string(fileContents), "\n")

			if lib.ValidVersionFormat(tfversion) { //check if version is correct
				lib.Install(string(tfversion), *custBinPath)
			} else {
				fmt.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				os.Exit(1)
			}
		} else if len(args) == 1 { //if tf version is provided in command line
			fmt.Println("ARGs NO conf file - delete")
			if lib.ValidVersionFormat(args[0]) {

				requestedVersion := args[0]
				listAll := true                                     //set list all true - all versions including beta and rc will be displayed
				tflist, _ := lib.GetTFList(hashiURL, listAll)       //get list of versions
				exist := lib.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

				if exist {
					lib.Install(requestedVersion, *custBinPath)
				} else {
					fmt.Println("The provided terraform version does not exist. Try `tfswitch -l` to see all available versions.")
				}

			} else {
				fmt.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				fmt.Println("Args must be a valid terraform version")
				usageMessage()
			}

		} else if *listAllFlag {
			listAll := true //set list all true - all versions including beta and rc will be displayed
			installOption(listAll, custBinPath)

		} else if len(args) == 0 { //if there are no commmand line arguments

			listAll := false //set list all false - only official release will be displayed
			installOption(listAll, custBinPath)

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
func installOption(listAll bool, custBinPath *string) {

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
	tfversion = strings.Trim(tfversion, " *recent") //trim versions with the string " *recent" appended

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	lib.Install(tfversion, *custBinPath)
}
