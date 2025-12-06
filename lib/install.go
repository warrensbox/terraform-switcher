//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/mattn/go-isatty"

	"github.com/hashicorp/go-version"
)

var installLocation = "/tmp"

// initialize : removes existing symlink to terraform binary based on provided binPath
//
//nolint:unused // FIXME: Function is not used 10-Mar-2025
func initialize(binPath string) {
	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand(binPath)
	next := cmd.Find()

	/* override installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		binPath = path
	}

	/* check if current symlink to terraform binary exists */
	symlinkExist := CheckSymlink(binPath)

	/* remove current symlink if exists */
	if symlinkExist {
		if err := RemoveSymlink(binPath); err != nil {
			logger.Errorf("Error removing symlink: %v", err)
		}
	}
}

// GetInstallLocation : get location where the terraform binary will be installed,
// will create the installDir if it does not exist
func GetInstallLocation(installPath string) string {
	/* set installation location */
	installLocation = filepath.Join(installPath, InstallDir)

	/* Create local installation directory if it does not exists */
	createDirIfNotExist(installLocation)
	return installLocation
}

// install : install the provided version in the argument
//
//nolint:gocyclo
func install(product Product, dryRun, showRequiredFlag bool, tfversion, binPath, installPath, mirrorURL, goarch string) error {
	installLocation := GetInstallLocation(installPath)
	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, product.GetVersionPrefix()+tfversion))

	// check to see if the requested version has been downloaded before
	if CheckFileExist(installFileVersionPath) && !showRequiredFlag {
		if dryRun {
			logger.Infof("[DRY-RUN] Would have attempted to switch %s to version %q", product.GetName(), tfversion)
			return nil
		}
		return switchToVersion(product, tfversion, binPath, installPath, installFileVersionPath)
	}

	// If the requested version had not been downloaded before,
	// set list all true - all versions including beta and rc will be displayed
	tflist, errTFList := getTFList(mirrorURL, true) // Get list of versions
	if errTFList != nil {
		return fmt.Errorf("Error getting list of %s versions from %q: %v", product.GetName(), mirrorURL, errTFList)
	}

	// Check if version exists before downloading it
	if !versionExist(tfversion, tflist) {
		return fmt.Errorf("Provided %s version does not exist: %q.\n\tTry `tfswitch -l` to see all available versions", product.GetName(), tfversion)
	}

	if goarch != runtime.GOARCH {
		logger.Warnf("Installing for %q CPU architecture on %q!", goarch, runtime.GOARCH)
	}

	goos := runtime.GOOS

	// Terraform darwin arm64 comes with 1.0.2 and next version
	tfver, tfverErr := version.NewVersion(tfversion)
	if tfverErr != nil {
		return fmt.Errorf("Error parsing %q version: %v", tfversion, tfverErr)
	}
	tf102, tf102Err := version.NewVersion(tfDarwinArm64StartVersion)
	if tf102Err != nil {
		return fmt.Errorf("Error parsing %q version: %v", tfDarwinArm64StartVersion, tf102Err)
	}
	if goos == "darwin" && goarch == "arm64" && tfver.LessThan(tf102) {
		goarch = "amd64"
	}

	logArgs := fmt.Sprintf("%s version %q for %q", product.GetName(), tfversion, goos+"_"+goarch)
	switch {
	case dryRun:
		logger.Infof("[DRY-RUN] Would have attempted to install %s", logArgs)
		return nil
	case showRequiredFlag:
		logger.Infof("Showing required %s", logArgs)
		fmt.Printf("%s\n", tfversion)
		return nil
	default:
		logger.Infof("Installing %s", logArgs)
	}

	var wg sync.WaitGroup

	// Create exclusive lock to prevent multiple concurrent installations
	// Put lockfile in temp directory to get it cleaned up on reboot
	lockFile := filepath.Join(os.TempDir(), ".tfswitch."+product.GetId()+".lock")
	// 90 attempts * 2 seconds = 3 minutes to acquire lock, otherwise bail out
	lockedFH, err := acquireLock(lockFile, 90, 2*time.Second)
	if err != nil {
		logger.Fatal(err)
	}
	// Release lock when done
	defer releaseLock(lockFile, lockedFH)

	// If selected version doesn't already exist, proceed to download it
	zipFile, errDownload := DownloadProductFromURL(product, installLocation, product.GetArtifactUrl(mirrorURL, tfversion), tfversion, product.GetArchivePrefix(), goos, goarch)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		// logger.Fatal doesn't invoke deferred functions,
		// so need to release the lock explicitly
		releaseLock(lockFile, lockedFH)
		logger.Fatalf("Error downloading: %s", errDownload)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := Unzip(zipFile, installLocation, product.GetExecutableName())
	if errUnzip != nil {
		// logger.Fatal doesn't invoke deferred functions,
		// so need to release the lock explicitly
		releaseLock(lockFile, lockedFH)
		logger.Fatalf("Unable to unzip %q file: %v", zipFile, errUnzip)
	}

	logger.Debug("Waiting for deferred functions")
	wg.Wait()
	/* rename unzipped file to terraform version name - terraform_x.x.x */
	installFilePath := ConvertExecutableExt(filepath.Join(installLocation, product.GetExecutableName()))
	RenameFile(installFilePath, installFileVersionPath)

	/* remove zipped file to clear clutter */
	RemoveFiles(zipFile)

	return switchToVersion(product, tfversion, binPath, installPath, installFileVersionPath)
}

func switchToVersion(product Product, tfversion string, binPath string, installPath string, installFileVersionPath string) error {
	err := ChangeProductSymlink(product, installFileVersionPath, binPath)
	if err != nil {
		return err
	}

	logger.Infof("Switched %s to version %q", product.GetName(), tfversion)

	// add to recent file for faster lookup
	addRecent(tfversion, installPath, product)
	return nil
}

// ConvertExecutableExt : convert executable with local OS extension
func ConvertExecutableExt(fpath string) string {
	const winExt = ".exe"
	switch runtime.GOOS {
	case windows:
		if filepath.Ext(fpath) == winExt {
			return fpath
		}
		return fpath + winExt
	default:
		return fpath
	}
}

// installableBinLocation : Checks if terraform is installable in the location provided by the user.
// If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH and return $HOME/bin as install location
//
// Deprecated: This function has been deprecated and will be removed in v2.0.0
//
//nolint:unused // Function is deprecated
func installableBinLocation(product Product, userBinPath string) string {
	homedir := GetHomeDirectory()         // get user's home directory
	binDir := Path(userBinPath)           // get path directory from binary path
	binPathExist := CheckDirExist(binDir) // the default is /usr/local/bin but users can provide custom bin locations

	if binPathExist { // if bin path exist - check if we can write to it

		binPathWritable := false // assume bin path is not writable
		if runtime.GOOS != windows {
			binPathWritable = CheckDirWritable(binDir) // check if is writable on ( only works on LINUX)
		}

		// IF: "/usr/local/bin" or `custom bin path` provided by user is non-writable, (binPathWritable == false),
		// we will attempt to install terraform at the ~/bin location. See ELSE
		if !binPathWritable {
			homeBinDir := filepath.Join(homedir, "bin")
			if !CheckDirExist(homeBinDir) { // if ~/bin exist, install at ~/bin/terraform
				logger.Noticef("Unable to write to %q", userBinPath)
				logger.Infof("Creating bin directory at %q", homeBinDir)
				createDirIfNotExist(homeBinDir) // create ~/bin
				logger.Warnf("Run `export PATH=\"$PATH:%s\"` to append bin to $PATH", homeBinDir)
			}
			logger.Infof("Installing %s at %q", product.GetName(), homeBinDir)
			return filepath.Join(homeBinDir, product.GetExecutableName())

		}
		// ELSE: the "/usr/local/bin" or custom path provided by user is writable, we will return installable location
		return userBinPath
	}

	logger.Fatalf("Binary path (%q) does not exist. Manually create bin directory %q and try again", userBinPath, binDir)
	os.Exit(1)
	return ""
}

// InstallLatestVersion install latest stable tf version
//
// Deprecated: This function has been deprecated in favor of InstallLatestProductVersion and will be removed in v2.0.0
func InstallLatestVersion(dryRun, showRequiredFlag bool, customBinaryPath, installPath, mirrorURL, arch string) {
	product := getLegacyProduct()
	//nolint:errcheck // Function is deprecated
	InstallLatestProductVersion(product, dryRun, showRequiredFlag, customBinaryPath, installPath, mirrorURL, arch)
}

// InstallLatestProductVersion install latest stable tf version
func InstallLatestProductVersion(product Product, dryRun, showRequiredFlag bool, customBinaryPath, installPath, mirrorURL, arch string) error {
	logger.Infof("Latest version requested explicitly")
	tfversion, err := getTFLatest(mirrorURL)
	if err != nil {
		return fmt.Errorf("Error getting latest %s version from %q: %v", product.GetName(), mirrorURL, err)
	}

	return install(product, dryRun, showRequiredFlag, tfversion, customBinaryPath, installPath, mirrorURL, arch)
}

// InstallLatestImplicitVersion install latest - argument (version) must be provided
//
// Deprecated: This function has been deprecated in favor of InstallLatestProductImplicitVersion and will be removed in v2.0.0
func InstallLatestImplicitVersion(dryRun, showRequiredFlag bool, requestedVersion, customBinaryPath, installPath, mirrorURL, arch string, preRelease bool) {
	product := getLegacyProduct()
	//nolint:errcheck // Function is deprecated
	InstallLatestProductImplicitVersion(product, dryRun, showRequiredFlag, requestedVersion, customBinaryPath, installPath, mirrorURL, arch, preRelease)
}

// InstallLatestProductImplicitVersion install latest - argument (version) must be provided
func InstallLatestProductImplicitVersion(product Product, dryRun, showRequiredFlag bool, requestedVersion, customBinaryPath, installPath, mirrorURL, arch string, preRelease bool) error {
	logger.Infof("Latest implicit version matching %q requested explicitly", requestedVersion)
	_, err := version.NewConstraint(requestedVersion)
	if err != nil {
		// @TODO Should this return an error?
		logger.Errorf("Error parsing constraint %q: %v", requestedVersion, err)
	}
	tfversion, err := getTFLatestImplicit(mirrorURL, preRelease, requestedVersion)
	if err == nil && tfversion != "" {
		if errInstall := install(product, dryRun, showRequiredFlag, tfversion, customBinaryPath, installPath, mirrorURL, arch); errInstall != nil {
			return fmt.Errorf("Error installing %s version %q: %v", product.GetName(), tfversion, errInstall)
		}
		return nil
	}
	PrintInvalidMinorTFVersion()
	return fmt.Errorf("error parsing constraint %q: %v", requestedVersion, err)
}

// InstallVersion install Terraform product
//
// Deprecated: This function has been deprecated in favor of InstallProductVersion and will be removed in v2.0.0
func InstallVersion(dryRun, showRequiredFlag bool, version, customBinaryPath, installPath, mirrorURL, arch string) {
	product := getLegacyProduct()
	//nolint:errcheck // Function is deprecated
	InstallProductVersion(product, dryRun, showRequiredFlag, version, customBinaryPath, installPath, mirrorURL, arch)
}

// InstallProductVersion install with provided version as argument
func InstallProductVersion(product Product, dryRun, showRequiredFlag bool, version, customBinaryPath, installPath, mirrorURL, arch string) error {
	var logPrefix string
	if dryRun {
		logPrefix = "[DRY-RUN] "
	}
	logger.Debugf("%sTargeting %s version %q", logPrefix, product.GetName(), version)

	if validVersionFormat(version) {
		return install(product, dryRun, showRequiredFlag, version, customBinaryPath, installPath, mirrorURL, arch)
	}
	PrintInvalidTFVersion()
	return fmt.Errorf("Argument must be a valid %s version", product.GetName())
}

// InstallProductOption displays & installs tf version
// listAll = true - all versions including beta and rc will be displayed */
// listAll = false - only official stable release are displayed */
//
// Deprecated: This function has been deprecated in favor of InstallProductOption and will be removed in v2.0.0
func InstallOption(listAll, dryRun, showRequiredFlag bool, customBinaryPath, installPath, mirrorURL, arch string) {
	product := getLegacyProduct()
	//nolint:errcheck // Function is deprecated
	InstallProductOption(product, listAll, dryRun, showRequiredFlag, customBinaryPath, installPath, mirrorURL, arch)
}

type VersionSelector struct {
	Version string
	Label   string
}

// InstallProductOption displays & installs tf version
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func InstallProductOption(product Product, listAll, dryRun, showRequiredFlag bool, customBinaryPath, installPath, mirrorURL, arch string) error {
	var selectedVersion string

	if !showRequiredFlag {
		var selectVersions []VersionSelector

		// Get all available versions from remote
		tfList, errTFList := getTFList(mirrorURL, listAll)
		if errTFList != nil {
			return fmt.Errorf("Error getting list of %s versions from %q: %v", product.GetName(), mirrorURL, errTFList)
		}

		if len(tfList) == 0 {
			return fmt.Errorf("[%s] Remote returned empty versions list: %s", product.GetName(), mirrorURL)
		}

		versionMap := make(map[string]bool)

		// Add recent versions
		//nolint:errcheck // getRecentVersions always returns nil error %-/
		recentVersions, _ := getRecentVersions(installPath, product)
		for _, version := range recentVersions {
			selectVersions = append(selectVersions, VersionSelector{
				Version: version,
				Label:   version + " *recent",
			})
			versionMap[version] = true
		}

		for _, version := range tfList {
			if !versionMap[version] {
				selectVersions = append(selectVersions, VersionSelector{
					Version: version,
					Label:   version,
				})
			}
		}

		if len(selectVersions) == 0 {
			return fmt.Errorf("%s version list is empty: %s", product.GetName(), mirrorURL)
		}

		/* prompt user to select version of product */
		if !isatty.IsTerminal(os.Stdin.Fd()) || !isatty.IsTerminal(os.Stdout.Fd()) {
			return errors.New("Interactive prompt isn't meant for non-interactive terminal")
		}
		prompt := promptui.Select{
			Label: fmt.Sprintf("Select %s version", product.GetName()),
			Items: selectVersions,
			Templates: &promptui.SelectTemplates{
				// Use templates from defaults in promptui, but specifying
				// the Label attribute of the VersionSelectors
				Label:    fmt.Sprintf("%s {{.Label}}: ", promptui.IconInitial),
				Active:   fmt.Sprintf("%s {{ .Label | underline }}", promptui.IconSelect),
				Inactive: "  {{.Label}}",
				Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Label | faint }}`, promptui.IconGood),
			},
		}

		selectedIdx, _, errPrompt := prompt.Run()
		if errPrompt != nil {
			if errPrompt.Error() == "^C" {
				// Cancel execution
				return fmt.Errorf("Interrupted by user")
			}
			return fmt.Errorf("PromptUI unexpectedly failed: %v", errPrompt)
		}

		selectedVersion = selectVersions[selectedIdx].Version
		logger.Infof("Selected %s version: %s", product.GetName(), selectedVersion)
	} else {
		var errGetLatest error
		selectedVersion, errGetLatest = getTFLatest(mirrorURL)
		if errGetLatest != nil {
			return fmt.Errorf("Error getting latest %s version from %q: %v", product.GetName(), mirrorURL, errGetLatest)
		}
		logger.Warnf("No required or otherwise explicitly requested version found: defaulting to latest %s version (%s)", product.GetName(), selectedVersion)
	}

	return install(product, dryRun, showRequiredFlag, selectedVersion, customBinaryPath, installPath, mirrorURL, arch)
}
