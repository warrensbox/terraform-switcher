package lib

import "fmt"

const legacyProductId = "terraform"

type ProductDetails struct {
	ID                    string
	Name                  string
	DefaultMirror         string
	VersionPrefix         string
	DefaultDownloadMirror string
	ExecutableName        string
	ArchivePrefix         string
	PublicKeyId           string
	PublicKeyUrl          string
}

type TerraformProduct struct {
	ProductDetails
}
type OpenTofuProduct struct {
	ProductDetails
}

type Product interface {
	GetId() string
	GetName() string
	GetDefaultMirrorUrl() string
	GetVersionPrefix() string
	GetExecutableName() string
	GetArchivePrefix() string
	GetPublicKeyId() string
	GetPublicKeyUrl() string
	GetShaSignatureSuffix() string
	GetArtifactUrl(mirrorURL string, version string) string
}

// Terraform Product
func (p TerraformProduct) GetId() string {
	return p.ID
}
func (p TerraformProduct) GetName() string {
	return p.Name
}
func (p TerraformProduct) GetDefaultMirrorUrl() string {
	return p.DefaultMirror
}
func (p TerraformProduct) GetVersionPrefix() string {
	return p.VersionPrefix
}
func (p TerraformProduct) GetExecutableName() string {
	return p.ExecutableName
}
func (p TerraformProduct) GetArchivePrefix() string {
	return p.ArchivePrefix
}
func (p TerraformProduct) GetArtifactUrl(mirrorURL string, version string) string {
	return fmt.Sprintf("%s%s", mirrorURL, version)
}
func (p TerraformProduct) GetPublicKeyId() string {
	return p.PublicKeyId
}
func (p TerraformProduct) GetPublicKeyUrl() string {
	return p.PublicKeyUrl
}
func (p TerraformProduct) GetShaSignatureSuffix() string {
	return p.GetPublicKeyId() + ".sig"
}

// OpenTofu methods
func (p OpenTofuProduct) GetId() string {
	return p.ID
}
func (p OpenTofuProduct) GetName() string {
	return p.Name
}
func (p OpenTofuProduct) GetDefaultMirrorUrl() string {
	return p.DefaultMirror
}
func (p OpenTofuProduct) GetVersionPrefix() string {
	return p.VersionPrefix
}
func (p OpenTofuProduct) GetExecutableName() string {
	return p.ExecutableName
}
func (p OpenTofuProduct) GetArchivePrefix() string {
	return p.ArchivePrefix
}
func (p OpenTofuProduct) GetArtifactUrl(mirrorURL string, version string) string {
	return fmt.Sprintf("%s/v%s", p.DefaultDownloadMirror, version)
}
func (p OpenTofuProduct) GetPublicKeyId() string {
	return p.PublicKeyId
}
func (p OpenTofuProduct) GetPublicKeyUrl() string {
	return p.PublicKeyUrl
}
func (p OpenTofuProduct) GetShaSignatureSuffix() string {
	return "gpgsig"
}

// Factory methods
var products = []Product{
	TerraformProduct{
		ProductDetails{
			ID:             "terraform",
			Name:           "Terraform",
			DefaultMirror:  "https://releases.hashicorp.com/terraform",
			VersionPrefix:  "terraform_",
			ExecutableName: "terraform",
			ArchivePrefix:  "terraform_",
			PublicKeyId:    "72D7468F",
			PublicKeyUrl:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
		},
	},
	OpenTofuProduct{
		ProductDetails{
			ID:                    "opentofu",
			Name:                  "OpenTofu",
			DefaultMirror:         "https://get.opentofu.org/tofu",
			DefaultDownloadMirror: "https://github.com/opentofu/opentofu/releases/download",
			VersionPrefix:         "opentofu_",
			ExecutableName:        "opentofu",
			ArchivePrefix:         "tofu_",
			PublicKeyId:           "0C0AF313E5FD9F80",
			PublicKeyUrl:          "https://get.opentofu.org/opentofu.asc",
		},
	},
}

func GetProductById(id string) Product {
	for _, product := range products {
		if product.GetId() == id {
			return product
		}
	}
	return nil
}

func GetAllProducts() []Product {
	return products
}

func getLegacyProduct() Product {
	product := GetProductById(legacyProductId)
	if product == nil {
		logger.Fatalf("Default product could not be found")
	}
	return product
}
