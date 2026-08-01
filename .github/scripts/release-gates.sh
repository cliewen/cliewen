#!/usr/bin/env bash
#
# The release gates, in one place.
#
# These checks used to live only in tag-on-merge.yml, which runs after the
# human merge. That is too late: a release PR that could never be tagged still
# looked mergeable, and the failure surfaced as a red run on main with the
# version stamp already committed. Tagging fires only on the merge that raises
# the stamp, so that merge cannot be retried — the version has to be abandoned
# and a new one cut. 0.11.1 was lost exactly that way.
#
# So the same script now runs on the pull request, where a missing changelog
# section costs a red check the author can fix, and on the merge, where it
# remains the authority that guards the tag. One implementation, because two
# copies of a gate drift and the drift is invisible until a release is already
# half-made.
#
# Usage: release-gates.sh <base-ref> <head-ref>
#
# Exits 0 and does nothing when the stamp did not change between the two refs:
# an ordinary change is not a release and owes none of this.
set -euo pipefail

base="${1:?usage: release-gates.sh <base-ref> <head-ref>}"
head="${2:?usage: release-gates.sh <base-ref> <head-ref>}"

tmpl=internal/skills/source/shared/frontmatter.md.tmpl

# The skills' frontmatter template is the repository's single bump site, and
# already the input the release's drift gate judges the tag against. Reading
# the version from anywhere else would create a second source that could
# disagree with it.
if git diff --quiet "$base" "$head" -- "$tmpl"; then
  echo "The version stamp did not change — not a release, so the release gates do not apply."
  exit 0
fi

# Everything the gates judge is read out of <head-ref> rather than the working
# tree. Both callers happen to have <head-ref> checked out, so this changes
# nothing in CI — but a script that takes a ref and then reads whatever is on
# disk answers about the wrong release the moment it is run anywhere else, and
# reports a pass while doing it. Replaying a past merge to check this gate would
# have silently graded the current checkout. A gate that quietly grades the
# wrong thing is the failure this whole change exists to remove.
v=$(git show "${head}:${tmpl}" | sed -n 's/^version:[[:space:]]*//p' | head -n 1)
if ! printf '%s' "$v" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "::error title=Unreadable version stamp::${tmpl} does not carry a bare semver version; got '${v}'" >&2
  exit 1
fi
echo "This change raises the version stamp to ${v}; checking it could be released."

notes=$(mktemp)
trap 'rm -f "$notes"' EXIT
git show "${head}:CHANGELOG.md" | awk -v ver="$v" '
  BEGIN { header = "## [" ver "]" }
  !found { if (index($0, header) == 1) found = 1; next }
  /^## \[/ { exit }
  { print }
' > "$notes"

# A tag with no user-facing notes cannot be released at all: the release
# workflow fails at its extraction step (ADR-012). Failing here instead means a
# missing section costs a red run, not a spent version number.
if ! grep -q '[^[:space:]]' "$notes"; then
  echo "::error title=No release notes::CHANGELOG.md has no '## [${v}]' section with content; the stamp was raised without user-facing notes (ADR-012)" >&2
  exit 1
fi

# A release that changes the migration machinery has something an adopter must
# act on, and they will not learn it from a Fixed entry describing internals.
if ! git diff --quiet "$base" "$head" -- internal/migrate; then
  if ! awk '
    /^### Migration$/ { found = 1; next }
    found && /^### / { exit }
    found && /[^[:space:]]/ { content = 1 }
    END { exit !content }
  ' "$notes"; then
    echo "::error title=No migration guidance::The release changes internal/migrate but its CHANGELOG section has no non-empty ### Migration section (ADR-039)" >&2
    exit 1
  fi
fi

# A release must record what its generated carriers look like, or every adopter
# already on it is told their pristine skills are locally edited and `clue
# migrate` refuses the whole write set. 0.10.0 shipped without this and the gap
# was invisible until an adopter tried to upgrade a release later.
#
# This half of the check only makes sense while a release is being cut: between
# releases the working tree is meant to differ from the last published one.
#
# It asks whether each digest appears in the manifest at all, not whether it sits
# in this version's entry. That is deliberate: a partial written by hand is
# caught by the Go guard, which verifies every entry against the tag it names,
# and a carrier whose bytes did not change between releases is still recognized
# through the older entry it already matches. What this catches is the failure
# that actually shipped — a release recorded nowhere.
missing=0
for f in $(cd .agents/skills && find . -name '*.md' | sed 's|^\./||' | sort); do
  want=$(sha256sum ".agents/skills/${f}" | cut -d' ' -f1)
  if ! grep -Fq "\"${f}\":" internal/migrate/migrate.go; then
    echo "::error title=Carrier missing from the release manifest::${f} is not listed in legacyDigests (internal/migrate/migrate.go); add it to the ${v} entry with digest ${want}" >&2
    missing=1
  elif ! grep -Fq "\"${want}\"" internal/migrate/migrate.go; then
    echo "::error title=Release manifest is stale::${f} has digest ${want} but no entry in legacyDigests carries it; the ${v} entry must record this release's carriers" >&2
    missing=1
  fi
done
if [ "$missing" -ne 0 ]; then
  echo "::error title=No manifest entry for ${v}::Cutting ${v} without its carrier digests leaves every adopter on ${v} unable to migrate off it (ADR-028)" >&2
  exit 1
fi

echo "Release gates pass for ${v}."
