#!/usr/bin/env bash
# Do not use this script manually, Use makefile

golangci-lint run --skip-files service/rice-box.go
