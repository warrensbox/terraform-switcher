// Suppressing linter warnings for this package:
// - revive: FIXME: don't use an underscore in package name
// - staticcheck: ST1005: error strings should not be capitalized (staticcheck)
//
//nolint:revive,staticcheck
package param_parsing

import (
	"fmt"
	"path/filepath"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	"github.com/warrensbox/terraform-switcher/lib"
)

const (
	paramTypeVersionTF      = "Terraform/OpenTofu module"
	terraformBlockType      = "terraform"
	requiredVersionAttrName = "required_version"
)

// getRequiredVersionsFromFile parses a single hcl file and extracts all required_version constraints
func getRequiredVersionsFromFile(filePath string) ([]string, error) {
	var versions []string

	parser := hclparse.NewParser()
	hclFile, diagnostics := parser.ParseHCLFile(filePath)
	if diagnostics.HasErrors() {
		logger.Errorf("Unable to parse HCL file %q: %v", filePath, diagnostics.Error())
		return nil, fmt.Errorf("Could not parse HCL file %q: %v", filePath, diagnostics.Error())
	}

	// Define the schema for the terraform block
	terraformBlockSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: terraformBlockType},
		},
	}

	content, _, diags := hclFile.Body.PartialContent(terraformBlockSchema)
	if diags.HasErrors() {
		logger.Debugf("No %s blocks found in %q", terraformBlockType, filePath)
		return versions, nil
	}

	for _, block := range content.Blocks {
		if block.Type == terraformBlockType {
			// Extract required_version from the terraform block
			terraformAttributesSchema := &hcl.BodySchema{
				Attributes: []hcl.AttributeSchema{
					{Name: requiredVersionAttrName},
				},
			}
			blockContent, _, attrDiags := block.Body.PartialContent(terraformAttributesSchema)
			if attrDiags.HasErrors() {
				logger.Debugf("Error getting attributes from %q block in %q: %v", terraformBlockType, filePath, attrDiags.Error())
				continue
			}

			if attr, exists := blockContent.Attributes[requiredVersionAttrName]; exists {
				val, valDiags := attr.Expr.Value(nil)
				if valDiags.HasErrors() {
					logger.Debugf("Error evaluating %q in %q: %v", requiredVersionAttrName, filePath, valDiags.Error())
					continue
				}
				if !val.IsKnown() || !val.Type().Equals(cty.String) {
					logger.Debugf("Skipping not known or non-string value of %q at %q: %q", requiredVersionAttrName, filePath, val)
					continue
				}
				versionStr := val.AsString()
				if versionStr != "" {
					logger.Debugf("Found %q %q in %q", requiredVersionAttrName, versionStr, filePath)
					versions = append(versions, versionStr)
				}
			}
		}
	}

	return versions, nil
}

func getConstraintFromVersionsTF(params Params) (Params, error) {
	var tfConstraints []string

	relPath, err := lib.GetRelativePath(params.ChDirPath)
	if err != nil {
		return params, err
	}

	logger.Infof("Reading version constraint from %s at %q", paramTypeVersionTF, relPath)

	extensionsPerProduct := lib.GetProductById(params.Product).GetFileExtensions()
	var hclFiles []string
	for _, ext := range extensionsPerProduct {
		files, globErr := filepath.Glob(filepath.Join(relPath, fmt.Sprintf("*.%s", ext)))
		if globErr != nil {
			return params, fmt.Errorf("Could not list %s files in %q: %v", ext, relPath, globErr)
		}
		hclFiles = append(hclFiles, files...)
	}

	if len(hclFiles) == 0 {
		logger.Debugf("No %s files found in %q", strings.Join(extensionsPerProduct, ", "), relPath)
		return params, nil
	}

	// Parse each file and collect required_version constraints
	for _, hclFile := range hclFiles {
		if !lib.CheckFileExist(hclFile) {
			continue
		}

		versions, parseErr := getRequiredVersionsFromFile(hclFile)
		if parseErr != nil {
			logger.Errorf("Error parsing %s file %q: %v", paramTypeVersionTF, hclFile, parseErr)
			return params, parseErr
		}

		for _, v := range versions {
			// Check if the version constraint is valid
			constraint, constraintErr := semver.NewConstraint(v)
			if constraintErr != nil {
				logger.Errorf("Invalid version constraint found: %q", v)
				return params, constraintErr
			}
			tfConstraints = append(tfConstraints, constraint.String())
		}
	}

	if len(tfConstraints) == 0 {
		logger.Debugf("No version constraint found in %s configuration", paramTypeVersionTF)
		return params, nil
	}

	params.VersionRequirement = strings.Join(tfConstraints, ", ")
	logger.Debugf("Using version constraint from %s at %q: %q", paramTypeVersionTF, relPath, params.VersionRequirement)
	return params, nil
}

func GetVersionFromVersionsTF(params Params) (Params, error) {
	params, err := getConstraintFromVersionsTF(params)
	if err != nil {
		return params, err
	}

	// If parsing was successful but no version constraint was found, return as is
	if params.VersionRequirement == "" {
		return params, nil
	}

	if params.MatchVersionRequirement == "" {
		version, err2 := lib.GetSemver(params.VersionRequirement, params.MirrorURL)
		if err2 != nil {
			logger.Errorf("No version found matching %q", params.VersionRequirement)
			return params, err2
		}
		params.Version = version
		logger.Debugf("Using version from %s: %q", paramTypeVersionTF, params.Version)
	}
	return params, nil
}
