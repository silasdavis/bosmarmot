#
# |_  _  _ _ _  _  _ _ _  _ _|_
# |_)(_)_\| | |(_|| | | |(_) |
#
# Bosmarmot Makefile
#
# Requires go version 1.8 or later.
#

SHELL := /bin/bash
REPO := $(shell pwd)
GO_FILES := $(shell go list -f "{{.Dir}}" ./...)
GOPACKAGES_NOVENDOR := $(shell go list ./...)
COMMIT := $(shell git rev-parse --short HEAD)
BURROW_PACKAGE := github.com/hyperledger/burrow

### Integration test binaries
# We make the relevant targets for building/fetching these depend on the Makefile itself - if unnecessary rebuilds
# when changing the Makefile become a problem we can move these values into individual files elsewhere and make those
# files specific targets for their respective binaries

### Tests and checks
# Run goimports (also checks formatting) first display output first, then check for success
.PHONY: check
check:
	@go get golang.org/x/tools/cmd/goimports
	@goimports -l -d ${GO_FILES}
	@goimports -l ${GO_FILES} | read && echo && \
	echo "Your marmot has found a problem with the formatting style of the code."\
	 1>&2 && exit 1 || true

# Just fix it
.PHONY: fix
fix:
	@goimports -l -w ${GO_FILES}

.PHONY: test_js
test_js:
	@cd legacy-contracts.js && npm test
	#re-enable after fixing a few things
	#@cd legacy-db.js && npm test

# Run tests
.PHONY:	test_bos
test_bos: check bin/solc
	@tests/scripts/bin_wrapper.sh go test ${GOPACKAGES_NOVENDOR}

.PHONY:	test
test: test_bos

# Run tests for development (noisy)
.PHONY:	test_dev
test_dev:
	@go test -v ${GOPACKAGES_NOVENDOR}

# Install dependency and make legacy-contracts depend on legacy-db by relative path
.PHONY: npm_install
npm_install:
	@cd legacy-db.js && npm install
	@cd legacy-contracts.js && npm install --save ../legacy-db.js
	@cd legacy-contracts.js && npm install

# Run tests including integration tests
.PHONY:	test_integration_bos
test_integration_bos: build_bin bin/solc bin/burrow
	@tests/scripts/bin_wrapper.sh tests/run_pkgs_tests.sh

.PHONY:	test_integration_js
test_integration_js: build_bin bin/solc bin/burrow
	@cd legacy-contracts.js && TEST=record ../tests/scripts/bin_wrapper.sh npm test

.PHONY:	test_integration
test_integration: test_integration_bos test_integration_js

# Use a provided/local Burrow
.PHONY:	test_integration_js_no_burrow
test_integration_js_no_burrow: build_bin bin/solc
	@cd legacy-contracts.js && TEST=record ../tests/scripts/bin_wrapper.sh npm test

.PHONY:	test_integration_bos_no_burrow
test_integration_bos_no_burrow: build_bin bin/solc
	@tests/scripts/bin_wrapper.sh tests/run_pkgs_tests.sh

PHONY:	test_integration_no_burrow
test_integration_no_burrow: test_integration_bos_no_burrow test_integration_js_no_burrow

### Vendoring

# erase vendor wipes the full vendor directory
.PHONY: erase_vendor
erase_vendor:
	rm -rf ${REPO}/vendor/

# install vendor uses dep to install vendored dependencies
.PHONY: reinstall_vendor
reinstall_vendor: erase_vendor
	@dep ensure -v

# delete the vendor directy and pull back using dep lock and constraints file
# will exit with an error if the working directory is not clean (any missing files or new
# untracked ones)
.PHONY: ensure_vendor
ensure_vendor: reinstall_vendor
	@tests/scripts/is_checkout_dirty.sh

### Builds

.PHONY: build_bin
build_bin:
	@go build -ldflags "-X github.com/monax/bosmarmot/project.commit=${COMMIT}" -o bin/bos ./bos/cmd/bos

bin/solc: ./tests/scripts/deps/solc.sh
	@mkdir -p bin
	@tests/scripts/deps/solc.sh bin/solc
	@touch bin/solc

tests/scripts/deps/burrow.sh: Gopkg.lock
	@go get -u github.com/golang/dep/cmd/dep
	@tests/scripts/deps/burrow-gen.sh > tests/scripts/deps/burrow.sh
	@chmod +x tests/scripts/deps/burrow.sh

.PHONY: burrow_local
burrow_local:
	@rm -rf .gopath_burrow
	@mkdir -p .gopath_burrow/src/${BURROW_PACKAGE}
	@cp -r ${GOPATH}/src/${BURROW_PACKAGE}/. .gopath_burrow/src/${BURROW_PACKAGE}

bin/burrow: ./tests/scripts/deps/burrow.sh
	mkdir -p bin
	@GOPATH="${REPO}/.gopath_burrow" \
	tests/scripts/go_get_revision.sh \
	https://github.com/hyperledger/burrow.git \
	${BURROW_PACKAGE} \
	$(shell ./tests/scripts/deps/burrow.sh) \
	"make build_db" && \
	cp .gopath_burrow/src/${BURROW_PACKAGE}/bin/burrow ./bin/burrow


# Build all the things
.PHONY: build
build:	build_bin

# Build binaries for all architectures
.PHONY: build_dist
build_dist:
	@goreleaser --rm-dist --skip-publish --skip-validate

# Do all available tests and checks then build
.PHONY: build_ci
build_ci: check test build

### Release and versioning

# Print version
.PHONY: version
version:
	@go run ./project/cmd/version/main.go

# Generate full changelog of all release notes
CHANGELOG.md: ./project/history.go ./project/cmd/changelog/main.go
	@go run ./project/cmd/changelog/main.go > CHANGELOG.md

# Generated release notes for this version
NOTES.md:  ./project/history.go ./project/cmd/notes/main.go
	@go run ./project/cmd/notes/main.go > NOTES.md

# Tag the current HEAD commit with the current release defined in
# ./release/release.go
.PHONY: tag_release
tag_release: test check CHANGELOG.md build_bin
	@tests/scripts/tag_release.sh

# If the checked out commit is tagged with a version then release to github
.PHONY: release
release: NOTES.md
	@tests/scripts/release.sh
