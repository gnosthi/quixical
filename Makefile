# QUIXICAL Makefile
# default is make all.
GOFILES_BUILD           := $(shell find . -type f -name '*.go' -not -name '*_test.go')
DATE                    := $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" '+%FT%T%z' 2>/dev/null || date -u '+%FT%T%z')
QUIXICAL_VERSION        ?= $(shell cat VERSION)
QUIXICAL_REVISION       := $(shell cat COMMIT 2>/dev/null || git rev-parse --short=8 HEAD)
QUIXICAL_OUTPUT         ?= quixical
PWD                     := $(shell pwd)
PREFIX                  ?= $(GOPATH)
BUILDFLAGS              := -ldflags="-s -w -X main.version=$(QUIXICAL_VERSION) -X main.commit=$(QUIXICAL_REVISION) -X main.date=$(DATE)" -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)" -buildmode=pie
BINDIR                  ?= $(PREFIX)/bin
GO                      := go
GOOS                    ?= $(shell go version | cut -d ' ' -f4 | cut -d '/' -f1)
GOARCH                  ?= $(shell go version | cut -d ' ' -f4 | cut -d '/' -f2)

OK                      := $(shell tput setaf 6; echo ' [OK]'; tput sgr0;)

all: build
build: $(QUIXICAL_OUTPUT)
travis: sysinfo crosscompile build install test

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo -n "     PLATFORM: $(shell uname -a)"
	@printf '%s\n' '$(OK)'
	@echo -n "     PWD:      : $(shell pwd)"
	@printf "%s\n" '$(OK)'
	@echo -n "     GO        : $(shell go version)"
	@printf '%s\n' '$(OK)'
	@echo -n "      BUILDFLAGS: $(BUILDFLAGS)"
	@printf '%s\n'  '$(OK)'
	@echo -n "      GIT       : $(shell git version)"
	@printf '%s\n'  '$(OK)'

clean:
	@echo -n ">> CLEAN"
	@$(GO) clean -i

$(QUIXICAL_OUTPUT): $(GOFILES_BUILD)
    ## Removed $(QUIXICAL_REVISION) as it was causing weird issues.
    ## Readded $(QUIXICAL_REVISION) forgot to define $(BUILDFLAGS).
	@echo -n ">> BUILD, version = $(QUIXICAL_VERSION/$(QUIXICAL_REVISION)), output = $@)"
	@$(GO) build -o $@ $(BUILDFLAGS)
	@printf '%s\n' '$(OK)'

install: all
	@echo -n ">> INSTALL, version = $(QUIXICAL_VERSION)"
	@install -m 0755 -d $(DESTDIR)$(BINDIR)
	@install -m 0755 $(QUIXICAL_OUTPUT) $(DESTDIR)$(BINDIR)/quixical
	@printf '%s\n' '$(OK)'

test: $(QUIXICAL_OUTPUT)
	@echo ">> TEST, \"fast-mode\": race detector off"
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg, $(PKGS),\
		echo -n "    ";\
		$(GO) test $(BUILDFLAGS) -coverprofile=coverage.out -covermode=count $(pkg)||exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	@$(GO) tool cover -html=coverage-all.out -o coverage-all.html

crosscompile:
	@echo -n ">> CROSSCOMPILE linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(QUIXICAL_OUTPUT)-linux-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(QUIXICAL_OUTPUT)-macOS-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(QUIXICAL_OUTPUT)-windows-amd64
	@printf '%s\n' '$(OK)'

check-release-env:
	ifndef GITHUB_TOKEN
		$(error GITHUB_TOKEN is undefined)
	endif

release: goreleaser

goreleaser: check-release-env travis clean
	@echo ">> RELEASE, goreleaser"
	@goreleaser

docker-test:
    docker build -t quixical:$(QUIXICAL_REVISION) .
    docker run --rm quixical:$(QUIXICAL_REVISION) make test

