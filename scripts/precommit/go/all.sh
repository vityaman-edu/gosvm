#!/usr/bin/env bash
set -euo pipefail

bash scripts/precommit/go/lint.sh

bash scripts/precommit/go/build.sh

bash scripts/precommit/go/test.sh

bash scripts/precommit/go/scan.sh
