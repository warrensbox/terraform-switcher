package lib_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/warrensbox/terraform-switcher/lib"
)

// TestRenameFile : Create a file, check filename exist,
// rename file, check new filename exit
func TestRenameFile(t *testing.T) {
	installFile := lib.ConvertExecutableExt("terraform")
	installVersion := "terraform_"
	installPath := "/.terraform.versions_test/"
	version := "0.0.7"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)

	createFile(installFilePath)

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("File exist %v", installFilePath)
	} else {
		t.Logf("File does not exist %v", installFilePath)
		t.Error("Missing file")
	}

	installVersionFilePath := lib.ConvertExecutableExt(filepath.Join(installLocation, installVersion+version))

	lib.RenameFile(installFilePath, installVersionFilePath)

	if exist := checkFileExist(installVersionFilePath); exist {
		t.Logf("New file exist %v", installVersionFilePath)
	} else {
		t.Logf("New file does not exist %v", installVersionFilePath)
		t.Error("Missing new file")
	}

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("Old file should not exist %v", installFilePath)
		t.Error("Did not rename file")
	} else {
		t.Logf("Old file does not exist %v", installFilePath)
	}

	cleanUp(installLocation)
}

// TestRemoveFiles : Create a file, check file exist,
// remove file, check file does not exist
func TestRemoveFiles(t *testing.T) {
	installFile := lib.ConvertExecutableExt("terraform")
	installPath := "/.terraform.versions_test/"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)

	createFile(installFilePath)

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("File exist %v", installFilePath)
	} else {
		t.Logf("File does not exist %v", installFilePath)
		t.Error("Missing file")
	}

	lib.RemoveFiles(installFilePath)

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("Old file should not exist %v", installFilePath)
		t.Error("Did not remove file")
	} else {
		t.Logf("Old file does not exist %v", installFilePath)
	}

	cleanUp(installLocation)
}

// TestUnzip : Create a file, check file exist,
// remove file, check file does not exist
func TestUnzip(t *testing.T) {
	installPath := "/.terraform.versions_test/"
	absPath, _ := filepath.Abs("../test-data/test-data.zip")

	fmt.Println(absPath)

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	files, errUnzip := lib.Unzip(absPath, installLocation)

	if errUnzip != nil {
		fmt.Println("Unable to unzip zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	tst := strings.Join(files, "")

	if exist := checkFileExist(tst); exist {
		t.Logf("File exist %v", tst)
	} else {
		t.Logf("File does not exist %v", tst)
		t.Error("Missing file")
	}

	cleanUp(installLocation)
}

// TestCreateDirIfNotExist : Create a directory, check directory exist
func TestCreateDirIfNotExist(t *testing.T) {
	installPath := "/.terraform.versions_test/"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	cleanUp(installLocation)

	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		t.Logf("Directory should not exist %v (expected)", installLocation)
	} else {
		t.Logf("Directory already exist %v (unexpected)", installLocation)
		t.Error("Directory should not exist")
	}

	lib.CreateDirIfNotExist(installLocation)
	t.Logf("Creating directory %v", installLocation)

	if _, err := os.Stat(installLocation); err == nil {
		t.Logf("Directory exist %v (expected)", installLocation)
	} else {
		t.Logf("Directory should exist %v (unexpected)", installLocation)
		t.Error("Directory should exist")
	}

	cleanUp(installLocation)
}

//TestWriteLines : write to file, check readline to verify
func TestWriteLines(t *testing.T) {
	installPath := "/.terraform.versions_test/"
	recentFile := "RECENT"
	semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}(-\w+\d*)?\z`)
	//semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	recentFilePath := filepath.Join(installLocation, recentFile)
	var test_array = []*lib.Release{
		ReleaseConstructor("0.1.1"),
		ReleaseConstructor("0.0.2"),
		ReleaseConstructor("0.0.3"),
		ReleaseConstructor("0.12.0-rc1"),
		ReleaseConstructor("0.12.0-beta1"),
	}

	errWrite := lib.WriteLines(test_array, recentFilePath)

	if errWrite != nil {
		t.Logf("Write should work %v (unexpected)", errWrite)
		log.Fatal(errWrite)
	} else {
		var (
			file             *os.File
			part             []byte
			prefix           bool
			errOpen, errRead error
			localReleases    []*lib.Release
		)
		if file, errOpen = os.Open(recentFilePath); errOpen != nil {
			log.Fatal(errOpen)
		}

		reader := bufio.NewReader(file)
		buffer := bytes.NewBuffer(make([]byte, 0))
		for {
			if part, prefix, errRead = reader.ReadLine(); errRead != nil {
				break
			}
			buffer.Write(part)
			if !prefix {
				var release *lib.Release
				if err := json.Unmarshal(buffer.Bytes(), &release); err != nil {
					log.Fatalf("%q: %s", err, buffer.Bytes())
				}
				localReleases = append(localReleases, release)
				buffer.Reset()
			}
		}
		if errRead == io.EOF {
			errRead = nil
		}

		if errRead != nil {
			log.Fatalf("Error: %s\n", errRead)
		}

		for _, release := range localReleases {
			if !semverRegex.MatchString(release.Version.String()) {
				fmt.Println(release.Version)
				log.Fatalf("Write to file is not invalid: %v\n", release)
				break
			}
		}

		file.Close()
		t.Log("Write versions exist (expected)")
	}

	cleanUp(installLocation)
}

// TestReadLines : read from file, check write to verify
func TestReadLines(t *testing.T) {
	installPath := "/.terraform.versions_test/"
	recentFile := "RECENT"
	semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}(-\w+\d*)?\z`)

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	recentFilePath := filepath.Join(installLocation, recentFile)
	var test_array = []*lib.Release{
		ReleaseConstructor("0.1.1"),
		ReleaseConstructor("0.0.2"),
		ReleaseConstructor("0.0.3"),
		ReleaseConstructor("0.12.0-rc1"),
		ReleaseConstructor("0.12.0-beta1"),
	}
	err := lib.WriteLines(test_array, recentFilePath)
	if err != nil {
		log.Fatalf("Error writing releases: %q", err)
	}

	localReleases, errRead := lib.ReadLines(recentFilePath)

	if errRead != nil {
		log.Fatalf("Error reading Releases from file: %s\n", errRead)
	}

	for _, release := range localReleases {
		if !semverRegex.MatchString(release.Version.String()) {
			fmt.Println(release.Version)
			log.Fatalf("Write to file is not invalid: %v\n", release)
			break
		}
	}

	t.Log("Read versions exist (expected)")

	cleanUp(installLocation)
}

// TestIsDirEmpty : create empty directory, check if empty
func TestIsDirEmpty(t *testing.T) {
	current := time.Now()

	installPath := "/.terraform.versions_test/"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	test_dir := current.Format("2006-01-02")
	test_dir_path := filepath.Join(installLocation, test_dir)
	t.Logf("Create test dir: %v \n", test_dir)

	createDirIfNotExist(installLocation)

	createDirIfNotExist(test_dir_path)

	empty := lib.IsDirEmpty(test_dir_path)

	t.Logf("Expected directory to be empty %v [expected]", test_dir_path)

	if empty == true {
		t.Logf("Directory empty")
	} else {
		t.Error("Directory not empty")
	}

	cleanUp(test_dir_path)
	cleanUp(installLocation)
}

// TestCheckDirHasTGBin : create tg file in directory, check if exist
func TestCheckDirHasTFBin(t *testing.T) {
	goarch := runtime.GOARCH
	goos := runtime.GOOS
	installPath := "/.terraform.versions_test/"
	installFilePrefix := "terraform"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, installFilePrefix+"_"+goos+"_"+goarch))
	createFile(installFileVersionPath)

	empty := lib.CheckDirHasTGBin(installLocation, installFilePrefix)

	t.Logf("Expected directory to have tf file %v [expected]", installFileVersionPath)

	if empty == true {
		t.Logf("Directory empty")
	} else {
		t.Error("Directory not empty")
	}

	cleanUp(installLocation)
}

// TestPath : create file in directory, check if path exist
func TestPath(t *testing.T) {
	installPath := "/.terraform.versions_test"
	installFile := lib.ConvertExecutableExt("terraform")

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := filepath.Join(usr.HomeDir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)
	createFile(installFilePath)

	path := lib.Path(installFilePath)

	t.Logf("Path created %s\n", installFilePath)
	t.Logf("Path expected %s\n", installLocation)
	t.Logf("Path from library %s\n", path)
	if path == installLocation {
		t.Logf("Path exist (expected)")
	} else {
		t.Error("Path does not exist (unexpected)")
	}

	cleanUp(installLocation)
}

// TestGetFileName : remove file ext.  .tfswitch.config returns .tfswitch
func TestGetFileName(t *testing.T) {

	fileNameWithExt := "file.toml"

	fileName := lib.GetFileName(fileNameWithExt)

	if fileName == "file" {
		t.Logf("File removed extension (expected)")
	} else {
		t.Error("File did not remove extension (unexpected)")
	}
}

// TestConvertExecutableExt : convert executable binary with extension
func TestConvertExecutableExt(t *testing.T) {
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	installPath := "/.terraform.versions_test/"
	test_array := []string{
		"terraform",
		"terraform.exe",
		filepath.Join(usr.HomeDir, installPath, "terraform"),
		filepath.Join(usr.HomeDir, installPath, "terraform.exe"),
	}

	for _, fpath := range test_array {
		fpathExt := lib.ConvertExecutableExt((fpath))
		outputMsg := fpath + " converted to " + fpathExt + " on " + runtime.GOOS

		switch runtime.GOOS {
		case "windows":
			if filepath.Ext(fpathExt) != ".exe" {
				t.Error(outputMsg + " (unexpected)")
				continue
			}

			if filepath.Ext(fpath) == ".exe" {
				if fpathExt != fpath {
					t.Error(outputMsg + " (unexpected)")
				} else {
					t.Logf(outputMsg + " (expected)")
				}
				continue
			}

			if fpathExt != fpath+".exe" {
				t.Error(outputMsg + " (unexpected)")
				continue
			}

			t.Logf(outputMsg + " (expected)")
		default:
			if fpath != fpathExt {
				t.Error(outputMsg + " (unexpected)")
				continue
			}

			t.Logf(outputMsg + " (expected)")
		}
	}
}
