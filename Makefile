# go option
GO         ?= go
PKG        := go mod vendor
LDFLAGS    := -w -s
GOFLAGS    :=
TAGS       := 
BINDIR     := $(CURDIR)/bin
APP_VERSION := v0.0.0

# Required for globs to work correctly
SHELL=/bin/bash


.PHONY: all
all: build

.PHONY: dep
dep:
	@$(PKG)

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=1 $(GO) build -o $(BINDIR)/ -v $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./

.PHONY: release
release:
	git tag -d "$(APP_VERSION)" || true
	git tag "$(APP_VERSION)"
	./scripts/gobin.sh github.com/goreleaser/goreleaser release --skip-publish --rm-dist

.PHONY: fmt
fmt:
	goimports -w ./
