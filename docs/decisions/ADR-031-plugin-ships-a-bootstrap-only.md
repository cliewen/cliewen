---
id: ADR-031
type: decision
status: verified
links: [CAP-001, CAP-004, C-015, ADR-011, ADR-022, ADR-030]
title: The Claude Code plugin ships a bootstrap, and the managed skills never ride in it
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-031 — The plugin ships a bootstrap, not the skills

## Context and problem statement

Cliewen is built for people working inside a coding agent, and Claude Code has a native distribution channel — a plugin marketplace, added with `/plugin marketplace add owner/repo` and backed by nothing more than a JSON file in a repository the user can already read. [ADR-030](ADR-030-verified-install-scripts.md) removed the six manual install steps for a human at a terminal. It did nothing for the reader already in a session, who must still leave it, find the guide, and run a shell command against [C-015](../constraints/C-015-onboarding-under-30-minutes.md)'s thirty-minute budget.

The obvious plugin to build is the one that arrives with Cliewen's five lifecycle skills already loaded. That is the shape of nearly every plugin in circulation, and it is the shape this decision refuses.

Cliewen's skills are not portable helpers. `clue init` writes them into `.agents/skills/` as committed repository files, each stamped with the version of the binary that wrote it ([ADR-022](ADR-022-skill-ownership-marker.md)), and a released `clue` fails validation when its own version and its marked skills disagree ([ADR-011](ADR-011-version-stamping.md)). That drift check is scoped to `.agents/skills/`, and it is what makes a green `clue validate` mean something: the judge, the guidance, and the corpus conventions are the same generation.

A plugin's components do not live there. Claude Code copies an installed plugin into a per-user cache and enables it through a settings file, so a plugin-shipped skill is invisible to `checkSkillVersions` no matter which installation scope is chosen — the scope decides who sees the plugin, not where its files land. It would also be one copy spanning every repository the user opens, while the whole point of a stamped skill is that it is pinned per repository.

So the failure is not that bundled skills would be redundant. It is that they would be a second, unversioned set of Cliewen instructions sitting beside the committed set, silently overriding nothing and contradicting whichever repository happened to pin a different release — and the drift check that exists to catch exactly that would never see them.

## Decision outcome

**The marketplace ships one plugin, the plugin ships one skill, and that skill's whole job is to get a real `clue` onto the machine and then stop.**

- **The bootstrap is thin on purpose.** It detects the host, installs `clue` through the channels ADR-030 published, confirms the binary reports a release version rather than `dev`, and asks before running `clue init`. Everything after that comes from the repository.
- **The five lifecycle skills are never plugin components.** `clue init` is the only supported writer of Cliewen skills, and `.agents/skills/` is the only place they are valid. This is the load-bearing half of the decision; the rest is packaging.
- **The bootstrap pins no `clue` version.** It installs the newest release. A version literal inside a channel that nobody remembers to bump is precisely the drift ADR-011 exists to detect, and it would hand new users a stale binary at the exact moment they have no way to know it.
- **The plugin manifest omits `version`.** Claude Code then treats the commit as the version, so the plugin has no hand-maintained stamp to forget. The cost is accepted openly: users are offered an update whenever this repository's default branch moves, even for changes that do not touch the plugin. A bootstrap this small is a fair trade against a second number that can contradict the binary, which is the failure ADR-011 was written after.
- **Installing the plugin is not consent to scaffold.** `clue init` writes into the user's repository, so the bootstrap asks first. A plugin install authorizes the plugin, nothing further.
- **The plugin tree is hand-authored and lives under `plugins/`,** outside both directories the skill generator owns (`.agents/skills/` and `internal/scaffold/templates/skills/`), so `go generate ./internal/skills` neither writes it nor reports it as an unexpected file.

The negative rule is the one that needs a carrier, because nothing about a plugin manifest announces it and the pressure to bundle will return with every new contributor. It is held by a criterion with tests over the committed plugin tree, and stated in prose on a guide page whose reason to exist is saying what the plugin does not install.

**Carrier:** `.claude-plugin/marketplace.json` and `plugins/cliewen/`; `guide/plugin.md`; AC-039 on [CAP-001](../capabilities/CAP-001-onboarding/criteria.md), whose tests read the committed plugin tree so a bundled lifecycle skill or a pinned version fails here rather than in a stranger's session.

### Rejected: ship the five lifecycle skills in the plugin

The default shape of a plugin, and the reason this record exists. It is rejected because it breaks the drift guarantee: the copies land in the plugin cache where `checkSkillVersions` cannot reach them, they carry one version across every repository the user opens, and a repository pinned to a different release would be advised by skills it never installed and cannot validate. The convenience is real and the failure is silent, which is the worst pairing available.

### Rejected: ship them at project scope instead of user scope

Project scope looks like the fix, because it is checked into a repository and reaches collaborators. It changes nothing that matters here: the plugin is still copied to the cache rather than into `.agents/skills/`, still invisible to the drift check, and still redundant, since a repository at project scope has already run `clue init` and owns the real, stamped set.

### Rejected: run `clue init` automatically after installing

It would make the first-run path a single action and shave a prompt off the thirty minutes. It is rejected because `clue init` creates a documentation taxonomy, an `AGENTS.md`, skills, and a CI workflow in a repository the user has not yet said should become a Cliewen repository. Scaffolding is a decision, not a step.

### Deferred: publishing to a wider plugin catalog

Listing in a community or organization-level catalog is additive and needs nothing this decision forecloses; the marketplace manifest is already the artifact such a catalog would point at. It stays out until the bootstrap has been used by someone who did not write it.
