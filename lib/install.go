package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"

	"github.com/hashicorp/go-version"
)

var (
	installLocation = "/tmp"
)

// initialize : removes existing symlink to terraform binary based on provided binPath
func initialize(binPath string) {
	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand(binPath)
	next := cmd.Find()

	/* override installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		binPath = path
	}

	/* check if current symlink to terraform binary exist */
	symlinkExist := CheckSymlink(binPath)

	/* remove current symlink if exist*/
	if symlinkExist {
		RemoveSymlink(binPath)
	}
}

// GetInstallLocation : get location where the terraform binary will be installed,
// will create the installDir if it does not exist
func GetInstallLocation(installPath string) string {
	/* set installation location */
	installLocation = filepath.Join(installPath, InstallDir)

	/* Create local installation directory if it does not exist */
	createDirIfNotExist(installLocation)
	return installLocation
}

// install : install the provided version in the argument
func install(tfversion string, binPath string, installPath string, mirrorURL string) {
	var wg sync.WaitGroup
	/* Check to see if user has permission to the default bin location which is  "/usr/local/bin/terraform"
	 * If user does not have permission to default bin location, proceed to create $HOME/bin and install the tfswitch there
	 * Inform user that they don't have permission to default location, therefore tfswitch was installed in $HOME/bin
	 * Tell users to add $HOME/bin to their path
	 */
	binPath = installableBinLocation(binPath)

	initialize(binPath)                               //initialize path
	installLocation = GetInstallLocation(installPath) //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// Terraform darwin arm64 comes with 1.0.2 and next version
	tfver, _ := version.NewVersion(tfversion)
	tf102, _ := version.NewVersion(tfDarwinArm64StartVersion)
	if goos == "darwin" && goarch == "arm64" && tfver.LessThan(tf102) {
		goarch = "amd64"
	}

	/* check if selected version already downloaded */
	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, TerraformPrefix+tfversion))
	fileExist := CheckFileExist(installFileVersionPath)

	/* if selected version already exist, */
	if fileExist {

		/* remove current symlink if exist*/
		symlinkExist := CheckSymlink(binPath)

		if symlinkExist {
			RemoveSymlink(binPath)
		}

		/* set symlink to desired version */
		CreateSymlink(installFileVersionPath, binPath)
		logger.Infof("Switched terraform to version %q", tfversion)
		addRecent(tfversion, installPath, distributionTerraform) //add to recent file for faster lookup
		return
	}

	//if does not have slash - append slash
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	if !hasSlash {
		mirrorURL = fmt.Sprintf("%s/", mirrorURL)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	zipFile, errDownload := DownloadFromURL(installLocation, mirrorURL, tfversion, TerraformPrefix, goos, goarch)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		logger.Fatalf("Error downloading: %s", errDownload)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		logger.Fatalf("Unable to unzip %q file: %v", zipFile, errUnzip)
	}

	logger.Debug("Waiting for deferred functions.")
	wg.Wait()
	/* rename unzipped file to terraform version name - terraform_x.x.x */
	installFilePath := ConvertExecutableExt(filepath.Join(installLocation, installFile))
	RenameFile(installFilePath, installFileVersionPath)

	/* remove zipped file to clear clutter */
	RemoveFiles(zipFile)

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)

	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(installFileVersionPath, binPath)
	logger.Infof("Switched terraform to version %q", tfversion)
	addRecent(tfversion, installPath, distributionTerraform) //add to recent file for faster lookup
	return
}

// ConvertExecutableExt : convert executable with local OS extension
func ConvertExecutableExt(fpath string) string {
	switch runtime.GOOS {
	case "windows":
		if filepath.Ext(fpath) == ".exe" {
			return fpath
		}
		return fpath + ".exe"
	default:
		return fpath
	}
}

// installableBinLocation : Checks if terraform is installable in the location provided by the user.
// If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH and return $HOME/bin as install location
func installableBinLocation(userBinPath string) string {

	homedir := GetHomeDirectory()         //get user's home directory
	binDir := Path(userBinPath)           //get path directory from binary path
	binPathExist := CheckDirExist(binDir) //the default is /usr/local/bin but users can provide custom bin locations

	if binPathExist { //if bin path exist - check if we can write to it

		binPathWritable := false //assume bin path is not writable
		if runtime.GOOS != "windows" {
			binPathWritable = CheckDirWritable(binDir) //check if is writable on ( only works on LINUX)
		}

		// IF: "/usr/local/bin" or `custom bin path` provided by user is non-writable, (binPathWritable == false), we will attempt to install terraform at the ~/bin location. See ELSE
		if !binPathWritable {

			homeBinExist := CheckDirExist(filepath.Join(homedir, "bin")) //check to see if ~/bin exist
			if homeBinExist {                                            //if ~/bin exist, install at ~/bin/terraform
				logger.Infof("Installing terraform at %q", filepath.Join(homedir, "bin"))
				return filepath.Join(homedir, "bin", "terraform")
			} else { //if ~/bin directory does not exist, create ~/bin for terraform installation
				logger.Noticef("Unable to write to %q", userBinPath)
				logger.Infof("Creating bin directory at %q", filepath.Join(homedir, "bin"))
				createDirIfNotExist(filepath.Join(homedir, "bin")) //create ~/bin
				logger.Warnf("Run `export PATH=\"$PATH:%s\"` to append bin to $PATH", filepath.Join(homedir, "bin"))
				return filepath.Join(homedir, "bin", "terraform")
			}
		} else { // ELSE: the "/usr/local/bin" or custom path provided by user is writable, we will return installable location
			return filepath.Join(userBinPath)
		}
	}

	logger.Fatalf("Binary path (%q) does not exist. Manually create bin directory %q and try again.", userBinPath, binDir)
	os.Exit(1)
	return ""
}

// InstallLatestVersion install latest stable tf version
func InstallLatestVersion(dryRun bool, customBinaryPath, installPath string, mirrorURL string) {
	tfversion, _ := getTFLatest(mirrorURL)
	if !dryRun {
		install(tfversion, customBinaryPath, installPath, mirrorURL)
	}
}

// InstallLatestImplicitVersion install latest - argument (version) must be provided
func InstallLatestImplicitVersion(dryRun bool, requestedVersion, customBinaryPath, installPath string, mirrorURL string, preRelease bool) {
	_, err := version.NewConstraint(requestedVersion)
	if err != nil {
		logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	}
	tfversion, err := getTFLatestImplicit(mirrorURL, preRelease, requestedVersion)
	if err == nil && tfversion != "" && !dryRun {
		install(tfversion, customBinaryPath, installPath, mirrorURL)
	}
	logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	PrintInvalidMinorTFVersion()
}

// InstallVersion install with provided version as argument
func InstallVersion(dryRun bool, version, customBinaryPath, installPath, mirrorURL string) {
	logger.Debugf("Install version %s. Dry run: %s", version, strconv.FormatBool(dryRun))
	if !dryRun {
		if validVersionFormat(version) {
			requestedVersion := version

			//check to see if the requested version has been downloaded before
			installLocation := GetInstallLocation(installPath)
			installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, TerraformPrefix+requestedVersion))
			recentDownloadFile := CheckFileExist(installFileVersionPath)
			if recentDownloadFile {
				ChangeSymlink(installFileVersionPath, customBinaryPath)
				logger.Infof("Switched terraform to version %q", requestedVersion)
				addRecent(requestedVersion, installPath, distributionTerraform) //add to recent file for faster lookup
				return
			}

			// If the requested version had not been downloaded before
			// Set list all true - all versions including beta and rc will be displayed
			tflist, _ := getTFList(mirrorURL, true)         // Get list of versions
			exist := versionExist(requestedVersion, tflist) // Check if version exists before downloading it

			if exist {
				install(requestedVersion, customBinaryPath, installPath, mirrorURL)
			} else {
				logger.Fatal("The provided terraform version does not exist.\n Try `tfswitch -l` to see all available versions")
			}
		} else {
			PrintInvalidTFVersion()
			logger.Error("Args must be a valid terraform version")
			UsageMessage()
			os.Exit(1)
		}
	}
}

// InstallOption displays & installs tf version
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func InstallOption(listAll, dryRun bool, customBinaryPath, installPath string, mirrorURL string) {
	tflist, _ := getTFList(mirrorURL, listAll)                                 // Get list of versions
	recentVersions, _ := getRecentVersions(installPath, distributionTerraform) // Get recent versions from RECENT file
	tflist = append(recentVersions, tflist...)                                 // Append recent versions to the top of the list
	tflist = removeDuplicateVersions(tflist)                                   // Remove duplicate version

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
		if errPrompt.Error() == "^C" {
			// Cancel execution
			os.Exit(1)
		} else {
			logger.Fatalf("Prompt failed %v", errPrompt)
		}
	}
	if !dryRun {
		install(tfversion, customBinaryPath, installPath, mirrorURL)
	}
	os.Exit(0)
}
