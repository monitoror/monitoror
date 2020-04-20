#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

UI_PORT=${UI_PORT:-"8000"}
PROXY_PORT=${PROXY_PORT:-"8100"}

revproxy --port="${PROXY_PORT}" --prefix="monitoror" "http://localhost:${UI_PORT}"
