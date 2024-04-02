package param_parsing

import (
	"fmt"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
	"os"
	"strings"
)

const terraformVersionFileName = ".terraform-version"

func GetParamsFromTerraformVersion(params Params) Params {
	filePath := params.ChDirPath + "/" + terraformVersionFileName
	if lib.CheckFileExist(filePath) {
		fmt.Printf("Reading configuration from %s\n", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("Could not read file content", filePath, err)
		}
		params.Version = strings.TrimSpace(string(content))
	}
	return params
}

func terraformVersionFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + terraformVersionFileName
	return lib.CheckFileExist(filePath)
}
