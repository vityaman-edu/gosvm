#!/usr/bin/env bash
set -euo pipefail

build_directory=$(mktemp -d)
trap 'rm -rf "$build_directory"' EXIT

gotestsum --format standard-verbose -- \
  -race \
  -shuffle=on \
  -count=1 \
  -covermode=atomic \
  -coverpkg=./... \
  -coverprofile="$build_directory/coverage.out" \
  ./...

go tool cover -func="$build_directory/coverage.out" > "$build_directory/coverage.txt"

awk '
  $1 == "total:" {
    coverage = $3 + 0
    printf "Total statement coverage: %.1f%% (minimum: 90%%)\n", coverage
    found = 1

    if (coverage < 90) {
      print "ERROR: statement coverage is below 90%." > "/dev/stderr"
      exit 1
    }
  }

  END {
    if (!found) {
      print "ERROR: no total coverage result was produced." > "/dev/stderr"
      exit 1
    }
  }
' "$build_directory/coverage.txt"
