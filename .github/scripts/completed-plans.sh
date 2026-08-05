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

base_ref="${1:?usage: completed-plans.sh <base-ref> <head-ref>}"
head="${2:?usage: completed-plans.sh <base-ref> <head-ref>}"

# The merge base, not the base branch tip. `pull_request.base.sha` follows the
# base branch, so a two-dot diff also lists plans that only *main* changed —
# and the first digest to close a campaign would turn every open pull request
# red for a file none of them touched.
base=$(git merge-base "$base_ref" "$head")

# Captured rather than piped: a failure here must stop the script, and `set -e`
# does not reach into a process substitution feeding a loop. A guard that goes
# quiet when it cannot resolve its own base is worse than no guard.
changed=$(git diff --name-only "$base" "$head" -- 'docs/plans/*.md')

failed=0
while IFS= read -r file; do
  [ -n "$file" ] || continue
  # A plan added by this change has no base-side file, and `git show` failing
  # is the ordinary way to learn that rather than an error to report.
  if ! before=$(git show "$base:$file" 2>/dev/null); then
    continue
  fi
  # Frontmatter only, and only when the file opens with a fence: line 2 through
  # the closing `---`. A `status:` in the body is prose about a plan, not the
  # plan's own status, and a file with no frontmatter has no status to read.
  case "$before" in
    ---*) ;;
    *) continue ;;
  esac
  # The value is read the way YAML writes it, not the way this repository
  # happens to: `status: "completed"` and `status: 'completed'` are the same
  # plan to `clue validate`, and a guard that matched only the bare word would
  # let a quoted one through in silence — the one failure mode a freeze may not
  # have. A trailing comment is part of YAML's line and not part of the value.
  if printf '%s\n' "$before" | tr -d '\r' | sed -n '2,/^---$/p' |
    grep -Eq "^status:[[:space:]]*[\"']?completed[\"']?[[:space:]]*(#.*)?$"; then
    echo "FAIL: $file was completed on $base — a completed plan is frozen and never deleted (C-008)"
    failed=1
  fi
done <<EOF
$changed
EOF

if [ "$failed" -ne 0 ]; then
  echo "A finished campaign is the project's record of what it achieved. Correct it in a successor plan, not in place."
  exit 1
fi

echo "No completed plan was modified."
