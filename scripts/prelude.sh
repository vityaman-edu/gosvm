#!/usr/bin/env bash
set -euo pipefail

# Populate the named array with matching, non-ignored repository files.
collect_files() {
  local array_name="$1"
  local pattern="$2"
  local -n collected_files="$array_name"
  local index
  local file_list
  local status=0

  file_list=$(mktemp)

  git ls-files --cached --others --exclude-standard -z -- "$pattern" > "$file_list" || status=$?

  if ((status != 0)); then
    rm -f "$file_list"

    return "$status"
  fi

  mapfile -d '' "$array_name" < "$file_list" || status=$?

  rm -f "$file_list"

  if ((status != 0)); then
    return "$status"
  fi

  for index in "${!collected_files[@]}"; do
    if [[ -f "${collected_files[index]}" ]]; then
      continue
    fi

    unset 'collected_files[index]'
  done
}
