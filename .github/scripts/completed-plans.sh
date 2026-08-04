#!/usr/bin/env bash
#
# C-008's machine: a completed plan is frozen.
#
# The plans index doubles as the project's achievement overview, so rewriting a
# finished campaign rewrites history. Nothing in `clue validate` can hold this
# — the rule compares a file against what it used to be, and the judge reads a
# state rather than a transition (ADR-044). A workflow, on the other hand,
# knows exactly what it is merging into, and comparing against the base is
# ordinary work for one.
#
# The status is read from <base-ref>, never from the working tree. That is what
# makes the change that *closes* a plan legal: it sets `completed` in its
# digest, and on the base that plan is still `active`.
#
# Usage: completed-plans.sh <base-ref> <head-ref>
set -euo pipefail

base="${1:?usage: completed-plans.sh <base-ref> <head-ref>}"
head="${2:?usage: completed-plans.sh <base-ref> <head-ref>}"

failed=0
while IFS= read -r file; do
  [ -n "$file" ] || continue
  # A plan added by this change has no base-side file, and `git show` failing
  # is the ordinary way to learn that rather than an error to report.
  if ! before=$(git show "$base:$file" 2>/dev/null); then
    continue
  fi
  # Frontmatter only: line 2 through the closing fence. A `status:` in the body
  # is prose about a plan, not the plan's own status.
  if printf '%s\n' "$before" | sed -n '2,/^---$/p' | grep -q '^status: completed'; then
    echo "FAIL: $file was completed on $base — a completed plan is frozen and never deleted (C-008)"
    failed=1
  fi
done < <(git diff --name-only "$base" "$head" -- 'docs/plans/*.md')

if [ "$failed" -ne 0 ]; then
  echo "A finished campaign is the project's record of what it achieved. Correct it in a successor plan, not in place."
  exit 1
fi

echo "No completed plan was modified."
