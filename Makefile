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
	@./scripts/test/unit.sh

.PHONY: test-coverage
test-coverage: ## run test coverage
	@./scripts/test/coverage.sh

.PHONY: test-coverage-html
test-coverage-html: ## run test coverage and generate cover.html
	@./scripts/test/coverage-html.sh

# ============= LINT =============
.PHONY: lint
lint: ## run linter
	@./scripts/test/lint.sh

# ============= MOCKS =============
.PHONY: mocks
mocks: ## generate mocks
	@./scripts/mocks.sh

# ============= BUILDS =============
.PHONY: build
build: package-front ## build executable for current environment
	@./scripts/build.sh

.PHONY: build-cross
build-cross: package-front ## build all executables
	@./scripts/build.sh linux/amd64
	@./scripts/build.sh linux/ARMv5
	@./scripts/build.sh windows
	@./scripts/build.sh macos

.PHONY: build-linux-amd64
build-linux-amd64: package-front ## build executable for Linux
	@./scripts/build.sh linux/amd64

.PHONY: build-linux-ARMv5
build-linux-ARMv5: package-front ## build executable for Raspberry Pi (ARM V5)
	@./scripts/build.sh linux/ARMv5

.PHONY: build-windows
build-windows: package-front ## build executable for Windows
	@./scripts/build.sh windows

.PHONY: build-macos
build-macos: package-front ## build executable for MacOs
	@./scripts/build.sh macos

.PHONY: build-faker-linux-amd64
build-faker-linux-amd64: package-front ## build faker executable linux amd64 (only for demo)
	@MB_GO_TAGS="faker" ./scripts/build.sh linux/amd64

# ============= PACKAGE =============
.PHONY: package-front
package-front: ## package front directory ui/dist into go source
	@./scripts/package/front.sh

.PHONY: package-docker
package-docker: ## package linux amd64 into docker image
	@./scripts/package/docker.sh

# ============= RUN =============
.PHONY: run
run: ## run monitoror
	@./scripts/run.sh

.PHONY: run-faker
run-faker: ## run monitoror in faker mode
	@MB_GO_TAGS="faker"  ./scripts/run.sh

# ============= VERSION =============
.PHONY: version
version: ## bump version of monitoror
	@./scripts/bump-version.sh

# ============= TOOLING =============
.PHONY: clean
clean: ## remove build artifacts
	rm -rf ./binaries/*

.PHONY: install
install: ## installing tools / dependencies
	@./scripts/install.sh

.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
