---
id: CH-111-tasks
type: tasks
status: open
links: [CH-111-proposal]
title: Tasks for CH-111
---

# CH-111 — tasks

## Decision first

- [ ] Write ADR-046 replacing ADR-019's single-line-entry clause: the predefined prefix, the writer-owned and adopter-owned split, the extraction rule, and why no date column exists

## Criteria before implementation

- [ ] Revise AC-024 and AC-025 for the table contract
- [ ] Add AC-096…AC-099 to CAP-001 and CAP-002 with `Test-type: Unit`, each with positive and negative evidence

## The regeneration engine

- [ ] Detect a table by a header and separator as the first two non-empty lines inside the markers, leaving every other block byte-identical to today (AC-096)
- [ ] Preserve the header and separator across regeneration, which is the decapitation fix (AC-096)
- [ ] Key each row by the link target in its first cell, reusing `coverTarget` so subfolder coverage is unchanged
- [ ] Rewrite cells one through four from the artifact and preserve cells five onward verbatim (AC-097)
- [ ] Append a missing target as a row with a derived prefix and empty extras (AC-098)
- [ ] Extract the description: lede beneath the H1 when present, otherwise the first sentence under the first heading, truncated at a sentence boundary
- [ ] Report every artifact whose description came from the fallback

## The judge

- [ ] `checkIndexes` requires the predefined four headers in order when the block is a table (AC-099)
- [ ] Confirm the judge reads no git and gains no new dependency

## Templates and migration

- [ ] Ship the table header in the eight `internal/scaffold/templates/docs/**/README.md` index blocks
- [ ] Confirm a fresh `init` reports nothing regenerated on its first run, as ADR-019 requires
- [ ] MIG-008 converts an emitted list block to the table shape: preview by default, idempotent, prose outside the markers untouched
- [ ] Add MIG-008 to the migration inventory

## Extraction guidance

- [ ] Teach `clue-extract` to fold a foreign index file into its sibling README block, map columns onto the prefix, drop the ones restating `links:`, and delete the absorbed file
- [ ] Add the index-file row to `mappings/openspec.md` and `mappings/madr.md`
- [ ] Edit `internal/skills/source/` only, then run `go generate ./internal/skills` and confirm no drift

## Carriers

- [ ] CAP-001 criteria and design, CAP-002's index rule, and the migration inventory state the same contract
- [ ] `[Unreleased]` in `CHANGELOG.md`, including the `### Migration` section the release gate requires

## Verification

- [ ] Run the new engine against the adopter repository that prompted this change and confirm the header survives `init`
- [ ] `go build ./...`, `go test ./... -coverprofile=coverage.out`, coverage at or above 80%
- [ ] `go run ./cmd/clue validate --forbid-changes` and `git diff --check`
- [ ] Digest the workspace into the corpus and set M-051 `done` with its evidence
