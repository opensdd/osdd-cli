#!/bin/bash
set -euo pipefail

go run internal/cmd/osdd/main.go recipe execute -f testRecipe.yaml -i claude -p "claude=claude" -l tmp/launch.json test
