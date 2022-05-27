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
	fmt.Printf("Downloading to: %s\n", installLocation)

	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("[Error] : Error while downloading from %s\n%q", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		//Sometimes hashicorp terraform file names are not consistent
		//For example 0.12.0-alpha4 naming convention in the release repo is not consistent
		return "", fmt.Errorf("[Error] : Unable to download from %s", url)
	}

	zipFile := filepath.Join(installLocation, fileName)
	output, err := os.Create(zipFile)
	if err != nil {
		return "", fmt.Errorf("[Error] : Error while creating file: %s\n%q", zipFile, err)
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return "", fmt.Errorf("[Error] : Error while parsing response body from %s\n%q", url, err)
	}

	fmt.Println(n, "bytes downloaded")
	return zipFile, nil
}
