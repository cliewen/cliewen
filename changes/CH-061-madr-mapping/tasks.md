---
id: CH-061-tasks
type: tasks
status: open
links: [CH-061]
title: Tasks for CH-061
---

# Tasks

- [ ] Read Tank Royale's converted `docs/decisions` records and its extraction report to fix the worked case this mapping cites, and re-read `mappings/openspec.md` for the shape and voice a mapping file uses.
- [ ] Write `internal/skills/source/resources/clue-extract/mappings/madr.md`: layout line covering MADR 3.x/4.x frontmatter, the template headings, and the Nygard form the same folders mix in.
- [ ] Add the status-vocabulary rows: where `proposed`, `rejected`, `accepted`, `deprecated`, and `superseded by …` each go, given that every converted decision is born `inferred` and the vocabulary has only `inferred` and `verified`.
- [ ] Add the `accepted-by` rule for acceptance that predates Cliewen: what may be written into `accepted-by:`, what stays body prose, and preservation of the original acceptance dates.
- [ ] Add the ID-preservation rows: numeric filename prefix survives as `ADR-xxx`, zero-padding and gaps, duplicate numbers across mixed folders, and records carrying no number (minted per contract item 3).
- [ ] Add the "watch for" line for the traps the worked case actually hit.
- [ ] Update the **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl` to name both mappings instead of one current mapping.
- [ ] Update the extraction-mapping paragraph in `guide/adoption.md`, which says one mapping ships today.
- [ ] Regenerate the skills (`go generate ./internal/skills`) and confirm no drift (`TestSanity_CommittedSkillsMatchCanonicalSources`).
- [ ] Run `go test ./...` and confirm green.
- [ ] Run `clue validate` and confirm green.
- [ ] Run the strict guide build (`npm run guide:build`) and confirm green.
- [ ] Update `docs/plans/P-006-first-adoption.md`: mark M-023 `done` with an evidence cell citing CH-061 and the concrete files and checks.
- [ ] Add a `CHANGELOG.md` entry under `[Unreleased]` describing the new mapping from an adopter's perspective.
- [ ] Record the status-collapse and `accepted-by` rules as a `docs/decisions/log.md` row citing CH-061.
- [ ] Run `clue-verify`'s local checks and agentic review loop (in-context fallback, disclosed) before opening the PR.
