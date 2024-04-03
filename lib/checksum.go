package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getChecksumFromFile Extract the checksum from the signature file
func getChecksumFromHashFile(signatureFilePath string, terraformFileName string) (string, error) {
	readFile, err := os.Open(signatureFilePath)
	if err != nil {
		fmt.Println("[Error] : Could not open ", signatureFilePath)
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
		fmt.Println("[Error] : Could not get expected checksum from file: " + err.Error())
		return false
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, targetFile); err != nil {
		closeFileHandlers(fileHandlersToClose)
		fmt.Println("[Error] : Calculating Checksum failed: " + err.Error())
		return false
	}
	checksum := hex.EncodeToString(hash.Sum(nil))
	if expectedChecksum != checksum {
		closeFileHandlers(fileHandlersToClose)
		fmt.Println("[Error] : Checksum mismatch. Expected: ", expectedChecksum, " got ", checksum)
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

	log.Println("Verifying signature of checksum file...")
	keyring, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		closeFileHandlers(fileHandlersToClose)
		log.Fatal("[Error] : Read armored key ring: " + err.Error())
		return false
	}

	_, err = openpgp.CheckDetachedSignature(keyring, hashFile, signatureFile)
	if err != nil {
		closeFileHandlers(fileHandlersToClose)
		log.Fatal("[Error] : Checking detached signature: " + err.Error())
		return false
	}
	log.Println("Verification successful.")
	closeFileHandlers(fileHandlersToClose)
	return true
}
