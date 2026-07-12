---
id: CAP-004-design
type: design
status: active
links: [CAP-004, ADR-011]
title: Design for CAP-004 clue ships
---

# Design — CAP-004 `clue ships`

## Version surface

`cmd/clue` carries `var version = "dev"`. Release builds override it with `-ldflags "-X main.version=<semver>"`; source and local builds keep `dev`. `clue version` and `clue --version` print `clue <version>` (via `runVersion(io.Writer)`, so the AC test can capture it). `runValidate` threads the same package var into `corpus.Options.Version` so the drift rule can compare.

## Skill stamps and the drift rule

Each skill lives at `.agents/skills/<name>/skill.md` and declares `version:` in its frontmatter. `corpus.checkSkillVersions` (wired into `Validate`) reads them directly from disk — skills sit outside the `docs/`+`changes/` scan, like `checkFolderReadmes` already reaches past the artifact graph. Three checks, ordered by what each build can see:

1. **Stamp present (AC-020)** — a `skill.md` with no frontmatter `version` fails. Runs on every build; this is the "skills carry a version" half of M-004, enforced by cliewen's own `dev` CI.
2. **Set consistency (AC-021)** — the skills must agree on one version. The first skill (sorted by path) is the reference; any that disagrees is reported. Per-skill markers, kept a set by the rule ([ADR-011](../../decisions/ADR-011-version-stamping.md)). Also runs on every build.
3. **Drift vs the binary (AC-022)** — only when the skills already agree *and* the binary is a real release (`version` is neither `""` nor `dev`): each skill version must equal the binary's, else it is drift. `dev`/source builds skip this — the preview contract is "OK (skips drift)". This is what an adopted repo's installed (released) `clue` uses to answer "are my skills current?".

Issues sort into the same `path: message` stream as every other rule.

## Release pipeline

`.github/workflows/release.yml` triggers on `v*` tags (and `workflow_dispatch`). One ubuntu runner cross-compiles Go for linux/darwin/windows × amd64/arm64 with `GOOS`/`GOARCH`, stamping `main.version` from the tag; `softprops/action-gh-release` attaches the binaries. Git tags carry the conventional `v` prefix (Go module tagging); the stamped and skill-frontmatter version is bare semver, so the workflow strips the `v` (`${GITHUB_REF_NAME#v}`) — the drift rule then compares like against like.

Distribution targets the private-repo install story (P-002): `go install github.com/cliewen/cliewen/cmd/clue@vX.Y.Z` builds from source (reports `dev` unless the tag is stamped in), and `gh release download vX.Y.Z` fetches a stamped prebuilt binary.

## Deliberate limits (doors, not gaps)

- **Skills outside `.agents/skills/`** — the rule looks only there. An adopter who relocates skills gets no drift check until the path is made configurable.
- **`go install` reports `dev`** — ldflags stamping happens in the release job, not on a bare `go install`; the stamped provenance lives in the downloadable release artifacts. A future `-ldflags` recipe or `go generate` step could stamp installs too.
- **No `clue`-driven bump** — cutting a release edits the five skill frontmatter versions by hand and tags. A `clue bump` command is post-M-004.
- **Operational proof of the release** — the cross-platform build is guarded structurally by `TestSanity_ReleaseWorkflowIsCrossPlatform`; the real proof is the first tagged release producing downloadable binaries.
