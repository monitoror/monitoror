#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

go clean testdata ./...
rm -f coverage.txt coverage.html
gotestsum -- -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v mocks)
go tool cover -html=coverage.txt -o coverage.html
