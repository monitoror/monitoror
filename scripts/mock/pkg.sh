#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

######################################################
# This script is to generate mocks for pkg directory #
######################################################

dirs=("pkg/gogithub" "pkg/gojenkins" "pkg/gopingdom" "pkg/gotravis" "pkg/validator")

# Generate mocks of interfaces find inside directory listed on dirs
for dir in "${dirs[@]}"; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock-pkg"
done

# Generate mocks for external libs with interface
# Azure DevOps
rm -rf pkg/goazuredevops/build/mocks pkg/goazuredevops/release/mocks
mockery -name Client -output pkg/goazuredevops/build/mocks  -dir $(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/build
mockery -name Client -output pkg/goazuredevops/release/mocks  -dir $(go list -m -f '{{ .Dir }}' github.com/jsdidierlaurent/azure-devops-go-api/azuredevops)/release
