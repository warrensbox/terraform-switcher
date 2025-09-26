package main

/*
* https://tfswitch.warrensbox.com/
* A command line tool to switch between different versions of terraform
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
	var err error
	switch {
	case parameters.VersionFlag:
		fmt.Printf("Version: ")
		if version != "" {
			fmt.Println(version)
		} else {
			fmt.Println("not defined during build")
		}
		os.Exit(0)
	case parameters.HelpFlag:
		lib.UsageMessage()
		os.Exit(0)
	case parameters.ListAllFlag:
		/* show all terraform version including betas and RCs*/
		err = lib.InstallProductOption(parameters.ProductEntity, true, parameters.DryRun, parameters.ShowRequiredFlag, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.LatestPre != "":
		/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		err = lib.InstallLatestProductImplicitVersion(parameters.ProductEntity, parameters.DryRun, parameters.ShowRequiredFlag, parameters.LatestPre, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch, true)
	case parameters.ShowLatestPre != "":
		/* show latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestPre, parameters.MirrorURL, true)
	case parameters.LatestStable != "":
		/* latest implicit version. Ex: tfswitch --latest-stable 0.13 downloads 0.13.5 (latest) */
		err = lib.InstallLatestProductImplicitVersion(parameters.ProductEntity, parameters.DryRun, parameters.ShowRequiredFlag, parameters.LatestStable, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch, false)
	case parameters.ShowLatestStable != "":
		/* show latest implicit stable version. Ex: tfswitch --show-latest-stable 0.13 downloads 0.13.5 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestStable, parameters.MirrorURL, false)
	case parameters.LatestFlag:
		/* latest stable version */
		err = lib.InstallLatestProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.ShowRequiredFlag, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.ShowLatestFlag:
		/* show latest stable version */
		lib.ShowLatestVersion(parameters.MirrorURL)
	case parameters.Version != "":
		err = lib.InstallProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.ShowRequiredFlag, parameters.Version, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	case parameters.DefaultVersion != "":
		/* if default version is provided - Pick this instead of going for prompt */
		err = lib.InstallProductVersion(parameters.ProductEntity, parameters.DryRun, parameters.ShowRequiredFlag, parameters.DefaultVersion, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	default:
		// Set list all false - only official release will be displayed
		err = lib.InstallProductOption(parameters.ProductEntity, false, parameters.DryRun, parameters.ShowRequiredFlag, parameters.CustomBinaryPath, parameters.InstallPath, parameters.MirrorURL, parameters.Arch)
	}
	if err != nil {
		logger.Fatal(err)
	}
}
