package lib_test

import (
	"reflect"
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

// TestNewCommand : pass value and check if returned value is a pointer
func TestNewCommand(t *testing.T) {

	testCmd := "terraform"
	cmd := lib.NewCommand(testCmd)

	if reflect.ValueOf(cmd).Kind() == reflect.Ptr {
		t.Logf("Value returned is a pointer %v [expected]", cmd)
	} else {
		t.Errorf("Value returned is not a pointer %v [expected", cmd)
	}
}

// TestPathList : check if bin path exist
func TestPathList(t *testing.T) {

	testCmd := ""
	cmd := lib.NewCommand(testCmd)
	listBin := cmd.PathList()

	if listBin == nil {
		t.Error("No bin path found [unexpected]")
	} else {
		t.Logf("Found bin path [expected]")
	}
}

// TestFind : check common "cd" command exist
// This is assuming that Windows and linux has the "cd" command
func TestFind(t *testing.T) {

	testCmd := "cd"
	cmd := lib.NewCommand(testCmd)

	next := cmd.Find()
	for path := next(); len(path) > 0; path = next() {
		if path != "" {
			t.Logf("Found installation path: %v [expected]\n", path)
		} else {
			t.Errorf("Unable to find '%v' command in this operating system [unexpected]", testCmd)
		}
	}
}
