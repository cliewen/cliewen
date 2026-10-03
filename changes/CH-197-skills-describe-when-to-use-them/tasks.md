---
id: CH-197-tasks
type: tasks
status: open
links: [CH-197]
title: CH-197 tasks
---

# CH-197 tasks

- [x] Probe one skill with `name` and `description` beside the existing keys in a fresh Claude Code session; stop and report if it is renamed, dropped, or loses its references
- [x] Add AC-217 (CAP-004): every generated entry point carries `name` equal to its directory and a distinct `description` of when to use it, within the format's length limit and never the generated-file comment; positive and negative evidence
- [x] Add a when-to-use clause to each skill definition, emit `name` and `description` from the generator, and run `go generate ./internal/skills` (AC-217)
- [x] Update the CAP-004 design where it describes skill frontmatter
