#!/usr/bin/env bash
#
# Fails when total statement coverage falls below a threshold.
#
# Usage: check_coverage.sh <coverage-profile> <minimum-percent>
#
# The threshold is a ratchet, not a target: raise it as coverage rises so a
# change can never quietly give ground. It is deliberately the last of the
# acceptance criteria in TESTING_PLAN.md — the point is covering the use cases,
# and the percentage is what follows from that.

set -euo pipefail

profile="${1:-coverage.out}"
minimum="${2:-90}"

if [ ! -f "$profile" ]; then
  echo "coverage profile not found: $profile" >&2
  exit 1
fi

total=$(go tool cover -func="$profile" | awk '/^total:/ {print $3}' | tr -d '%')

if [ -z "$total" ]; then
  echo "could not read a total from $profile" >&2
  exit 1
fi

echo "total coverage: ${total}% (minimum ${minimum}%)"

if awk -v total="$total" -v minimum="$minimum" 'BEGIN { exit !(total < minimum) }'; then
  echo "coverage ${total}% is below the ${minimum}% minimum" >&2
  echo "" >&2
  echo "packages furthest from full coverage:" >&2
  go tool cover -func="$profile" \
    | grep -v '^total:' \
    | awk '{ gsub(/%/, "", $NF); print $NF, $1 }' \
    | sort -n \
    | head -15 >&2
  exit 1
fi

echo "coverage check passed"
