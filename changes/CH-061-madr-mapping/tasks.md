---
id: CH-061-tasks
type: tasks
status: open
links: [CH-061]
title: Tasks for CH-061
---

# Tasks

- [ ] Read Tank Royale's converted `docs/decisions` records and its extraction report to fix the worked case this mapping cites, and re-read `mappings/openspec.md` for the shape and voice a mapping file uses.
- [ ] Write the `accepted-by` ADR at the next free ID: what an extraction may write into `accepted-by:` on a record that stays `inferred`, what stays body prose, and why pre-Cliewen acceptance is preserved without promoting to `verified`. Born `status: inferred`, `author: agent`.
- [ ] Add the `accepted-by:` clause to `internal/skills/source/shared/decision-records.md.tmpl` so the shared rule that binds the field to promotion admits the shape the ADR sanctions — this block renders into five skills, so read the result in each.
- [ ] Record the status-collapse rule as a `docs/decisions/log.md` row citing CH-061: what a `rejected` or `superseded` MADR record becomes when the target vocabulary has neither status.
- [ ] Write `internal/skills/source/resources/clue-extract/mappings/madr.md`: layout line covering MADR 3.x/4.x frontmatter, the template headings, and the Nygard form the same folders mix in.
- [ ] Add the status-vocabulary rows, stating the rule the log row settles: where `proposed`, `rejected`, `accepted`, `deprecated`, and `superseded by …` each go, given that every converted decision is born `inferred` and the vocabulary has only `inferred` and `verified`.
- [ ] Add the `accepted-by` rows, stating the rule the ADR settles for MADR's `decision-makers`, `consulted`, and `informed`, plus preservation of the original acceptance dates.
- [ ] Add the ID-preservation rows: numeric filename prefix survives as `ADR-xxx`, zero-padding and gaps, duplicate numbers across mixed folders, and records carrying no number (minted per contract item 3).
- [ ] Add the "watch for" line for the traps the worked case actually hit.
- [ ] Replace the `Nygard/MADR ADRs in docs/decisions` row in `mappings/openspec.md` with a pointer to `madr.md`, so one MADR conversion rule has one home.
- [ ] Make `internal/skills/generate.go` emit every file under `source/resources/clue-extract/mappings/` instead of reading `openspec.md` by name — without this the new mapping ships to nobody.
- [ ] Add a test in `internal/skills/generate_test.go` that fails when a mapping source has no generated counterpart; confirm it fails against the hardcoded-read generator before the fix.
- [ ] Update the **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl` to name both mappings instead of one current mapping.
- [ ] Update the extraction-mapping paragraph in `guide/adoption.md`, which says one mapping ships today.
- [ ] Regenerate the skills (`go generate ./internal/skills`) and confirm no drift (`TestSanity_CommittedSkillsMatchCanonicalSources`), including the new `mappings/madr.md` under both generated trees.
- [ ] Confirm `clue init` writes the new mapping into a fresh repository — the scaffold embeds `all:templates`, so verify rather than assume.
- [ ] Run `go test ./...` and confirm green.
- [ ] Run `clue validate` and confirm green.
- [ ] Run the strict guide build (`npm run guide:build`) and confirm green.
- [ ] Update `docs/plans/P-006-first-adoption.md`: mark M-023 `done` with an evidence cell citing CH-061, the ADR and log row, the mapping and generator files, and the checks.
- [ ] Add a `CHANGELOG.md` entry under `[Unreleased]` describing the new mapping from an adopter's perspective.
- [ ] Run `clue-verify`'s local checks and agentic review loop (in-context fallback, disclosed) before opening the PR.
