package param_parsing

import (
	"fmt"
	"strings"

	"github.com/gookit/slog"
	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
)

type Params struct {
	ChDirPath        string
	CustomBinaryPath string
	DefaultVersion   string
	DryRun           bool
	HelpFlag         bool
	InstallPath      string
	LatestFlag       bool
	LatestPre        string
	LatestStable     string
	ListAllFlag      bool
	LogLevel         string
	MirrorURL        string
	ShowLatestFlag   bool
	ShowLatestPre    string
	ShowLatestStable string
	Product          string
	Version          string
	VersionFlag      bool
}

var logger *slog.Logger

func GetParameters() Params {
	var params Params
	params = initParams(params)

	var productIds []string
	var defaultMirrors []string
	for _, product := range lib.GetAllProducts() {
		productIds = append(productIds, product.GetId())
		defaultMirrors = append(defaultMirrors, fmt.Sprintf("%s: %s", product.GetName(), product.GetDefaultMirrorUrl()))
	}

	getopt.StringVarLong(&params.ChDirPath, "chdir", 'c', "Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory")
	getopt.StringVarLong(&params.CustomBinaryPath, "bin", 'b', "Custom binary path. Ex: tfswitch -b "+lib.ConvertExecutableExt("/Users/username/bin/terraform"))
	getopt.StringVarLong(&params.DefaultVersion, "default", 'd', "Default to this version in case no other versions could be detected. Ex: tfswitch --default 1.2.4")
	getopt.BoolVarLong(&params.DryRun, "dry-run", 'r', "Only show what tfswitch would do. Don't download anything.")
	getopt.BoolVarLong(&params.HelpFlag, "help", 'h', "Displays help message")
	getopt.StringVarLong(&params.InstallPath, "install", 'i', "Custom install path. Ex: tfswitch -i /Users/username. The binaries will be in the sub installDir directory e.g. /Users/username/"+lib.InstallDir)
	getopt.BoolVarLong(&params.LatestFlag, "latest", 'u', "Get latest stable version")
	getopt.StringVarLong(&params.LatestPre, "latest-pre", 'p', "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.LatestStable, "latest-stable", 's', "Latest implicit version based on a constraint. Ex: tfswitch --latest-stable 0.13.0 downloads 0.13.7 and 0.13 downloads 0.15.5 (latest)")
	getopt.BoolVarLong(&params.ListAllFlag, "list-all", 'l', "List all versions of terraform - including beta and rc")
	getopt.StringVarLong(&params.LogLevel, "log-level", 'g', "Set loglevel for tfswitch. One of (INFO, NOTICE, DEBUG, TRACE)")
	getopt.StringVarLong(&params.MirrorURL, "mirror", 'm', "install from a remote API other than the default. Default (based on product):\n"+strings.Join(defaultMirrors, "\n"))
	getopt.BoolVarLong(&params.ShowLatestFlag, "show-latest", 'U', "Show latest stable version")
	getopt.StringVarLong(&params.ShowLatestPre, "show-latest-pre", 'P', "Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.ShowLatestStable, "show-latest-stable", 'S', "Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)")
	getopt.StringVarLong(&params.Product, "product", 'q', fmt.Sprintf("Specifies which product to use. Ex: `tfswitch --product opentofu` will install Terraform. Options: (%s). Default: %s", strings.Join(productIds, ", "), lib.DefaultProductId))
	getopt.BoolVarLong(&params.VersionFlag, "version", 'v', "Displays the version of tfswitch")

	// Parse the command line parameters to fetch stuff like chdir
	getopt.Parse()

	oldLogLevel := params.LogLevel
	logger = lib.InitLogger(params.LogLevel)
	var err error
	// Read configuration files
	if tomlFileExists(params) {
		params, err = getParamsTOML(params)
	} else if tfSwitchFileExists(params) {
		params, err = GetParamsFromTfSwitch(params)
	} else if terraformVersionFileExists(params) {
		params, err = GetParamsFromTerraformVersion(params)
	} else if isTerraformModule(params) {
		params, err = GetVersionFromVersionsTF(params)
	} else if terraGruntFileExists(params) {
		params, err = GetVersionFromTerragrunt(params)
	} else {
		params = GetParamsFromEnvironment(params)
	}
	if err != nil {
		logger.Fatalf("Error parsing configuration file: %q", err)
	}

	// Set defaults based on product
	product := lib.GetProductById(params.Product)
	if product == nil {
		logger.Fatalf("Invalid \"product\" configuration value: %q", params.Product)
	} else { // Use else as there is a warning that params maybe nil, as it does not see Fatalf as a break condition
		params.MirrorURL = product.GetDefaultMirrorUrl()
	}

	// Logger config was changed by the config files. Reinitialise.
	if params.LogLevel != oldLogLevel {
		logger = lib.InitLogger(params.LogLevel)
	}

	// Parse again to overwrite anything that might by defined on the cli AND in any config file (CLI always wins)
	getopt.Parse()
	args := getopt.Args()
	if len(args) == 1 {
		/* version provided on command line as arg */
		params.Version = args[0]
	}

	return params
}

func initParams(params Params) Params {
	params.ChDirPath = lib.GetCurrentDirectory()
	params.CustomBinaryPath = lib.ConvertExecutableExt(lib.GetDefaultBin())
	params.DefaultVersion = lib.DefaultLatest
	params.DryRun = false
	params.HelpFlag = false
	params.InstallPath = lib.GetHomeDirectory()
	params.LatestFlag = false
	params.LatestPre = lib.DefaultLatest
	params.LatestStable = lib.DefaultLatest
	params.ListAllFlag = false
	params.LogLevel = "INFO"
	params.MirrorURL = ""
	params.ShowLatestFlag = false
	params.ShowLatestPre = lib.DefaultLatest
	params.ShowLatestStable = lib.DefaultLatest
	params.Version = lib.DefaultLatest
	params.Product = lib.DefaultProductId
	params.VersionFlag = false
	return params
}
