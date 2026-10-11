#!/usr/bin/env bash
set -euo pipefail

bash scripts/precommit/nix/lint.sh

bash scripts/precommit/nix/test.sh
