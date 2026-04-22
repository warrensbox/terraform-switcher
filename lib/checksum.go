//nolint:staticcheck //ST1005: error strings should not be capitalized (staticcheck)
package lib

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
	// Along with that, skip the part[0] altogether.
	parts := strings.Split(armored, pgpPublicKeyBegin)
	keys := make([]*crypto.Key, 0, len(parts)-1)
	for i, part := range parts[1:] {
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
func checkSignatureOfChecksums(keyFileContent []byte, hashFileContent []byte, signatureContent []byte) bool {
	keys, err := parsePublicKeys(string(keyFileContent))
	if err != nil {
		logger.Errorf("Could not parse PGP keys: %v", err)
		return false
	}

	var verificationErrors []error
	for idx, key := range keys {
		logger.Debugf(
			"Trying to verify PGP signature using key №%d (out of %d) with fingerprint %q",
			idx+1, len(keys), key.GetFingerprint(),
		)

		verifier, err := crypto.PGP().Verify().VerificationKey(key).New()
		if err != nil {
			verificationErrors = append(verificationErrors, fmt.Errorf(
				"Could not read PGP signing key №%d (out of %d): %v",
				idx+1, len(keys), err,
			))
			continue
		}

		verifyRes, err := verifier.VerifyDetached(hashFileContent, signatureContent, crypto.Auto)
		if err != nil {
			verificationErrors = append(verificationErrors, fmt.Errorf(
				"Could not verify detached signature PGP message using key №%d (out of %d): %v",
				idx+1, len(keys), err,
			))
			continue
		}

		if err := verifyRes.SignatureError(); err != nil {
			verificationErrors = append(verificationErrors, fmt.Errorf(
				"Could not verify PGP signature using key №%d (out of %d): %v",
				idx+1, len(keys), err,
			))
			continue
		}

		logger.Info("Checksum file PGP signature verification successful")
		return true
	}

	// Print errors once (if any)
	for verificationError := range slices.Values(verificationErrors) {
		logger.Error(verificationError)
	}
	return false
}
