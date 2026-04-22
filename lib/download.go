//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// DownloadFromURL : Downloads the terraform binary and its hash from the source url
func DownloadFromURL(installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch string) (string, error) {
	product := getLegacyProduct()
	return DownloadProductFromURL(product, installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch)
}

//nolint:gocyclo
func DownloadProductFromURL(product Product, installLocation, mirrorURL, tfversion, versionPrefix, goos, goarch string) (string, error) {
	var wg sync.WaitGroup
	defer wg.Done()
	// nolint:revive // FIXME: var-naming: var zipUrl should be zipURL (revive)
	zipUrl := mirrorURL + "/" + versionPrefix + tfversion + "_" + goos + "_" + goarch + ".zip"
	// nolint:revive // FIXME: var-naming: var hashUrl should be hashURL (revive)
	hashUrl := mirrorURL + "/" + versionPrefix + tfversion + "_SHA256SUMS"
	// nolint:revive // FIXME: var-naming: var hashSignatureUrl should be hashSignatureURL (revive)
	hashSignatureUrl := mirrorURL + "/" + versionPrefix + tfversion + "_SHA256SUMS." + product.GetShaSignatureSuffix()

	match := false

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
	defer func() {
		if !match {
			os.Remove(zipFilePath)
		}
	}()

	logger.Infof("Downloading %q", hashUrl)
	hashFilePath, err := downloadFromURL(installLocation, hashUrl, &wg)
	if err != nil {
		logger.Error("Could not download hash file")
		return "", err
	}
	defer os.Remove(hashFilePath)

	logger.Infof("Downloading %q", hashSignatureUrl)
	hashSigFilePath, err := downloadFromURL(installLocation, hashSignatureUrl, &wg)
	if err != nil {
		logger.Error("Could not download hash signature file")
		return "", err
	}
	defer os.Remove(hashSigFilePath)

	publicKeyFile, err := os.Open(pubKeyFilename)
	if err != nil {
		logger.Errorf("Could not open public key %q: %v", pubKeyFilename, err)
		return "", err
	}
	defer publicKeyFile.Close()

	signatureFile, err := os.Open(hashSigFilePath)
	if err != nil {
		logger.Errorf("Could not open hash signature file %q: %v", hashSigFilePath, err)
		return "", err
	}
	defer signatureFile.Close()

	targetFile, err := os.Open(zipFilePath)
	if err != nil {
		logger.Errorf("Could not open zip file %q: %v", zipFilePath, err)
		return "", err
	}
	defer targetFile.Close()

	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		logger.Errorf("Could not open hash file %q: %v", hashFilePath, err)
		return "", err
	}
	defer hashFile.Close()

	var verifySucceed bool
	if verifySucceed, err = verifySignature(product, publicKeyFile, hashFile, signatureFile); err != nil {
		return "", err
	} else if !verifySucceed {
		return "", fmt.Errorf("Unable to verify checksum signature against PGP key")
	}

	match = checkChecksumMatches(hashFilePath, targetFile)
	if !match {
		return "", errors.New("Checksums did not match")
	}

	return zipFilePath, err
}

// verifySignature: Verify signature of hash file.
// If
func verifySignature(product Product, publicKeyFile, hashFile, signatureFile *os.File) (bool, error) {
	// CAUTION: Skip PGP signature verification of checksum file if TF_SKIP_SIGNATURE_VERIFICATION
	// environment variable is set to true-ish value: 1, t, T, TRUE, true, True
	// THIS IS NOT RECOMMENDED AND SHOULD ONLY BE USED FOR TESTING PURPOSES!
	skipSignatureVerificationStr, exists := os.LookupEnv("TF_SKIP_SIGNATURE_VERIFICATION")
	if !exists {
		skipSignatureVerificationStr = "false"
	}
	skipSignatureVerification, err := strconv.ParseBool(skipSignatureVerificationStr)
	if err != nil {
		logger.Warnf(
			"Unable to parse \"TF_SKIP_SIGNATURE_VERIFICATION\" env var value %q, defaulting to \"false\"",
			skipSignatureVerificationStr,
		)
		return true, nil
	}

	if skipSignatureVerification {
		logger.Warn(
			"Skipping PGP signature verification of checksum file due to " +
				"\"TF_SKIP_SIGNATURE_VERIFICATION\" environment variable being set",
		)
		logger.Warn("!!! THIS IS NOT RECOMMENDED AND SHOULD ONLY BE USED FOR TESTING PURPOSES !!!")
		return true, nil
	}

	logger.Infof("Verifying PGP signature of checksum file: %q", hashFile.Name())

	keyFileContent, err := io.ReadAll(publicKeyFile)
	if err != nil {
		return false, fmt.Errorf("Could not read PGP key file %q: %v", publicKeyFile.Name(), err)
	}

	hashFileContent, err := io.ReadAll(hashFile)
	if err != nil {
		return false, fmt.Errorf("Could not read hash file %q: %v", hashFile.Name(), err)
	}

	signatureContent, err := io.ReadAll(signatureFile)
	if err != nil {
		return false, fmt.Errorf("Could not read PGP signature file %q: %v", signatureFile.Name(), err)
	}

	// Verify signature using key
	verified := checkSignatureOfChecksums(keyFileContent, hashFileContent, signatureContent)
	if verified {
		return true, nil
	}

	// Fail fast if there is no legacy builtin PGP public key to fall back to
	if product.GetPublicKeyLegacyLiteral() == "" {
		return false, errors.New("Signature of checksum file could not be verified and fallback does not exist")
	}

	legacyBuiltinKeyIdentifier := "legacy builtin PGP public key"
	logger.Warnf(
		"Checksum file PGP signature verification failed with public key from %q file. "+
			"Falling back to %s", publicKeyFile.Name(), legacyBuiltinKeyIdentifier,
	)

	verified = checkSignatureOfChecksums([]byte(product.GetPublicKeyLegacyLiteral()), hashFileContent, signatureContent)
	if !verified {
		logger.Errorf("Signature of checksum file could not be verified with %s either", legacyBuiltinKeyIdentifier)
	}
	return verified, nil
}

func downloadFromURL(installLocation string, url string, wg *sync.WaitGroup) (string, error) {
	wg.Add(1)
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	logger.Infof("Downloading to %q", filepath.Join(installLocation, "/", fileName))

	response, err := http.Get(url) // nolint:gosec // `url' is expected to be variable
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
	logger.Debugf("Looking up public PGP-key file at %q", pubKeyFilePath)
	publicKeyFileExists := FileExistsAndIsNotDir(pubKeyFilePath)
	if !publicKeyFileExists {
		// Public PGP-key does not exist. Let's grab it
		publicKeyURLs := product.GetPublicKeyURLs()
		var pubKeyFile string
		var errDl error
		var errsDl []string
		for idx, publicKeyURL := range publicKeyURLs {
			logger.Debugf("Attempting to download public PGP-key from %q", publicKeyURL)
			pubKeyFile, errDl = downloadFromURL(installLocation, publicKeyURL, wg)
			if errDl != nil {
				errsDl = append(errsDl, errDl.Error())
				logger.Errorf("Failed to fetch public PGP-key from %q", publicKeyURL)

				// Return all failures if all URLs have been tried so far
				if idx+1 == len(publicKeyURLs) {
					return "", errors.New(strings.Join(errsDl, "; "))
				}

				// Try the next URL
				continue
			}
			// Download succeeded, break out of the loop
			break
		}

		logger.Debugf("Renaming public PGP-key file from %q to %q", pubKeyFile, pubKeyFilePath)
		errRename := os.Rename(pubKeyFile, pubKeyFilePath)
		if errRename != nil {
			logger.Errorf("Error renaming public PGP-key file from %q to %q", pubKeyFile, pubKeyFilePath)
			return "", errRename
		}
	}
	return pubKeyFilePath, nil
}
