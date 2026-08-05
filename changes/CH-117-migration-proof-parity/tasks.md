---
id: CH-117-tasks
type: tasks
status: open
links: [CH-117]
title: Tasks for migration proof parity
---

- [x] Write [ADR-049](../../docs/decisions/ADR-049-migration-parity-manifests.md): the source/target manifest schema and the five required failure classes
- [x] Add AC-109..AC-114 to `docs/capabilities/CAP-003-extract/criteria.md`: clean parity pass, missing criterion, orphaned tag, changed direction/location, stale source fingerprint, unjustified `@draft`/`Human`/retirement disposition
- [ ] Extend `internal/corpus` AC-evidence harvest to expose per-criterion evidence locations and direction(s) (needed by target-manifest derivation), with unit coverage that existing `checkACTests`/`Coverage` behavior is unchanged
- [ ] Implement `internal/parity/` (manifest types, `LoadSourceManifest`, `DeriveTargetManifest`, `Compare`, `Report`) with positive/negative unit tests naming AC-109..AC-114
- [ ] Implement `clue parity <source-manifest> [root]` (+ `--out`) in `cmd/clue/main.go`, wired through `internal/parity/`
- [ ] Update the canonical `clue-extract` skill source (`internal/skills/source/resources/clue-extract/`) and the OpenSpec mapping to require a source manifest during rehearsal and a clean `clue parity` run before mutation, then `go generate ./internal/skills`
- [ ] Update `docs/capabilities/CAP-003-extract/design.md` and `README.md` with the parity command's role
- [ ] Add `[Unreleased]` CHANGELOG entry: `clue parity` compares a source mapping's manifest against the derived corpus state
- [ ] Set M-053 to `done` in `docs/plans/P-011-truthful-brownfield-migration.md`, citing this change's evidence
- [ ] Digest: delete this change workspace
