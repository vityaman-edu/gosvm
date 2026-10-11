#!/usr/bin/env bash
set -euo pipefail

source scripts/prelude.sh

files=()
collect_files files '*.nix'

if ((${#files[@]} == 0)); then
  exit 0
fi

nixfmt --strict "${files[@]}"
