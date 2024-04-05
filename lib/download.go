package lib

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	PubKeyId     = "72D7468F"
	PubKeyPrefix = "hashicorp_"
	PubKeyUri    = "https://www.hashicorp.com/.well-known/pgp-key.txt"
)

const (
	pubKeySuffix = ".asc"
)

// DownloadFromURL : Downloads the terraform binary and its hash from the source url
func DownloadFromURL(installLocation string, mirrorURL string, tfversion string, versionPrefix string, goos string, goarch string) (string, error) {
	pubKeyFilename := filepath.Join(installLocation, "/", PubKeyPrefix+PubKeyId+pubKeySuffix)
	zipUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_" + goos + "_" + goarch + ".zip"
	hashUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_SHA256SUMS"
	hashSignatureUrl := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_SHA256SUMS." + PubKeyId + ".sig"

	err := downloadPublicKey(installLocation, pubKeyFilename)
	if err != nil {
		logger.Error("Could not download public key file")
		return "", err
	}

	logger.Infof("Downloading %q", zipUrl)
	zipFilePath, err := downloadFromURL(installLocation, zipUrl)
	if err != nil {
		logger.Error("Could not download zip file")
		return "", err
	}

	logger.Infof("Downloading %q", hashUrl)
	hashFilePath, err := downloadFromURL(installLocation, hashUrl)
	if err != nil {
		logger.Error("Could not download hash file")
		return "", err
	}

	logger.Infof("Downloading %q", hashSignatureUrl)
	hashSigFilePath, err := downloadFromURL(installLocation, hashSignatureUrl)
	if err != nil {
		logger.Error("Could not download hash signature file")
		return "", err
	}

	publicKeyFile, err := os.Open(pubKeyFilename)
	if err != nil {
		logger.Error("Could not open the public key")
		return "", err
	}

	signatureFile, err := os.Open(hashSigFilePath)
	if err != nil {
		logger.Error("Could not open the public key")
		return "", err
	}

	targetFile, err := os.Open(zipFilePath)
	if err != nil {
		logger.Error("Could not open the terraform binary for signature verification.")
		return "", err
	}

	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		logger.Error("Could not open the terraform binary for signature verification.")
		return "", err
	}

	var filesToCleanup []string
	filesToCleanup = append(filesToCleanup, hashFilePath)
	filesToCleanup = append(filesToCleanup, hashSigFilePath)

	verified := checkSignatureOfChecksums(publicKeyFile, hashFile, signatureFile)
	if !verified {
		cleanup(filesToCleanup)
		return "", errors.New("signature of checksum files could not be verified")
	}
	match := checkChecksumMatches(hashFilePath, targetFile)
	if !match {
		cleanup(filesToCleanup)
		return "", errors.New("checksums did not match")
	}
	cleanup(filesToCleanup)
	return zipFilePath, err
}

func downloadFromURL(installLocation string, url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	logger.Infof("Downloading to: %s", filepath.Join(installLocation, "/", fileName))

	response, err := http.Get(url)
	if err != nil {
		logger.Errorf("Error while downloading %q - %v", url, err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		//Sometimes hashicorp terraform file names are not consistent
		//For example 0.12.0-alpha4 naming convention in the release repo is not consistent
		return "", errors.New("Unable to download from " + url)
	}

	filePath := filepath.Join(installLocation, fileName)
	output, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("Error while creating %q: %v", filePath, err)
		return "", err
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		logger.Errorf("Error while downloading %s: %v", url, err)
		return "", err
	}

	logger.Info(n, "bytes downloaded")
	return filePath, nil
}

func downloadPublicKey(installLocation string, targetFileName string) error {
	logger.Debugf("Looking up public key file at %q", targetFileName)
	publicKeyFileExists := FileExists(targetFileName)
	if !publicKeyFileExists {
		// Public key does not exist. Let's grab it from hashicorp
		pubKeyFile, errDl := downloadFromURL(installLocation, PubKeyUri)
		if errDl != nil {
			logger.Error("Error while fetching the public key file from ", PubKeyUri)
			return errDl
		}
		errRename := os.Rename(pubKeyFile, targetFileName)
		if errRename != nil {
			logger.Error("Error while renaming the public key file from ", pubKeyFile, " to ", targetFileName)
			return errRename
		}
	}
	return nil
}

func cleanup(paths []string) {
	for _, path := range paths {
		logger.Infof("Deleting %q", path)
		err := os.Remove(path)
		if err != nil {
			logger.Error("Error deleting %q - %q", path, err)
		}
	}
}
