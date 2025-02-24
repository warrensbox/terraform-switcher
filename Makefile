EXE  := tfswitch
PKG  := github.com/warrensbox/terraform-switcher
VER := $(shell git ls-remote --tags --sort=version:refname git@github.com:warrensbox/terraform-switcher.git | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | tail -1)
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o $@ $(PKG)

.PHONY: release
release: $(EXE) darwin linux

.PHONY: darwin linux
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o $(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -f $(EXE) $(EXE)-*-*-*

.PHONY: test
test: vet $(EXE)
	mv $(EXE) build
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
	#cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..
	cd www && mkdocs gh-deploy --force

