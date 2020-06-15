#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

source ./scripts/setup-variables.sh

##########################################
# This script is used to build monitoror #
##########################################

# Override OS/ARCH
if [[ $# -eq 1 ]]; then
  if [[ $1 == "linux/amd64" ]]; then
    export GOOS=linux
    export GOARCH=amd64
  elif [[ $1 == "linux/ARMv5" ]]; then
    export GOOS=linux
    export GOARCH=arm
    export GOARM=5
  elif [[ $1 == "windows" ]]; then
    export GOOS=windows
    export GOARCH=amd64
  elif [[ $1 == "macos" ]]; then
    export GOOS=darwin
    export GOARCH=amd64
  fi
fi

# Define target base name
targetBaseName="$MB_BINARIES_PATH/monitoror"

# Define target os/arch decorator
targetOsArch="-$GOOS-$GOARCH"
if [[ $GOOS == "darwin" ]]; then
  targetOsArch="-macos"
fi

# Define target version decorator
targetVersion="-$MB_VERSION"

# Define target tags decorator
targetTags=""
if [[ $MB_GO_TAGS != "" ]]; then
  targetTags="-${MB_GO_TAGS/,/-}"
fi

# Define target extension
ext=""
if [[ $GOOS == "windows" ]]; then
  ext=".exe"
fi

# Target
target=$(printf %s%s%s%s%s "$targetBaseName" "$targetOsArch" "$targetVersion" "$targetTags" "$ext")

echo "Building statically linked $target"
go build -o "$target" --ldflags "$MB_GO_LDFLAGS" --tags "$MB_GO_TAGS" "${MB_SOURCE_PATH}"
