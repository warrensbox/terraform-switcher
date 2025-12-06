//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
)

// RenameFile : rename file name
func RenameFile(src string, dest string) {
	logger.Debugf("Renaming file %q to %q", src, dest)
	err := os.Rename(src, dest)
	if err != nil {
		logger.Fatal(err)
	}
}

// RemoveFiles : remove file
func RemoveFiles(src string) {
	// Keep both identical functions for backward compatibility
	// FIXME: need to plan deprecation of either of them
	// 09-Mar-2025
	removeFiles(src)
}

// CheckFileExist : check if file exist in directory
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
// fileToUnzip (parameter 3) specifies the file within the zipfile to be extracted.
// This is optional and default to "terraform"
func Unzip(src string, dest string, fileToUnzipSlice ...string) ([]string, error) {
	logger.Debugf("Unzipping file %q", src)

	// Handle old signature of method, where fileToUnzip did not exist
	legacyProduct := getLegacyProduct()
	fileToUnzip := legacyProduct.GetExecutableName()
	if len(fileToUnzipSlice) == 1 {
		fileToUnzip = fileToUnzipSlice[0]
	} else if len(fileToUnzipSlice) > 1 {
		logger.Fatal("Too many args passed to Unzip")
	}

	var filenames []string

	reader, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer reader.Close()
	destination, err := filepath.Abs(dest)
	if err != nil {
		logger.Fatalf("Could not open destination: %v", err)
	}
	var unzipWaitGroup sync.WaitGroup
	for _, f := range reader.File {
		// Only extract the main binary
		// from the archive, ignoring LICENSE and other files
		if f.Name != ConvertExecutableExt(fileToUnzip) {
			continue
		}

		unzipWaitGroup.Add(1)
		unzipErr := unzipFile(f, destination, &unzipWaitGroup)
		if unzipErr != nil {
			return nil, fmt.Errorf("Error unzipping: %v", unzipErr)
		}
		// nolint:gosec // The "G305: File traversal when extracting zip/tar archive" is handled by unzipFile()
		filenames = append(filenames, filepath.Join(destination, f.Name))
	}
	logger.Debug("Waiting for deferred functions")
	unzipWaitGroup.Wait()

	if len(filenames) < 1 {
		logger.Fatalf("Could not find %s file in release archive to unzip", fileToUnzip)
	} else if len(filenames) > 1 {
		logger.Fatal("Extracted more files than expected in release archive")
	}

	return filenames, nil
}

// createDirIfNotExist : create directory if directory does not exist
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Infof("Creating %q directory", dir)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			logger.Panicf("Unable to create %q directory: %v", dir, err)
		}
	}
}

// WriteLines : writes into file
//
// Deprecated: This method has been deprecated and will be removed in v2.0.0
func WriteLines(lines []string, path string) (err error) {
	var file *os.File

	if file, err = os.Create(path); err != nil {
		return err
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			logger.Error(err)
			break
		}
	}

	return nil
}

// ReadLines : Read a whole file into the memory and store it as array of lines
//
// Deprecated: This method has been deprecated and will be removed in v2.0.0
func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// IsDirEmpty : check if directory is empty (TODO UNIT TEST)
func IsDirEmpty(dir string) bool {
	exist := false

	f, err := os.Open(dir)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		exist = true
	}
	return exist // Either not empty or error, suits both cases
}

// CheckDirHasTGBin : // check binary exist (TODO UNIT TEST)
func CheckDirHasTGBin(dir, prefix string) bool {
	exist := false

	files, err := os.ReadDir(dir)
	if err != nil {
		logger.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
			exist = true
		}
	}
	return exist
}

// CheckDirExist : check if directory exist
// dir=path to file
// return bool
func CheckDirExist(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Debugf("Directory %q doesn't exist", dir)
		return false
	}

	return true
}

// CheckIsDir: check if is directory
// dir=path to file
// return bool
func CheckIsDir(dir string) bool {
	fi, err := os.Stat(dir)

	if err != nil {
		logger.Debugf("Error checking %q: %v", dir, err)
		return false
	} else if !fi.IsDir() {
		logger.Debugf("The %q is not a directory", dir)
		return false
	}

	return true
}

// CheckDirIsReadable : check if directory is readable
func CheckDirIsReadable(dir string) bool {
	if !CheckDirExist(dir) || !CheckIsDir(dir) {
		return false
	}

	_, err := os.ReadDir(dir)
	if err != nil {
		logger.Debugf("Failed to read directory %q: %v", dir, err)
		return false
	}

	return true
}

// Path : returns path of directory
// value=path to file
func Path(value string) string {
	return filepath.Dir(value)
}

// GetFileName : remove file ext.  .tfswitch.config returns .tfswitch
func GetFileName(configfile string) string {
	return strings.TrimSuffix(configfile, filepath.Ext(configfile))
}

// GetCurrentDirectory : return the current directory
func GetCurrentDirectory() string {
	dir, err := os.Getwd() // get current directory
	if err != nil {
		logger.Fatalf("Failed to get current directory: %v", err)
	}
	return dir
}

// GetHomeDirectory : return the user's home directory
func GetHomeDirectory() string {
	homedir, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Failed to get user's home directory: %v", err)
	}
	return homedir
}

func unzipFile(f *zip.File, destination string, wg *sync.WaitGroup) error {
	defer wg.Done()
	// 1. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name) // nolint:gosec // The "G305: File traversal when extracting zip/tar archive" is handled below
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("Invalid file path: %q", filePath)
	}

	// 2. Create directory tree
	if f.FileInfo().IsDir() {
		logger.Debugf("Extracting directory %q", filePath)
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 3. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	defer func(destinationFile *os.File) {
		logger.Debugf("Closing destination file handler %q", destinationFile.Name())
		_ = destinationFile.Close()
	}(destinationFile)
	if err != nil {
		return err
	}

	// 4. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	defer func(zippedFile io.ReadCloser) {
		logger.Debugf("Closing zipped file handler %q", f.Name)
		_ = zippedFile.Close()
	}(zippedFile)
	if err != nil {
		return err
	}

	logger.Debugf("Extracting file %q to %q", f.Name, destinationFile.Name())
	// Prevent the "G110: Potential DoS vulnerability via decompression bomb (gosec)"
	totalCopied := int64(0)
	maxSize := int64(1024 * 1024 * 1024) // 1 GB
	for {
		copied, err := io.CopyN(destinationFile, zippedFile, 1024*1024)
		totalCopied += copied
		if totalCopied%(10*1024*1024) == 0 { // Print stats every 10 MB
			logger.Debugf("Size copied so far: %3.d MB\r", totalCopied/1024/1024)
		}
		if err != nil {
			if err == io.EOF {
				logger.Debugf("Total size copied: %4.d MB\r", totalCopied/1024/1024)
				break
			}
			return err
		}
		if totalCopied > maxSize {
			return fmt.Errorf("file %q is too large (> %d MB)", f.Name, maxSize/1024/1024)
		}
	}
	return nil
}
