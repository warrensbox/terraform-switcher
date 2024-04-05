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
		logger.Errorf("Could not open %q", signatureFilePath)
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

// checkChecksumMatches This will calculate and compare the check sum of the downloaded zip file.
func checkChecksumMatches(hashFile string, targetFile *os.File) bool {
	var fileHandlersToClose []*os.File
	fileHandlersToClose = append(fileHandlersToClose, targetFile)

	_, fileName := filepath.Split(targetFile.Name())
	expectedChecksum, err := getChecksumFromHashFile(hashFile, fileName)
	if err != nil {
		closeFileHandlers(fileHandlersToClose)
		logger.Errorf("Could not get expected checksum from file: %q", err.Error())
		return false
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, targetFile); err != nil {
		closeFileHandlers(fileHandlersToClose)
		logger.Errorf("Calculating Checksum failed: %q", err.Error())
		return false
	}
	checksum := hex.EncodeToString(hash.Sum(nil))
	if expectedChecksum != checksum {
		closeFileHandlers(fileHandlersToClose)
		logger.Errorf("Checksum mismatch. Expected: %q, expected %v", expectedChecksum, checksum)
		return false
	}
	closeFileHandlers(fileHandlersToClose)
	return true
}

// checkSignatureOfChecksums THis will verify the signature of the file containing the hash sums
func checkSignatureOfChecksums(keyRingReader *os.File, hashFile *os.File, signatureFile *os.File) bool {
	var fileHandlersToClose []*os.File
	fileHandlersToClose = append(fileHandlersToClose, keyRingReader)
	fileHandlersToClose = append(fileHandlersToClose, hashFile)
	fileHandlersToClose = append(fileHandlersToClose, signatureFile)

	logger.Info("Verifying signature of checksum file...")
	keyring, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		closeFileHandlers(fileHandlersToClose)
		logger.Errorf("Read armored key ring: %q", err.Error())
		return false
	}

	_, err = openpgp.CheckDetachedSignature(keyring, hashFile, signatureFile)
	if err != nil {
		closeFileHandlers(fileHandlersToClose)
		logger.Errorf("Checking detached signature: %q", err.Error())
		return false
	}
	logger.Info("Verification successful.")
	closeFileHandlers(fileHandlersToClose)
	return true
}
