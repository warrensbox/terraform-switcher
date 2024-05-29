package lib

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
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

	openpgpConfig := packet.Config{
		RSABits:                1024,
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
	}
	gpgKeyEntity, err := openpgp.NewEntity("TestProductSign", "Signing key for test product", "example@localhost.com", &openpgpConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	gpgFingerprint := hex.EncodeToString(gpgKeyEntity.PrimaryKey.Fingerprint[:])[:8]

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

	var publicKeySerialiseBuffer bytes.Buffer
	err = gpgKeyEntity.Serialize(&publicKeySerialiseBuffer)
	if err != nil {
		t.Fatal(err)
	}
	var publicKey bytes.Buffer
	publicKeyWriter, err := armor.Encode(&publicKey, openpgp.PublicKeyType, nil)
	if err != nil {
		t.Fatal(err)
	}
	publicKeyWriter.Close()

	zipFileBytes := zipFileBuffer.Bytes()

	// Calculate SHA256 sum of ZIP file
	sha256HashWriter := sha256.New()
	sha256Hash := sha256HashWriter.Sum(zipFileBytes)

	// Create checksum file
	checksumFileContent := string(sha256Hash) + "  " + "my_product_download_2.1.0_linux_amd64.zip"
	checksumFileReader := bytes.NewBuffer([]byte(checksumFileContent))

	// Create signature of checksum file
	var sigFile bytes.Buffer
	err = openpgp.DetachSign(&sigFile, gpgKeyEntity, checksumFileReader, &openpgpConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/testproduct/gpg-key.txt":
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write(publicKey.Bytes())
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
			w.Write(sigFile.Bytes())
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
	if expectedZipPath := ""; zipFilePath != expectedZipPath {
		t.Errorf("Returned zipFile not expected path. Expected: %q, actual: %q", expectedZipPath, zipFilePath)
	}
}
