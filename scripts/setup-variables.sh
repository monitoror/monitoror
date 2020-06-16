#!/usr/bin/env bash
# Do not use this script manually, Use makefile

############################################
## Imported script with 'source' command. ##
############################################

MB_SOURCE_PATH="./cmd/monitoror"
MB_VERSION_PATH="./VERSION"
MB_BINARIES_PATH="./binaries"

# ENVIRONMENT (can be "development" or "production")
MB_ENVIRONMENT=${MB_ENVIRONMENT:-"development"}

# VERSION (based on "VERSION" with prefix for dev)
MB_VERSION=${MB_VERSION:-$(cat $MB_VERSION_PATH 2>/dev/null)}
if [[ $MB_ENVIRONMENT == "development" ]]; then
  MB_VERSION=$(printf %s-dev "$MB_VERSION")
fi

# Git commit short hash
MB_GITCOMMIT=${MB_GITCOMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}

# Build time
MB_BUILDTIME=${MB_BUILDTIME:-$(TZ=GMT date +"%F %H:%M:%S+00:00" 2> /dev/null | sed -e 's/ /T/')}

# Default OS/ARCH
GOOS="${GOOS:-$(go env GOHOSTOS)}"
GOARCH="${GOARCH:-$(go env GOHOSTARCH)}"

# Default TAGS/LDFLAGS
MB_GO_TAGS=${MB_GO_TAGS:-""}
MB_GO_LDFLAGS="-w  \
-X \"github.com/monitoror/monitoror/cli/version.Version=${MB_VERSION}\" \
-X \"github.com/monitoror/monitoror/cli/version.GitCommit=${MB_GITCOMMIT}\" \
-X \"github.com/monitoror/monitoror/cli/version.BuildTime=${MB_BUILDTIME}\" \
-X \"github.com/monitoror/monitoror/cli/version.BuildTags=${MB_GO_TAGS}\" \
${MB_GO_LDFLAGS:-} \
"
