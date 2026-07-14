---
version: 0.1.0
---

# clue-delta

The change loop. Use for every mutation of `main` — features, fixes, docs, plans. Two tiers: the **full loop** below, and a **light tier** (bottom of this skill) for changes that decide nothing and change no meaning.

1. **Branch** off `main`, named `ch-xxx-slug`. Take the next free CH number (grep `git log` and `/changes/` for the highest used).
2. **Propose**: create `/changes/CH-xxx-slug/` with
   - `proposal.md` — what & why; frontmatter `links` MUST reference the plan item it serves (P/M-IDs) or state plan-less explicitly. No fake plan items.
   - `tasks.md` — task breakdown as `- [ ]` checkboxes; sub-tasks nest one level. Order tasks so dependencies come first: a task that needs another done lists after it. Mark `[x]` **the moment a task completes** — never in batch at the end; the unticked list is what tells you (and anyone watching the branch) what is actually left. Mark `[-]` for addressed-but-not-feasible, and a `[-]` must carry its reason on the same line (e.g. a third-party library turned out to lack the feature) — a skip without a reason is an ignored task. A task that changes behavior names the AC-IDs it serves. Tests never bind to tasks — tasks are transient and die at the digest; tests trace to ACs. A behavior-changing task with no AC to name means a criterion is missing in `criteria.md` — add it there first.
   - `open-questions.md` — blocking questions. When one appears, write it and **stop**; the human's answer is recorded as a decision — an ADR when it constrains future changes, a decision-log row when reversing it later is cheap and local.

   Commit this folder before implementing: the branch is the proposal.
3. **Implement** against the corpus: capabilities get/update their folder (README / criteria / design), ACs are Gherkin with `@AC-xxx` tags, every AC gets a positive + negative test pair. An AC too complex to verify with a focused test pair is too big: split it in `criteria.md` before writing tests — granularity is flexible, the anchor is mandatory. ACs and their tests land in the same change; a capability whose ACs cannot be tested yet stays `draft`. Every test declares exactly one purpose: `AC<digits>` for criterion verification, `Unit` for coverage backstops, `Sanity` for environment/repo invariants, `Arch` for structural checks — as a framework tag where tags exist, as the name prefix in Go. When a requirement change alters what an AC *means*, retire it and mint a new ID — never redefine an existing ID: add `@retired` to its tag line (`@AC-xxx @retired`), keep the tombstone in `criteria.md`, and remove or re-tag its tests (the linter forces this).
4. **Digest** (final commits before PR): update permanent `/docs` (including README index blocks), apply any repo-local digest conventions the repo's AGENTS.md declares (e.g. recording user-visible impact in a changelog) — local conventions extend this loop, never override it; when AGENTS.md contradicts a skill, that is a blocking question for `open-questions.md`, not a choice to make silently — record decisions made — expensive-to-reverse ones (architecture, methodology mechanics, public contracts, AC semantics) as ADRs, the rest as rows in the decision log (`docs/decisions/log.md`, columns `Date | Decision | Why | Change/PR`); the litmus test: if reversing it later is cheap and local, it's a log row; if it constrains future changes, it's an ADR — update plan bookkeeping, then **delete `/changes/CH-xxx-slug/`**. Digest precondition: every task in `tasks.md` is `[x]` or `[-]`-with-reason — an open `[ ]` means the change is not done, so finish it or downgrade it honestly.
5. **PR + merge**: the PR is the review gate (human verifies meaning), the merge commit is the acceptance. Run `clue-verify` first. `main` must never contain `/changes/`.

Keep deltas small — git merges text, not meaning; small deltas plus a post-rebase consistency check are the mitigation.

## Light tier

A change is **light** when ALL of these hold: no new decision is needed (no ADR, no decision-log row), no acceptance criterion or capability meaning is added/changed/retired, no semantic plan mutation (milestone-status bookkeeping is fine), and no methodology carrier is touched (skills, AGENTS.md rules, lint rules). Examples: typos, doc clarity, dependency bumps, pure refactors, CI plumbing.

A light change skips the transient workspace: branch `ch-xxx-slug` → commits → PR, where the **PR description is the proposal** — it states what and why, and names the plan item it serves or declares itself plan-less, exactly as a proposal.md would. CH numbering stays global across tiers (same next-free-number rule), so provenance grep works identically.

**Escalate immediately:** the moment an open question, a decision, or an AC change appears mid-work, create `/changes/CH-xxx-slug/` and continue the full loop from there. The tier is judged by what the change *became*, not what it was expected to be.
