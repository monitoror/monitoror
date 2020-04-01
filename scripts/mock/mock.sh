#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

#######################################################################
# This script is to generate mocks for monitorables and api directory #
#######################################################################

# monitorables/*
for dir in monitorables/*/api/; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock"
done

# api/config/*
for dir in api/*/; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock"
done

# cli
rm -rf cli/mocks
mockery \
  -dir cli \
  -output cli/mocks \
  -all \
  -note "If you want to rebuild this file, make mock"

# service
rm -rf service/mocks
mockery \
  -dir service \
  -output service/mocks \
  -all \
  -note "If you want to rebuild this file, make mock"

# pkg/*
dirs=("pkg/gogithub" "pkg/gojenkins" "pkg/gopingdom" "pkg/gotravis" "pkg/validator")
for dir in "${dirs[@]}"; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock"
done

# Generate mocks for external libs with interface
# Azure DevOps
rm -rf pkg/goazuredevops/build/mocks pkg/goazuredevops/release/mocks
mockery -name Client -output pkg/goazuredevops/build/mocks  -dir $(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/build
mockery -name Client -output pkg/goazuredevops/release/mocks  -dir $(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/release
