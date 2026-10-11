#!/usr/bin/env bash
set -euo pipefail

build_directory=$(mktemp -d)
trap 'rm -rf "$build_directory"' EXIT

go build -o "$build_directory/" ./...
