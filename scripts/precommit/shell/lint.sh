#!/usr/bin/env bash
set -euo pipefail

source scripts/prelude.sh

files=()
collect_files files '*.sh'

if ((${#files[@]} == 0)); then
  exit 0
fi

shfmt -d -i 2 -ci -sr "${files[@]}"

shellcheck --enable=all --severity=info "${files[@]}"
