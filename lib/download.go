package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFromURL : Downloads the binary from the source url
func DownloadFromURL(installLocation string, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)
	fmt.Println("Downloading ...")

	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		//Sometimes hashicorp terraform file names are not consistent
		//For example 0.12.0-alpha4 naming convention in the release repo is not consistent
		return "", fmt.Errorf("Unable to download from %s\nPlease download manually from https://releases.hashicorp.com/terraform/", url)
	}

	zipFile := filepath.Join(installLocation, fileName)
	output, err := os.Create(zipFile)
	if err != nil {
		fmt.Println("Error while creating", zipFile, "-", err)
		return "", err
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}

	fmt.Println(n, "bytes downloaded.")
	return zipFile, nil
}
