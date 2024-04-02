package param_parsing

import (
	"fmt"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
	"os"
	"strings"
)

const tfSwitchFileName = ".tfswitchrc"

func GetParamsFromTfSwitch(params Params) Params {
	filePath := params.ChDirPath + "/" + tfSwitchFileName
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

func tfSwitchFileExists(params Params) bool {
	filePath := params.ChDirPath + "/" + tfSwitchFileName
	return lib.CheckFileExist(filePath)
}
