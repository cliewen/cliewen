---
id: C-006
type: constraint
status: active
links: []
title: Decision records are timeless prose; a method decision names its carrier
source: docs/decisions/README.md
enforcement: agent
---

# C-006 — Decision records are timeless prose; method decisions name their carrier

A record's context states the problem, not the episode: a motivating incident earns at most one sentence. This holds for ADRs and PDRs alike. A decision that changes the methodology for adopting projects must name its carrier — the `clue` rule (machine), the skill text (agent), or the init template (default) that ships it; a method decision without a carrier has not reached anyone.

**Promotion trigger:** timelessness is meaning and stays human-reviewed; the carrier half promotes when `clue` can classify method decisions and require a carrier marker in them — then `enforcement: machine` for that half.
