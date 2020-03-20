#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

golangci-lint run --skip-files service/rice-box.go
