package lib_test

import (
	"testing"

	"github.com/warrensbox/terraform-switcher/lib"
)

var versionsRaw = []string{
	"1.1",
	"1.2.1",
	"1.2.2",
	"1.2.3",
	"1.3",
	"1.1.4",
	"0.7.1",
	"1.4-beta",
	"1.4",
	"2",
}

// TestSemverParser1 : Test to see if SemVerParser parses valid version
// Test version 1.1
func TestSemverParserCase1(t *testing.T) {
	tfconstraint := "1.1"
	tfversion, semVerErr := lib.SemVerParser(&tfconstraint, versionsRaw)
	if semVerErr != nil {
		t.Errorf("Error parsing version %q: %v", tfconstraint, semVerErr)
	}
	expected := "1.1.0"
	if tfversion == expected {
		t.Logf("Version %q exists in list [expected]", expected)
	} else {
		t.Logf("Version %q does not exist in list [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %q", expected)
	}
}

// TestSemverParserCase2 : Test to see if SemVerParser parses valid version
// Test version ~> 1.1 should return  1.1.4
func TestSemverParserCase2(t *testing.T) {
	tfconstraint := "~> 1.1.0"
	tfversion, semVerErr := lib.SemVerParser(&tfconstraint, versionsRaw)
	if semVerErr != nil {
		t.Errorf("Error parsing version %q: %v", tfconstraint, semVerErr)
	}
	expected := "1.1.4"
	if tfversion == expected {
		t.Logf("Version %q exist in list [expected]", expected)
	} else {
		t.Logf("Version %q does not exist in list [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %q", expected)
	}
}

// TestSemverParserCase3 : Test to see if SemVerParser parses valid version
// Test version ~> 1.1 should return  1.1.4
func TestSemverParserCase3(t *testing.T) {
	tfconstraint := "~> 1.A.0"
	_, err := lib.SemVerParser(&tfconstraint, versionsRaw)
	if err != nil {
		t.Logf("This test is supposed to error on %q [expected]", tfconstraint)
	} else {
		t.Errorf("This test is supposed to error on %q but passed [unexpected]", tfconstraint)
	}
}

// TestSemverParserCase4 : Test to see if SemVerParser parses valid version
// Test version ~> >= 1.0, < 1.4 should return  1.3.0
func TestSemverParserCase4(t *testing.T) {
	tfconstraint := ">= 1.0, < 1.4"
	tfversion, semVerErr := lib.SemVerParser(&tfconstraint, versionsRaw)
	if semVerErr != nil {
		t.Errorf("Error parsing version %q: %v", tfconstraint, semVerErr)
	}
	expected := "1.3.0"
	if tfversion == expected {
		t.Logf("Version %q exist in list [expected]", expected)
	} else {
		t.Logf("Version %q does not exist in list [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %q", expected)
	}
}

// TestSemverParserCase5 : Test to see if SemVerParser parses valid version
// Test version ~> >= 1.0 should return  2.0.0
func TestSemverParserCase5(t *testing.T) {
	tfconstraint := ">= 1.0"
	tfversion, semVerErr := lib.SemVerParser(&tfconstraint, versionsRaw)
	if semVerErr != nil {
		t.Errorf("Error parsing version %q: %v", tfconstraint, semVerErr)
	}
	expected := "2.0.0"
	if tfversion == expected {
		t.Logf("Version %q exist in list [expected]", expected)
	} else {
		t.Logf("Version %q does not exist in list [unexpected]", tfconstraint)
		t.Errorf("This is unexpected. Parsing failed. Expected: %q", expected)
	}
}

func TestSemverCheckFoss(t *testing.T) {
	tests := map[string]struct {
		version string
		expected   bool
	}{
		"FOSS": {
			version: "1.5.1",
			expected: true,
		},
		"BSL": {
			version: "1.7.1",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := lib.SemVerCheckFoss(test.version)
			if err != nil {
				t.Errorf("Error checking for Foss licensed version %q: %v", test.version, err)
			}
			if actual != test.expected {
				t.Errorf("%s: Version %q returned %v. Expected: %v", name, test.version, actual, test.expected)
			}
		})
	}
}

