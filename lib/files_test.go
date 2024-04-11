package lib

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
)

// TestRenameFile : Create a file, check filename exist,
// rename file, check new filename exit
func TestRenameFile(t *testing.T) {
	installFile := ConvertExecutableExt("terraform")
	installVersion := "terraform_"
	installPath := "/.terraform.versions_test/"
	version := "0.0.7"

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)

	createFile(installFilePath)

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("File exist %v", installFilePath)
	} else {
		t.Logf("File does not exist %v", installFilePath)
		t.Error("Missing file")
	}

	installVersionFilePath := ConvertExecutableExt(filepath.Join(installLocation, installVersion+version))

	RenameFile(installFilePath, installVersionFilePath)

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
	installFile := ConvertExecutableExt("terraform")
	installPath := "/.terraform.versions_test/"

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)

	createFile(installFilePath)

	if exist := checkFileExist(installFilePath); exist {
		t.Logf("File exist %v", installFilePath)
	} else {
		t.Logf("File does not exist %v", installFilePath)
		t.Error("Missing file")
	}

	RemoveFiles(installFilePath)

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

	t.Logf("Absolute Path: %q", absPath)

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	files, errUnzip := Unzip(absPath, installLocation)

	if errUnzip != nil {
		t.Errorf("Unable to unzip %q file: %v", absPath, errUnzip)
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

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	cleanUp(installLocation)

	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		t.Logf("Directory should not exist %v (expected)", installLocation)
	} else {
		t.Logf("Directory already exist %v (unexpected)", installLocation)
		t.Error("Directory should not exist")
	}

	createDirIfNotExist(installLocation)
	t.Logf("Creating directory %v", installLocation)

	if _, err := os.Stat(installLocation); err == nil {
		t.Logf("Directory exist %v (expected)", installLocation)
	} else {
		t.Logf("Directory should exist %v (unexpected)", installLocation)
		t.Error("Directory should exist")
	}

	cleanUp(installLocation)
}

// TestWriteLines : write to file, check readline to verify
func TestWriteLines(t *testing.T) {
	installPath := "/.terraform.versions_test/"
	recentFile := "RECENT"
	semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}(-\w+\d*)?\z`)
	//semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	recentFilePath := filepath.Join(installLocation, recentFile)
	test_array := []string{"0.1.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	errWrite := WriteLines(test_array, recentFilePath)

	if errWrite != nil {
		t.Logf("Write should work %v (unexpected)", errWrite)
		t.Error(errWrite)
	} else {
		var (
			file             *os.File
			part             []byte
			prefix           bool
			errOpen, errRead error
			lines            []string
		)
		if file, errOpen = os.Open(recentFilePath); errOpen != nil {
			t.Error(errOpen)
		}

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
			t.Errorf("Error: %s", errRead)
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				t.Errorf("Write to file is not invalid: %s", line)
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

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	recentFilePath := filepath.Join(installLocation, recentFile)
	test_array := []string{"0.0.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	var (
		file      *os.File
		errCreate error
	)

	if file, errCreate = os.Create(recentFilePath); errCreate != nil {
		t.Errorf("Error: %s", errCreate)
	}

	for _, item := range test_array {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			t.Errorf("Error: %s", err)
			break
		}
	}

	lines, errRead := ReadLines(recentFilePath)

	if errRead != nil {
		t.Errorf("Error: %s", errRead)
	}

	for _, line := range lines {
		if !semverRegex.MatchString(line) {
			t.Errorf("Write to file is not invalid: %s", line)
		}
	}

	file.Close()
	t.Log("Read versions exist (expected)")

	cleanUp(installLocation)
}

// TestIsDirEmpty : create empty directory, check if empty
func TestIsDirEmpty(t *testing.T) {
	current := time.Now()

	installPath := "/.terraform.versions_test/"

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	test_dir := current.Format("2006-01-02")
	test_dir_path := filepath.Join(installLocation, test_dir)
	t.Logf("Create test dir: %v \n", test_dir)

	createDirIfNotExist(installLocation)

	createDirIfNotExist(test_dir_path)

	empty := IsDirEmpty(test_dir_path)

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

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, installFilePrefix+"_"+goos+"_"+goarch))
	createFile(installFileVersionPath)

	empty := CheckDirHasTGBin(installLocation, installFilePrefix)

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
	installFile := ConvertExecutableExt("terraform")

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	installFilePath := filepath.Join(installLocation, installFile)
	createFile(installFilePath)

	path := Path(installFilePath)

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

	fileName := GetFileName(fileNameWithExt)

	if fileName == "file" {
		t.Logf("File removed extension (expected)")
	} else {
		t.Error("File did not remove extension (unexpected)")
	}
}

// TestConvertExecutableExt : convert executable binary with extension
func TestConvertExecutableExt(t *testing.T) {
	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}

	installPath := "/.terraform.versions_test/"
	test_array := []string{
		"terraform",
		"terraform.exe",
		filepath.Join(homedir, installPath, "terraform"),
		filepath.Join(homedir, installPath, "terraform.exe"),
	}

	for _, fpath := range test_array {
		fpathExt := ConvertExecutableExt((fpath))
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
