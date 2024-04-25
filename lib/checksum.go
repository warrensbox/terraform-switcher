package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		if split[1] == terraformFileName {
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
func checkSignatureOfChecksums(keyRingReader *os.File, hashFile *os.File, signatureFile *os.File) bool {
	var fileHandlersToClose []*os.File
	fileHandlersToClose = append(fileHandlersToClose, keyRingReader)
	fileHandlersToClose = append(fileHandlersToClose, hashFile)
	fileHandlersToClose = append(fileHandlersToClose, signatureFile)
	defer closeFileHandlers(fileHandlersToClose)

	logger.Info("Verifying signature of checksum file...")
	keyring, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		logger.Errorf("Could not read armored key ring: %v", err)
		return false
	}

	_, err = openpgp.CheckDetachedSignature(keyring, hashFile, signatureFile)
	if err != nil {
		logger.Errorf("Could not check detached signature: %v", err)
		return false
	}
	logger.Info("Checksum file signature verification successful.")
	return true
}
