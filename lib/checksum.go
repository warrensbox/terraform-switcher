package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// pgpPublicKeyBegin is the ASCII armor marker that delimits every public
// key block in a concatenated armored file (such as HashiCorp's
// pgp-key.txt, which publishes more than one key around rotations).
const pgpPublicKeyBegin = "-----BEGIN PGP PUBLIC KEY BLOCK-----"

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

// parsePublicKeys extracts every armored PGP public key from the given
// blob. HashiCorp's pgp-key.txt contains more than one key around key
// rotations, and the key that signed a release may not be the first
// block. crypto.NewKeyFromArmored refuses multi-entity input, so the
// caller must feed it one armored block at a time.
//
// Blocks that fail to parse are logged at Debug and skipped; the
// function returns an error only when no usable keys remain.
func parsePublicKeys(armored string) ([]*crypto.Key, error) {
	// strings.Split leaves the text preceding the first BEGIN marker in
	// parts[0] (always non-key content); each subsequent entry is the
	// body of one block. Re-prepending the marker rather than trimming
	// preserves the mandatory blank line between the marker/headers and
	// the base64 payload that RFC 4880 armor requires.
	parts := strings.Split(armored, pgpPublicKeyBegin)
	keys := make([]*crypto.Key, 0, len(parts))
	for i, part := range parts {
		if i == 0 {
			continue
		}
		block := pgpPublicKeyBegin + part
		key, err := crypto.NewKeyFromArmored(block)
		if err != nil {
			logger.Debugf("Skipping unparsable PGP public key block №%d: %v", i, err)
			continue
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil, errors.New("no parsable PGP public keys found in key file")
	}
	return keys, nil
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
		logger.Errorf("Could not read PGP key file %q: %v", keyFile.Name(), err)
		return false
	}

	keys, err := parsePublicKeys(string(keyFileContent))
	if err != nil {
		logger.Errorf("Could not parse PGP keys from %q: %v", keyFile.Name(), err)
		return false
	}

	verifyBuilder := crypto.PGP().Verify()
	// Parse every armored block from the key file and register each parsed key
	// with verify handle builder. Successive VerificationKey() calls append
	// to the builder's internal keyring, and the key matching the signature's
	// KeyID is picked automatically.
	for key := range slices.Values(keys) {
		verifyBuilder = verifyBuilder.VerificationKey(key)
	}
	signingKey, err := verifyBuilder.New()
	if err != nil {
		logger.Errorf("Could not read PGP signing key: %v", err)
		return false
	}

	hashFileContent, err := io.ReadAll(hashFile)
	if err != nil {
		logger.Errorf("Could not read hash file %q: %v", hashFile.Name(), err)
		return false
	}

	signatureContent, err := io.ReadAll(signatureFile)
	if err != nil {
		logger.Errorf("Could not read PGP signature file %q: %v", signatureFile.Name(), err)
		return false
	}

	verifyRes, err := signingKey.VerifyDetached(hashFileContent, signatureContent, crypto.Auto)
	if err != nil {
		logger.Errorf("Could not verify detached signature PGP message: %v", err)
		return false
	}

	if err := verifyRes.SignatureError(); err != nil {
		logger.Errorf("Could not verify PGP signature (tried %d keys): %v", len(keys), err)
		return false
	}

	logger.Info("Checksum file PGP signature verification successful")
	return true
}
