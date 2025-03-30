package lib

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// TestDownloadFromURL_FileNameMatch : Check expected filename exist when downloaded
func TestDownloadFromURL_FileNameMatch(t *testing.T) {
	logger = InitLogger("DEBUG")
	hashiURL := "https://releases.hashicorp.com/terraform/"
	installVersion := "terraform_"
	tempDir := t.TempDir()
	installLocation := filepath.Join(tempDir, ".terraform.versions_test")
	macOS := "_darwin_amd64.zip"

	t.Logf("install Location: %v", installLocation)

	// create /.terraform.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		t.Logf("Creating directory for terraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0o755)

		t.Cleanup(func() {
			defer os.Remove(tempDir)
			t.Logf("Cleanup temporary directory %q", tempDir)
		})

		if err != nil {
			t.Logf("Unable to create directory for terraform: %v", installLocation)
			t.Error("Test fail")
		}
	}

	/* test download old terraform version */
	lowestVersion := "0.11.0"
	var wg sync.WaitGroup
	defer wg.Done()
	urlToDownload := hashiURL + lowestVersion + "/" + installVersion + lowestVersion + macOS
	expectedFile := filepath.Join(installLocation, installVersion+lowestVersion+macOS)
	installedFile, errDownload := downloadFromURL(installLocation, urlToDownload, &wg)

	if errDownload != nil {
		t.Logf("Expected file name %v to be downloaded", expectedFile)
		t.Error("Download not possible (unexpected)")
	}

	if installedFile == expectedFile {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Log("Download file matches expected file")
	} else {
		t.Logf("Expected file %v", expectedFile)
		t.Logf("Downloaded file %v", installedFile)
		t.Error("Download file mismatches expected file (unexpected)")
	}

	// check file name is what is expected
	_, checkErr := os.Stat(expectedFile)
	if checkErr != nil {
		t.Logf("Expected file does not exist %v", expectedFile)
	}
}

// TestDownloadFromURL_Valid : Test if https://releases.hashicorp.com/terraform/ is still valid
func TestDownloadFromURL_Valid(t *testing.T) {
	logger = InitLogger("DEBUG")
	hashiURL := "https://releases.hashicorp.com/terraform/"

	url, err := url.ParseRequestURI(hashiURL)
	if err != nil {
		t.Errorf("Invalid URL %v [unexpected]", err)
	} else {
		t.Logf("Valid URL from %v [expected]", url)
	}
}

type DownloadProductTestConfig struct {
	GpgFingerprint      string
	ZipFileContent      []byte
	ZipFileChecksum     string
	ChecksumFileContent string
	PublicKey           string
}

//nolint:gocyclo
func setupTestDownloadServer(t *testing.T, downloadProductTestConfig *DownloadProductTestConfig) *httptest.Server {
	logger = InitLogger("DEBUG")
	gpgKey, err := crypto.GenerateKey("TestProductSign", "example@localhost.com", "RSA", 1024)
	if err != nil {
		t.Fatal(err)
	}
	if downloadProductTestConfig.GpgFingerprint == "" {
		downloadProductTestConfig.GpgFingerprint = gpgKey.GetFingerprint()[0:8]
	}

	if len(downloadProductTestConfig.ZipFileContent) == 0 {
		zipFileBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipFileBuffer)

		executableBytes := []byte("This is the main executable")
		zipFileContentWriter, err := zipWriter.Create("myprod")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := zipFileContentWriter.Write(executableBytes); err != nil {
			t.Fatal(err)
		}
		zipWriter.Flush()
		zipWriter.Close()
		downloadProductTestConfig.ZipFileContent = zipFileBuffer.Bytes()
	}

	if downloadProductTestConfig.PublicKey == "" {
		downloadProductTestConfig.PublicKey, err = gpgKey.GetArmoredPublicKey()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Calculate SHA256 sum of ZIP file
	if downloadProductTestConfig.ZipFileChecksum == "" {
		sha256HashWriter := sha256.New()
		zipFileReadBuffer := new(bytes.Buffer)
		_, err := zipFileReadBuffer.Write(downloadProductTestConfig.ZipFileContent)
		if err != nil {
			t.Fatal(err)
		}

		if _, err = io.Copy(sha256HashWriter, zipFileReadBuffer); err != nil {
			t.Fatal(err)
		}
		sha256Hash := sha256HashWriter.Sum(nil)
		downloadProductTestConfig.ZipFileChecksum = hex.EncodeToString(sha256Hash)
	}

	// Create checksum file
	if downloadProductTestConfig.ChecksumFileContent == "" {
		downloadProductTestConfig.ChecksumFileContent = downloadProductTestConfig.ZipFileChecksum + "  " + "my_product_download_2.1.0_linux_amd64.zip"
	}

	// Create signature of checksum file
	binMessage := crypto.NewPlainMessageFromFile([]byte(downloadProductTestConfig.ChecksumFileContent), "my_product_download_2.1.0_SHA256SUMS", uint32(crypto.GetUnixTime()))
	keyRing, err := crypto.NewKeyRing(gpgKey)
	if err != nil {
		t.Fatal(err)
	}

	signatureObj, err := keyRing.SignDetached(binMessage)
	if err != nil {
		t.Fatal(err)
	}

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/testproduct/gpg-key.txt":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(downloadProductTestConfig.PublicKey)); err != nil {
				t.Error(err)
			}
		case "/productdownload/2.1.0/my_product_download_2.1.0_linux_amd64.zip":
			w.Header().Set("Content-Type", "application/zip")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(downloadProductTestConfig.ZipFileContent); err != nil {
				t.Error(err)
			}
		case "/productdownload/2.1.0/my_product_download_2.1.0_SHA256SUMS":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(downloadProductTestConfig.ChecksumFileContent)); err != nil {
				t.Error(err)
			}
		case "/productdownload/2.1.0/my_product_download_2.1.0_SHA256SUMS." + downloadProductTestConfig.GpgFingerprint + ".sig":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(signatureObj.GetBinary()); err != nil {
				t.Error(err)
			}
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	return mockServer
}

// TestDownloadProductFromURL : Test DownloadProductFromURL
func TestDownloadProductFromURL(t *testing.T) {
	downloadProductTestConfig := DownloadProductTestConfig{}
	mockServer := setupTestDownloadServer(t, &downloadProductTestConfig)
	defer mockServer.Close()

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    downloadProductTestConfig.GpgFingerprint,
			PublicKeyUrl:   mockServer.URL + "/testproduct/gpg-key.txt",
		},
	}

	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	zipFilePath, err := DownloadProductFromURL(mockProduct, tempDir, mockProduct.GetArtifactUrl(mockServer.URL+"/productdownload", "2.1.0"), "2.1.0", mockProduct.GetArchivePrefix(), "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
	if expectedZipPath := filepath.Join(tempDir, "my_product_download_2.1.0_linux_amd64.zip"); zipFilePath != expectedZipPath {
		t.Errorf("Returned zipFile not expected path. Expected: %q, actual: %q", expectedZipPath, zipFilePath)
	}
}

// TestDownloadProductFromURL_invalid_checksum : Test DownloadProductFromURL with invalid zipfile checksum
func TestDownloadProductFromURL_invalid_checksum(t *testing.T) {
	downloadProductTestConfig := DownloadProductTestConfig{
		ZipFileChecksum: "abcdef",
	}
	mockServer := setupTestDownloadServer(t, &downloadProductTestConfig)
	defer mockServer.Close()

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    downloadProductTestConfig.GpgFingerprint,
			PublicKeyUrl:   mockServer.URL + "/testproduct/gpg-key.txt",
		},
	}

	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = DownloadProductFromURL(mockProduct, tempDir, mockProduct.GetArtifactUrl(mockServer.URL+"/productdownload", "2.1.0"), "2.1.0", mockProduct.GetArchivePrefix(), "linux", "amd64")
	if err == nil {
		t.Fatal("DownloadProductFromURL did not throw error")
	} else if expectedError := "Checksums did not match"; !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("DownloadProductFromURL did not throw expected error. Expected: %q, Actual: %q", expectedError, err)
	}
}

// TestDownloadProductFromURL_zip_file_not_present_in_checksum : Test DownloadProductFromURL with zipfile not present in checksum
func TestDownloadProductFromURL_zip_file_not_present_in_checksum(t *testing.T) {
	downloadProductTestConfig := DownloadProductTestConfig{
		ChecksumFileContent: "e64d27cf0fd05eaa0deab98756f4d533d14d467a7198b82a885be260e1ec4885  doesnotexist.txt",
	}
	mockServer := setupTestDownloadServer(t, &downloadProductTestConfig)
	defer mockServer.Close()

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    downloadProductTestConfig.GpgFingerprint,
			PublicKeyUrl:   mockServer.URL + "/testproduct/gpg-key.txt",
		},
	}

	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = DownloadProductFromURL(mockProduct, tempDir, mockProduct.GetArtifactUrl(mockServer.URL+"/productdownload", "2.1.0"), "2.1.0", mockProduct.GetArchivePrefix(), "linux", "amd64")
	if err == nil {
		t.Fatal("DownloadProductFromURL did not throw error")
	} else if expectedError := "Checksums did not match"; !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("DownloadProductFromURL did not throw expected error. Expected: %q, Actual: %q", expectedError, err)
	}
}

// TestDownloadProductFromURL_unable_to_download_public_key : Test DownloadProductFromURL unable to download public key
func TestDownloadProductFromURL_unable_to_download_public_key(t *testing.T) {
	downloadProductTestConfig := DownloadProductTestConfig{}
	mockServer := setupTestDownloadServer(t, &downloadProductTestConfig)
	defer mockServer.Close()

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    downloadProductTestConfig.GpgFingerprint,
			PublicKeyUrl:   mockServer.URL + "/testproduct/invalid-public-key",
		},
	}

	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = DownloadProductFromURL(mockProduct, tempDir, mockProduct.GetArtifactUrl(mockServer.URL+"/productdownload", "2.1.0"), "2.1.0", mockProduct.GetArchivePrefix(), "linux", "amd64")
	errorMatch := regexp.MustCompile("Unable to download from .*/testproduct/invalid-public-key")
	if err == nil {
		t.Fatal("DownloadProductFromURL did not throw error")
	} else if !errorMatch.MatchString(err.Error()) {
		t.Fatalf("DownloadProductFromURL did not throw expected error. Expected: %q, Actual: %q", errorMatch, err)
	}
}

// TestDownloadProductFromURL_invalid_public_key : Test DownloadProductFromURL with invalid public key content
func TestDownloadProductFromURL_invalid_public_key(t *testing.T) {
	downloadProductTestConfig := DownloadProductTestConfig{
		PublicKey: "thisisinvalid",
	}
	mockServer := setupTestDownloadServer(t, &downloadProductTestConfig)
	defer mockServer.Close()

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    downloadProductTestConfig.GpgFingerprint,
			PublicKeyUrl:   mockServer.URL + "/testproduct/gpg-key.txt",
		},
	}

	// Create temp location
	tempDir, err := os.MkdirTemp("", "addRecentVersion")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	_, err = DownloadProductFromURL(mockProduct, tempDir, mockProduct.GetArtifactUrl(mockServer.URL+"/productdownload", "2.1.0"), "2.1.0", mockProduct.GetArchivePrefix(), "linux", "amd64")
	if err == nil {
		t.Fatal("DownloadProductFromURL did not throw error")
	} else if expectedError := "Signature of checksum file could not be verified"; !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("DownloadProductFromURL did not throw expected error. Expected: %q, Actual: %q", expectedError, err)
	}
}
