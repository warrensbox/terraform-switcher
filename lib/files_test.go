package lib_test

import (
	"log"
	"os/user"
	"testing"

	"github.com/warren-veerasingam/terraform-switcher/lib"
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
