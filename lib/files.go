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
	files, err := filepath.Glob(src)
	if err != nil {

		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

// CheckFileExist : check if file exist in directory
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	logger.Debugf("Unzipping file %q", src)

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
	var wg sync.WaitGroup
	for _, f := range reader.File {
		wg.Add(1)
		unzipErr := unzipFile(f, destination, &wg)
		if unzipErr != nil {
			logger.Fatalf("Error unzipping %v", unzipErr)
		} else {
			filenames = append(filenames, filepath.Join(destination, f.Name))
		}
	}
	return filenames, nil
}

// createDirIfNotExist : create directory if directory does not exist
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Infof("Creating directory for terraform binary at %q", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Panicf("Unable to create %q directory for terraform: %v", dir, err)
		}
	}
}

// WriteLines : writes into file
func WriteLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

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
func IsDirEmpty(name string) bool {

	exist := false

	f, err := os.Open(name)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
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
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
			res = append(res, filepath.Join(dir, f.Name()))
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

	dir, err := os.Getwd() //get current directory
	if err != nil {
		logger.Fatalf("Failed to get current directory %v", err)
	}
	return dir
}

func unzipFile(f *zip.File, destination string, wg *sync.WaitGroup) error {
	defer wg.Done()
	// 1. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
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
		_ = destinationFile.Close()
	}(destinationFile)
	if err != nil {
		return err
	}

	// 4. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	defer func(zippedFile io.ReadCloser) {
		_ = zippedFile.Close()
	}(zippedFile)
	if err != nil {
		return err
	}

	logger.Debugf("Extracting File %q", destinationFile.Name())
	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	logger.Debugf("Closing destination file handler %q", destinationFile.Name())
	_ = destinationFile.Close()
	logger.Debugf("Closing zipped file handler %q", f.Name)
	_ = zippedFile.Close()
	return nil
}
