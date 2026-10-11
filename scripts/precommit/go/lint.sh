#!/usr/bin/env bash
set -euo pipefail

go mod tidy -diff
go mod verify

golangci-lint config verify

golangci-lint fmt --diff

golangci-lint run ./...

go vet ./...
