//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/gookit/slog"
	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
)

type Params struct {
	Arch             string
	ChDirPath        string
	CustomBinaryPath string
	DefaultVersion   string
	DryRun           bool
	ForceColor       bool
	HelpFlag         bool
	InstallPath      string
	LatestFlag       bool
	LatestPre        string
	LatestStable     string
	ListAllFlag      bool
	LogLevel         string
	MirrorURL        string
	NoColor          bool
	ProductEntity    lib.Product
	Product          string
	ShowLatestFlag   bool
	ShowLatestPre    string
	ShowLatestStable string
	TomlDir          string
	VersionFlag      bool
	Version          string
}

// This is used to automatically instate Environment variables and TOML keys
var paramMappings = []struct {
	description string
	env         string
	param       string
	ptype       reflect.Kind
	toml        string
}{
	{param: "Arch", ptype: reflect.String, env: "TF_ARCH", toml: "arch", description: "CPU architecture"},
	{param: "CustomBinaryPath", ptype: reflect.String, env: "TF_BINARY_PATH", toml: "bin", description: "Custom binary path"},
	{param: "DefaultVersion", ptype: reflect.String, env: "TF_DEFAULT_VERSION", toml: "default-version", description: "Default version"},
	{param: "ForceColor", ptype: reflect.Bool, env: "FORCE_COLOR", toml: "force-color", description: "Force color output if terminal supports it"},
	{param: "InstallPath", ptype: reflect.String, env: "TF_INSTALL_PATH", toml: "install", description: "Custom install path"},
	{param: "LogLevel", ptype: reflect.String, env: "TF_LOG_LEVEL", toml: "log-level", description: "Log level"},
	{param: "NoColor", ptype: reflect.Bool, env: "NO_COLOR", toml: "no-color", description: "Disable color output"},
	{param: "Product", ptype: reflect.String, env: "TF_PRODUCT", toml: "product", description: "Product"},
	{param: "Version", ptype: reflect.String, env: "TF_VERSION", toml: "version", description: "Version"},
}

var logger *slog.Logger

func GetParameters() Params {
	var params Params
	params = initParams(params)
	params = populateParams(params)
	return params
}

//nolint:gocyclo
func populateParams(params Params) Params {
	var productIds []string
	var defaultMirrors []string
	for _, product := range lib.GetAllProducts() {
		productIds = append(productIds, product.GetId())
		defaultMirrors = append(defaultMirrors, fmt.Sprintf("%s: %s", product.GetName(), product.GetDefaultMirrorUrl()))
	}

	// String params
	getopt.StringVarLong(&params.Arch, "arch", 'A', fmt.Sprintf("Override CPU architecture type for downloaded binary. Ex: `tfswitch --arch amd64` will attempt to download the amd64 version of the binary. Default: %s", runtime.GOARCH))
	getopt.StringVarLong(&params.ChDirPath, "chdir", 'c', "Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory")
	getopt.StringVarLong(&params.CustomBinaryPath, "bin", 'b', "Custom binary path. Ex: tfswitch -b "+lib.ConvertExecutableExt("/Users/username/bin/terraform"))
	getopt.StringVarLong(&params.DefaultVersion, "default", 'd', "Default to this version in case no other versions could be detected. Ex: tfswitch --default 1.2.4")
	getopt.StringVarLong(&params.InstallPath, "install", 'i', "Custom install path. Ex: tfswitch -i /Users/username. The binaries will be in the sub installDir directory e.g. /Users/username/"+lib.InstallDir)
	getopt.StringVarLong(&params.LatestPre, "latest-pre", 'p', "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.LatestStable, "latest-stable", 's', "Latest implicit version based on a constraint. Ex: tfswitch --latest-stable 0.13.0 downloads 0.13.7 and 0.13 downloads 0.15.5 (latest)")
	getopt.StringVarLong(&params.LogLevel, "log-level", 'g', "Set loglevel for tfswitch. One of (ERROR, INFO, NOTICE, DEBUG, TRACE)")
	getopt.StringVarLong(&params.MirrorURL, "mirror", 'm', "install from a remote API other than the default. Default (based on product):\n"+strings.Join(defaultMirrors, "\n"))
	getopt.StringVarLong(&params.ShowLatestPre, "show-latest-pre", 'P', "Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.ShowLatestStable, "show-latest-stable", 'S', "Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)")
	getopt.StringVarLong(&params.Product, "product", 't', fmt.Sprintf("Specifies which product to use. Ex: `tfswitch --product opentofu` will install OpenTofu. Options: (%s). Default: %s", strings.Join(productIds, ", "), lib.DefaultProductId))

	// Bool params
	getopt.BoolVarLong(&params.DryRun, "dry-run", 'r', "Only show what tfswitch would do. Don't download anything")
	getopt.BoolVarLong(&params.ForceColor, "force-color", 'K', "Force color output if terminal supports it")
	getopt.BoolVarLong(&params.HelpFlag, "help", 'h', "Displays help message")
	getopt.BoolVarLong(&params.LatestFlag, "latest", 'u', "Get latest stable version")
	getopt.BoolVarLong(&params.ListAllFlag, "list-all", 'l', "List all versions of terraform - including beta and rc")
	getopt.BoolVarLong(&params.NoColor, "no-color", 'k', "Disable color output. Useful for piping output to a file or when the terminal does not support colors")
	getopt.BoolVarLong(&params.ShowLatestFlag, "show-latest", 'U', "Show latest stable version")
	getopt.BoolVarLong(&params.VersionFlag, "version", 'v', "Displays the version of tfswitch")

	// Parse the command line parameters to fetch stuff like chdir
	getopt.Parse()

	isNotShortRun := !params.VersionFlag && !params.HelpFlag

	if isNotShortRun {
		if params.ForceColor && params.NoColor {
			os.Setenv("NO_COLOR", "true") // Prefer no color in this case
			logger = lib.InitLogger(params.LogLevel)
			logger.Fatal("(init) Cannot force color and disable color at the same time. Please choose either of them.")
		}

		if params.ForceColor {
			os.Setenv("FORCE_COLOR", "true")
		} else if params.NoColor {
			os.Setenv("NO_COLOR", "true")
		}

		oldForceColor := params.ForceColor
		oldNoColor := params.NoColor

		oldLogLevel := params.LogLevel
		logger = lib.InitLogger(params.LogLevel)

		var err error
		// Read configuration files
		// TOML from Homedir
		if tomlFileExists(params) {
			params, err = getParamsTOML(params)
			if err != nil {
				logger.Fatalf("Failed to obtain settings from TOML config in home directory: %v", err)
			}

			if params.ForceColor && params.NoColor {
				logger.Fatal("(toml) Cannot force color and disable color at the same time. Please choose either of them.")
			}

			if !oldForceColor && params.ForceColor {
				os.Setenv("FORCE_COLOR", "true")
			} else if !oldNoColor && params.NoColor {
				os.Setenv("NO_COLOR", "true")
			}
			// Re-init logger
			logger = lib.InitLogger(params.LogLevel)

			// Update assignment in case values changed by TOML
			oldForceColor = params.ForceColor
			oldNoColor = params.NoColor
		}

		// First pass to obtain environment variables to override product
		params = GetParamsFromEnvironment(params)

		if params.ForceColor && params.NoColor {
			logger.Fatal("(env) Cannot force color and disable color at the same time. Please choose either of them.")
		}

		if !oldForceColor && params.ForceColor {
			os.Setenv("FORCE_COLOR", "true")
		} else if !oldNoColor && params.NoColor {
			os.Setenv("NO_COLOR", "true")
		}
		// Re-init logger
		logger = lib.InitLogger(params.LogLevel)

		// Set defaults based on product
		// This must be performed after TOML file, to obtain product.
		// But the mirror URL, if set to default product URL,
		// is used by some of the version getter methods, to
		// obtain list of versions.
		product := lib.GetProductById(params.Product)
		if product == nil {
			logger.Fatalf("Invalid \"product\" configuration value: %q", params.Product)
		} else { // Use else as there is a warning that params maybe nil, as it does not see Fatalf as a break condition
			if params.MirrorURL == "" {
				params.MirrorURL = product.GetDefaultMirrorUrl()
				logger.Debugf("Default mirror URL: %q", params.MirrorURL)
			}

			// Set default bin directory, if not configured
			if params.CustomBinaryPath == "" {
				if runtime.GOOS == "windows" {
					params.CustomBinaryPath = filepath.Join(lib.GetHomeDirectory(), "bin", lib.ConvertExecutableExt(product.GetExecutableName()))
				} else {
					params.CustomBinaryPath = filepath.Join("/usr/local/bin", product.GetExecutableName())
				}
			}
			params.ProductEntity = product
		}

		if tfSwitchFileExists(params) {
			params, err = GetParamsFromTfSwitch(params)
			if err != nil {
				logger.Fatalf("Failed to obtain settings from \".tfswitch\" file: %v", err)
			}
		}

		if terraformVersionFileExists(params) {
			params, err = GetParamsFromTerraformVersion(params)
			if err != nil {
				logger.Fatalf("Failed to obtain settings from \".terraform-version\" file: %v", err)
			}
		}

		if isTerraformModule(params) {
			params, err = GetVersionFromVersionsTF(params)
			if err != nil {
				logger.Fatalf("Failed to obtain settings from Terraform module: %v", err)
			}
		}

		if terraGruntFileExists(params) {
			params, err = GetVersionFromTerragrunt(params)
			if err != nil {
				logger.Fatalf("Failed to obtain settings from Terragrunt configuration: %v", err)
			}
		}

		params = GetParamsFromEnvironment(params)

		// Logger config was changed by the config files. Reinitialise.
		if params.LogLevel != oldLogLevel {
			logger = lib.InitLogger(params.LogLevel)
		}
	}

	// Parse again to overwrite anything that might by defined on the cli AND in any config file (CLI always wins)
	getopt.Parse()
	args := getopt.Args()
	if len(args) == 1 {
		/* version provided on command line as arg */
		params.Version = args[0]
	}

	if isNotShortRun {
		if params.DryRun {
			logger.Info("[DRY-RUN] No changes will be made")
		} else {
			logger.Debugf("Resolved dry-run: %t", params.DryRun)
		}

		logger.Debugf("Resolved CPU architecture: %q", params.Arch)
		if params.DefaultVersion != "" {
			logger.Debugf("Resolved fallback version: %q", params.DefaultVersion)
		}
		logger.Debugf("Resolved binary path: %q", params.CustomBinaryPath)
		logger.Debugf("Resolved force color: %t", params.ForceColor)
		logger.Debugf("Resolved install path: %q", filepath.Join(params.InstallPath, lib.InstallDir))
		logger.Debugf("Resolved install version: %q", params.Version)
		logger.Debugf("Resolved log level: %q", params.LogLevel)
		logger.Debugf("Resolved mirror URL: %q", params.MirrorURL)
		logger.Debugf("Resolved no color: %t", params.NoColor)
		logger.Debugf("Resolved product name: %q", params.Product)
		logger.Debugf("Resolved working directory: %q", params.ChDirPath)
	}

	return params
}

func initParams(params Params) Params {
	params.Arch = runtime.GOARCH
	params.ChDirPath = lib.GetCurrentDirectory()
	params.CustomBinaryPath = ""
	params.DefaultVersion = lib.DefaultLatest
	params.DryRun = false
	params.ForceColor = false
	params.HelpFlag = false
	params.InstallPath = lib.GetHomeDirectory()
	params.LatestFlag = false
	params.LatestPre = lib.DefaultLatest
	params.LatestStable = lib.DefaultLatest
	params.ListAllFlag = false
	params.LogLevel = "INFO"
	params.MirrorURL = ""
	params.NoColor = false
	params.ShowLatestFlag = false
	params.ShowLatestPre = lib.DefaultLatest
	params.ShowLatestStable = lib.DefaultLatest
	params.TomlDir = lib.GetHomeDirectory()
	params.Version = lib.DefaultLatest
	params.Product = lib.DefaultProductId
	params.VersionFlag = false
	return params
}
