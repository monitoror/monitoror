#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

go clean testdata ./...
gotestsum -- $(go list ./... | grep -v mocks)
