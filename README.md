# Cliewen

> A verifiable thread from goal to test — docs as the system of record
> for agent-driven development.

**Cliewen** (Old English *cliewen*, "ball of thread" — the word that
became *clue*) is a methodology and toolchain for agent-driven software
development: the documentation corpus is the system-of-record, and the
chain **goal → capability → acceptance criterion → test** is
mechanically enforced. CLI binary: `clue`.

SDD frameworks document the *change*; Cliewen documents the *system*.
Changes are transient deltas digested into the permanent corpus at
merge — `git log docs/` is the provenance archive.

## Status

Elaboration baseline in progress
([P-001](docs/plans/P-001-elaboration-baseline.md)). This repo dogfoods
its own conventions from commit one — start reading at
[docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
