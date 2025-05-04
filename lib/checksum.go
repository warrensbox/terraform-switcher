package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// getChecksumFromFile Extract the checksum from the signature file
func getChecksumFromHashFile(signatureFilePath string, terraformFileName string) (string, error) {
	readFile, err := os.Open(signatureFilePath)
	if err != nil {
		logger.Errorf("Could not open %q: %v", signatureFilePath, err)
		return "", err
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "  ")
		if len(split) == 2 && split[1] == terraformFileName {
			return split[0], nil
		}
	}
	return "", nil
}

// checkChecksumMatches This will calculate and compare the check sum of the downloaded zip file
func checkChecksumMatches(hashFile string, targetFile *os.File) bool {
	logger.Debugf("Checksum comparison for %q", targetFile.Name())
	var fileHandlersToClose []*os.File
	fileHandlersToClose = append(fileHandlersToClose, targetFile)
	defer closeFileHandlers(fileHandlersToClose)

	_, fileName := filepath.Split(targetFile.Name())
	expectedChecksum, err := getChecksumFromHashFile(hashFile, fileName)
	if err != nil {
		logger.Errorf("Could not get checksum from file %q: %v", hashFile, err)
		return false
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, targetFile); err != nil {
		logger.Errorf("Checksum calculation failed for %q: %v", fileName, err)
		return false
	}
	checksum := hex.EncodeToString(hash.Sum(nil))
	if expectedChecksum != checksum {
		logger.Errorf("Checksum mismatch for %q. Expected: %q, calculated: %v", fileName, expectedChecksum, checksum)
		return false
	}
	return true
}

// checkSignatureOfChecksums This will verify the signature of the file containing the hash sums
func checkSignatureOfChecksums(keyFile *os.File, hashFile *os.File, signatureFile *os.File) bool {
	var fileHandlersToClose []*os.File
	fileHandlersToClose = append(fileHandlersToClose, keyFile)
	fileHandlersToClose = append(fileHandlersToClose, hashFile)
	fileHandlersToClose = append(fileHandlersToClose, signatureFile)
	defer closeFileHandlers(fileHandlersToClose)

	logger.Infof("Verifying PGP signature of checksum file: %q", hashFile.Name())

	keyFileContent, err := io.ReadAll(keyFile)
	if err != nil {
		logger.Errorf("Could not read PGP key file %q: %v", keyFile, err)
		return false
	}

	keyFromArmored, err := crypto.NewKeyFromArmored(string(keyFileContent))
	if err != nil {
		logger.Errorf("Could not read PGP armored key: %v", err)
		return false
	}

	signingKey, err := crypto.PGP().Verify().VerificationKey(keyFromArmored).New()
	if err != nil {
		logger.Errorf("Could not read PGP signing key: %v", err)
		return false
	}

	hashFileContent, err := io.ReadAll(hashFile)
	if err != nil {
		logger.Errorf("Could not read hash file %q: %v", hashFile, err)
		return false
	}

	signatureContent, err := io.ReadAll(signatureFile)
	if err != nil {
		logger.Errorf("Could not read PGP signature file %q: %v", signatureFile, err)
		return false
	}

	verifyRes, err := signingKey.VerifyDetached(hashFileContent, signatureContent, crypto.Auto)
	if err != nil {
		logger.Errorf("Could not verify detached signature PGP message: %v", err)
		return false
	}

	if err := verifyRes.SignatureError(); err != nil {
		logger.Errorf("Could not verify PGP signature: %v", err)
		return false
	}

	logger.Info("Checksum file PGP signature verification successful")
	return true
}
