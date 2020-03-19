#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

source ./scripts/setup-variables.sh

#########################################################
# This script is used build docker images from binaries #
#########################################################

docker build -t "monitoror/monitoror:${MB_VERSION}" --build-arg VERSION="${MB_VERSION}" .
