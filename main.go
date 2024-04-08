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

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/mitchellh/go-homedir"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	"github.com/spf13/viper"

	lib "github.com/warrensbox/terraform-switcher/lib"
)

const (
	defaultMirror = "https://releases.hashicorp.com/terraform"
	defaultLatest = ""
	tfvFilename   = ".terraform-version"
	rcFilename    = ".tfswitchrc"
	tomlFilename  = ".tfswitch.toml"
	tgHclFilename = "terragrunt.hcl"
	versionPrefix = "terraform_"
)

var logger = lib.InitLogger()
var version = "0.12.0"

func main() {
	dir := lib.GetCurrentDirectory()
	custBinPath := getopt.StringLong("bin", 'b', lib.ConvertExecutableExt(lib.GetDefaultBin()), "Custom binary path. Ex: tfswitch -b "+lib.ConvertExecutableExt("/Users/username/bin/terraform"))
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	latestPre := getopt.StringLong("latest-pre", 'p', defaultLatest, "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)")
	showLatestPre := getopt.StringLong("show-latest-pre", 'P', defaultLatest, "Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)")
	latestStable := getopt.StringLong("latest-stable", 's', defaultLatest, "Latest implicit version based on a constraint. Ex: tfswitch --latest-stable 0.13.0 downloads 0.13.7 and 0.13 downloads 0.15.5 (latest)")
	showLatestStable := getopt.StringLong("show-latest-stable", 'S', defaultLatest, "Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)")
	latestFlag := getopt.BoolLong("latest", 'u', "Get latest stable version")
	showLatestFlag := getopt.BoolLong("show-latest", 'U', "Show latest stable version")
	mirrorURL := getopt.StringLong("mirror", 'm', defaultMirror, "Install from a remote API other than the default. Default: "+defaultMirror)
	chDirPath := getopt.StringLong("chdir", 'c', dir, "Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory")
	versionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	defaultVersion := getopt.StringLong("default", 'd', defaultLatest, "Default to this version in case no other versions could be detected. Ex: tfswitch --default 1.2.4")

	getopt.StringVarLong(&lib.PubKeyId, "public-key-id", 'k', "The ID of the public key to check the checksums against")
	getopt.StringVarLong(&lib.PubKeyPrefix, "public-key-prefix", 'x', "The prefix of the public key. i.e. \"hashicorp_\"")
	getopt.StringVarLong(&lib.PubKeyUri, "public-key-uri", 'y', "The URI to download the public key from")

	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	homedir, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Unable to get home directory: %v", err)
		os.Exit(1)
	}

	TFVersionFile := filepath.Join(*chDirPath, tfvFilename)    //settings for .terraform-version file in current directory (tfenv compatible)
	RCFile := filepath.Join(*chDirPath, rcFilename)            //settings for .tfswitchrc file in current directory (backward compatible purpose)
	TOMLConfigFile := filepath.Join(*chDirPath, tomlFilename)  //settings for .tfswitch.toml file in current directory (option to specify bin directory)
	HomeTOMLConfigFile := filepath.Join(homedir, tomlFilename) //settings for .tfswitch.toml file in home directory (option to specify bin directory)
	TGHACLFile := filepath.Join(*chDirPath, tgHclFilename)     //settings for terragrunt.hcl file in current directory (option to specify bin directory)

	switch {
	case *versionFlag:
		//if *versionFlag {
		logger.Infof("Version: %s", version)
	case *helpFlag:
		//} else if *helpFlag {
		usageMessage()
	/* Checks if the .tfswitch.toml file exist in home or current directory
	 * This block checks to see if the tfswitch toml file is provided in the current path.
	 * If the .tfswitch.toml file exist, it has a higher precedence than the .tfswitchrc file
	 * You can specify the custom binary path and the version you desire
	 * If you provide a custom binary path with the -b option, this will override the bin value in the toml file
	 * If you provide a version on the command line, this will override the version value in the toml file
	 */
	case lib.FileExists(TOMLConfigFile) || lib.FileExists(HomeTOMLConfigFile):
		version := ""
		binPath := *custBinPath
		if lib.FileExists(TOMLConfigFile) { //read from toml from current directory
			version, binPath = getParamsTOML(binPath, *chDirPath)
		} else { // else read from toml from home directory
			version, binPath = getParamsTOML(binPath, homedir)
		}

		switch {
		/* GIVEN A TOML FILE, */
		/* show all terraform version including betas and RCs*/
		case *listAllFlag:
			listAll := true //set list all true - all versions including beta and rc will be displayed
			installOption(listAll, &binPath, mirrorURL)
		/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
		case *latestPre != "":
			preRelease := true
			installLatestImplicitVersion(*latestPre, &binPath, mirrorURL, preRelease)
		/* latest implicit version. Ex: tfswitch --latest 0.13 downloads 0.13.5 (latest) */
		case *latestStable != "":
			preRelease := false
			installLatestImplicitVersion(*latestStable, &binPath, mirrorURL, preRelease)
		/* latest stable version */
		case *latestFlag:
			installLatestVersion(&binPath, mirrorURL)
		/* version provided on command line as arg */
		case len(args) == 1:
			installVersion(args[0], &binPath, mirrorURL)
		/* provide an tfswitchrc file (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(RCFile) && len(args) == 0:
			readingFileMsg(rcFilename)
			tfversion := retrieveFileContents(RCFile)
			installVersion(tfversion, &binPath, mirrorURL)
		/* if .terraform-version file found (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(TFVersionFile) && len(args) == 0:
			readingFileMsg(tfvFilename)
			tfversion := retrieveFileContents(TFVersionFile)
			installVersion(tfversion, &binPath, mirrorURL)
		/* if versions.tf file found (IN ADDITION TO A TOML FILE) */
		case checkTFModuleFileExist(*chDirPath) && len(args) == 0:
			installTFProvidedModule(*chDirPath, &binPath, mirrorURL)
		/* if Terraform Version environment variable is set */
		case checkTFEnvExist() && len(args) == 0 && version == "":
			tfversion := os.Getenv("TF_VERSION")
			logger.Infof("\"TF_VERSION\" environment variable value: %q", tfversion)
			installVersion(tfversion, &binPath, mirrorURL)
		/* if terragrunt.hcl file found (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile) && len(args) == 0:
			installTGHclFile(&TGHACLFile, &binPath, mirrorURL)
		// if no arg is provided - but toml file is provided
		case version != "":
			installVersion(version, &binPath, mirrorURL)
		default:
			listAll := false //set list all false - only official release will be displayed
			installOption(listAll, &binPath, mirrorURL)
		}

	/* show all terraform version including betas and RCs*/
	case *listAllFlag:
		installWithListAll(custBinPath, mirrorURL)

	/* latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
	case *latestPre != "":
		preRelease := true
		installLatestImplicitVersion(*latestPre, custBinPath, mirrorURL, preRelease)

	/* show latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest) */
	case *showLatestPre != "":
		preRelease := true
		showLatestImplicitVersion(*showLatestPre, custBinPath, mirrorURL, preRelease)

	/* latest implicit version. Ex: tfswitch --latest 0.13 downloads 0.13.5 (latest) */
	case *latestStable != "":
		preRelease := false
		installLatestImplicitVersion(*latestStable, custBinPath, mirrorURL, preRelease)

	/* show latest implicit stable version. Ex: tfswitch --latest 0.13 downloads 0.13.5 (latest) */
	case *showLatestStable != "":
		preRelease := false
		showLatestImplicitVersion(*showLatestStable, custBinPath, mirrorURL, preRelease)

	/* latest stable version */
	case *latestFlag:
		installLatestVersion(custBinPath, mirrorURL)

	/* show latest stable version */
	case *showLatestFlag:
		showLatestVersion(custBinPath, mirrorURL)

	/* version provided on command line as arg */
	case len(args) == 1:
		installVersion(args[0], custBinPath, mirrorURL)

	/* provide an tfswitchrc file */
	case lib.FileExists(RCFile) && len(args) == 0:
		readingFileMsg(rcFilename)
		tfversion := retrieveFileContents(RCFile)
		installVersion(tfversion, custBinPath, mirrorURL)

	/* if .terraform-version file found */
	case lib.FileExists(TFVersionFile) && len(args) == 0:
		readingFileMsg(tfvFilename)
		tfversion := retrieveFileContents(TFVersionFile)
		installVersion(tfversion, custBinPath, mirrorURL)

	/* if versions.tf file found */
	case checkTFModuleFileExist(*chDirPath) && len(args) == 0:
		installTFProvidedModule(*chDirPath, custBinPath, mirrorURL)

	/* if terragrunt.hcl file found */
	case lib.FileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile) && len(args) == 0:
		installTGHclFile(&TGHACLFile, custBinPath, mirrorURL)

	/* if Terraform Version environment variable is set */
	case checkTFEnvExist() && len(args) == 0:
		tfversion := os.Getenv("TF_VERSION")
		logger.Infof("\"TF_VERSION\" environment variable value: %q", tfversion)
		installVersion(tfversion, custBinPath, mirrorURL)

	/* if default version is provided - Pick this instead of going for prompt */
	case *defaultVersion != "":
		installVersion(*defaultVersion, custBinPath, mirrorURL)

	// if no arg is provided
	default:
		listAll := false //set list all false - only official release will be displayed
		installOption(listAll, custBinPath, mirrorURL)
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
	logger.Infof("Reading %q from %s", tomlFilename, path)
	configfileName := lib.GetFileName(tomlFilename) // get the config file
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
	curDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Could not get current working directory: %v", err)
		os.Exit(1)
	}

	absPath := dir
	if !filepath.IsAbs(dir) {
		absPath, err = filepath.Abs(dir)
		if err != nil {
			logger.Fatalf("Could not derive absolute path to %q: %v", dir, err)
			os.Exit(1)
		}
	}

	relPath, err := filepath.Rel(curDir, absPath)
	if err != nil {
		logger.Fatalf("Could not derive relative path to %q: %v", dir, err)
		os.Exit(1)
	}

	logger.Infof("Reading required version from terraform module at %q", relPath)
	module, _ := tfconfig.LoadModule(dir)
	tfconstraint := module.RequiredCore[0] //we skip duplicated definitions and use only first one
	installFromConstraint(&tfconstraint, custBinPath, mirrorURL)
}

// install using a version constraint
func installFromConstraint(tfconstraint *string, custBinPath, mirrorURL *string) {

	tfversion, err := lib.GetSemver(tfconstraint, mirrorURL)
	if err == nil {
		lib.Install(tfversion, *custBinPath, *mirrorURL)
	}
	logger.Fatalf("No version found to match constraint: %v.\n Follow the README.md instructions for setup: https://github.com/warrensbox/terraform-switcher/blob/master/README.md", err)
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
