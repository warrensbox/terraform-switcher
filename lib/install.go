package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	installFile    = "terraform"
	installVersion = "terraform_"
	installPath    = ".terraform.versions"
	recentFile     = "RECENT"
	defaultBin     = "/usr/local/bin/terraform" //default bin installation dir
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

// getInstallLocation : get location where the terraform binary will be installed,
// will create a directory in the home location if it does not exist
func getInstallLocation() string {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	userCommon := usr.HomeDir

	/* For snapcraft users, SNAP_USER_COMMON environment variable is set by default.
	 * tfswitch does not have permission to save to $HOME/.terraform.versions for snapcraft users
	 * tfswitch will save binaries into $SNAP_USER_COMMON/.terraform.versions */
	if os.Getenv("SNAP_USER_COMMON") != "" {
		userCommon = os.Getenv("SNAP_USER_COMMON")
	}

	/* set installation location */
	installLocation = filepath.Join(userCommon, installPath)

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

	return installLocation

}

//Install : Install the provided version in the argument
func Install(tfversion string, binPath string, mirrorURL string) {

	if !ValidVersionFormat(tfversion) {
		fmt.Printf("The provided terraform version format does not exist - %s. Try `tfswitch -l` to see all available versions.\n", tfversion)
		os.Exit(1)
	}

	pathDir := Path(binPath) //get path directory from binary path
	//binDirExist := CheckDirExist(pathDir) //check bin path exist
	/* Check to see if user has permission to the default bin location which is  "/usr/local/bin/terraform"
	 * If user does not have permission to default bin location, proceed to create $HOME/bin and install the tfswitch there
	 * Inform user that they dont have permission to default location, therefore tfswitch was installed in $HOME/bin
	 * Tell users to add $HOME/bin to their path
	 */
	binPath = InstallableBinLocation(pathDir)

	initialize()                           //initialize path
	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// TODO: Workaround for macos arm64 since terraform doesn't have a binary for it yet
	if goos == "darwin" && goarch == "arm64" {
		goarch = "amd64"
	}

	/* check if selected version already downloaded */
	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, installVersion+tfversion))
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
		fmt.Printf("Switched terraform to version %q \n", tfversion)
		AddRecent(tfversion) //add to recent file for faster lookup
		os.Exit(0)
	}

	//if does not have slash - append slash
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	if !hasSlash {
		mirrorURL = fmt.Sprintf("%s/", mirrorURL)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := mirrorURL + tfversion + "/" + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip"
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
	fmt.Printf("Switched terraform to version %q \n", tfversion)
	AddRecent(tfversion) //add to recent file for faster lookup
	os.Exit(0)
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file
	versionFile := filepath.Join(installLocation, recentFile)

	fileExist := CheckFileExist(versionFile)
	if fileExist {
		lines, errRead := ReadLines(versionFile)

		if errRead != nil {
			fmt.Printf("[Error] : %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !ValidVersionFormat(line) {
				fmt.Println("File dirty. Recreating cache file.")
				RemoveFiles(versionFile)
				CreateRecentFile(requestedVersion)
				return
			}
		}

		versionExist := VersionExist(requestedVersion, lines)

		if !versionExist {
			if len(lines) >= 3 {
				_, lines = lines[len(lines)-1], lines[:len(lines)-1]

				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, versionFile)
			} else {
				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, versionFile)
			}
		}

	} else {
		CreateRecentFile(requestedVersion)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions() ([]string, error) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file
	versionFile := filepath.Join(installLocation, recentFile)

	fileExist := CheckFileExist(versionFile)
	if fileExist {

		lines, errRead := ReadLines(versionFile)
		outputRecent := []string{}

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			/* 	checks if versions in the recent file are valid.
			If any version is invalid, it will be consider dirty
			and the recent file will be removed
			*/
			if !ValidVersionFormat(line) {
				RemoveFiles(versionFile)
				return nil, errRead
			}

			/* 	output can be confusing since it displays the 3 most recent used terraform version
			append the string *recent to the output to make it more user friendly
			*/
			outputRecent = append(outputRecent, fmt.Sprintf("%s *recent", line))
		}

		return outputRecent, nil
	}

	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	WriteLines([]string{requestedVersion}, filepath.Join(installLocation, recentFile))
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
//If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH
//Return $HOME/bin as install location
func InstallableBinLocation(userBin string) string {

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	/* Setup for SNAPCRAFT Users */
	SNAP := os.Getenv("SNAP_REAL_HOME")
	if SNAP != "" { //if SNAP_USER_COMMON env is set, install
		snapHomePath := filepath.Join(SNAP, "bin")
		fmt.Println(snapHomePath)
		CreateDirIfNotExist(snapHomePath)
		fmt.Printf("Installing terraform at %s\n", snapHomePath)
		fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", snapHomePath)
		return filepath.Join(snapHomePath, "terraform")
	}

	existDefaultBin := CheckDirExist(userBin) //the default is /usr/local/bin but users can provide custom bin locations

	if existDefaultBin { //if exist - now see if we can write to to it

		writableToDefault := false
		if runtime.GOOS != "windows" {
			writableToDefault = CheckDirWritable(userBin) //check if is writable on ( only works on LINUX)
		}

		if !writableToDefault {
			exisHomeBin := CheckDirExist(filepath.Join(usr.HomeDir, "bin"))
			if exisHomeBin {
				fmt.Printf("Installing terraform at %s\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terraform")
			}
			PrintCreateDirStmt(userBin, filepath.Join(usr.HomeDir, "bin"))
			CreateDirIfNotExist(filepath.Join(usr.HomeDir, "bin"))
			return filepath.Join(usr.HomeDir, "bin", "terraform")
		}
		return filepath.Join(userBin, "terraform")
	}
	fmt.Printf("[Error] : Binary path does not exist: %s\n", userBin)
	fmt.Printf("[Error] : Manually create bin directory at: %s and try again.\n", userBin)
	os.Exit(1)
	return ""
}

func PrintCreateDirStmt(unableDir string, writable string) {
	fmt.Printf("Unable to write to: %s\n", unableDir)
	fmt.Printf("Creating bin directory at: %s\n", writable)
	fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", writable)
}
