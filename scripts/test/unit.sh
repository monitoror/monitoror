#!/usr/bin/env bash
# Do not use this script manually, Use makefile

go clean testdata ./...
gotestsum -- $(go list ./... | grep -v mocks)
