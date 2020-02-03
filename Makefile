#
# github.com/monitoror/monitoror
#

DEFAULT: build

MAKEFLAGS = --silent

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

# ============= LINT =============
.PHONY: lint
lint: ## run linter
	@./scripts/lint

# ============= MOCKS =============
.PHONY: mock
mock: mock-monitorable mock-pkg

.PHONY: mock-monitorable
mock-monitorable: ## generate mocks of monitorable sub-directories
	@./scripts/mock/mock-monitorable

.PHONY: mock-pkg
mock-pkg: ## generate mocks of pkg directory listed in scripts/mock/mock-pkg
	@./scripts/mock/mock-pkg

# ============= BUILDS =============
.PHONY: build
build: ## build executable for current environment
	@./scripts/build/rice
	@./scripts/build/build

.PHONY: build-cross
build-cross: ## build all executables
	@./scripts/build/rice
	@./scripts/build/build linux
	@./scripts/build/build windows
	@./scripts/build/build macos
	@./scripts/build/build raspberrypi

.PHONY: build-linux
build-linux: ## build executable for Linux
	@./scripts/build/rice
	@./scripts/build/build linux

.PHONY: build-windows
build-windows: ## build executable for Windows
	@./scripts/build/rice
	@./scripts/build/build windows

.PHONY: build-macos
build-macos: ## build executable for MacOs
	@./scripts/build/rice
	@./scripts/build/build macos

.PHONY: build-raspberrypi
build-raspberrypi: ## build executable for Raspberry Pi
	@./scripts/build/rice
	@./scripts/build/build raspberrypi

.PHONY: build-faker
build-faker: ## build faker executable for current environment
	@./scripts/build/rice
	@./scripts/build/faker

.PHONY: build-faker-cross
build-faker-cross: ## build all faker executables
	@./scripts/build/rice
	@./scripts/build/faker linux
	@./scripts/build/faker windows
	@./scripts/build/faker macos
	@./scripts/build/faker raspberrypi

.PHONY: build-faker-linux
build-faker-linux: ## build faker executable for Linux
	@./scripts/build/rice
	@./scripts/build/faker linux

.PHONY: build-faker-windows
build-faker-windows: ## build faker executable for Windows
	@./scripts/build/rice
	@./scripts/build/faker windows

.PHONY: build-faker-macos
build-faker-macos: ## build faker executable for MacOs
	@./scripts/build/rice
	@./scripts/build/faker macos

.PHONY: build-faker-raspberrypi
build-faker-raspberrypi: ## build faker executable for Raspberry Pi
	@./scripts/build/rice
	@./scripts/build/faker raspberrypi

# ============= RUN =============
.PHONY: run
run: ## run monitoror
	@./scripts/run/run

.PHONY: run-faker
run-faker: ## run monitoror in faker mode
	@./scripts/run/faker

# ============= VERSION =============
.PHONY: version
version: ## bump version of monitoror
	@./scripts/version/bump

.PHONY: release
release: ## publish version of monitoror
	@./scripts/version/release

# ============= TOOLING =============
.PHONY: clean
clean: ## remove build artifacts
	rm -rf ./binaries/*
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
