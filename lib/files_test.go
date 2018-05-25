package lib_test

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

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

// TestCreateDirIfNotExist : Create a directory, check directory exist,
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
