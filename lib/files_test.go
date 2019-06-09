package lib_test

import (
	"bufio"
	"bytes"
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

	installFile := "terraform"
	installVersion := "terraform_"
	installPath := "/.terraform.versions_test/"
	version := "0.0.7"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	createFile(installLocation + installFile)

	if exist := checkFileExist(installLocation + installFile); exist {
		t.Logf("File exist %v", installLocation+installFile)
	} else {
		t.Logf("File does not exist %v", installLocation+installFile)
		t.Error("Missing file")
	}

	lib.RenameFile(installLocation+installFile, installLocation+installVersion+version)

	if exist := checkFileExist(installLocation + installVersion + version); exist {
		t.Logf("New file exist %v", installLocation+installVersion+version)
	} else {
		t.Logf("New file does not exist %v", installLocation+installVersion+version)
		t.Error("Missing new file")
	}

	if exist := checkFileExist(installLocation + installFile); exist {
		t.Logf("Old file should not exist %v", installLocation+installFile)
		t.Error("Did not rename file")
	} else {
		t.Logf("Old file does not exist %v", installLocation+installFile)
	}

	cleanUp(installLocation)
}

// TestRemoveFiles : Create a file, check file exist,
// remove file, check file does not exist
func TestRemoveFiles(t *testing.T) {

	installFile := "terraform"
	installPath := "/.terraform.versions_test/"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	createFile(installLocation + installFile)

	if exist := checkFileExist(installLocation + installFile); exist {
		t.Logf("File exist %v", installLocation+installFile)
	} else {
		t.Logf("File does not exist %v", installLocation+installFile)
		t.Error("Missing file")
	}

	lib.RemoveFiles(installLocation + installFile)

	if exist := checkFileExist(installLocation + installFile); exist {
		t.Logf("Old file should not exist %v", installLocation+installFile)
		t.Error("Did not remove file")
	} else {
		t.Logf("Old file does not exist %v", installLocation+installFile)
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
	installLocation := usr.HomeDir + installPath

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
	installLocation := usr.HomeDir + installPath

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
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	test_array := []string{"0.1.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	errWrite := lib.WriteLines(test_array, installLocation+recentFile)

	if errWrite != nil {
		t.Logf("Write should work %v (unexpected)", errWrite)
		log.Fatal(errWrite)
	} else {

		var (
			file             *os.File
			part             []byte
			prefix           bool
			errOpen, errRead error
			lines            []string
		)
		if file, errOpen = os.Open(installLocation + recentFile); errOpen != nil {
			log.Fatal(errOpen)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		buffer := bytes.NewBuffer(make([]byte, 0))
		for {
			if part, prefix, errRead = reader.ReadLine(); errRead != nil {
				break
			}
			buffer.Write(part)
			if !prefix {
				lines = append(lines, buffer.String())
				buffer.Reset()
			}
		}
		if errRead == io.EOF {
			errRead = nil
		}

		if errRead != nil {
			log.Fatalf("Error: %s\n", errRead)
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				log.Fatalf("Write to file is not invalid: %s\n", line)
				break
			}
		}

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
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	test_array := []string{"0.0.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	var (
		file      *os.File
		errCreate error
	)

	if file, errCreate = os.Create(installLocation + recentFile); errCreate != nil {
		log.Fatalf("Error: %s\n", errCreate)
	}
	defer file.Close()

	for _, item := range test_array {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			log.Fatalf("Error: %s\n", err)
			break
		}
	}

	lines, errRead := lib.ReadLines(installLocation + recentFile)

	if errRead != nil {
		log.Fatalf("Error: %s\n", errRead)
	}

	for _, line := range lines {
		if !semverRegex.MatchString(line) {
			log.Fatalf("Write to file is not invalid: %s\n", line)
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
	installLocation := usr.HomeDir + installPath

	test_dir := current.Format("2006-01-02")
	t.Logf("Create test dir: %v \n", test_dir)

	createDirIfNotExist(installLocation)

	createDirIfNotExist(installLocation + "/" + test_dir)

	empty := lib.IsDirEmpty(installLocation + "/" + test_dir)

	t.Logf("Expected directory to be empty %v [expected]", installLocation+"/"+test_dir)

	if empty == true {
		t.Logf("Directory empty")
	} else {
		t.Error("Directory not empty")
	}

	cleanUp(installLocation + "/" + test_dir)

	cleanUp(installLocation)

}

// TestCheckDirHasTGBin : create tg file in directory, check if exist
func TestCheckDirHasTFBin(t *testing.T) {

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	installPath := "/.terraform.versions_test/"
	installFile := "terraform"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	createFile(installLocation + installFile + "_" + goos + "_" + goarch)

	empty := lib.CheckDirHasTGBin(installLocation, installFile)

	t.Logf("Expected directory to have tf file %v [expected]", installLocation+installFile+"_"+goos+"_"+goarch)

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
	installFile := "terraform"

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	installLocation := usr.HomeDir + installPath

	createDirIfNotExist(installLocation)

	createFile(installLocation + "/" + installFile)

	path := lib.Path(installLocation + "/" + installFile)

	t.Logf("Path created %s\n", installLocation+installFile)
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
