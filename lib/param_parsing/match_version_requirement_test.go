//nolint:revive // FIXME: don't use an underscore in package name
package param_parsing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pborman/getopt"
)

func TestMatchVersionRequirement_match(t *testing.T) {
	var actual bool
	var err error
	expected := true
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Log("Testing match with no requirement")
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatal("Version requirement not matched (unexpected)")
	}
	t.Log("Version requirement matched (expected)")

	t.Log("Testing match with default fallback version")
	params.DefaultVersion = "1.0.0"
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatal("Version requirement not matched (unexpected)")
	}
	t.Log("Version requirement matched (expected)")

	t.Log("Testing match with explicit version")
	params.Version = "1.0.0"
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatal("Version requirement not matched (unexpected)")
	}
	t.Log("Version requirement matched (expected)")
}

func TestMatchVersionRequirement_mismatch(t *testing.T) {
	var actual bool
	var err error
	expected := false
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Log("Testing mismatch with default fallback version")
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.DefaultVersion = "1.0.1"
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatal("Version requirement not mismatched (unexpected)")
	}
	t.Log("Version requirement mismatched (expected)")

	t.Log("Testing mismatch with explicit version")
	params.Version = "1.0.2"
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatal("Version requirement not mismatched (unexpected)")
	}
	t.Log("Version requirement mismatched (expected)")
}

func TestMatchVersionRequirement_arg_validation(t *testing.T) {
	var err error
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Log("Testing argument validation error")
	os.Args = []string{"cmd", "--match-version-requirement=incorrect_version_string"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params = populateParams(params)
	_, err = MatchVersionRequirement(params)

	if err == nil {
		t.Fatal("Argument validation error not caught (not expected)")
	}
	t.Log("Argument validation error caught (expected)")
}

func TestMatchVersionRequirement_match_toml(t *testing.T) {
	var actual bool
	var err error
	expected := true
	path := "../../test-data/integration-tests/test_tfswitchtoml/.tfswitch.toml"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing match with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.6.2"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.TomlDir = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not matched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement matched (expected)", kind)
}

func TestMatchVersionRequirement_mismatch_toml(t *testing.T) {
	var actual bool
	var err error
	expected := false
	path := "../../test-data/integration-tests/test_tfswitchtoml/.tfswitch.toml"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing mismatch with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.TomlDir = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not mismatched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement mismatched (expected)", kind)
}

func TestMatchVersionRequirement_match_terraform_version(t *testing.T) {
	var actual bool
	var err error
	expected := true
	path := "../../test-data/integration-tests/test_terraform-version/.terraform-version"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing match with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=0.11.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not matched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement matched (expected)", kind)
}

func TestMatchVersionRequirement_mismatch_terraform_version(t *testing.T) {
	var actual bool
	var err error
	expected := false
	path := "../../test-data/integration-tests/test_terraform-version/.terraform-version"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing mismatch with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not mismatched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement mismatched (expected)", kind)
}

func TestMatchVersionRequirement_match_terragrunt(t *testing.T) {
	var actual bool
	var err error
	expected := true
	path := "../../test-data/integration-tests/test_terragrunt_hcl/terragrunt.hcl"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing match with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=0.13.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not matched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement matched (expected)", kind)
}

func TestMatchVersionRequirement_mismatch_terragrunt(t *testing.T) {
	var actual bool
	var err error
	expected := false
	path := "../../test-data/integration-tests/test_terragrunt_hcl/terragrunt.hcl"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing mismatch with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not mismatched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement mismatched (expected)", kind)
}

func TestMatchVersionRequirement_match_tfswitchrc(t *testing.T) {
	var actual bool
	var err error
	expected := true
	path := "../../test-data/integration-tests/test_tfswitchrc/.tfswitchrc"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing match with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=0.10.5"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not matched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement matched (expected)", kind)
}

func TestMatchVersionRequirement_mismatch_tfswitchrc(t *testing.T) {
	var actual bool
	var err error
	expected := false
	path := "../../test-data/integration-tests/test_tfswitchrc/.tfswitchrc"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing mismatch with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.0.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not mismatched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement mismatched (expected)", kind)
}

func TestMatchVersionRequirement_match_module(t *testing.T) {
	var actual bool
	var err error
	expected := true
	path := "../../test-data/integration-tests/test_versiontf/version.tf"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing match with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.0.1"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not matched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement matched (expected)", kind)
}

func TestMatchVersionRequirement_mismatch_module(t *testing.T) {
	var actual bool
	var err error
	expected := false
	path := "../../test-data/integration-tests/test_versiontf/version.tf"
	kind := filepath.Base(path)
	params := Params{}

	t.Cleanup(func() {
		getopt.CommandLine = getopt.New()
	})

	t.Logf("Testing mismatch with %s", kind)
	os.Args = []string{"cmd", "--match-version-requirement=1.1.0"}
	params = initParams(params)
	params.LogLevel = "INFO"
	params.ChDirPath = filepath.Dir(path)
	params = populateParams(params)
	actual, err = MatchVersionRequirement(params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual != expected {
		t.Fatalf("[%s] Version requirement not mismatched (unexpected)", kind)
	}
	t.Logf("[%s] Version requirement mismatched (expected)", kind)
}
