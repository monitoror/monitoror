#
# github.com/jsdidierlaurent/monitoror
#

DEFAULT: build-cross

MAKEFLAGS = --silent --ignore-errors

# ============= TESTS =============
.PHONY: test
test: test-unit ## run tests

.PHONY: test-unit
test-unit: ## run unit tests, to change the output format use: GOTESTSUM_FORMAT=(dots|short|standard-quiet|short-verbose|standard-verbose) make test-unit
	@./scripts/test/test-unit

.PHONY: test-coverage
test-coverage: ## run test coverage
	@./scripts/test/test-coverage

.PHONY: test-coverage-html
test-coverage-html: ## run test coverage and generate cover.html
	@./scripts/test/test-coverage-html

# ============= MOCKS =============
.PHONY: mock
mock: mock-monitorable

.PHONY: mock-full
mock: mock-monitorable mock-pkg ## generate all mocks

.PHONY: mock-monitorable
mock-monitorable: ## generate mocks of monitorable sub-directories
	@./scripts/mock/mock-monitorable

.PHONY: mock-pkg
mock-pkg: ## generate mocks of pkg sub-directories
	@./scripts/mock/mock-pkg

# ============= BUILDS =============
.PHONY: build
build: build-cross

.PHONY: build-all
build-all: build-binary build-windows build-osx ## build all executables

.PHONY: build-binary
build-binary: ## build executable for Linux
	@./scripts/build/binary

.PHONY: build-windows
build-windows: ## build executable for Windows
	@./scripts/build/windows

.PHONY: build-osx
build-osx: ## build executable for Mac Osx
	@./scripts/build/osx

.PHONY: build-cross
build-cross: ## build executable for current environment
	@./scripts/build/cross

# ============= RUN =============
.PHONY: run
run: ## run monitoror
	@./scripts/run/run

.PHONY: run-faker
run-faker: ## run monitoror in faker mode
	@./scripts/run/faker

# ============= TOOLINGS =============
.PHONY: clean
clean: ## remove build artifacts
	rm -rf ./build/*
	@go clean ./...

.PHONY: install
install: ## installing tools / dependencies
	@./scripts/install

.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: fmt
fmt:
	go fmt ./...
