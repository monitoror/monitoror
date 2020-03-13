#!/usr/bin/env bash
# Do not use this script manually, Use makefile

##############################################################
# This script is to generate mocks for monitorable directory #
##############################################################

for dir in monitorable/*/; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock-monitorable"
done
