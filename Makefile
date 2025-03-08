EXE    := tfswitch
PKG    := github.com/warrensbox/terraform-switcher
PATH   := build:$(PATH)
VER    ?= $(shell git ls-remote --tags --sort=version:refname git@github.com:warrensbox/terraform-switcher.git | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | tail -1)
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o $@ $(PKG)

.PHONY: release
release: $(EXE) darwin linux windows

.PHONY: darwin linux windows
darwin linux windows:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o $(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -vrf $(EXE) $(EXE)-*-*-* build/

.PHONY: test
test: vet $(EXE)
	mkdir -p build
	mv $(EXE) build/ # can't figure what's this for (also `PATH' var) (c) @yermulnik 01-Mar-2025
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: install
install: $(EXE)
	mkdir -p ~/bin
	mv $(EXE) ~/bin

.PHONY: docs
docs:
	@#cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..
	cd www && mkdocs gh-deploy --force

.PHONY: goreleaser-release-snapshot
goreleaser-release-snapshot:
	RELEASE_VERSION=$(VER) goreleaser release --config ./.goreleaser.yml --snapshot --clean
