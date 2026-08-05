---
links: [P-011]
---

# CH-116 — Implement the corpus-wide identity ledger

## What

[ADR-048](../../docs/decisions/ADR-048-corpus-wide-id-ledger.md) decided the design; this change implements it, closing M-052 of [P-011](../../docs/plans/P-011-truthful-brownfield-migration.md):

- `internal/ledger/` — a `Ledger` type over `.clue/id-ledger.yaml`: `byID` map (exact canonical ID → entry) and a per-prefix `counters` map, both populated by one `Load()`. Entries carry `id`, `kind` (`numeric | opaque`), `state` (`reserved | live | retired`), and — for numeric — `prefix`/`component`; an imported entry additionally carries `source-revision`/`source-location`.
- `clue id next <prefix>` — a new CLI subcommand that allocates the next numeric ID for a prefix as an O(1) counter increment, reserves it in the ledger, persists, and prints the new ID.
- `checkLedger` — a new `clue validate` rule in `internal/corpus/rules.go` cross-checking every live corpus ID against the ledger: rejects a live artifact whose ID is `retired` or `reserved` under a different reservation, rejects a numeric entry with no valid decimal component and an opaque entry with one, and rejects an opaque ID not preserved verbatim. Gated so a corpus without a ledger file yet is unaffected.
- `clue migrate` backfill step — a new idempotent migration that seeds the ledger from the current corpus scan (one `live` entry per currently-live ID, `counters` seeded at each prefix's current max) so existing history is not renumbered.

## Why

Without a persisted ledger, a retired native ID's number can be silently re-minted, colliding with historic evidence, commit history, or external references that still carry the deleted artifact's original meaning ([ADR-048](../../docs/decisions/ADR-048-corpus-wide-id-ledger.md)). ADR-048 records the design; M-052 remains `todo` until the ledger, allocator, validator rule, and backfill actually exist.

## Scope

In scope: the four build items above, focused positive/negative fixtures proving an archived numeric ID above the live maximum and a UUID-like opaque ID cannot be reused (M-052's stated exit fixtures), and closing M-052 in the P-011 plan table in the digest.

Out of scope: wiring `clue id next` into the `clue-delta` skill's allocation guidance as the *required* path (the skill still permits Git-history scanning as a fallback) — that migration of agent guidance is a follow-up once the ledger has run in this repository for a while. Also out of scope: an opaque-namespace source-owned generator implementation for any specific brownfield source — this change proves the *contract* (verbatim preservation, generator-checked-against-`byID`) with a fixture generator, not a production one (that belongs to M-054/M-056).
