package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-version"
)

const (
	installFile               = "terraform"
	versionPrefix             = "terraform_"
	installPath               = ".terraform.versions"
	recentFile                = "RECENT"
	tfDarwinArm64StartVersion = "1.0.2"
)

var (
	installLocation = "/tmp"
)

// initialize : removes existing symlink to terraform binary// I Don't think this is needed
func initialize() {

	/* Step 1 */
	/* initilize default binary path for terraform */
	/* assumes that terraform is installed here */
	/* we will find the terraform path instalation later and replace this variable with the correct installed bin path */
	installedBinPath := "/usr/local/bin/terraform"

	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand("terraform")
	next := cmd.Find()

	/* overrride installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		installedBinPath = path
	}

	/* check if current symlink to terraform binary exist */
	symlinkExist := CheckSymlink(installedBinPath)

	/* remove current symlink if exist*/
	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}

}

// GetInstallLocation : get location where the terraform binary will be installed,
// will create a directory in the home location if it does not exist
func GetInstallLocation() string {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	userCommon := usr.HomeDir

	/* set installation location */
	installLocation = filepath.Join(userCommon, installPath)

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

	return installLocation

}

//Install : Install the provided version in the argument
func Install(tfRelease *Release, binPath string) {

	/* Check to see if user has permission to the default bin location which is  "/usr/local/bin/terraform"
	 * If user does not have permission to default bin location, proceed to create $HOME/bin and install the tfswitch there
	 * Inform user that they dont have permission to default location, therefore tfswitch was installed in $HOME/bin
	 * Tell users to add $HOME/bin to their path
	 */
	binPath = InstallableBinLocation(binPath)

	initialize()                           //initialize path
	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// Terraform darwin arm64 comes with 1.0.2 and next version
	tfver, err := version.NewVersion(tfRelease.Version)
	if err != nil {
		log.Fatalf("Error generating terraform version for %q: %s", tfRelease.Version, err)
	}
	tf102, err := version.NewVersion(tfDarwinArm64StartVersion)
	if err != nil {
		log.Fatalf("Error generating terraform version for %q: %s", tfDarwinArm64StartVersion, err)
	}
	if goos == "darwin" && goarch == "arm64" && tfver.LessThan(tf102) {
		goarch = "amd64"
	}

	/* check if selected version already downloaded */
	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+tfRelease.Version))
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
		fmt.Printf("Switched terraform to version %q \n", tfRelease.Version)
		AddRecent(tfRelease) //add to recent file for faster lookup
		os.Exit(0)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := ""
	for _, build := range tfRelease.Builds {
		if build.OS == goos && build.Arch == goarch {
			url = build.URL
		}
	}
	if url == "" {
		log.Fatalln("Couldn't determine download url from release")
	}
	zipFile, errDownload := DownloadFromURL(installLocation, url)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		fmt.Println(errDownload)
		os.Exit(1)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("[Error] : Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

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
	fmt.Printf("Switched terraform to version %q \n", tfRelease.Version)
	AddRecent(tfRelease) //add to recent file for faster lookup
	os.Exit(0)
}

// AddRecent : add to recent file
func AddRecent(requestedRelease *Release) {

	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file
	versionFile := filepath.Join(installLocation, recentFile)

	fileExist := CheckFileExist(versionFile)
	if fileExist {
		releases, errRead := ReadLines(versionFile)

		if errRead != nil {
			fmt.Println("File dirty or encountered issue while parsing Release metadata. Recreating cache file.")
			RemoveFiles(versionFile)
			CreateRecentFile(requestedRelease)
			fmt.Printf("[Error] : %s\n", errRead)
			return
		}

		versionExist := VersionExist(requestedRelease, releases)

		if !versionExist {
			if len(releases) >= 3 {
				_, releases = releases[len(releases)-1], releases[:len(releases)-1]

				releases = append([]*Release{requestedRelease}, releases...)
				err := WriteLines(releases, versionFile)
				if err != nil {
					log.Fatalf("Encountered error while updating versions file: %s\n", err)
				}
			} else {
				releases = append([]*Release{requestedRelease}, releases...)
				err := WriteLines(releases, versionFile)
				if err != nil {
					log.Fatalf("Encountered error while updating versions file: %s\n", err)
				}
			}
		}

	} else {
		CreateRecentFile(requestedRelease)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions(mirrorURL string) ([]*Release, error) {

	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file
	versionFile := filepath.Join(installLocation, recentFile)

	fileExist := CheckFileExist(versionFile)
	if fileExist {

		localReleases, errRead := ReadLines(versionFile)
		outputRecent := []*Release{}

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, release := range localReleases {
			/* 	checks if versions in the recent file are valid.
			If any version is invalid, it will be consider dirty
			and the recent file will be removed
			*/
			if !ValidVersionFormat(release.Version) {
				RemoveFiles(versionFile)
				return nil, errRead
			}

			/* 	output can be confusing since it displays the 3 most recent used terraform version
			append the string *recent to the output to make it more user friendly
			*/
			release.Version = fmt.Sprintf("%s *recent", release.Version)
			outputRecent = append(outputRecent, release)
		}

		return outputRecent, nil
	}

	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion *Release) {

	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file

	err := WriteLines([]*Release{requestedVersion}, filepath.Join(installLocation, recentFile))
	if err != nil {
		log.Fatalf("Encountered error while updating versions file: %s\n", err)
	}
}

//ConvertExecutableExt : convert excutable with local OS extension
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

//InstallableBinLocation : Checks if terraform is installable in the location provided by the user.
//If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH and return $HOME/bin as install location
func InstallableBinLocation(userBinPath string) string {

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	binDir := Path(userBinPath)           //get path directory from binary path
	binPathExist := CheckDirExist(binDir) //the default is /usr/local/bin but users can provide custom bin locations

	if binPathExist { //if bin path exist - check if we can write to to it

		binPathWritable := false //assume bin path is not writable
		if runtime.GOOS != "windows" {
			binPathWritable = CheckDirWritable(binDir) //check if is writable on ( only works on LINUX)
		}

		// IF: "/usr/local/bin" or `custom bin path` provided by user is non-writable, (binPathWritable == false), we will attempt to install terraform at the ~/bin location. See ELSE
		if !binPathWritable {

			homeBinExist := CheckDirExist(filepath.Join(usr.HomeDir, "bin")) //check to see if ~/bin exist
			if homeBinExist {                                                //if ~/bin exist, install at ~/bin/terraform
				fmt.Printf("Installing terraform at %s\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terraform")
			} else { //if ~/bin directory does not exist, create ~/bin for terraform installation
				fmt.Printf("Unable to write to: %s\n", userBinPath)
				fmt.Printf("Creating bin directory at: %s\n", filepath.Join(usr.HomeDir, "bin"))
				CreateDirIfNotExist(filepath.Join(usr.HomeDir, "bin")) //create ~/bin
				fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terraform")
			}
		} else { // ELSE: the "/usr/local/bin" or custom path provided by user is writable, we will return installable location
			return filepath.Join(userBinPath)
		}
	}
	fmt.Printf("[Error] : Binary path does not exist: %s\n", userBinPath)
	fmt.Printf("[Error] : Manually create bin directory at: %s and try again.\n", binDir)
	os.Exit(1)
	return ""
}
