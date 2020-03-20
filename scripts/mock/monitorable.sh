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
      -note "If you want to rebuild this file, make mock-monitorable"
done

# api/config/*
for dir in api/*/; do
    rm -rf ${dir}/mocks
    mockery \
      -dir ${dir} \
      -output ${dir}/mocks \
      -all \
      -note "If you want to rebuild this file, make mock-monitorable"
done
