# Claude Code entry point

This repository's agent instructions live in [`AGENTS.md`](AGENTS.md), which is the cross-agent standard and the file to edit. This file exists only because Claude Code loads `CLAUDE.md` and nothing else, so without it none of the routing hub reaches the session.

@AGENTS.md

## Why this file is a pointer and not a copy

The routing rules are already stated once, and a second copy would drift silently — the two would disagree for exactly as long as it takes someone to notice, which for agent instructions is usually one wrong change. If you are adding a rule, add it to `AGENTS.md`.

## Where methodology rules actually live

`AGENTS.md` routes; it does not restate the method. The lifecycle rules an agent must follow are in the generated skills under `.agents/skills/`, mirrored to `.claude/skills/` by a symlink so Claude Code lists them. **Read the skill before deciding how a change is shaped** — a rule found in `/docs` alone may be this repository's own bookkeeping rather than something the method binds, and a rule missing from the skills never reaches an adopter at all.

The skills are generated. Edit `internal/skills/source/`, then run `go generate ./internal/skills`; never edit a file under `.agents/skills/` directly.
