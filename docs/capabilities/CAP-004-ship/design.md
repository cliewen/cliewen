---
id: CAP-004-design
type: design
status: active
links: [CAP-004, ADR-011, ADR-012, ADR-021]
title: Design for CAP-004 clue ships
---

# Design — CAP-004 `clue ships`

## Version surface

`cmd/clue` carries `var version = "dev"`. Release builds override it with `-ldflags "-X main.version=<semver>"`; when no stamp was injected, `main` falls back to the module version Go embeds in `go install module@vX.Y.Z` builds (`debug.ReadBuildInfo`, `v` trimmed to bare semver) — `(devel)`, a pseudo-version (branch/commit install, or the VCS-derived version a checkout build embeds since Go 1.24), and a `+dirty` suffix (local modifications) all stay `dev`. `clue version` and `clue --version` print `clue <version>` (via `runVersion(io.Writer)`, so the AC test can capture it). `runValidate` threads the same package var into `corpus.Options.Version` so the drift rule can compare.

## Skill stamps and the drift rule

Each skill lives at `.agents/skills/<name>/skill.md` and declares `version:` in its frontmatter. `corpus.checkSkillVersions` (wired into `Validate`) reads them directly from disk — skills sit outside the `docs/`+`changes/` scan, like `checkFolderReadmes` already reaches past the artifact graph. Three checks, ordered by what each build can see:

1. **Stamp present (AC-020)** — a `skill.md` with no frontmatter `version` fails. Runs on every build; this is the "skills carry a version" half of M-004, enforced by cliewen's own `dev` CI.
2. **Set consistency (AC-021)** — the skills must agree on one version. The reference is the version the majority of skills carry (ties go to the earliest-sorted holder), so the report names the actual outlier. Per-skill markers, kept a set by the rule ([ADR-011](../../decisions/ADR-011-version-stamping.md)). Also runs on every build. Malformed frontmatter and non-string versions (YAML reads `1.0` as a number) are named as such, not as missing stamps.
3. **Drift vs the binary (AC-022)** — only when the skills already agree *and* the binary is a real release (`version` is neither `""` nor `dev`): each skill version must equal the binary's, else it is drift. `dev`/source builds skip this — the preview contract is "OK (skips drift)". This is what an adopted repo's installed (released) `clue` uses to answer "are my skills current?".

Issues sort into the same `path: message` stream as every other rule.

## Skill authoring and generation

The complete skill files remain the installed, versioned artifacts described above, but their authoring source is centralized under `internal/skills/source/` ([ADR-021](../../decisions/ADR-021-generated-standalone-skills.md)). Skill-specific templates compose shared instruction fragments at generation time; no generated skill has a runtime dependency on another skill or on the source tree.

`go generate ./internal/skills` writes both `.agents/skills/` and `internal/scaffold/templates/skills/`. The generator owns the five `clue-*` directories in both trees, including `clue-extract` resources, and normalizes output to LF for deterministic bytes. AC-028 tests exercise generation into an empty repository and reject missing, changed, or unexpected files. A Sanity test checks the committed outputs against the same in-memory rendering, so editing a generated skill instead of its source fails the build.

## Release pipeline

`.github/workflows/release.yml` triggers on `v*` tags (and `workflow_dispatch`, which a guard restricts to tag refs — dispatching from a branch fails before any build, so nothing can stamp binaries as a branch name or publish a branch-named release). The release body is the tag's `## [X.Y.Z]` section of `CHANGELOG.md`, extracted verbatim; a missing or empty section fails the release before anything is built, and `TestSanity_ReleaseNotesComeFromChangelog` keeps auto-generated notes from returning ([ADR-012](../../decisions/ADR-012-release-notes-from-changelog.md)). The job runs `go test ./...` first — a tag can land on any commit, so nothing ships untested. One ubuntu runner then cross-compiles Go for linux/darwin/windows × amd64/arm64 with `GOOS`/`GOARCH`, stamping `main.version` from the tag, writes a `SHA256SUMS` file next to the binaries, and `softprops/action-gh-release` (pinned by commit SHA — it runs with `contents: write`) attaches everything. Git tags carry the conventional `v` prefix (Go module tagging); the stamped and skill-frontmatter version is bare semver, so the workflow strips the `v` (`${GITHUB_REF_NAME#v}`) — the drift rule then compares like against like.

Distribution targets the private-repo install story (P-002): `go install github.com/cliewen/cliewen/cmd/clue@vX.Y.Z` builds from source and self-stamps from Go's build info, and `gh release download vX.Y.Z` fetches a stamped prebuilt binary.

## Deliberate limits (doors, not gaps)

- **Skills outside `.agents/skills/`** — the rule looks only there. An adopter who relocates skills gets no drift check until the path is made configurable.
- **Untagged installs report `dev`** — `go install module@vX.Y.Z` self-stamps from Go's build info, but a checkout build or an install of a branch/commit (a pseudo-version) has no release to report and stays `dev`, exempt from drift.
- **No `clue`-driven bump** — cutting a release edits the shared generated-skill frontmatter source, regenerates the five skill artifacts, and tags. A `clue bump` command is post-M-004.
- **Operational proof of the release** — the cross-platform build is guarded structurally by `TestSanity_ReleaseWorkflowIsCrossPlatform`; the real proof is the first tagged release producing downloadable binaries.
