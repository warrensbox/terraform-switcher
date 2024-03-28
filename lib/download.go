package lib

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	pubKeyId     = "72D7468F"
	pubKeyPrefix = "hashicorp_"
	pubKeySuffix = ".asc"
	pubKeyUri    = "https://www.hashicorp.com/.well-known/pgp-key.txt"
)

// DownloadFromURL : Downloads the terraform binary and its hash from the source url
func DownloadFromURL(installLocation string, mirrorURL string, tfversion string, versionPrefix string, goos string, goarch string) (string, error) {
	pubKeyFilename := filepath.Join(installLocation, "/", pubKeyPrefix+pubKeyId+pubKeySuffix)
	zipUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_" + goos + "_" + goarch + ".zip"
	hashUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_SHA256SUMS"
	hashSignatureUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_SHA256SUMS." + pubKeyId + ".sig"

	err := downloadPublicKey(installLocation, pubKeyFilename)
	if err != nil {
		log.Fatal("[Error]: Could not download public key file")
	}

	log.Println("Downloading ", zipUrl)
	zipFilePath, err := downloadFromURL(installLocation, zipUrl)
	if err != nil {
		log.Fatal("[Error]: Could not download zip file")
	}

	log.Println("Downloading ", hashUrl)
	hashFilePath, err := downloadFromURL(installLocation, hashUrl)
	if err != nil {
		log.Fatal("[Error]: Could not download hash file")
	}

	log.Println("Downloading ", hashSignatureUrl)
	hashSigFilePath, err := downloadFromURL(installLocation, hashSignatureUrl)
	if err != nil {
		log.Fatal("[Error]: Could not download hash signature file")
	}

	publicKeyFile, err := os.Open(pubKeyFilename)
	if err != nil {
		log.Fatal("[Error]: Could not open the public key.", pubKeyFilename)
	}

	signatureFile, err := os.Open(hashSigFilePath)
	if err != nil {
		log.Fatal("[Error]: Could not open the signature file.", hashSigFilePath)
	}

	targetFile, err := os.Open(zipFilePath)
	if err != nil {
		log.Fatal("[Error]: Could not open the terraform binary for checksum verification.", zipFilePath)
	}

	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		log.Fatal("[Error]: Could not open the terraform checksum file for signature verification.", hashFilePath)
	}
	verified := checkSignatureOfChecksums(publicKeyFile, hashFile, signatureFile)
	if !verified {
		return "", errors.New("signature of checksum files could not be verified")
	}
	match := checkChecksumMatches(hashFilePath, targetFile)
	if !match {
		return "", errors.New("checksums did not match")
	}
	return zipFilePath, err
}

func downloadFromURL(installLocation string, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	log.Printf("Downloading to: %s\n", filepath.Join(installLocation, "/", fileName))

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("[Error] : Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		//Sometimes hashicorp terraform file names are not consistent
		//For example 0.12.0-alpha4 naming convention in the release repo is not consistent
		log.Fatalf("[Error] : Unable to download from %s", url)
	}

	zipFile := filepath.Join(installLocation, fileName)
	output, err := os.Create(zipFile)
	if err != nil {
		log.Fatal("[Error] : Error while creating", zipFile, "-", err)
		return "", err
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Fatal("[Error] : Error while writing file", url, "-", err)
		return "", err
	}

	log.Println(n, "bytes downloaded")
	return zipFile, nil
}

func downloadPublicKey(installLocation string, targetFileName string) error {
	fmt.Println("Looking up public key file at ", targetFileName)
	publicKeyFileExists := FileExists(targetFileName)
	if !publicKeyFileExists {
		// Public key does not exist. Let's grab it from hashicorp
		pubKeyFile, errDl := downloadFromURL(installLocation, pubKeyUri)
		if errDl != nil {
			log.Fatal("[Error]: Error while fetching the public key file from ", pubKeyUri)
			return errDl
		}
		errRename := os.Rename(pubKeyFile, targetFileName)
		if errRename != nil {
			log.Fatal("[Error]: Error while renaming the public key file from ", pubKeyFile, " to ", targetFileName)
			return errRename
		}
	}
	return nil
}
