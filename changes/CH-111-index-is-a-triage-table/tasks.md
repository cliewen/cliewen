---
id: CH-111-tasks
type: tasks
status: open
links: [CH-111-proposal]
title: Tasks for CH-111
---

# CH-111 — tasks

## Decision first

- [ ] Write ADR-046 extending ADR-041: an appended row seeds a description from the body, the description is curated thereafter, extraction seeds and never asserts, and rows without one are counted rather than failed on
- [ ] Amend C-016 with the description clause, keeping its existing row-opening rule intact

## Criteria before implementation

- [ ] Add AC-096 to CAP-005: an appended row carries a description extracted from the artifact body
- [ ] Add AC-097 to CAP-005: an artifact with no extractable sentence keeps the existing row shape and never emits an empty tail
- [ ] Add AC-098 to CAP-002: the judge counts rows that state their record but carry no description, and lists them on request
- [ ] Each with `Test-type: Unit` and both positive and negative evidence

## The generator

- [ ] Extract the description: a lede paragraph beneath the H1 where one exists, otherwise the first sentence of the first paragraph under the first heading
- [ ] Skip tables, lists, blockquotes and fenced blocks when looking for that sentence
- [ ] Truncate at a sentence boundary and keep the row on one line, since a multi-line entry is not one `checkIndexes` recognizes
- [ ] Append `- [<id> — <title>](<file>) · \`<status>\` — <description>` (AC-096)
- [ ] Degrade to today's row when no sentence can be read (AC-097)
- [ ] Confirm regeneration still rewrites nothing that already exists

## The judge

- [ ] Extend the index-row backlog in `internal/corpus/index.go` with the description-less population (AC-098)
- [ ] Surface it on the OK line and behind a list flag, as a count and never an `Issue`
- [ ] Confirm the judge reads no new source and gains no dependency

## Extraction guidance

- [ ] Teach `clue-extract` to absorb a foreign index file into its sibling README block, map its description column onto the row tail, drop columns restating `links:`, and delete the absorbed file
- [ ] Add the index-file row to `mappings/openspec.md` and `mappings/madr.md`
- [ ] Edit `internal/skills/source/` only, then run `go generate ./internal/skills` and confirm no drift

## Carriers

- [ ] C-016, CAP-002 and CAP-005 criteria and design, and the CLI surface state the same contract
- [ ] `[Unreleased]` in `CHANGELOG.md`; no `### Migration` section, because no corpus obligation is added or narrowed

## Verification

- [ ] Run the generator against the adopter repository that prompted this change and read the rows it produces
- [ ] `go build ./...`, `go test ./... -coverprofile=coverage.out`, coverage at or above 80%
- [ ] `go run ./cmd/clue validate --forbid-changes` and `git diff --check`
- [ ] Digest the workspace into the corpus and set M-051 `done` with its evidence
