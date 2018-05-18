package main

/*
* Version 0.0.1
* Compatible with Mac OS X ONLY
 */

/*** WORKFLOW ***/
/*
* 1- Check if user has sudo permission
* 2- Ask password to run sudo commands
* 3- Create /usr/local/terraform directory if does not exist
* 4- Download zip file from url to /usr/local/terraform
* 5- Unzip the file to /usr/local/terraform
* 6- Rename the file from `terraform` to `terraform_version`
* 7- Remove the downloaded zip file
* 8- Read the existing symlink for terraform (Check if it's a homebrew symlink)
* 9- Remove that symlink (Check if it's a homebrew symlink)
* 10- Create new symlink to binary  `terraform_version`
 */

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	cmd "github.com/warren-veerasingam/lib"
)

type tfVersionList struct {
	tflist []string
}

const (
	hashiURL = "https://releases.hashicorp.com/terraform/"
	//installLocation = "/usr/local/terraform/"
	//installLocation = "~/.terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	binLocation    = "/usr/local/bin/terraform"
)

func main() {

	usr, err1 := user.Current()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println(usr.HomeDir)

	installLocation := usr.HomeDir + "/.terraform/"

	/* check if terraform is already installed */
	cmd := cmd.NewCommand("terraform")
	next := cmd.Find()
	installedPath := binLocation
	existed := false
	for path := next(); len(path) > 0; path = next() {
		log.Printf("Found installation path: %v", path)
		installedPath = path
		existed = true
	}
	if !existed {
		installedPath = binLocation
		log.Printf("Installation path created: %v", installedPath)
	}
	/* 3- Create /usr/local/terraform directory if does not exist*/
	CreateDirIfNotExist(installLocation)

	resp, errURL := http.Get(hashiURL)
	if errURL != nil {
		log.Printf("Error getting url: %v", errURL)
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Printf("Error reading body: %v", errBody)
		return
	}

	bodyString := string(body)
	result := strings.Split(bodyString, "\n")

	var tfVersionList tfVersionList

	for i := range result {
		//getting versions from body
		r, _ := regexp.Compile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)

		if r.MatchString(result[i]) {
			str := r.FindString(result[i])
			trimstr := strings.Trim(str, "/")
			tfVersionList.tflist = append(tfVersionList.tflist, trimstr)
		}
	}

	prompt := promptui.Select{
		Label: "Select Version",
		Items: tfVersionList.tflist,
	}

	_, version, errPrompt := prompt.Run()

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		return
	}

	log.Printf("You picked %q\n", version)

	/* check if version exist locally*/
	fileExist := CheckFileExist(installLocation + installVersion + version)

	if fileExist {
		removeSymlink(binLocation)
		CreateSymlink(installLocation+installVersion+version, binLocation)
		log.Printf("Exiting early")
		os.Exit(0)
	}

	log.Printf("Still working")
	url := hashiURL + version + "/terraform_" + version + "_darwin_amd64.zip"

	zipFile, _ := DownloadFromURL(installLocation, url)

	log.Printf("ZipFile: " + zipFile)

	files, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		log.Fatal(errUnzip)
	}

	log.Printf("Unzipped:\n" + strings.Join(files, "\n"))

	RenameFile(installLocation+installFile, installLocation+installVersion+version)
	RemoveFiles(installLocation + installVersion + version + "_darwin_amd64.zip")
	//ReadlinkI(binLocation)
	//readLink(binLocation)
	removeSymlink(binLocation)
	CreateSymlink(installLocation+installVersion+version, binLocation)

}

// DownloadFromURL : Downloads the binary from the source url
func DownloadFromURL(installLocation string, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(installLocation + fileName)
	if err != nil {
		fmt.Println("Error while creating", installLocation+fileName, "-", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}

	fmt.Println(n, "bytes downloaded.")
	return installLocation + fileName, nil
}

//RenameFile : rename file name
func RenameFile(src string, dest string) {

	err := os.Rename(src, dest)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func getUser() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello " + user.Name)
	fmt.Println("Hello GroupId" + user.Gid)
	fmt.Println("====")
	fmt.Println("Id: " + user.Uid)
	fmt.Println("Username: " + user.Username)
	fmt.Println("Home Dir: " + user.HomeDir)
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	fmt.Println("src: " + src)

	fmt.Println("dest: " + dest)
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
			os.MkdirAll(fpath, os.ModePerm)

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

// RemoveFiles : remove file
func RemoveFiles(src string) {

	fmt.Println(src)

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

func ReadlinkI(path string) {
	fmt.Println("PATH")
	// fmt.Println(path)
	// exit_code := 0
	// defer os.Exit(exit_code)
	ln, err := os.Readlink(path)
	if err != nil {
		fmt.Println("[ERR]", err)
		//exit_code = 1
		return
	}
	fmt.Println("[FOUND]", ln)
}

func readLink(path string) {
	// exit_code := 0
	// defer os.Exit(exit_code)
	ln, err := filepath.EvalSymlinks(path)
	if err != nil {
		fmt.Println("[ERR]", err)
		//exit_code = 1
		return
	}
	fmt.Println("[FOUND]", ln)
}

func CreateSymlink(cwd string, dir string) error {

	if err := os.Symlink(cwd, dir); err != nil {
		return err
	}
	return nil
}

func removeSymlink(symlinkPath string) error {

	if _, err := os.Lstat(symlinkPath); err != nil {
		return err
	}
	os.Remove(symlinkPath)
	return nil
}

func CreateDirIfNotExist(dir string) {
	log.Printf("entering here")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Creating directory for teraform: %v", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("Unable to create directory for teraform: %v", dir)
			panic(err)
		}
	}
}

func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}
