package lib

import (
	"os"
	"testing"
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
