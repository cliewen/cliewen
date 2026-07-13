---
version: 0.1.0
---

# clue-delta

The change loop. Use for every mutation of `main` — features, fixes, docs, plans.

1. **Branch** off `main`, named `ch-xxx-slug`. Take the next free CH number (grep `git log` and `/changes/` for the highest used).
2. **Propose**: create `/changes/CH-xxx-slug/` with
   - `proposal.md` — what & why; frontmatter `links` MUST reference the plan item it serves (P/M-IDs) or state plan-less explicitly. No fake plan items.
   - `tasks.md` — task breakdown as `- [ ]` checkboxes; sub-tasks nest one level. Order tasks so dependencies come first: a task that needs another done lists after it. Mark `[x]` **the moment a task completes** — never in batch at the end; the unticked list is what tells you (and anyone watching the branch) what is actually left. Mark `[-]` for addressed-but-not-feasible, and a `[-]` must carry its reason on the same line (e.g. a third-party library turned out to lack the feature) — a skip without a reason is an ignored task. A task that changes behavior names the AC-IDs it serves. Tests never bind to tasks — tasks are transient and die at the digest; tests trace to ACs. A behavior-changing task with no AC to name means a criterion is missing in `criteria.md` — add it there first.
   - `open-questions.md` — blocking questions. When one appears, write it and **stop**; the human's answer becomes an ADR.

   Commit this folder before implementing: the branch is the proposal.
3. **Implement** against the corpus: capabilities get/update their folder (README / criteria / design), ACs are Gherkin with `@AC-xxx` tags, every AC gets a positive + negative test pair. An AC too complex to verify with a focused test pair is too big: split it in `criteria.md` before writing tests — granularity is flexible, the anchor is mandatory. ACs and their tests land in the same change; a capability whose ACs cannot be tested yet stays `draft`. Every test declares exactly one purpose: `AC<digits>` for criterion verification, `Unit` for coverage backstops, `Sanity` for environment/repo invariants, `Arch` for structural checks — as a framework tag where tags exist, as the name prefix in Go. When a requirement change alters what an AC *means*, retire it and mint a new ID — never redefine an existing ID: add `@retired` to its tag line (`@AC-xxx @retired`), keep the tombstone in `criteria.md`, and remove or re-tag its tests (the linter forces this).
4. **Digest** (final commits before PR): update permanent `/docs` (including README index blocks), apply any repo-local digest conventions the repo's AGENTS.md declares (e.g. recording user-visible impact in a changelog) — local conventions extend this loop, never override it; when AGENTS.md contradicts a skill, that is a blocking question for `open-questions.md`, not a choice to make silently — write ADRs for decisions made, update plan bookkeeping, then **delete `/changes/CH-xxx-slug/`**. Digest precondition: every task in `tasks.md` is `[x]` or `[-]`-with-reason — an open `[ ]` means the change is not done, so finish it or downgrade it honestly.
5. **PR + merge**: the PR is the review gate (human verifies meaning), the merge commit is the acceptance. Run `clue-verify` first. `main` must never contain `/changes/`.

Keep deltas small — git merges text, not meaning; small deltas plus a post-rebase consistency check are the mitigation.
