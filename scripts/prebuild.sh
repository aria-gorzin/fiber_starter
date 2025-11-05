#!/usr/bin/env bash
set -euo pipefail

STATE_DIR="scripts/.state"
STATE_FILE="$STATE_DIR/prebuild.lastrun"

mkdir -p "$STATE_DIR"

run_sqlc=false
run_proto=false
run_swag=false

if [ ! -f "$STATE_FILE" ]; then
  echo "First run → running all pre-steps"
  run_sqlc=true
  run_proto=true
  run_swag=true
  : > "$STATE_FILE"
else
  echo "Changed files:"
  changed=$(find . -type f \
    ! -path "./.git/*" \
    ! -path "./tmp/*" \
    ! -path "./vendor/*" \
    ! -path "./db/sqlc/*" \
    ! -path "./pb/*" \
    ! -path "./doc/*" \
    ! -path "./docs/*" \
    -newer "$STATE_FILE" \
    -print 2>/dev/null)

  if [ -z "$changed" ]; then
    echo "  (none since last run)"
  else
    echo "$changed" | sed 's/^/  /'
  fi

  # Check changed files and decide what to run
  if printf '%s\n' "$changed" | grep -qE '^\./db/query/'; then
     run_sqlc=true
   fi
   if printf '%s\n' "$changed" | grep -qE '^\./proto/'; then
     run_proto=true
   fi
   if printf '%s\n' "$changed" | grep -qE '^\./api/'; then
     run_swag=true
   fi
fi

# Run only what’s needed
if [ "$run_sqlc" = true ]; then
  echo "▶ make sqlc"
  make sqlc
fi

if [ "$run_proto" = true ]; then
  echo "▶ make proto"
  make proto
fi

if [ "$run_swag" = true ]; then
  echo "▶ swag init"
  swag init
fi

# Update timestamp
touch "$STATE_FILE"
