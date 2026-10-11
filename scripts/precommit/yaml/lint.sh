#!/usr/bin/env bash
set -euo pipefail

yamllint --strict .

actionlint -color
