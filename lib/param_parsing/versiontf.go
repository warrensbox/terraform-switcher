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

func getVersionConstraintsFromHCLFile(fileName string, hclFile *hcl.File) ([]string, error) {
	var constraints []string

	// Define the schema for the terraform block
	terraformBlockSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: terraformBlockType},
		},
	}

	content, _, diags := hclFile.Body.PartialContent(terraformBlockSchema)
	if diags.HasErrors() {
		logger.Debugf("No %s blocks found in %q", terraformBlockType, fileName)
		return constraints, nil
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
				logger.Debugf("Error getting attributes from %q block in %q: %v", terraformBlockType, fileName, attrDiags.Error())
				continue
			}

			if attr, exists := blockContent.Attributes[requiredVersionAttrName]; exists {
				val, valDiags := attr.Expr.Value(nil)
				if valDiags.HasErrors() {
					logger.Debugf("Error evaluating %q in %q: %v", requiredVersionAttrName, fileName, valDiags.Error())
					continue
				}
				if !val.IsKnown() || !val.Type().Equals(cty.String) {
					logger.Debugf("Skipping not known or non-string value of %q in %q: %#v", requiredVersionAttrName, fileName, val)
					continue
				}
				versionStr := val.AsString()
				if versionStr == "" {
					logger.Debugf("Skipping empty %q in %q", requiredVersionAttrName, fileName)
					continue
				}
				constraint, constraintErr := semver.NewConstraint(versionStr)
				if constraintErr != nil {
					logger.Errorf("Invalid version constraint found in %q: %q", fileName, versionStr)
					return nil, constraintErr
				}
				logger.Debugf("Found %q %q in %q", requiredVersionAttrName, constraint.String(), fileName)
				constraints = append(constraints, constraint.String())
			}
		}
	}

	return constraints, nil
}

func getVersionConstraintsFromFiles(filesPath []string) (string, error) {
	parser := hclparse.NewParser()
	for _, filePath := range filesPath {
		_, diagnostics := parser.ParseHCLFile(filePath)
		if diagnostics.HasErrors() {
			return "", fmt.Errorf("Could not parse HCL file %q: %v", filePath, diagnostics.Error())
		}
	}

	var constraints []string
	for fileName, hclFile := range parser.Files() {
		parsedConstraints, err := getVersionConstraintsFromHCLFile(fileName, hclFile)
		if err != nil {
			return "", err
		}
		constraints = append(constraints, parsedConstraints...)
	}

	return strings.Join(constraints, ", "), nil
}

func getConstraintFromVersionsTF(params Params) (Params, error) {
	relPath, err := lib.GetRelativePath(params.ChDirPath)
	if err != nil {
		return params, err
	}

	logger.Infof("Reading version constraint from %s at %q", paramTypeVersionTF, relPath)

	extensionsPerProduct := lib.GetProductById(params.Product).GetFileExtensions()
	var hclFiles []string
	var fileGlobs []string
	for _, ext := range extensionsPerProduct {
		globPattern := fmt.Sprintf("*.%s", ext)
		fileGlobs = append(fileGlobs, globPattern)
		files, globErr := filepath.Glob(filepath.Join(relPath, globPattern))
		if globErr != nil {
			return params, fmt.Errorf("Could not list %s files in %q: %v", globPattern, relPath, globErr)
		}
		hclFiles = append(hclFiles, files...)
	}

	if len(hclFiles) == 0 {
		logger.Debugf("No %s files found in %q", strings.Join(fileGlobs, ", "), relPath)
		return params, nil
	}

	versionRequirements, err := getVersionConstraintsFromFiles(hclFiles)
	if err != nil {
		return params, err
	}

	if versionRequirements == "" {
		logger.Debugf("No version requirements found in %s files in %q", strings.Join(fileGlobs, ", "), relPath)
		return params, nil
	}

	params.VersionRequirement = versionRequirements
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
