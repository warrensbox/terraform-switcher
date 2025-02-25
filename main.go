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
	"os"

	lib "github.com/warrensbox/terraform-switcher/lib"
	"github.com/warrensbox/terraform-switcher/lib/param_parsing"
)

var (
	parameters = param_parsing.GetParameters()
	logger     = lib.InitLogger(parameters.LogLevel)
	version    string
)

func main() {
	var err error = nil
	switch {
	case parameters.VersionFlag:
		if version != "" {
			fmt.Printf("Version: %s\n", version)
		} else {
			fmt.Println("Version not defined during build")
		}
		os.Exit(0)
	case parameters.HelpFlag:
		lib.UsageMessage()
		os.Exit(0)
	case parameters.ListAllFlag:
		/* show all terraform version including betas and RCs*/
		err = lib.InstallProductOption(parameters.ProductEntity, true, parameters.DryRun, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.LatestPre != "":
		/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		err = lib.InstallLatestProductImplicitVersion(parameters.ProductEntity, parameters.DryRun, parameters.LatestPre, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch, true)
	case parameters.ShowLatestPre != "":
		/* show latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestPre, parameters.MirrorURL, true)
	case parameters.LatestStable != "":
		/* latest implicit version. Ex: tfswitch --latest-stable 0.13 downloads 0.13.5 (latest) */
		err = lib.InstallLatestProductImplicitVersion(parameters.ProductEntity, parameters.DryRun, parameters.LatestStable, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch, false)
	case parameters.ShowLatestStable != "":
		/* show latest implicit stable version. Ex: tfswitch --show-latest-stable 0.13 downloads 0.13.5 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestStable, parameters.MirrorURL, false)
	case parameters.LatestFlag:
		/* latest stable version */
		err = lib.InstallLatestProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.ShowLatestFlag:
		/* show latest stable version */
		lib.ShowLatestVersion(parameters.MirrorURL)
	case parameters.Version != "":
		err = lib.InstallProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.Version, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.DefaultVersion != "":
		/* if default version is provided - Pick this instead of going for prompt */
		err = lib.InstallProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.DefaultVersion, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	default:
		// Set list all false - only official release will be displayed
		err = lib.InstallProductOption(parameters.ProductEntity, false, parameters.DryRun, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	}
	if err != nil {
		logger.Fatal(err)
	}
}
