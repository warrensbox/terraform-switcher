package lib

import (
	"fmt"
	"strings"
)

// nolint:revive // FIXME: var-naming: const legacyProductId should be legacyProductID (revive)
const legacyProductId = "terraform"

// nolint:revive // FIXME: var-naming: struct field PublicKeyId should be PublicKeyID (revive)
type ProductDetails struct {
	ID                    string
	Name                  string
	DefaultMirror         string
	VersionPrefix         string
	DefaultDownloadMirror string
	ExecutableName        string
	ArchivePrefix         string
	PublicKeyId           string
	PublicKeyURLs         []string
	FileExtensions        []string
}

type TerraformProduct struct {
	ProductDetails
}
type OpenTofuProduct struct {
	ProductDetails
}

// nolint:revive // FIXME: var-naming: method GetId should be GetID (revive)
// nolint:revive // FIXME: var-naming: method GetDefaultMirrorUrl should be GetDefaultMirrorURL (revive)
// nolint:revive // FIXME: var-naming: method GetArtifactUrl should be GetArtifactURL (revive)
// nolint:revive // FIXME: var-naming: method GetPublicKeyId should be GetPublicKeyID (revive)
type Product interface {
	GetId() string
	GetName() string
	GetDefaultMirrorUrl() string
	GetVersionPrefix() string
	GetExecutableName() string
	GetArchivePrefix() string
	GetPublicKeyId() string
	GetPublicKeyURLs() []string
	GetShaSignatureSuffix() string
	GetArtifactUrl(mirrorURL string, version string) string
	GetRecentVersionProduct(recentFile *RecentFile) []string
	SetRecentVersionProduct(recentFile *RecentFile, versions []string)
	GetFileExtensions() []string
}

// Terraform Product
// nolint:revive // FIXME: var-naming: method GetId should be GetID (revive)
func (p TerraformProduct) GetId() string {
	return p.ID
}

func (p TerraformProduct) GetName() string {
	return p.Name
}

// nolint:revive // FIXME: var-naming: method GetDefaultMirrorUrl should be GetDefaultMirrorURL (revive)
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

// nolint:revive // FIXME: var-naming: method GetArtifactUrl should be GetArtifactURL (revive)
func (p TerraformProduct) GetArtifactUrl(mirrorURL string, version string) string {
	mirrorURL = strings.TrimRight(mirrorURL, "/")
	return fmt.Sprintf("%s/%s", mirrorURL, version)
}

// nolint:revive // FIXME: var-naming: method GetPublicKeyId should be GetPublicKeyID (revive)
func (p TerraformProduct) GetPublicKeyId() string {
	return p.PublicKeyId
}

func (p TerraformProduct) GetPublicKeyURLs() []string {
	return p.PublicKeyURLs
}

func (p TerraformProduct) GetShaSignatureSuffix() string {
	return p.GetPublicKeyId() + ".sig"
}

func (p TerraformProduct) GetRecentVersionProduct(recentFile *RecentFile) []string {
	return recentFile.Terraform
}

func (p TerraformProduct) SetRecentVersionProduct(recentFile *RecentFile, versions []string) {
	recentFile.Terraform = versions
}

func (p TerraformProduct) GetFileExtensions() []string {
	return p.FileExtensions
}

// OpenTofu methods
// nolint:revive // FIXME: var-naming: method GetId should be GetID (revive)
func (p OpenTofuProduct) GetId() string {
	return p.ID
}

func (p OpenTofuProduct) GetName() string {
	return p.Name
}

// nolint:revive // FIXME: var-naming: method GetDefaultMirrorUrl should be GetDefaultMirrorURL (revive)
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

// nolint:revive // FIXME: parameter 'mirrorURL' is not used (custom Mirror URL is not implemented for OpenTofu? 10-Mar-2025)
// nolint:revive // FIXME: var-naming: method GetArtifactUrl should be GetArtifactURL (revive)
func (p OpenTofuProduct) GetArtifactUrl(mirrorURL string, version string) string {
	return fmt.Sprintf("%s/v%s", p.DefaultDownloadMirror, version)
}

// nolint:revive // FIXME: var-naming: method GetPublicKeyId should be GetPublicKeyID (revive)
func (p OpenTofuProduct) GetPublicKeyId() string {
	return p.PublicKeyId
}

func (p OpenTofuProduct) GetPublicKeyURLs() []string {
	return p.PublicKeyURLs
}

func (p OpenTofuProduct) GetShaSignatureSuffix() string {
	return "gpgsig"
}

func (p OpenTofuProduct) GetRecentVersionProduct(recentFile *RecentFile) []string {
	return recentFile.OpenTofu
}

func (p OpenTofuProduct) SetRecentVersionProduct(recentFile *RecentFile, versions []string) {
	recentFile.OpenTofu = versions
}

func (p OpenTofuProduct) GetFileExtensions() []string {
	return p.FileExtensions
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
			PublicKeyURLs:  []string{"https://www.hashicorp.com/.well-known/pgp-key.txt", "https://keybase.io/hashicorp/pgp_keys.asc"},
			FileExtensions: []string{"tf"},
		},
	},
	OpenTofuProduct{
		ProductDetails{
			ID:                    "opentofu",
			Name:                  "OpenTofu",
			DefaultMirror:         "https://get.opentofu.org/tofu",
			DefaultDownloadMirror: "https://github.com/opentofu/opentofu/releases/download",
			VersionPrefix:         "opentofu_",
			ExecutableName:        "tofu",
			ArchivePrefix:         "tofu_",
			PublicKeyId:           "0C0AF313E5FD9F80",
			PublicKeyURLs:         []string{"https://get.opentofu.org/opentofu.asc"},
			FileExtensions:        []string{"tf", "tofu"},
		},
	},
}

// nolint:revive // FIXME: var-naming: func GetProductById should be GetProductByID (revive)
func GetProductById(id string) Product {
	for _, product := range products {
		if strings.EqualFold(product.GetId(), id) {
			return product
		}
	}
	return nil
}

func GetAllProducts() []Product {
	return products
}

// Obtain produced used by deprecated public methods that
// now expect a product to be called.
// Once these public methods are removed, this function can be removed
func getLegacyProduct() Product {
	product := GetProductById(legacyProductId)
	if product == nil {
		logger.Fatal("Default product could not be found")
	}
	return product
}
