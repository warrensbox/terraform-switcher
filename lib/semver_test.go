package lib_test

import (
	"github.com/warrensbox/terraform-switcher/lib"
	"testing"
)

var releasesRaw = []lib.Release{
	{Version: "1.1"},
	{Version: "1.2.1"},
	{Version: "1.2.2"},
	{Version: "1.2.3"},
	{Version: "1.3"},
	{Version: "1.1.4"},
	{Version: "0.7.1"},
	{Version: "1.4-beta"},
	{Version: "1.4"},
	{Version: "2"},
}

// TestSemverParser1 : Test to see if SemVerParser parses valid version
// Test version 1.1
func TestSemverParserCase1(t *testing.T) {

	tfconstraint := "1.1"
	tfversion, _ := lib.SemVerParser(&tfconstraint, releasesRaw)
	expected := "1.1.0"
	if tfversion.Version == expected {
		t.Logf("Version exist in list %v [expected]", expected)
	} else {
		t.Logf("Version does not exist in list %v [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %v", expected)
	}
}

// TestSemverParserCase2 : Test to see if SemVerParser parses valid version
// Test version ~> 1.1 should return  1.1.4
func TestSemverParserCase2(t *testing.T) {

	tfconstraint := "~> 1.1.0"
	tfversion, _ := lib.SemVerParser(&tfconstraint, releasesRaw)
	expected := "1.1.4"
	if tfversion.Version == expected {
		t.Logf("Version exist in list %v [expected]", expected)
	} else {
		t.Logf("Version does not exist in list %v [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %v", expected)
	}
}

// TestSemverParserCase3 : Test to see if SemVerParser parses valid version
// Test version ~> 1.1 should return  1.1.4
func TestSemverParserCase3(t *testing.T) {

	tfconstraint := "~> 1.A.0"
	_, err := lib.SemVerParser(&tfconstraint, releasesRaw)
	if err != nil {
		t.Logf("This test is suppose to error %v [expected]", tfconstraint)
	} else {
		t.Errorf("This test is suppose to error but passed %v [expected]", tfconstraint)
	}
}

// TestSemverParserCase4 : Test to see if SemVerParser parses valid version
// Test version ~> >= 1.0, < 1.4 should return  1.3.0
func TestSemverParserCase4(t *testing.T) {

	tfconstraint := ">= 1.0, < 1.4"
	tfversion, _ := lib.SemVerParser(&tfconstraint, releasesRaw)
	expected := "1.3.0"
	if tfversion.Version == expected {
		t.Logf("Version exist in list %v [expected]", expected)
	} else {
		t.Logf("Version does not exist in list %v [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %v", expected)
	}
}

// TestSemverParserCase5 : Test to see if SemVerParser parses valid version
// Test version ~> >= 1.0 should return  2.0.0
func TestSemverParserCase5(t *testing.T) {

	tfconstraint := ">= 1.0"
	tfversion, _ := lib.SemVerParser(&tfconstraint, releasesRaw)
	expected := "2.0.0"
	if tfversion.Version == expected {
		t.Logf("Version exist in list %v [expected]", expected)
	} else {
		t.Logf("Version does not exist in list %v [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %v", expected)
	}
}
