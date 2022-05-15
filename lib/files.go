package lib

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// RenameFile : rename file name
func RenameFile(src string, dest string) {
	err := os.Rename(src, dest)
	if err != nil {
		fmt.Println(err)
		return
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
	return err == nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return nil, err
			}

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}

//CreateDirIfNotExist : create directory if directory does not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Creating directory for terraform binary at: %v\n", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terraform binary at: %v", dir)
			panic(err)
		}
	}
}

//WriteLines : writes into file
func WriteLines(releases []*Release, path string) error {
	var (
		file *os.File
	)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range releases {

		b, err := json.Marshal(item)
		if err != nil {
			fmt.Println(err)
			break
		}
		_, err = file.Write(b)
		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = file.WriteString("\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return nil
}

// ReadLines : Read a whole file into the memory and store it as array of lines
func ReadLines(path string) (lines []*Release, err error) {
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
			var release *Release
			if err := json.Unmarshal(buffer.Bytes(), &release); err != nil {
				return nil, fmt.Errorf("%s: %s", err, buffer.Bytes())
			}
			if !ValidVersionFormat(release.Version) {
				return nil, nil
				return nil, fmt.Errorf("Invalid version parsed: %s", release.Version)
			}
			lines = append(lines, release)
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

//IsDirEmpty : check if directory is empty (TODO UNIT TEST)
func IsDirEmpty(name string) bool {

	exist := false

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		exist = true
	}
	return exist // Either not empty or error, suits both cases
}

//CheckDirHasTGBin : // check binary exist (TODO UNIT TEST)
func CheckDirHasTGBin(dir, prefix string) bool {

	exist := false
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
			exist = true
		}
	}
	return exist
}

//CheckDirExist : check if directory exist
//dir=path to file
//return bool
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
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}
	return dir
}

// GetHomeDirectory : return the home directory
func GetHomeDirectory() string {

	homedir, errHome := homedir.Dir()
	if errHome != nil {
		log.Printf("Failed to get home directory %v\n", errHome)
		os.Exit(1)
	}

	return homedir
}
