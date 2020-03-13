#!/usr/bin/env bash
# Do not use this script manually, Use makefile

source ./scripts/setup-variables.sh

###############################################
# This script is used to start monitoror core #
###############################################

# Force dev environment
MO_ENV=${MO_ENV:-$MB_ENVIRONMENT}

go run --ldflags "$MB_GO_LDFLAGS" --tags "$MB_GO_TAGS" "${MB_SOURCE_PATH}"
