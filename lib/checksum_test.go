package lib

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/constants"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

func Test_getChecksumFromHashFile(t *testing.T) {
	expected := "3ff056b5e8259003f67fd0f0ed7229499cfb0b41f3ff55cc184088589994f7a5"
	got, err := getChecksumFromHashFile("../test-data/terraform_1.7.5_SHA256SUMS", "terraform_1.7.5_linux_amd64.zip")
	if err != nil {
		t.Errorf("getChecksumFromHashFile() error = %v", err)
		return
	}
	if got != expected {
		t.Errorf("getChecksumFromHashFile() got = %v, expected %v", got, expected)
	}
}

func Test_checkChecksumMatches(t *testing.T) {
	InitLogger("TRACE")
	targetFile, err := os.Open("../test-data/checksum-check-file")
	if err != nil {
		t.Errorf("[Error]: Could not open testfile for signature verification.")
	}

	if got := checkChecksumMatches("../test-data/terraform_1.7.5_SHA256SUMS", targetFile); got != true {
		t.Errorf("checkChecksumMatches() = %v, want %v", got, true)
	}
}

// Key-generation cache. PGP keygen (RSA 2048, standard security) is the
// slow part of these tests; generating one keyset per package-test run
// keeps the full suite in the single-digit-seconds range.
//
// Failures inside the sync.Once closure are captured into sharedTestKeysErr
// rather than calling t.Fatalf directly, because Once.Do marks itself done
// regardless of how its closure exits: aborting the first caller with
// t.Fatalf would leave every subsequent caller looking at nil keys and
// panicking instead of reporting the real cause.
var (
	sharedTestKeysOnce sync.Once
	sharedTestKeyA     *crypto.Key
	sharedTestKeyB     *crypto.Key
	sharedTestKeyC     *crypto.Key
	sharedTestKeysErr  error
)

func sharedTestKeys(t *testing.T) (*crypto.Key, *crypto.Key, *crypto.Key) {
	t.Helper()
	sharedTestKeysOnce.Do(func() {
		pgp := crypto.PGPWithProfile(profile.RFC4880())
		for target := range slices.Values([]**crypto.Key{&sharedTestKeyA, &sharedTestKeyB, &sharedTestKeyC}) {
			gen := pgp.KeyGeneration().AddUserId("tfswitch-test", "tfswitch-test@example.invalid").New()
			k, err := gen.GenerateKeyWithSecurity(constants.StandardSecurity)
			if err != nil {
				sharedTestKeysErr = err
				return
			}
			*target = k
		}
	})
	if sharedTestKeysErr != nil {
		t.Fatalf("PGP key generation failed: %v", sharedTestKeysErr)
	}
	return sharedTestKeyA, sharedTestKeyB, sharedTestKeyC
}

func armoredPublicKey(t *testing.T, key *crypto.Key) string {
	t.Helper()
	armored, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey: %v", err)
	}
	return armored
}

func signDetached(t *testing.T, key *crypto.Key, payload []byte) []byte {
	t.Helper()
	signer, err := crypto.PGP().Sign().SigningKey(key).Detached().New()
	if err != nil {
		t.Fatalf("build signing handle: %v", err)
	}
	sig, err := signer.Sign(payload, crypto.Auto)
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}
	return sig
}

// writeReadableTempFile writes content to a fresh file under t.TempDir()
// and returns it opened read-only. checkSignatureOfChecksums closes every
// file it receives, so each test hands out its own handle. A t.Cleanup is
// registered to close the handle even when a test aborts before
// checkSignatureOfChecksums runs; os.File.Close on an already-closed file
// is harmless (returns an error we ignore), so the belt-and-braces is safe.
func writeReadableTempFile(t *testing.T, name string, content []byte) *os.File {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func Test_parsePublicKeys_singleBlock(t *testing.T) {
	InitLogger("DEBUG")
	keyA, _, _ := sharedTestKeys(t)
	keys, err := parsePublicKeys(armoredPublicKey(t, keyA))
	if err != nil {
		t.Fatalf("parsePublicKeys returned error: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("want 1 key, got %d", len(keys))
	}
	if keys[0].GetFingerprint() != keyA.GetFingerprint() {
		t.Errorf("fingerprint mismatch: want %s, got %s", keyA.GetFingerprint(), keys[0].GetFingerprint())
	}
}

func Test_parsePublicKeys_multipleBlocks(t *testing.T) {
	InitLogger("DEBUG")
	keyA, keyB, _ := sharedTestKeys(t)
	concatenated := armoredPublicKey(t, keyA) + "\n" + armoredPublicKey(t, keyB)

	keys, err := parsePublicKeys(concatenated)
	if err != nil {
		t.Fatalf("parsePublicKeys returned error: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("want 2 keys, got %d", len(keys))
	}
	if keys[0].GetFingerprint() != keyA.GetFingerprint() {
		t.Errorf("first key fingerprint mismatch: want %s, got %s", keyA.GetFingerprint(), keys[0].GetFingerprint())
	}
	if keys[1].GetFingerprint() != keyB.GetFingerprint() {
		t.Errorf("second key fingerprint mismatch: want %s, got %s", keyB.GetFingerprint(), keys[1].GetFingerprint())
	}
}

func Test_parsePublicKeys_emptyInput(t *testing.T) {
	InitLogger("DEBUG")
	for input := range slices.Values([]string{"", "   \n\n\t  "}) {
		_, err := parsePublicKeys(input)
		if err == nil {
			t.Errorf("parsePublicKeys(%q) expected error, got nil", input)
		}
	}
}

func Test_parsePublicKeys_mixedValidAndGarbage(t *testing.T) {
	InitLogger("DEBUG")
	keyA, _, _ := sharedTestKeys(t)
	garbage := pgpPublicKeyBegin + "\nnot-a-real-armored-body\n-----END PGP PUBLIC KEY BLOCK-----\n"
	concatenated := armoredPublicKey(t, keyA) + "\n" + garbage

	keys, err := parsePublicKeys(concatenated)
	if err != nil {
		t.Fatalf("parsePublicKeys returned error: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("want 1 key, got %d", len(keys))
	}
	if keys[0].GetFingerprint() != keyA.GetFingerprint() {
		t.Errorf("fingerprint mismatch: want %s, got %s", keyA.GetFingerprint(), keys[0].GetFingerprint())
	}
}

func Test_checkSignatureOfChecksums_singleKey(t *testing.T) {
	InitLogger("DEBUG")
	keyA, _, _ := sharedTestKeys(t)
	payload := []byte("abc123  terraform_1.7.5_linux_amd64.zip\n")
	signature := signDetached(t, keyA, payload)

	keyFile := writeReadableTempFile(t, "pubkey.asc", []byte(armoredPublicKey(t, keyA)))
	hashFile := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile := writeReadableTempFile(t, "SHA256SUMS.sig", signature)

	if !checkSignatureOfChecksums(keyFile, hashFile, sigFile) {
		t.Fatal("checkSignatureOfChecksums returned false for a single-key happy path")
	}
}

func Test_checkSignatureOfChecksums_signerIsFirst(t *testing.T) {
	InitLogger("DEBUG")
	keyA, keyB, _ := sharedTestKeys(t)
	payload := []byte("abc123  terraform_1.7.5_linux_amd64.zip\n")
	signature := signDetached(t, keyA, payload)
	keyring := armoredPublicKey(t, keyA) + "\n" + armoredPublicKey(t, keyB)

	keyFile := writeReadableTempFile(t, "pubkey.asc", []byte(keyring))
	hashFile := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile := writeReadableTempFile(t, "SHA256SUMS.sig", signature)

	if !checkSignatureOfChecksums(keyFile, hashFile, sigFile) {
		t.Fatal("checkSignatureOfChecksums returned false with signer at position 0")
	}
}

// Test_checkSignatureOfChecksums_signerIsSecond is the direct regression
// test for issue #746: HashiCorp's current pgp-key.txt ships the expired
// key first and the active signing key second, and every release signed
// with the successor key must still verify.
func Test_checkSignatureOfChecksums_signerIsSecond(t *testing.T) {
	InitLogger("DEBUG")
	keyA, keyB, _ := sharedTestKeys(t)
	payload := []byte("abc123  terraform_1.14.9_linux_amd64.zip\n")
	signature := signDetached(t, keyB, payload)
	keyring := armoredPublicKey(t, keyA) + "\n" + armoredPublicKey(t, keyB)

	keyFile := writeReadableTempFile(t, "pubkey.asc", []byte(keyring))
	hashFile := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile := writeReadableTempFile(t, "SHA256SUMS.sig", signature)

	if !checkSignatureOfChecksums(keyFile, hashFile, sigFile) {
		t.Fatal("checkSignatureOfChecksums returned false with signer at position 1 (GitHub issue #746 scenario)")
	}
}

func Test_checkSignatureOfChecksums_noMatchingKey(t *testing.T) {
	InitLogger("DEBUG")
	keyA, keyB, keyC := sharedTestKeys(t)
	payload := []byte("abc123  terraform_1.7.5_linux_amd64.zip\n")
	signature := signDetached(t, keyC, payload)
	keyring := armoredPublicKey(t, keyA) + "\n" + armoredPublicKey(t, keyB)

	keyFile := writeReadableTempFile(t, "pubkey.asc", []byte(keyring))
	hashFile := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile := writeReadableTempFile(t, "SHA256SUMS.sig", signature)

	if checkSignatureOfChecksums(keyFile, hashFile, sigFile) {
		t.Fatal("checkSignatureOfChecksums returned true for signature made by a key not in the file")
	}
}

func Test_checkSignatureOfChecksums_malformedKeyFile(t *testing.T) {
	InitLogger("DEBUG")
	keyA, _, _ := sharedTestKeys(t)
	payload := []byte("abc123  terraform_1.7.5_linux_amd64.zip\n")
	signature := signDetached(t, keyA, payload)

	keyFile := writeReadableTempFile(t, "pubkey.asc", []byte("this is definitely not a PGP key"))
	hashFile := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile := writeReadableTempFile(t, "SHA256SUMS.sig", signature)

	if checkSignatureOfChecksums(keyFile, hashFile, sigFile) {
		t.Fatal("checkSignatureOfChecksums returned true for a malformed key file")
	}
	// Guard against strings.Split silently "parsing" a file that contains the
	// BEGIN marker but nothing useful after it.
	junkWithMarker := pgpPublicKeyBegin + "\nstill garbage\n-----END PGP PUBLIC KEY BLOCK-----\n"
	if !strings.Contains(junkWithMarker, pgpPublicKeyBegin) {
		t.Fatal("test fixture lost BEGIN marker; check constant")
	}
	keyFile2 := writeReadableTempFile(t, "pubkey2.asc", []byte(junkWithMarker))
	hashFile2 := writeReadableTempFile(t, "SHA256SUMS", payload)
	sigFile2 := writeReadableTempFile(t, "SHA256SUMS.sig", signature)
	if checkSignatureOfChecksums(keyFile2, hashFile2, sigFile2) {
		t.Fatal("checkSignatureOfChecksums returned true for a file with only a malformed armored block")
	}
}
