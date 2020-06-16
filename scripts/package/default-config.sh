#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

#######################################################
# This script is used to package ui/dist in go source #
#######################################################

# Cleanup
rm -f cli/commands/init/rice-box.go
rm -rf cli/commands/init/default-files

# Copy default config
mkdir -p cli/commands/init/default-files
cp .env.example cli/commands/init/default-files/
cp config-example.json cli/commands/init/default-files/

# package config into go file
rice embed-go -i ./cli/commands/init

# cleanup
rm -rf cli/commands/init/default-files
