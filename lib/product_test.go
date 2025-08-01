package lib

import (
	"testing"
)

// Factory method tests
func Test_GetProductById(t *testing.T) {
	product := GetProductById("terraform")
	if product == nil {
		t.Errorf("Terraform product returned nil")
	} else {
		if expected := "terraform"; product.GetId() != expected {
			t.Errorf("Product ID does not match expected Id. Expected: %q, actual: %q", expected, product.GetId())
		}
	}

	product = GetProductById("opentofu")
	if product == nil {
		t.Errorf("Terraform product returned nil")
	} else {
		if expected := "opentofu"; product.GetId() != expected {
			t.Errorf("Product ID does not match expected Id. Expected: %q, actual: %q", expected, product.GetId())
		}
	}

	// Test case-insensitve match
	product = GetProductById("oPeNtOfU")
	if product == nil {
		t.Errorf("Terraform product returned nil")
	} else {
		if expected := "opentofu"; product.GetId() != expected {
			t.Errorf("Product ID does not match expected Id. Expected: %q, actual: %q", expected, product.GetId())
		}
	}

	product = GetProductById("doesnotexist")
	if product != nil {
		t.Errorf("Unknown product returned non-nil response")
	}
}

// Terraform Tests
func Test_GetId_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetId()
	if expected := "terraform"; actual != expected {
		t.Errorf("Product GetId does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetName_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetName()
	if expected := "Terraform"; actual != expected {
		t.Errorf("Product GetProductById does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetDefaultMirrorUrl_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetDefaultMirrorUrl()
	if expected := "https://releases.hashicorp.com/terraform"; actual != expected {
		t.Errorf("Product GetDefaultMirrorUrl does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetVersionPrefix_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetVersionPrefix()
	if expected := "terraform_"; actual != expected {
		t.Errorf("Product GetVersionPrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetExecutableName_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetExecutableName()
	if expected := "terraform"; actual != expected {
		t.Errorf("Product GetExecutableName does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetArchivePrefix_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetArchivePrefix()
	if expected := "terraform_"; actual != expected {
		t.Errorf("Product GetArchivePrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetArtifactUrl_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetArtifactUrl("https://example.com/terraform", "5.3.2")
	if expected := "https://example.com/terraform/5.3.2"; actual != expected {
		t.Errorf("Product GetArchivePrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetPublicKeyId_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetPublicKeyId()
	if expected := "72D7468F"; actual != expected {
		t.Errorf("Product GetPublicKeyId does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetPublicKeyUrl_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetPublicKeyUrl()
	if expected := "https://www.hashicorp.com/.well-known/pgp-key.txt"; actual != expected {
		t.Errorf("Product GetPublicKeyUrl does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetShaSignatureSuffix_Terraform(t *testing.T) {
	product := GetProductById("terraform")
	actual := product.GetShaSignatureSuffix()
	if expected := "72D7468F.sig"; actual != expected {
		t.Errorf("Product GetShaSignatureSuffix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetRecentVersionProduct_Terraform(t *testing.T) {
	recentFile := RecentFile{
		OpenTofu:  []string{"1.2.3", "3.2.1"},
		Terraform: []string{"5.4.3", "3.4.5"},
	}
	expected := []string{"5.4.3", "3.4.5"}

	product := GetProductById("terraform")
	actual := product.GetRecentVersionProduct(&recentFile)
	err := compareLists(actual, expected)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetRecentVersionProduct_Terraform(t *testing.T) {
	recentFile := RecentFile{
		OpenTofu:  []string{"1.2.3", "3.2.1"},
		Terraform: []string{"5.4.3", "3.4.5"},
	}
	expected := []string{"1.0.0", "1.0.1"}

	product := GetProductById("terraform")
	product.SetRecentVersionProduct(&recentFile, expected)
	err := compareLists(recentFile.Terraform, expected)
	if err != nil {
		t.Error(err)
	}

	err = compareLists(recentFile.OpenTofu, expected)
	if err == nil {
		t.Error("OpenTofu version list should not match version set for Terraform")
	}
}

// OpenTofu Tests
func Test_GetId_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetId()
	if expected := "opentofu"; actual != expected {
		t.Errorf("Product GetId does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetName_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetName()
	if expected := "OpenTofu"; actual != expected {
		t.Errorf("Product GetProductById does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetDefaultMirrorUrl_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetDefaultMirrorUrl()
	if expected := "https://get.opentofu.org/tofu"; actual != expected {
		t.Errorf("Product GetDefaultMirrorUrl does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetVersionPrefix_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetVersionPrefix()
	if expected := "opentofu_"; actual != expected {
		t.Errorf("Product GetVersionPrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetExecutableName_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetExecutableName()
	if expected := "tofu"; actual != expected {
		t.Errorf("Product GetExecutableName does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetArchivePrefix_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetArchivePrefix()
	if expected := "tofu_"; actual != expected {
		t.Errorf("Product GetArchivePrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetArtifactUrl_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetArtifactUrl("https://example.com/opentofu", "5.3.2")
	if expected := "https://github.com/opentofu/opentofu/releases/download/v5.3.2"; actual != expected {
		t.Errorf("Product GetArchivePrefix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetPublicKeyId_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetPublicKeyId()
	if expected := "0C0AF313E5FD9F80"; actual != expected {
		t.Errorf("Product GetPublicKeyId does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetPublicKeyUrl_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetPublicKeyUrl()
	if expected := "https://get.opentofu.org/opentofu.asc"; actual != expected {
		t.Errorf("Product GetPublicKeyUrl does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetShaSignatureSuffix_OpenTofu(t *testing.T) {
	product := GetProductById("opentofu")
	actual := product.GetShaSignatureSuffix()
	if expected := "gpgsig"; actual != expected {
		t.Errorf("Product GetShaSignatureSuffix does not match expected ID. Expected: %q, actual: %q", expected, actual)
	}
}

func Test_GetRecentVersionProduct_OpenTofu(t *testing.T) {
	recentFile := RecentFile{
		Terraform: []string{"1.2.3", "3.2.1"},
		OpenTofu:  []string{"5.4.3", "3.4.5"},
	}
	expected := []string{"5.4.3", "3.4.5"}

	product := GetProductById("opentofu")
	actual := product.GetRecentVersionProduct(&recentFile)
	err := compareLists(actual, expected)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetRecentVersionProduct_OpenTofu(t *testing.T) {
	recentFile := RecentFile{
		Terraform: []string{"1.2.3", "3.2.1"},
		OpenTofu:  []string{"5.4.3", "3.4.5"},
	}
	expected := []string{"1.0.0", "1.0.1"}

	product := GetProductById("opentofu")
	product.SetRecentVersionProduct(&recentFile, expected)
	err := compareLists(recentFile.OpenTofu, expected)
	if err != nil {
		t.Error(err)
	}

	err = compareLists(recentFile.Terraform, expected)
	if err == nil {
		t.Error("OpenTofu version list should not match version set for Terraform")
	}
}
