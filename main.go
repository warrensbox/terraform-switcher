package main

/*
* Version 0.12.0
* Compatible with Mac OS X AND other LINUX OS ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terraform directory if it does not exist
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
	"github.com/warrensbox/terraform-switcher/lib/param_parsing"
	"os"
	"path/filepath"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/manifoldco/promptui"
	lib "github.com/warrensbox/terraform-switcher/lib"
)

var logger = lib.InitLogger()
var version string

func main() {

	parameters := param_parsing.GetParameters()

	switch {
	case parameters.VersionFlag:
		if version != "" {
			fmt.Printf("Version: %s\n", version)
		} else {
			fmt.Println("Version not defined during build.")
		}
		os.Exit(0)
	case parameters.HelpFlag:
		param_parsing.UsageMessage()
		os.Exit(0)
	case parameters.ListAllFlag:
		/* show all terraform version including betas and RCs*/
		installOption(true, parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.LatestPre != "":
		/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		installLatestImplicitVersion(parameters.LatestPre, parameters.CustomBinaryPath, parameters.MirrorURL, true)
	case parameters.ShowLatestPre != "":
		/* show latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		showLatestImplicitVersion(parameters.ShowLatestPre, parameters.MirrorURL, true)
	case parameters.LatestStable != "":
		/* latest implicit version. Ex: tfswitch --latest-stable 0.13 downloads 0.13.5 (latest) */
		installLatestImplicitVersion(parameters.LatestStable, parameters.CustomBinaryPath, parameters.MirrorURL, false)
	case parameters.ShowLatestStable != "":
		/* show latest implicit stable version. Ex: tfswitch --show-latest-stable 0.13 downloads 0.13.5 (latest) */
		showLatestImplicitVersion(parameters.ShowLatestStable, parameters.MirrorURL, false)
	case parameters.LatestFlag:
		/* latest stable version */
		installLatestVersion(parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.ShowLatestFlag:
		/* show latest stable version */
		showLatestVersion(parameters.MirrorURL)
	case parameters.Version != "":
		installVersion(parameters.Version, parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.DefaultVersion != "":
		/* if default version is provided - Pick this instead of going for prompt */
		installVersion(parameters.DefaultVersion, parameters.CustomBinaryPath, parameters.MirrorURL)
	default:
		// Set list all false - only official release will be displayed
		installOption(false, parameters.CustomBinaryPath, parameters.MirrorURL)
	}
}

// install latest stable tf version
func installLatestVersion(customBinaryPath, mirrorURL string) {
	tfversion, _ := lib.GetTFLatest(mirrorURL)
	lib.Install(tfversion, customBinaryPath, mirrorURL)
}

// show install latest stable tf version
func showLatestVersion(mirrorURL string) {
	tfversion, _ := lib.GetTFLatest(mirrorURL)
	logger.Infof("%s", tfversion)
}

// install latest - argument (version) must be provided
func installLatestImplicitVersion(requestedVersion, customBinaryPath, mirrorURL string, preRelease bool) {
	_, err := semver.NewConstraint(requestedVersion)
	if err != nil {
		logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	}
	//if lib.ValidMinorVersionFormat(requestedVersion) {
	tfversion, err := lib.GetTFLatestImplicit(mirrorURL, preRelease, requestedVersion)
	if err == nil && tfversion != "" {
		lib.Install(tfversion, customBinaryPath, mirrorURL)
	}
	logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	lib.PrintInvalidMinorTFVersion()
}

// show latest - argument (version) must be provided
func showLatestImplicitVersion(requestedVersion, mirrorURL string, preRelease bool) {
	if lib.ValidMinorVersionFormat(requestedVersion) {
		tfversion, _ := lib.GetTFLatestImplicit(mirrorURL, preRelease, requestedVersion)
		if len(tfversion) > 0 {
			logger.Infof("%s", tfversion)
		} else {
			logger.Fatal("The provided terraform version does not exist.\n Try `tfswitch -l` to see all available versions")
			os.Exit(1)
		}
	} else {
		lib.PrintInvalidMinorTFVersion()
	}
}

// install with provided version as argument
func installVersion(arg, customBinaryPath, mirrorURL string) {
	if lib.ValidVersionFormat(arg) {
		requestedVersion := arg

		//check to see if the requested version has been downloaded before
		installLocation := lib.GetInstallLocation()
		installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, lib.VersionPrefix+requestedVersion))
		recentDownloadFile := lib.CheckFileExist(installFileVersionPath)
		if recentDownloadFile {
			lib.ChangeSymlink(installFileVersionPath, customBinaryPath)
			logger.Infof("Switched terraform to version %q", requestedVersion)
			lib.AddRecent(requestedVersion) //add to recent file for faster lookup
			os.Exit(0)
		}

		// If the requested version had not been downloaded before
		// Set list all true - all versions including beta and rc will be displayed
		tflist, _ := lib.GetTFList(mirrorURL, true)         // Get list of versions
		exist := lib.VersionExist(requestedVersion, tflist) // Check if version exists before downloading it

		if exist {
			lib.Install(requestedVersion, customBinaryPath, mirrorURL)
		} else {
			logger.Fatal("The provided terraform version does not exist.\n Try `tfswitch -l` to see all available versions")
			os.Exit(1)
		}

	} else {
		lib.PrintInvalidTFVersion()
		logger.Error("Args must be a valid terraform version")
		param_parsing.UsageMessage()
		os.Exit(1)
	}
}

/* installOption : displays & installs tf version */
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func installOption(listAll bool, customBinaryPath, mirrorURL string) {
	tflist, _ := lib.GetTFList(mirrorURL, listAll) // Get list of versions
	recentVersions, _ := lib.GetRecentVersions()   // Get recent versions from RECENT file
	tflist = append(recentVersions, tflist...)     // Append recent versions to the top of the list
	tflist = lib.RemoveDuplicateVersions(tflist)   // Remove duplicate version

	if len(tflist) == 0 {
		logger.Fatalf("Terraform version list is empty: %s", mirrorURL)
		os.Exit(1)
	}

	/* prompt user to select version of terraform */
	prompt := promptui.Select{
		Label: "Select Terraform version",
		Items: tflist,
	}

	_, tfversion, errPrompt := prompt.Run()
	tfversion = strings.Trim(tfversion, " *recent") //trim versions with the string " *recent" appended

	if errPrompt != nil {
		logger.Fatalf("Prompt failed %v", errPrompt)
		os.Exit(1)
	}

	lib.Install(tfversion, customBinaryPath, mirrorURL)
	os.Exit(0)
}
