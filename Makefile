EXE       := tfswitch
PKG       := github.com/warrensbox/terraform-switcher
BUILDPATH := build
PATH      := $(BUILDPATH):$(PATH)
VER       ?= $(shell git ls-remote --tags --sort=version:refname git@github.com:warrensbox/terraform-switcher.git | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | tail -1)
# Managing Go installations: Installing multiple Go versions
# https://go.dev/doc/manage-install
GOBINARY  ?= $(shell (egrep -m1 '^go[[:space:]]+[[:digit:]]+\.' go.mod | tr -d '[:space:]' | xargs which) || echo go)
GOOS      ?= $(shell $(GOBINARY) env GOOS)
GOARCH    ?= $(shell $(GOBINARY) env GOARCH)

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
test: vet $(EXE)
	$(GOBINARY) test -v ./...

.PHONY: test-single-function
test-single-function: vet
	@([ -z "$(TEST_FUNC_NAME)" ] && echo "TEST_FUNC_NAME is not set" && false) || true
	$(GOBINARY) test -v -run="$(TEST_FUNC_NAME)" ./...

.PHONY: vet
vet: version
	$(GOBINARY) vet ./...

.PHONY: version
version:
	@echo "Running $(GOBINARY) ($(shell $(GOBINARY) version))"

.PHONY: install
install: $(EXE)
	mkdir -p ~/bin
	mv "$(BUILDPATH)/$(EXE)" ~/bin/

.PHONY: docs
docs:
	@#cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..
	cd www && mkdocs gh-deploy --force

.PHONY: goreleaser-release-snapshot
goreleaser-release-snapshot:
	RELEASE_VERSION=$(VER) goreleaser release --config ./.goreleaser.yml --snapshot --clean
