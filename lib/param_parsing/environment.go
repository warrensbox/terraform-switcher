package param_parsing

import (
	"github.com/warrensbox/terraform-switcher/lib/types"
	"os"
)

func GetParamsFromEnvironment(params types.Params) types.Params {
	params.Version = os.Getenv("TF_VERSION")
	return params
}
