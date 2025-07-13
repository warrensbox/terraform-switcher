EXE       := tfswitch
PKG       := github.com/warrensbox/terraform-switcher
BUILDPATH := build
PATH      := $(BUILDPATH):$(PATH)
VER       ?= $(shell git ls-remote --tags --sort=version:refname git@github.com:warrensbox/terraform-switcher.git | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | tail -1)

# Managing Go installations: Installing multiple Go versions
# https://go.dev/doc/manage-install
GOBINARY ?= $(shell (egrep -m1 '^go[[:space:]]+[[:digit:]]+\.' go.mod | tr -d '[:space:]' | xargs which) || which go)

# Go vars
GOOS   ?= $(shell $(GOBINARY) env GOOS)
GOARCH ?= $(shell $(GOBINARY) env GOARCH)

# Vars for `super-linter` target
CONTAINER_ENGINE ?= $(shell command -v podman || command -v docker || echo "NONE")
# Git dirs for `super-linter` target
GIT_COMMON_DIR   := $(shell git rev-parse --git-common-dir)
GIT_DIR          := $(shell git rev-parse --git-dir)
GIT_TOPLEVEL_DIR := $(shell git rev-parse --show-toplevel)


$(EXE): version go.mod *.go lib/*.go
	mkdir -p "$(BUILDPATH)/"
	$(GOBINARY) build -v -ldflags "-X 'main.version=$(VER)'" -o "$(BUILDPATH)/$@" $(PKG)

.PHONY: release
release: $(EXE) darwin linux windows

.PHONY: darwin linux windows
darwin linux windows: version
	GOOS=$@ $(GOBINARY) build -ldflags "-X 'main.version=$(VER)'" -o "$(BUILDPATH)/$(EXE)-$(word 1, $(VER))-$@-$(GOARCH)" $(PKG)

.PHONY: clean
clean:
	rm -vrf "$(BUILDPATH)/"

.PHONY: test
test: vet vulncheck $(EXE)
	$(GOBINARY) test -v ./...

.PHONY: test-single-function
test-single-function: vet
ifndef TEST_FUNC_NAME
	$(error TEST_FUNC_NAME is not set)
endif
	$(eval TEST_FUNC_NAME := ^$(TEST_FUNC_NAME)$$)
	$(info Testing function matching "$(TEST_FUNC_NAME)" regexp)
	$(GOBINARY) test -v -run="$(TEST_FUNC_NAME)" ./...

.PHONY: vet
vet: version
	$(GOBINARY) vet ./...

.PHONY: version
version:
	@echo "Running $(GOBINARY) ($(shell $(GOBINARY) version))"

.PHONY: vulncheck
vulncheck:
	@command -v govulncheck >/dev/null 2>&1 && govulncheck -show color ./... || echo "govulncheck not found, skipping vulnerability check"

.PHONY: vulncheck-verbose
vulncheck-verbose:
	@command -v govulncheck >/dev/null 2>&1 && govulncheck -show traces,color,version,verbose ./... || echo "govulncheck not found, skipping vulnerability check"

.PHONY: install
install: $(EXE)
	mkdir -p ~/bin
	mv "$(BUILDPATH)/$(EXE)" ~/bin/

.PHONY: docs-build
docs-build:
	cd www && mkdocs build

.PHONY: docs-deploy
docs-deploy:
	cd www && mkdocs gh-deploy --force

.PHONY: docs-serve
docs-serve:
	cd www && mkdocs serve

.PHONY: docs-serve-public
docs-serve-public:
	cd www && mkdocs serve --dev-addr 0.0.0.0:8000

.PHONY: goreleaser-release-snapshot
goreleaser-release-snapshot:
	RELEASE_VERSION="$(VER)" goreleaser release --config ./.goreleaser.yml --snapshot --clean

.PHONY: lint
lint: super-linter

.PHONY: super-linter
super-linter:
ifeq ($(CONTAINER_ENGINE),NONE)
	$(error "No container engine found. Please install Podman or Docker.")
else
# Assume that if GIT_COMMON_DIR differs from GIT_DIR, then we're in Git worktree
# In this case, mount Git dir (`.git`) as a volume
ifneq ($(GIT_COMMON_DIR),$(GIT_DIR))
	$(eval GIT_COMMON_DIR_VOLUME := --volume "$(GIT_COMMON_DIR):$(GIT_COMMON_DIR)")
endif
	# Keep `--env' vars below the `VALIDATE_ALL_CODEBASE' in sync with .github/workflows/super-linter.yml
	$(CONTAINER_ENGINE) run \
		--name super-linter \
		--volume "$(GIT_TOPLEVEL_DIR):$(GIT_TOPLEVEL_DIR)" $(GIT_COMMON_DIR_VOLUME)\
		--workdir "$(GIT_TOPLEVEL_DIR)" \
		--env GITHUB_WORKSPACE="$(GIT_TOPLEVEL_DIR)" \
		--env RUN_LOCAL=true \
		--env VALIDATE_ALL_CODEBASE=false \
		--env FILTER_REGEX_EXCLUDE="^$(GIT_TOPLEVEL_DIR)/test-data/" \
		--env BASH_EXEC_IGNORE_LIBRARIES=true \
		--env VALIDATE_GO=false \
		--env VALIDATE_JSCPD=false \
		--rm ghcr.io/super-linter/super-linter:slim-latest
endif
