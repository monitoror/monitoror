#!/usr/bin/env bash
# Do not use this script manually, Use makefile

source ./scripts/setup-variables.sh

#############################################################
# This script is used for creating new version of monitoror #
#############################################################

MB_VERSION_REGEX="^([0-9]*)\\.([0-9]*)\\.([0-9]*)$"

function error {
  echo -e "$1" >&2
  exit 1
}

function bump_version () {
  if [[ "$1" =~ $MB_VERSION_REGEX ]]; then
    major=${BASH_REMATCH[1]}
    minor=${BASH_REMATCH[2]}
    patch=${BASH_REMATCH[3]}

    case "$2" in
      major) newVersion="$((major + 1)).0.0";;
      minor) newVersion="${major}.$((minor + 1)).0";;
      patch) newVersion="${major}.${minor}.$((patch + 1))";;
      *) error "Invalid command $2" ;;
    esac
  else
    error "Version $1 does not match the version scheme 'major.minor.patch'."
  fi

  echo "$newVersion"
}

#### START SCRIPT ####

if [[ -n $(git status -s) ]]; then
  error "You have changed files. Please clean up your git repository"
fi

version=$(cat "$MB_VERSION_PATH" 2>/dev/null)
echo "Current Version: $version"
read -p 'Command (major|minor|patch): ' input

case $input in
  major|minor|patch) command=$input;;
  *) error "Invalid command $input";;
esac

newVersion=$(bump_version "$version" "$command")

# Update VERSION file
echo "New version is $newVersion"
echo $newVersion > "$MB_VERSION_PATH"

# Commit !
git add "$MB_VERSION_PATH"
git commit -m "chore: bump version to $newVersion"





