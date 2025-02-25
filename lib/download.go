package lib

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// DownloadFromURL : Downloads the terraform binary and its hash from the source url
func DownloadFromURL(installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch string) (string, error) {
	product := getLegacyProduct()
	return DownloadProductFromURL(product, installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch)
}

func DownloadProductFromURL(product Product, installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch string) (string, error) {
	var wg sync.WaitGroup
	defer wg.Done()
	zipUrl := mirrorURL + "/" + versionPrefix + tfversion + "_" + goos + "_" + goarch + ".zip"
	hashUrl := mirrorURL + "/" + versionPrefix + tfversion + "_SHA256SUMS"
	hashSignatureUrl := mirrorURL + "/" + versionPrefix + tfversion + "_SHA256SUMS." + product.GetShaSignatureSuffix()

	pubKeyFilename, err := downloadPublicKey(product, installLocation, &wg)
	if err != nil {
		logger.Error("Could not download public PGP key file")
		return "", err
	}

	logger.Infof("Downloading %q", zipUrl)
	zipFilePath, err := downloadFromURL(installLocation, zipUrl, &wg)
	if err != nil {
		logger.Error("Could not download zip file")
		return "", err
	}

	logger.Infof("Downloading %q", hashUrl)
	hashFilePath, err := downloadFromURL(installLocation, hashUrl, &wg)
	if err != nil {
		logger.Error("Could not download hash file")
		return "", err
	}

	logger.Infof("Downloading %q", hashSignatureUrl)
	hashSigFilePath, err := downloadFromURL(installLocation, hashSignatureUrl, &wg)
	if err != nil {
		logger.Error("Could not download hash signature file")
		return "", err
	}

	// // Wait for wait group, as the file downloads are required for the below functionality
	// wg.Wait()

	publicKeyFile, err := os.Open(pubKeyFilename)
	if err != nil {
		logger.Errorf("Could not open public key %q: %v", pubKeyFilename, err)
		return "", err
	}

	signatureFile, err := os.Open(hashSigFilePath)
	if err != nil {
		logger.Errorf("Could not open hash signature file %q: %v", hashSigFilePath, err)
		return "", err
	}

	targetFile, err := os.Open(zipFilePath)
	if err != nil {
		logger.Errorf("Could not open zip file %q: %v", zipFilePath, err)
		return "", err
	}

	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		logger.Errorf("Could not open hash file %q: %v", hashFilePath, err)
		return "", err
	}

	var filesToCleanup []string
	filesToCleanup = append(filesToCleanup, hashFilePath)
	filesToCleanup = append(filesToCleanup, hashSigFilePath)
	defer cleanup(filesToCleanup, &wg)

	verified := checkSignatureOfChecksums(publicKeyFile, hashFile, signatureFile)
	if !verified {
		return "", errors.New("Signature of checksum file could not be verified")
	}
	match := checkChecksumMatches(hashFilePath, targetFile)
	if !match {
		return "", errors.New("Checksums did not match")
	}
	return zipFilePath, err
}

func downloadFromURL(installLocation string, url string, wg *sync.WaitGroup) (string, error) {
	wg.Add(1)
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	logger.Infof("Downloading to %q", filepath.Join(installLocation, "/", fileName))

	response, err := http.Get(url)
	if err != nil {
		logger.Errorf("Error downloading %s: %v", url, err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// Sometimes hashicorp terraform file names are not consistent
		// For example 0.12.0-alpha4 naming convention in the release repo is not consistent
		return "", errors.New("Unable to download from " + url)
	}

	filePath := filepath.Join(installLocation, fileName)
	output, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("Error creating %q: %v", filePath, err)
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

func downloadPublicKey(product Product, installLocation string, wg *sync.WaitGroup) (string, error) {
	pubKeyFilePath := filepath.Join(installLocation, "/", product.GetId()+"_"+product.GetPublicKeyId()+pubKeySuffix)
	logger.Debugf("Looking up public key file at %q", pubKeyFilePath)
	publicKeyFileExists := FileExistsAndIsNotDir(pubKeyFilePath)
	if !publicKeyFileExists {
		// Public key does not exist. Let's grab it from hashicorp
		pubKeyFile, errDl := downloadFromURL(installLocation, product.GetPublicKeyUrl(), wg)
		if errDl != nil {
			logger.Errorf("Error fetching public key file from %s", product.GetPublicKeyUrl())
			return "", errDl
		}
		errRename := os.Rename(pubKeyFile, pubKeyFilePath)
		if errRename != nil {
			logger.Errorf("Error renaming public key file from %q to %q", pubKeyFile, pubKeyFilePath)
			return "", errRename
		}
	}
	return pubKeyFilePath, nil
}

func cleanup(paths []string, wg *sync.WaitGroup) {
	for _, path := range paths {
		wg.Add(1)
		logger.Infof("Deleting %q", path)
		err := os.Remove(path)
		if err != nil {
			logger.Error("Error deleting %q: %v", path, err)
		}
	}
}
