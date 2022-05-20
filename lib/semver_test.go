package lib_test

import (
	"fmt"
	"testing"

	semver "github.com/hashicorp/go-version"

	"github.com/warrensbox/terraform-switcher/lib"
)

var tmp lib.Release

var releasesRaw = []*lib.Release{
	tmp.NewRelease("2"),
	tmp.NewRelease("1.4-beta"),
	tmp.NewRelease("1.4"),
	tmp.NewRelease("1.3"),
	tmp.NewRelease("1.2.3"),
	tmp.NewRelease("1.2.2"),
	tmp.NewRelease("1.2.1"),
	tmp.NewRelease("1.1.4"),
	tmp.NewRelease("1.1"),
	tmp.NewRelease("0.7.1"),
}

// TestSemverParser1 : Test to see if SemVerParser parses valid version
// Test version 1.1
func TestSemverParserCase1(t *testing.T) {
	tfconstraint := "1.1"
	tfversion, err := lib.SemVerParser(&tfconstraint, nil)
	if err != nil {
		t.Fatalf("%q", err)
	}
	expected, err := semver.NewVersion("1.1.0")
	if err != nil {
		t.Fatalf("%q", err)
	}
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
	tfconstraint := fmt.Sprintf("~> 1.1.0")
	tfversion, err := lib.SemVerParser(&tfconstraint, releasesRaw)
	if err != nil {
		t.Fatalf("%q", err)
	}
	expected, err := semver.NewVersion("1.1.4")
	if err != nil {
		t.Fatalf("%q", err)
	}
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
	tfversion, err := lib.SemVerParser(&tfconstraint, releasesRaw)
	if err != nil {
		t.Fatalf("%q", err)
	}
	expected, err := semver.NewVersion("1.3.0")
	if err != nil {
		t.Fatalf("%q", err)
	}
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
	expected, err := semver.NewVersion("2.0.0")
	if err != nil {
		t.Fatalf("%q", err)
	}
	if tfversion.Version == expected {
		t.Logf("Version exist in list %v [expected]", expected)
	} else {
		t.Logf("Version does not exist in list %v [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %v", expected)
	}
}
