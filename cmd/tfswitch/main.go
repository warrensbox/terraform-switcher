package main

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
)

type tfVersion struct {
	version string
	url     string
}

type tfList struct {
	tflist []tfVersion
}

type tfListA struct {
	tflist []string
}

func main() {

	resp, err := http.Get("https://releases.hashicorp.com/terraform/")
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle err
		log.Printf("Error reading body: %v", err)
		//http.Error(resp, "can't read body", http.StatusBadRequest)
		//test commit
		return
	}

	bodyString := string(body)
	//fmt.Println(bodyString)

	// r, _ := regexp.Compile("terraform")

	// fmt.Println(r.FindString(bodyString))

	//scanner := bufio.NewScanner(bodyString)
	result := strings.Split(bodyString, "\n")

	//var tfList tfList

	// Display all elements.

	var tfListA tfListA

	var tfList tfList

	for i := range result {

		//r, _ := regexp.Compile("terraform")
		r, _ := regexp.Compile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)
		//var re = regexp.MustCompile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)

		//fmt.Println(r.FindString(result[i]))
		//fmt.Println("u")
		//fmt.Println(r.MatchString("terraform"))

		var tfVersion tfVersion

		if r.MatchString(result[i]) {

			//fmt.Println(result[i])
			//fmt.Println(r.FindString(result[i]))
			str := r.FindString(result[i])
			trimstr := strings.Trim(str, "/")
			//fmt.Printf(trimstr)

			tfVersion.version = trimstr

			tfVersion.url = "https://releases.hashicorp.com/terraform/" + trimstr + "/terraform_" + trimstr + "_darwin_amd64.zip"

			tfListA.tflist = append(tfListA.tflist, trimstr)

			tfList.tflist = append(tfList.tflist, tfVersion)

		}

		//tfList = append(tfList[i].

	}

	//fmt.Println(tfListA)
	//fmt.Println(tfList)

	prompt := promptui.Select{
		Label: "Select Day",
		Items: tfListA.tflist,
		// Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		// 	"Saturday", "Sunday"},
	}

	_, result1, err1 := prompt.Run()

	if err1 != nil {
		fmt.Printf("Prompt failed %v\n", err1)
		return
	}
	//fmt.Printf(Items)
	fmt.Printf("You choose %q\n", result1)
	url := "https://releases.hashicorp.com/terraform/" + result1 + "/terraform_" + result1 + "_darwin_amd64.zip"
	zipFile, _ := downloadFromUrl(url)
	//getUser()

	fmt.Println("zipFile: " + zipFile)

	files, err2 := Unzip(zipFile, "/usr/local/terraform")
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	moveFile("/usr/local/terraform/terraform", "/usr/local/terraform/terraform"+"_"+result1)
	removeFiles("/usr/local/terraform/terraform_" + result1 + "_darwin_amd64.zip")
	ReadlinkI("/usr/local/bin/terraform")
	readLink("/usr/local/bin/terraform")
	removeSymlink("/usr/local/bin/terraform")
	CreateSymlink("/usr/local/terraform/terraform"+"_"+result1, "/usr/local/bin/terraform")

}

func downloadFromUrl(url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create("/usr/local/terraform/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", "/usr/local/terraform/"+fileName, "-", err)
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
	return "/usr/local/terraform/" + fileName, nil
}

func moveFile(src string, dest string) {

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

func removeFiles(src string) {

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
