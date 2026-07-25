---
id: CH-061-tasks
type: tasks
status: open
links: [CH-061]
title: Tasks for CH-061
---

# Tasks

- [x] Read Tank Royale's converted `docs/decisions` records and its extraction report to fix the worked case this mapping cites, and re-read `mappings/openspec.md` for the shape and voice a mapping file uses.
- [x] Write the `accepted-by` ADR at the next free ID: the field records approval given under Cliewen's merge boundary and nothing else, so acceptance that predates the corpus is preserved as body prose with its original dates. Name [PDR-004](../../docs/decisions/PDR-004-merge-binds-approval-signs.md) in `links:` and state that this extends its signing rule rather than replacing it. Born `status: inferred`, `author: agent`.
- [x] Add the `accepted-by:` clause to `internal/skills/source/shared/decision-records.md.tmpl` so the shared rule states the boundary the ADR settles — this block renders into five skills, so read the result in each.
- [x] List the new ADR in the `clue:index` block of `docs/decisions/README.md` — index completeness is machine-checked, so an unlisted record fails `clue validate`.
- [x] Read `docs/decisions/README.md` and `internal/scaffold/templates/docs/decisions/README.md` against the ADR's wording and amend each where it now says less than the rule does; the scaffold template ships into every repository `clue init` touches.
- [x] Record the status-collapse rule as a `docs/decisions/log.md` row citing CH-061: where `deprecated` and `superseded by …` land, and whether the source status survives as body prose, as a `links:` edge, or as both. The rule must be lossless — every converted record is kept and its source status survives as meaning, which is what makes the row reversible by a later mapping revision. For `rejected`, cite `docs/decisions/README.md` (rejected alternatives stay in the corpus, decisions are never deleted) rather than deciding it again.
- [x] Write `internal/skills/source/resources/clue-extract/mappings/madr.md`: layout line covering MADR 3.x/4.x frontmatter, the template headings, and the Nygard form the same folders mix in.
- [x] Add the status-vocabulary rows, stating the rule the log row settles: where `proposed`, `rejected`, `accepted`, `deprecated`, and `superseded by …` each go, given that every converted decision is born `inferred` and the vocabulary has only `inferred` and `verified`.
- [x] Add the `accepted-by` rows, stating the rule the ADR settles for MADR's `decision-makers`, `consulted`, and `informed`: body prose with the original acceptance dates, and the converted record carrying `accepted-by: []` — the exact shape an unsigned record uses, since the point of the rule is that `checkCoreFields` can one day be written against the field.
- [x] Add the ID-preservation rows: numeric filename prefix survives as `ADR-xxx`, zero-padding and gaps, duplicate numbers across mixed folders, and records carrying no number. State the minting rule for unnumbered records in the mapping itself — contract item 3 covers criterion IDs in `ac-prefix:` namespaces and does not reach decision records, and generalizing it is a contract edit past what M-022 settled.
- [x] Add the "watch for" line for the traps the worked case actually hit.
- [x] Replace the `Nygard/MADR ADRs in docs/decisions` row in `mappings/openspec.md` with a pointer to `madr.md`, so one MADR conversion rule has one home.
- [x] Make `internal/skills/generate.go` emit every file under `source/resources/clue-extract/mappings/` instead of reading `openspec.md` by name — without this the new mapping ships to nobody.
- [x] Add a test in `internal/skills/generate_test.go` that fails when a mapping source has no generated counterpart; confirm it fails against the hardcoded-read generator before the fix. Name it `TestSanity_…` — it asserts an invariant of this repo's own generator, and a test matching no purpose prefix fails `clue validate` ([ADR-006](../../docs/decisions/ADR-006-test-purpose-taxonomy.md)).
- [x] Update the **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl` to name both mappings instead of one current mapping.
- [x] Update the extraction-mapping paragraph in `guide/adoption.md`, which says one mapping ships today.
- [x] Regenerate the skills (`go generate ./internal/skills`) and confirm no drift (`TestSanity_CommittedSkillsMatchCanonicalSources`), including the new `mappings/madr.md` under both generated trees.
- [x] Confirm `clue init` writes the new mapping into a fresh repository — the scaffold embeds `all:templates`, so verify rather than assume.
- [x] Run `go test ./...` and confirm green.
- [x] Run `clue validate` and confirm green.
- [x] Run the strict guide build (`npm run guide:build`) and confirm green.
- [x] Update `docs/plans/P-006-first-adoption.md`: mark M-023 `done` with an evidence cell citing CH-061, the ADR and log row, the mapping and generator files, and the checks.
- [x] Add a `CHANGELOG.md` entry under `[Unreleased]` covering both user-visible parts from an adopter's perspective: the new mapping, and the `accepted-by` boundary that five shipped skills and the scaffolded decisions README now state.
- [ ] Run `clue-verify`'s local checks and agentic review loop (in-context fallback, disclosed) before opening the PR.
