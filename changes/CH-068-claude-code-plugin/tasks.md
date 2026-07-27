---
id: CH-068-tasks
type: tasks
status: open
links: [CH-068]
title: Tasks for CH-068
---

# CH-068 — tasks

- [x] Write ADR-031: the marketplace ships a bootstrap, never the managed skills; the bootstrap pins no `clue` version; the plugin manifest omits `version`
- [x] Add AC-039 to CAP-001's criteria — the published plugin bootstraps and pins nothing, and ships none of the managed lifecycle skills
- [x] Write `.claude-plugin/marketplace.json` at the repository root: marketplace `cliewen`, one plugin entry, relative source
- [x] Write `plugins/cliewen/.claude-plugin/plugin.json` — metadata only, no `version`, no component paths that reach outside the plugin directory
- [x] Write the bootstrap skill `plugins/cliewen/skills/setup/SKILL.md` — detect host, install via the M-030 channels, verify a release version, ask before `clue init`, and say what it does not install (AC-039)
- [x] Add the AC-039 positive and negative tests over the committed plugin tree (AC-039)
- [x] Add `guide/plugin.md` — the page whose reason to exist is stating what the plugin does not install — with its single Next action, and register it in the VitePress sidebar
- [x] Point `guide/getting-started.md` at the plugin page as the coding-agent route
- [x] Record CAP-001's design lesson: the plugin channel and why the managed skills stay out of it
- [x] Add the `[Unreleased]` changelog entry, written for a user
- [x] Run `go test ./...`, `clue validate`, `go generate ./internal/skills` with no drift, and `npm run guide:build`
