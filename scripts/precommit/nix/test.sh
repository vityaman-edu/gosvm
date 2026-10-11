#!/usr/bin/env bash
set -euo pipefail

nix flake check path:. \
  --no-build \
  --no-update-lock-file \
  --all-systems
