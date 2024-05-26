package lib

type ProductDetails struct {
	ID            string
	Name          string
	DefaultMirror string
}

var products = []ProductDetails{
	{
		ID:            "terraform",
		Name:          "Terraform",
		DefaultMirror: "https://releases.hashicorp.com/terraform",
	},
	{
		ID:            "opentofu",
		Name:          "OpenTofu",
		DefaultMirror: "https://get.opentofu.org/tofu",
	},
}

func GetProductById(id string) *ProductDetails {
	for _, product := range products {
		if product.ID == id {
			return &product
		}
	}
	return nil
}

func GetAllProducts() []ProductDetails {
	return products
}
