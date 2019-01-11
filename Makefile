# Go tools.
GO ?= GO111MODULE=off go
GO_MD2MAN ?= go-md2man

# Paths.
PROJECT := github.com/kubic-project/caasp-init
CMD := .

# We use Docker because Go is just horrific to deal with.
CASSPINIT_IMAGE := caasp-init_dev
DOCKER_RUN := docker run --rm -it --security-opt apparmor:unconfined --security-opt label:disable -v ${PWD}:/go/src/${PROJECT}

# Output directory.
BUILD_DIR ?= ./bin

# Release information.
GPG_KEYID ?=

# Version information.
VERSION := $(shell cat VERSION)
COMMIT_NO := $(shell git rev-parse HEAD 2> /dev/null || true)
COMMIT := $(if $(shell git status --porcelain --untracked-files=no),"${COMMIT_NO}-dirty","${COMMIT_NO}")

# Get current Version changelog
CHANGE := $(shell sed -e '1,/$(VERSION)/d;/v.*/Q' ./CHANGELOG.md)

BUILD_FLAGS ?=

BASE_FLAGS := ${BUILD_FLAGS} -tags "${BUILDTAGS}"

BASE_LDFLAGS := -X $(PROJECT)/cmd.version=$(VERSION)
BASE_LDFLAGS += -X $(PROJECT)/cmd.gitCommit=$(COMMIT)

DYN_BUILD_FLAGS := ${BASE_FLAGS} -buildmode=pie -ldflags "${BASE_LDFLAGS}"
TEST_BUILD_FLAGS := ${BASE_FLAGS} -buildmode=pie -ldflags "${BASE_LDFLAGS} -X ${PROJECT}/pkg/testutils.binaryType=test"
STATIC_BUILD_FLAGS := ${BASE_FLAGS} -ldflags "${BASE_LDFLAGS} -extldflags '-static'"

.DEFAULT: caasp-init

GO_SRC = $(shell find . -name \*.go)

# NOTE: If you change these make sure you also update local-validate-build.

caasp-init: $(GO_SRC)
	$(GO) build ${DYN_BUILD_FLAGS} -o $(BUILD_DIR)/$@ ${CMD}

caasp-init.static: $(GO_SRC)
	env CGO_ENABLED=0 $(GO) build ${STATIC_BUILD_FLAGS} -o $(BUILD_DIR)/$@ ${CMD}

install: $(GO_SRC)
	$(GO) install -v ${DYN_BUILD_FLAGS} ${CMD}

install.static: $(GO_SRC)
	$(GO) install -v ${STATIC_BUILD_FLAGS} ${CMD}

clean:
	rm -rf ./bin ./build ./release .tmp-validate coverage.txt
	rm -f $(MANPAGES)

local-validate: local-validate-git local-validate-go local-validate-reproducible

EPOCH_COMMIT ?= cd284b45372f452c56df31b2e105e805cdb8d55b
local-validate-git:
	@type git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found." && false)
ifdef TRAVIS_COMMIT_RANGE
	git-validation -q -run DCO,short-subject
else
	git-validation -q -run DCO,short-subject -range $(EPOCH_COMMIT)..HEAD
endif

local-validate-go:
	@type gofmt    >/dev/null 2>/dev/null || (echo "ERROR: gofmt not found." && false)
	test -z "$$(gofmt -s -l . | grep -vE '^vendor/|^third_party/' | tee /dev/stderr)"
	@type golint   >/dev/null 2>/dev/null || (echo "ERROR: golint not found." && false)
	test -z "$$(golint ./... | grep -v -E '^vendor/|^third_party/' | tee /dev/stderr)"
	@go doc cmd/vet >/dev/null 2>/dev/null || (echo "ERROR: go vet not found." && false)
	test -z "$$($(GO) vet $$($(GO) list $(PROJECT)/... | grep -vE '/vendor/|/third_party/') 2>&1 | tee /dev/stderr)"

# Make sure that our builds are reproducible even if you wait between them and
# the modified time of the files is different.
local-validate-reproducible:
	mkdir -p .tmp-validate
	make -B caasp-init && cp $(BUILD_DIR)/caasp-init .tmp-validate/caasp-init.a
	@echo sleep 10s
	@sleep 10s && touch $(GO_SRC)
	make -B caasp-init && cp $(BUILD_DIR)/caasp-init .tmp-validate/caasp-init.b
	diff -s .tmp-validate/caasp-init.{a,b}
	sha256sum .tmp-validate/caasp-init.{a,b}
	rm -r .tmp-validate/caasp-init.{a,b}

local-validate-build:
	$(GO) build ${DYN_BUILD_FLAGS} -o /dev/null ${CMD}
	env CGO_ENABLED=0 $(GO) build ${STATIC_BUILD_FLAGS} -o /dev/null ${CMD}
	$(GO) test -run nothing ${DYN_BUILD_FLAGS} $(PROJECT)/...

# Used for tests.
DOCKER_IMAGE :=kubic-project/amd64:tumbleweed

caasp-initimage:
	docker build -t $(CASSPINIT_IMAGE) .


test.unit: caasp-initimage
	$(DOCKER_RUN) $(CASSPINIT_IMAGE) make test

test: local-validate-go
	rm -rf /tmp/caasp-init
	$(GO) test -v ./...

cover:
	bash <scripts/cover.sh

dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist:
	rm -rf build/caasp-init/* release/*
	mkdir -p build/caasp-init/bin release/
	cp README.md LICENSE build/caasp-init
	GOOS=linux GOARCH=amd64 go build -o build/caasp-init/bin/caasp-init -ldflags="$(BASE_LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/caasp-init-linux.tgz caasp-init/
	GOOS=darwin GOARCH=amd64 go build -o build/caasp-init/bin/caasp-init -ldflags="$(BASE_LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/caasp-init-macos.tgz caasp-init/
	rm build/caasp-init/bin/caasp-init

release: dist
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif
	github-release release -u kubic-project -r caasp-init --tag $(VERSION)  --name $(VERSION) -s $(GITHUB_TOKEN) -d "$(CHANGE)"
	github-release upload -u kubic-project -r caasp-init --tag $(VERSION)  --name caasp-init-linux.tgz --file release/caasp-init-linux.tgz -s $(GITHUB_TOKEN)
	github-release upload -u kubic-project -r caasp-init --tag $(VERSION)  --name caasp-init-macos.tgz --file release/caasp-init-macos.tgz -s $(GITHUB_TOKEN)
	github-release upload -u kubic-project -r caasp-init --tag $(VERSION)  --name caasp-init-windows.tgz --file release/caasp-init-windows.tgz -s $(GITHUB_TOKEN)

MANPAGES_MD := $(wildcard doc/man/*.md)
MANPAGES    := $(MANPAGES_MD:%.md=%)

doc/man/%.1: doc/man/%.1.md
	@$(GO_MD2MAN) -in $< -out $@.out
	@go run doc/man/sanitize.go $@.out &> $@
	@rm $@.out

doc: $(MANPAGES)

.PHONY: caasp-init \
	caasp-init.static \
	install \
	install.static \
	clean \
	local-validate \
	local-validate-git \
	local-validate-go \
	local-validate-reproducible \
	local-validate-build \
	caasp-initimage \
	test.unit
	test \
	cover \
	bootstrap \
	dist \
	release \
	doc