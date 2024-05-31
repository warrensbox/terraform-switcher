package lib

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/mitchellh/go-homedir"
)

// TestDownloadFromURL_FileNameMatch : Check expected filename exist when downloaded
func TestDownloadFromURL_FileNameMatch(t *testing.T) {
	logger = InitLogger("DEBUG")
	hashiURL := "https://releases.hashicorp.com/terraform/"
	installVersion := "terraform_"
	tempDir := t.TempDir()
	installPath := fmt.Sprintf(tempDir + string(os.PathSeparator) + ".terraform.versions_test")
	macOS := "_darwin_amd64.zip"

	home, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Could not detect home directory")
	}

	t.Logf("Current home directory: %q", home)
	if runtime.GOOS != "windows" {
		installLocation = filepath.Join(home, installPath)
	} else {
		installLocation = installPath
	}
	t.Logf("install Location: %v", installLocation)

	// create /.terraform.versions_test/ directory to store code
	if _, err := os.Stat(installLocation); os.IsNotExist(err) {
		t.Logf("Creating directory for terraform: %v", installLocation)
		err = os.MkdirAll(installLocation, 0755)
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

	//check file name is what is expected
	_, err = os.Stat(expectedFile)
	if err != nil {
		t.Logf("Expected file does not exist %v", expectedFile)
	}

	t.Cleanup(func() {
		defer os.Remove(tempDir)
		t.Logf("Cleanup temporary directory %q", tempDir)
	})
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

// TestDownloadProductFromURL : Test DownloadProductFromURL
func TestDownloadProductFromURL(t *testing.T) {
	logger = InitLogger("DEBUG")
	gpgKey, err := crypto.GenerateKey("TestProductSign", "example@localhost.com", "RSA", 1024)
	if err != nil {
		t.Fatal(err)
	}
	gpgFingerprint := gpgKey.GetFingerprint()[0:8]

	zipFileBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipFileBuffer)

	mainExecutableBytes := []byte("This is the main executable")
	zipFileContentWriter, err := zipWriter.Create("myprod")
	if err != nil {
		t.Fatal(err)
	}
	zipFileContentWriter.Write(mainExecutableBytes)
	zipWriter.Flush()
	zipWriter.Close()
	zipFileBytes := zipFileBuffer.Bytes()

	publicKey, err := gpgKey.GetArmoredPublicKey()
	if err != nil {
		t.Fatal(err)
	}

	// Calculate SHA256 sum of ZIP file
	sha256HashWriter := sha256.New()
	if _, err = io.Copy(sha256HashWriter, zipFileBuffer); err != nil {
		t.Fatal(err)
	}
	sha256Hash := sha256HashWriter.Sum(nil)

	// Create checksum file
	checksumFileContent := hex.EncodeToString(sha256Hash) + "  " + "my_product_download_2.1.0_linux_amd64.zip"

	// Create signature of checksum file
	binMessage := crypto.NewPlainMessageFromFile([]byte(checksumFileContent), "my_product_download_2.1.0_SHA256SUMS", uint32(crypto.GetUnixTime()))
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
			w.Write([]byte(publicKey))
		case "/productdownload/2.1.0/my_product_download_2.1.0_linux_amd64.zip":
			w.Header().Set("Content-Type", "application/zip")
			w.WriteHeader(http.StatusOK)
			w.Write(zipFileBytes)
		case "/productdownload/2.1.0/my_product_download_2.1.0_SHA256SUMS":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(checksumFileContent))
		case "/productdownload/2.1.0/my_product_download_2.1.0_SHA256SUMS." + gpgFingerprint + ".sig":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write(signatureObj.GetBinary())
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	// Create mock product
	mockProduct := TerraformProduct{
		ProductDetails{
			ID:             "myproduct",
			Name:           "Mock Product",
			DefaultMirror:  mockServer.URL + "/productdownload",
			VersionPrefix:  "myprod_",
			ExecutableName: "myprod",
			ArchivePrefix:  "my_product_download_",
			PublicKeyId:    gpgFingerprint,
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
