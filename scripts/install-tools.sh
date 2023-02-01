#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

# gotestsum, used by `make test`. Test utilities
echo "Installing gotestsum"
GO111MODULE=off go get gotest.tools/gotestsum

# rice, used by `make build`. Embed UI dist into go binary
echo "Installing rice"
GO111MODULE=off go get github.com/GeertJohan/go.rice/rice

# mockery, used by `make mocks`. Generating mock for backend
echo "Installing mockery"
GO111MODULE=off go get github.com/vektra/mockery/.../

# revproxy, usef by `make proxy`. Start a proxy for test
echo "Installing revproxy"
GO111MODULE=off go get github.com/jsdidierlaurent/revproxy
