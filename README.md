# Cliewen

> A verifiable thread from goal to test — docs as the system of record for agent-driven development.

**Cliewen** (Old English *cliewen*, "ball of thread" — the word that became *clue*) is a methodology and toolchain for agent-driven software development: the documentation corpus is the system-of-record, and the chain **goal → capability → acceptance criterion → test** is mechanically enforced. CLI binary: `clue`.

SDD frameworks document the *change*; Cliewen documents the *system*. Changes are transient deltas digested into the permanent corpus at merge — `git log docs/` is the provenance archive.

## Install

`clue` is a single Go binary. Install from source, or download a stamped prebuilt binary from a release:

```sh
# From source — a release install stamps itself from Go's build info:
go install github.com/cliewen/cliewen/cmd/clue@latest

# A stamped prebuilt binary for your platform (linux/darwin/windows × amd64/arm64):
gh release download vX.Y.Z --repo cliewen/cliewen --pattern 'clue-*-<os>-<arch>*'
```

While the repo is private, `go install` additionally needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.

`clue version` reports the release it was built from — a checkout build (`go build ./cmd/clue`) or an install of an untagged commit reports `dev`. A tagged release (`vX.Y.Z`) builds the cross-platform binaries and stamps each with its version; the agent skills carry the same version, and `clue validate` flags drift between them ([CAP-004](docs/capabilities/CAP-004-ship/README.md), [ADR-011](docs/decisions/ADR-011-version-stamping.md)).

## Status

Baseline complete ([P-001](docs/plans/P-001-elaboration-baseline.md)); distribution and greenfield bootstrap under way ([P-002](docs/plans/P-002-leaves-home.md)). This repo dogfoods its own conventions from commit one — start reading at [docs/README.md](docs/README.md). Agents: see [AGENTS.md](AGENTS.md).

## License

[Apache 2.0](LICENSE)
