package main

/*
* Version 0.12.0
* Compatible with Mac OS X AND other LINUX OS ONLY
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
	"path/filepath"
	"sort"
	"strings"

	// original hashicorp upstream have broken dependencies, so using fork as workaround
	// TODO: move back to upstream
	"github.com/Masterminds/semver"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/kiranjthomas/terraform-config-inspect/tfconfig"
	"github.com/mitchellh/go-homedir"

	//	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	"github.com/spf13/viper"

	lib "github.com/warrensbox/terraform-switcher/lib"
)

const (
	defaultMirror = "https://releases.hashicorp.com/terraform"
	defaultBin    = "/usr/local/bin/terraform" //default bin installation dir
	defaultLatest = ""
	tfvFilename   = ".terraform-version"
	rcFilename    = ".tfswitchrc"
	tomlFilename  = ".tfswitch.toml"
	tgHclFilename = "terragrunt.hcl"
)

var (
	version             = "0.12.0\n"
	terraformBinaryPath = ""
)

func main() {
	custBinPath := getopt.StringLong("bin", 'b', lib.ConvertExecutableExt(defaultBin), "Custom binary path. Ex: "+lib.ConvertExecutableExt("/Users/username/bin/terraform"))
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	latestPre := getopt.StringLong("latest-pre", 'p', defaultLatest, "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)")
	latestStable := getopt.StringLong("latest-stable", 's', defaultLatest, "Latest implicit version. Ex: tfswitch --latest-stable 0.13 downloads 0.13.7 (latest)")
	explicitVersion := getopt.StringLong("explicit-version", 'e', "", "Explicit version. Ex: tfswitch --explicit-version 0.15.5 downloads 0.15.5")
	latestFlag := getopt.BoolLong("latest", 'u', "Get latest stable version")
	mirrorURL := getopt.StringLong("mirror", 'm', defaultMirror, "Install from a remote other than the default. Default: https://releases.hashicorp.com/terraform")
	versionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	dir, err := os.Getwd() //get current directory
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	homedir, errHome := homedir.Dir()
	if errHome != nil {
		log.Printf("Failed to get home directory %v\n", errHome)
		os.Exit(1)
	}

	TFVersionFile := filepath.Join(dir, tfvFilename)           //settings for .terraform-version file in current directory (tfenv compatible)
	RCFile := filepath.Join(dir, rcFilename)                   //settings for .tfswitchrc file in current directory (backward compatible purpose)
	TOMLConfigFile := filepath.Join(dir, tomlFilename)         //settings for .tfswitch.toml file in current directory (option to specify bin directory)
	HomeTOMLConfigFile := filepath.Join(homedir, tomlFilename) //settings for .tfswitch.toml file in home directory (option to specify bin directory)
	TGHACLFile := filepath.Join(dir, tgHclFilename)            //settings for terragrunt.hcl file in current directory (option to specify bin directory)

	switch {
	case *versionFlag:
		//if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
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
	case fileExists(TOMLConfigFile) || fileExists(HomeTOMLConfigFile):
		version := ""
		binPath := *custBinPath
		if fileExists(TOMLConfigFile) { //read from toml from current directory
			version, binPath = getParamsTOML(binPath, dir)
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
			installLatestImplicitVersion(*latestPre, custBinPath, mirrorURL, preRelease)
		/* latest implicit version. Ex: tfswitch --latest 0.13 downloads 0.13.5 (latest) */
		case *latestStable != "":
			preRelease := false
			installLatestImplicitVersion(*latestStable, custBinPath, mirrorURL, preRelease)
		/* latest stable version */
		case *latestFlag:
			installLatestVersion(custBinPath, mirrorURL)
		/* version provided on command line as arg */
		case len(args) == 1:
			installVersion(args[0], &binPath, mirrorURL)
		/* provide an tfswitchrc file (IN ADDITION TO A TOML FILE) */
		case fileExists(RCFile):
			readingFileMsg(rcFilename)
			tfversion := retrieveFileContents(RCFile)
			installVersion(tfversion, &binPath, mirrorURL)
		/* if .terraform-version file found (IN ADDITION TO A TOML FILE) */
		case fileExists(TFVersionFile):
			readingFileMsg(tfvFilename)
			tfversion := retrieveFileContents(TFVersionFile)
			installVersion(tfversion, &binPath, mirrorURL)
		/* if versions.tf file found (IN ADDITION TO A TOML FILE) */
		case checkTFModuleFileExist(dir):
			installTFProvidedModule(dir, &binPath, mirrorURL)
		/* if Terraform Version environment variable is set */
		case checkTFEnvExist() && version == "":
			tfversion := os.Getenv("TF_VERSION")
			fmt.Printf("Terraform version environment variable: %s\n", tfversion)
			installVersion(tfversion, custBinPath, mirrorURL)
		/* if terragrunt.hcl file found (IN ADDITION TO A TOML FILE) */
		case fileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile):
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

	/* latest implicit version. Ex: tfswitch --latest 0.13 downloads 0.13.5 (latest) */
	case *latestStable != "":
		preRelease := false
		installLatestImplicitVersion(*latestStable, custBinPath, mirrorURL, preRelease)

	/* latest stable version */
	case *latestFlag:
		installLatestVersion(custBinPath, mirrorURL)

	/* version provided on command line as arg */
	case *explicitVersion != "":
		installVersion(*explicitVersion, custBinPath, mirrorURL)

	/* provide an tfswitchrc file */
	case fileExists(RCFile):
		readingFileMsg(rcFilename)
		tfversion := retrieveFileContents(RCFile)
		installVersion(tfversion, custBinPath, mirrorURL)

	/* if .terraform-version file found */
	case fileExists(TFVersionFile):
		readingFileMsg(tfvFilename)
		tfversion := retrieveFileContents(TFVersionFile)
		installVersion(tfversion, custBinPath, mirrorURL)

	/* if versions.tf file found */
	case checkTFModuleFileExist(dir):
		installTFProvidedModule(dir, custBinPath, mirrorURL)

	/* if terragrunt.hcl file found */
	case fileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile):
		installTGHclFile(&TGHACLFile, custBinPath, mirrorURL)

	/* if Terraform Version environment variable is set */
	case checkTFEnvExist():
		tfversion := os.Getenv("TF_VERSION")
		fmt.Printf("Terraform version environment variable: %s\n", tfversion)
		installVersion(tfversion, custBinPath, mirrorURL)

	// if no arg is provided
	default:
		listAll := false //set list all false - only official release will be displayed
		installOption(listAll, custBinPath, mirrorURL)
	}

	os.Exit(0)
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
	Install(tfversion, *custBinPath, *mirrorURL)
}

// install latest - argument (version) must be provided
func installLatestImplicitVersion(requestedVersion string, custBinPath, mirrorURL *string, preRelease bool) {
	if lib.ValidMinorVersionFormat(requestedVersion) {
		tfversion, _ := lib.GetTFLatestImplicit(*mirrorURL, preRelease, requestedVersion)
		Install(tfversion, *custBinPath, *mirrorURL)
	} else {
		printInvalidMinorTFVersion()
	}
}

// install with provided version as argument
func installVersion(arg string, custBinPath *string, mirrorURL *string) {
	if lib.ValidVersionFormat(arg) {
		requestedVersion := arg
		listAll := true                                     //set list all true - all versions including beta and rc will be displayed
		tflist, _ := lib.GetTFList(*mirrorURL, listAll)     //get list of versions
		exist := lib.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

		if exist {
			Install(requestedVersion, *custBinPath, *mirrorURL)
		} else {
			fmt.Println("The provided terraform version does not exist. Try `tfswitch -l` to see all available versions.")
			os.Exit(1)
		}

	} else {
		printInvalidTFVersion()
		fmt.Println("Args must be a valid terraform version")
		usageMessage()
		os.Exit(1)
	}
}

// Print invalid TF version
func printInvalidTFVersion() {
	fmt.Println("Invalid terraform version format. Format should be #.#.# or #.#.#-@# where # are numbers and @ are word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
}

// Print invalid TF version
func printInvalidMinorTFVersion() {
	fmt.Println("Invalid minor terraform version format. Format should be #.# where # are numbers. For example, 0.11 is valid version")
}

//retrive file content of regular file
func retrieveFileContents(file string) string {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md\n", tfvFilename)
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	tfversion := strings.TrimSuffix(string(fileContents), "\n")
	return tfversion
}

// Print message reading file content of :
func readingFileMsg(filename string) {
	fmt.Printf("Reading file %s \n", filename)
}

// fileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
	path, _ := homedir.Dir()
	if dir == path {
		path = "home directory"
	} else {
		path = "current directory"
	}
	fmt.Printf("Reading configuration from %s\n", path+" for "+tomlFilename) //takes the default bin (defaultBin) if user does not specify bin path
	configfileName := lib.GetFileName(tomlFilename)                          //get the config file
	viper.SetConfigType("toml")
	viper.SetConfigName(configfileName)
	viper.AddConfigPath(dir)

	errs := viper.ReadInConfig() // Find and read the config file
	if errs != nil {
		fmt.Printf("Unable to read %s provided\n", tomlFilename) // Handle errors reading the config file
		fmt.Println(errs)
		os.Exit(1) // exit immediately if config file provided but it is unable to read it
	}

	bin := viper.Get("bin")                                            // read custom binary location
	if binPath == lib.ConvertExecutableExt(defaultBin) && bin != nil { // if the bin path is the same as the default binary path and if the custom binary is provided in the toml file (use it)
		binPath = os.ExpandEnv(bin.(string))
	}
	//fmt.Println(binPath) //uncomment this to debug
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
		fmt.Println("[ERROR] : List is empty")
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
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	Install(tfversion, *custBinPath, *mirrorURL)
}

// install when tf file is provided
func installTFProvidedModule(dir string, custBinPath, mirrorURL *string) {
	fmt.Printf("Reading required version from terraform file\n")
	module, _ := tfconfig.LoadModule(dir)
	tfconstraint := module.RequiredCore[0] //we skip duplicated definitions and use only first one
	installFromConstraint(&tfconstraint, custBinPath, mirrorURL)
}

// install using a version constraint
func installFromConstraint(tfconstraint *string, custBinPath, mirrorURL *string) {
	tfversion := ""
	listAll := true                                 //set list all true - all versions including beta and rc will be displayed
	tflist, _ := lib.GetTFList(*mirrorURL, listAll) //get list of versions
	fmt.Printf("Reading required version from constraint: %s\n", *tfconstraint)

	constrains, err := semver.NewConstraint(*tfconstraint) //NewConstraint returns a Constraints instance that a Version instance can be checked against
	if err != nil {
		fmt.Printf("Error parsing constraint: %s\nPlease check constrain syntax on terraform file.\n", err)
		fmt.Println()
		os.Exit(1)
	}
	versions := make([]*semver.Version, len(tflist))
	for i, tfvals := range tflist {
		version, err := semver.NewVersion(tfvals) //NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
		if err != nil {
			fmt.Printf("Error parsing version: %s", err)
			os.Exit(1)
		}

		versions[i] = version
	}

	sort.Sort(sort.Reverse(semver.Collection(versions)))

	for _, element := range versions {
		if constrains.Check(element) { // Validate a version against a constraint
			tfversion = element.String()
			fmt.Printf("Matched version: %s\n", tfversion)
			if lib.ValidVersionFormat(tfversion) { //check if version format is correct
				Install(tfversion, *custBinPath, *mirrorURL)
				return
			} else {
				printInvalidTFVersion()
				os.Exit(1)
			}
		}
	}

	fmt.Println("No version found to match constraint. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md")
	os.Exit(1)
}

func Install(tfversion string, binPath string, mirrorURL string) {
	terraformBinaryPath = lib.Install(tfversion, binPath, mirrorURL)
}

// Install using version constraint from terragrunt file
func installTGHclFile(tgFile *string, custBinPath, mirrorURL *string) {
	fmt.Printf("Terragrunt file found: %s\n", *tgFile)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
	if diags.HasErrors() {
		fmt.Println("Unable to parse HCL file")
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
		fmt.Println("Unable to parse HCL file")
		os.Exit(1)
	}
	var version terragruntVersionConstraints
	gohcl.DecodeBody(file.Body, nil, &version)
	if version == (terragruntVersionConstraints{}) {
		return false
	}
	return true
}
