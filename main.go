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
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/pborman/getopt"
	"github.com/spf13/viper"
	lib "github.com/warrensbox/terraform-switcher/lib"
	"github.com/warrensbox/terraform-switcher/lib/param_parsing"
)

var parameters = param_parsing.GetParameters()
var logger = lib.InitLogger(parameters.LogLevel)
var version string

func main() {

	switch {
	case parameters.VersionFlag:
		if version != "" {
			fmt.Printf("Version: %s\n", version)
		} else {
			fmt.Println("Version not defined during build.")
		}
		os.Exit(0)
	case parameters.HelpFlag:
		lib.UsageMessage()
		os.Exit(0)
	case parameters.ListAllFlag:
		/* show all terraform version including betas and RCs*/
		lib.InstallOption(true, parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.LatestPre != "":
		/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		lib.InstallLatestImplicitVersion(parameters.LatestPre, parameters.CustomBinaryPath, parameters.MirrorURL, true)
	case parameters.ShowLatestPre != "":
		/* show latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestPre, parameters.MirrorURL, true)
	case parameters.LatestStable != "":
		/* latest implicit version. Ex: tfswitch --latest-stable 0.13 downloads 0.13.5 (latest) */
		lib.InstallLatestImplicitVersion(parameters.LatestStable, parameters.CustomBinaryPath, parameters.MirrorURL, false)
	case parameters.ShowLatestStable != "":
		/* show latest implicit stable version. Ex: tfswitch --show-latest-stable 0.13 downloads 0.13.5 (latest) */
		lib.ShowLatestImplicitVersion(parameters.ShowLatestStable, parameters.MirrorURL, false)
	case parameters.LatestFlag:
		/* latest stable version */
		lib.InstallLatestVersion(parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.ShowLatestFlag:
		/* show latest stable version */
		lib.ShowLatestVersion(parameters.MirrorURL)
	case parameters.Version != "":
		lib.InstallVersion(parameters.Version, parameters.CustomBinaryPath, parameters.MirrorURL)
	case parameters.DefaultVersion != "":
		/* if default version is provided - Pick this instead of going for prompt */
		lib.InstallVersion(parameters.DefaultVersion, parameters.CustomBinaryPath, parameters.MirrorURL)
	default:
		// Set list all false - only official release will be displayed
		lib.InstallOption(false, parameters.CustomBinaryPath, parameters.MirrorURL)
	}
}

/* Helper functions */

// install with all possible versions, including beta and rc
func installWithListAll(custBinPath, mirrorURL *string) {
	listAll := true //set list all true - all versions including beta and rc will be displayed
	installOption(listAll, custBinPath, mirrorURL)
}

// install latest stable tf version
func installLatestVersion(custBinPath, mirrorURL *string) {
	tfversion, _ := lib.GetTFLatest(*mirrorURL)
	lib.Install(tfversion, *custBinPath, *mirrorURL)
}

// show install latest stable tf version
func showLatestVersion(custBinPath, mirrorURL *string) {
	tfversion, _ := lib.GetTFLatest(*mirrorURL)
	logger.Infof("%s", tfversion)
}

// install latest - argument (version) must be provided
func installLatestImplicitVersion(requestedVersion string, custBinPath, mirrorURL *string, preRelease bool) {
	_, err := semver.NewConstraint(requestedVersion)
	if err != nil {
		logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	}
	//if lib.ValidMinorVersionFormat(requestedVersion) {
	tfversion, err := lib.GetTFLatestImplicit(*mirrorURL, preRelease, requestedVersion)
	if err == nil && tfversion != "" {
		lib.Install(tfversion, *custBinPath, *mirrorURL)
	}
	logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	lib.PrintInvalidMinorTFVersion()
}

// show latest - argument (version) must be provided
func showLatestImplicitVersion(requestedVersion string, custBinPath, mirrorURL *string, preRelease bool) {
	if lib.ValidMinorVersionFormat(requestedVersion) {
		tfversion, _ := lib.GetTFLatestImplicit(*mirrorURL, preRelease, requestedVersion)
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
func installVersion(arg string, custBinPath *string, mirrorURL *string) {
	if lib.ValidVersionFormat(arg) {
		requestedVersion := arg

		//check to see if the requested version has been downloaded before
		installLocation := lib.GetInstallLocation()
		installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+requestedVersion))
		recentDownloadFile := lib.CheckFileExist(installFileVersionPath)
		if recentDownloadFile {
			lib.ChangeSymlink(installFileVersionPath, *custBinPath)
			logger.Infof("Switched terraform to version %q", requestedVersion)
			lib.AddRecent(requestedVersion) //add to recent file for faster lookup
			os.Exit(0)
		}

		//if the requested version had not been downloaded before
		listAll := true                                     //set list all true - all versions including beta and rc will be displayed
		tflist, _ := lib.GetTFList(*mirrorURL, listAll)     //get list of versions
		exist := lib.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

		if exist {
			lib.Install(requestedVersion, *custBinPath, *mirrorURL)
		} else {
			logger.Fatal("The provided terraform version does not exist.\n Try `tfswitch -l` to see all available versions")
			os.Exit(1)
		}

	} else {
		lib.PrintInvalidTFVersion()
		logger.Error("Args must be a valid terraform version")
		usageMessage()
		os.Exit(1)
	}
}

// retrive file content of regular file
func retrieveFileContents(file string) string {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		logger.Fatalf("Failed reading %q file: %v\n Follow the README.md instructions for setup: https://github.com/warrensbox/terraform-switcher/blob/master/README.md", tfvFilename, err)
		os.Exit(1)
	}
	tfversion := strings.TrimSuffix(string(fileContents), "\n")
	return tfversion
}

// Print message reading file content of :
func readingFileMsg(filename string) {
	logger.Infof("Reading file %q", filename)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func checkTFModuleFileExist(dir string) bool {

	module, _ := tfconfig.LoadModule(dir)
	if len(module.RequiredCore) >= 1 {
		return true
	}
	return false
}

// checkTFEnvExist - checks if the TF_VERSION environment variable is set
func checkTFEnvExist() bool {
	tfversion := os.Getenv("TF_VERSION")
	if tfversion != "" {
		return true
	}
	return false
}

/* parses everything in the toml file, return required version and bin path */
func getParamsTOML(binPath string, dir string) (string, string) {
	path, err := homedir.Dir()

	if err != nil {
		logger.Fatalf("Unable to get home directory: %v", err)
		os.Exit(1)
	}

	if dir == path {
		path = "home directory"
	} else {
		path = "current directory"
	}
	logger.Infof("Reading %q configuration from %q", tomlFilename, path) // Takes the default bin (defaultBin) if user does not specify bin path
	configfileName := lib.GetFileName(tomlFilename)                      //get the config file
	viper.SetConfigType("toml")
	viper.SetConfigName(configfileName)
	viper.AddConfigPath(dir)

	errs := viper.ReadInConfig() // Find and read the config file
	if errs != nil {
		logger.Fatalf("Failed to read %q: %v", tomlFilename, errs)
		os.Exit(1)
	}

	bin := viper.Get("bin")                                                     // read custom binary location
	if binPath == lib.ConvertExecutableExt(lib.GetDefaultBin()) && bin != nil { // if the bin path is the same as the default binary path and if the custom binary is provided in the toml file (use it)
		binPath = os.ExpandEnv(bin.(string))
	}
	//logger.Debug(binPath) // Uncomment this to debug
	version := viper.Get("version") //attempt to get the version if it's provided in the toml
	if version == nil {
		version = ""
	}

	return version.(string), binPath
}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}

/* installOption : displays & installs tf version */
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func installOption(listAll bool, custBinPath, mirrorURL *string) {
	tflist, _ := lib.GetTFList(*mirrorURL, listAll) //get list of versions
	recentVersions, _ := lib.GetRecentVersions()    //get recent versions from RECENT file
	tflist = append(recentVersions, tflist...)      //append recent versions to the top of the list
	tflist = lib.RemoveDuplicateVersions(tflist)    //remove duplicate version

	if len(tflist) == 0 {
		logger.Fatalf("Terraform version list is empty: %s", *mirrorURL)
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

	lib.Install(tfversion, *custBinPath, *mirrorURL)
	os.Exit(0)
}

// install when tf file is provided
func installTFProvidedModule(dir string, custBinPath, mirrorURL *string) {
	var tfconstraints []string
	var exactConstraints []string

	logger.Infof("Reading required version from terraform module at %q", dir)
	module, _ := tfconfig.LoadModule(dir)
	requiredVersions := module.RequiredCore

	for key := range requiredVersions {
		tfconstraint := requiredVersions[key]
		tfconstraintParts := strings.Fields(tfconstraint)

		if len(tfconstraintParts) > 2 {
			logger.Fatalf("Invalid version constraint found: %q", tfconstraint)
			os.Exit(1)
		} else if len(tfconstraintParts) == 1 {
			exactConstraints = append(exactConstraints, tfconstraint)
			tfconstraint = "= " + tfconstraintParts[0]
		}

		if tfconstraintParts[0] == "=" {
			exactConstraints = append(exactConstraints, tfconstraint)
		}

		tfconstraints = append(tfconstraints, tfconstraint)
	}

	if len(exactConstraints) > 0 && len(tfconstraints) > 1 {
		logger.Fatalf("Exact constraint (%q) cannot be combined with other conditions", strings.Join(exactConstraints, ", "))
		os.Exit(1)
	}

	tfconstraint := strings.Join(tfconstraints, ", ")
	installFromConstraint(&tfconstraint, custBinPath, mirrorURL)
}

// install using a version constraint
func installFromConstraint(tfconstraint *string, custBinPath, mirrorURL *string) {

	tfversion, err := lib.GetSemver(tfconstraint, mirrorURL)
	if err == nil {
		lib.Install(tfversion, *custBinPath, *mirrorURL)
	}
	logger.Fatalf("No version found to match constraint: %v\n Follow the README.md instructions for setup: https://github.com/warrensbox/terraform-switcher/blob/master/README.md", err)
	os.Exit(1)
}

// Install using version constraint from terragrunt file
func installTGHclFile(tgFile *string, custBinPath, mirrorURL *string) {
	logger.Infof("Terragrunt file found: %q", *tgFile)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
	if diags.HasErrors() {
		logger.Fatalf("Unable to parse %q file", *tgFile)
		os.Exit(1)
	}
	var version terragruntVersionConstraints
	gohcl.DecodeBody(file.Body, nil, &version)
	installFromConstraint(&version.TerraformVersionConstraint, custBinPath, mirrorURL)
}

type terragruntVersionConstraints struct {
	TerraformVersionConstraint string `hcl:"terraform_version_constraint"`
}

// check if version is defined in hcl file /* lazy-emergency fix - will improve later */
func checkVersionDefinedHCL(tgFile *string) bool {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
	if diags.HasErrors() {
		logger.Fatalf("Unable to parse %q file", *tgFile)
		os.Exit(1)
	}
	var version terragruntVersionConstraints
	gohcl.DecodeBody(file.Body, nil, &version)
	if version == (terragruntVersionConstraints{}) {
		return false
	}
	return true
}
