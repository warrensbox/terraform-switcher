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
	logger = InitLogger("DEBUG")
	installPath := "/.terraform.versions_test/"
	var pathToTestFile string
	switch runtime.GOOS {
	case windows:
		pathToTestFile = "../test-data/test-data_windows.zip"
	default:
		pathToTestFile = "../test-data/test-data.zip"
	}
	absPath, err := filepath.Abs(pathToTestFile)
	if err != nil {
		t.Error(err)
	}

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

	if fileCount := len(files); fileCount != 1 {
		t.Errorf("Expected extracted files size is %d, expected 1", fileCount)
	}

	// Ensure terraform file exists
	terraformFile := filepath.Join(installLocation, ConvertExecutableExt("terraform"))
	if terraformFileExists := checkFileExist(terraformFile); !terraformFileExists {
		t.Errorf("File does not exist %v", terraformFile)
	}

	terraformFileContent, err := os.ReadFile(terraformFile)
	if err != nil {
		t.Error(err)
	}
	// Ensure terraform file contains test content from
	if string(terraformFileContent) != "TerraformBinaryContent\n" {
		t.Errorf("Terraform test file content does not match expected value: %s", string(terraformFileContent))
	}

	// Ensure README and LICENSE files don't exist
	nonExistentFiles := []string{"README", "LICENSE"}
	for _, fileName := range nonExistentFiles {
		filePath := filepath.Join(installLocation, fileName)
		if fileExists := checkFileExist(filePath); fileExists {
			t.Errorf("Zip archive file should not exist: %v", fileExists)
		}
	}

	cleanUp(installLocation)
}

// TestUnzip_with_file_to_unzip : Test Unzip method with fileToUnzip argument
func TestUnzip_with_file_to_unzip(t *testing.T) {
	installPath := "/.terraform.versions_test/"
	var pathToTestFile string
	var expectedFilename string
	switch runtime.GOOS {
	case windows:
		pathToTestFile = "../test-data/test-data_windows.zip"
		expectedFilename = "another_file.exe"
	default:
		pathToTestFile = "../test-data/test-data.zip"
		expectedFilename = "another_file"
	}

	absPath, err := filepath.Abs(pathToTestFile)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Absolute Path: %q", absPath)

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	files, errUnzip := Unzip(absPath, installLocation, "another_file")

	if errUnzip != nil {
		t.Errorf("Unable to unzip %q file: %v", absPath, errUnzip)
	}

	if fileCount := len(files); fileCount != 1 {
		t.Errorf("Expected extracted files size is %d, expected 1", fileCount)
	}

	// Ensure terraform file exists
	expectedExtractFile := filepath.Join(installLocation, expectedFilename)
	if terraformFileExists := checkFileExist(expectedExtractFile); !terraformFileExists {
		t.Errorf("File does not exist %v", expectedExtractFile)
	}

	expectedFileContent, err := os.ReadFile(expectedExtractFile)
	if err != nil {
		t.Error(err)
	}
	// Ensure terraform file contains test content from
	if !strings.Contains(string(expectedFileContent), "This is another executable") {
		t.Errorf("Extract test file (%q) content does not match expected value: %s", expectedExtractFile, string(expectedFileContent))
	}

	// Ensure README and LICENSE files don't exist
	nonExistentFiles := []string{"terraform", "terraform.exe", "LICENSE"}
	for _, fileName := range nonExistentFiles {
		filePath := filepath.Join(installLocation, fileName)
		if fileExists := checkFileExist(filePath); fileExists {
			t.Errorf("Zip archive file should not exist: %v", fileExists)
		}
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
	// semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

	homedir, errCurr := homedir.Dir()
	if errCurr != nil {
		t.Error(errCurr)
	}
	installLocation := filepath.Join(homedir, installPath)

	createDirIfNotExist(installLocation)

	recentFilePath := filepath.Join(installLocation, recentFile)
	testArray := []string{"0.1.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	errWrite := WriteLines(testArray, recentFilePath)

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
	testArray := []string{"0.0.1", "0.0.2", "0.0.3", "0.12.0-rc1", "0.12.0-beta1"}

	var (
		file      *os.File
		errCreate error
	)

	if file, errCreate = os.Create(recentFilePath); errCreate != nil {
		t.Errorf("Error: %s", errCreate)
	}

	for _, item := range testArray {
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

	testDir := current.Format("2006-01-02")
	testDirPath := filepath.Join(installLocation, testDir)
	t.Logf("Create test dir: %v \n", testDir)

	createDirIfNotExist(installLocation)

	createDirIfNotExist(testDirPath)

	empty := IsDirEmpty(testDirPath)

	t.Logf("Expected directory to be empty %v [expected]", testDirPath)

	if empty == true {
		t.Logf("Directory empty")
	} else {
		t.Error("Directory not empty")
	}

	cleanUp(testDirPath)
	cleanUp(installLocation)
}

// TestCheckDirIsReadable : test dir readability check
func TestCheckDirIsReadable(t *testing.T) {
	var path string

	// Dir must be readable
	path = "../test-data/"
	if !CheckDirIsReadable(path) {
		t.Fatalf("Directory %q is not readable", path)
	}
	t.Logf("Directory %q is readable (expected)", path)

	// Dir must not exist
	path = "../test-data/non-existent-directory"
	if CheckDirIsReadable(path) {
		t.Fatalf("Directory %q must not exist", path)
	}
	t.Logf("Directory %q does not exist (expected)", path)

	// Path must not be a directory
	path = "../test-data/is-plain-file"
	if CheckDirIsReadable(path) {
		t.Fatalf("The %q must not be a directory", path)
	}
	t.Logf("The %q is not a directory (expected)", path)

	// Creating dir on Windows produces `drwxrwxrwx` permissions no matter
	// what perms are requested in `os.Mkdir()`, and even `os.Chmod()` doesn't help
	// So just skip this bit of the test on Windows. Windows is tricky ¯\_(ツ)_/¯
	// Informational reference: https://github.com/golang/go/issues/65377
	if runtime.GOOS != windows {
		// Dir must have no read permissions
		path = "../test-data/directory-without-read-permission"
		// 0200 is -w------- (user write only)
		if err := os.Mkdir(path, os.FileMode(0o200)); err != nil {
			t.Fatalf("Unexpected failure creating test directory %q: %v", path, err)
		}
		defer func() {
			// Add enough permissions to allow cleanup (user read/write/execute)
			if err := os.Chmod(path, 0o700); err != nil {
				// Don't fail, just log
				t.Logf("(cleanup) Unable to restore permissions on test directory %q: %v", path, err)
			}
			cleanUp(path)
		}()
		if CheckDirIsReadable(path) {
			t.Fatalf("Directory %q must not have read permission", path)
		}
		t.Logf("Directory %q has no read permission (expected)", path)
	}
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
	testArray := []string{
		"terraform",
		"terraform.exe",
		filepath.Join(homedir, installPath, "terraform"),
		filepath.Join(homedir, installPath, "terraform.exe"),
	}

	for _, fpath := range testArray {
		fpathExt := ConvertExecutableExt((fpath))
		outputMsg := fpath + " converted to " + fpathExt + " on " + runtime.GOOS

		switch runtime.GOOS {
		case windows:
			if filepath.Ext(fpathExt) != ".exe" {
				t.Errorf("%s (unexpected)", outputMsg)
				continue
			}

			if filepath.Ext(fpath) == ".exe" {
				if fpathExt != fpath {
					t.Errorf("%s (unexpected)", outputMsg)
				} else {
					t.Logf("%s (expected)", outputMsg)
				}
				continue
			}

			if fpathExt != fpath+".exe" {
				t.Errorf("%s (unexpected)", outputMsg)
				continue
			}

			t.Logf("%s (expected)", outputMsg)
		default:
			if fpath != fpathExt {
				t.Errorf("%s (unexpected)", outputMsg)
				continue
			}

			t.Logf("%s (expected)", outputMsg)
		}
	}
}
