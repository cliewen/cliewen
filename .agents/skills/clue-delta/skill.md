# clue-delta

The change loop. Use for every mutation of `main` — features, fixes, docs, plans.

1. **Branch** off `main`, named `ch-xxx-slug`. Take the next free CH number (grep `git log` and `/changes/` for the highest used).
2. **Propose**: create `/changes/CH-xxx-slug/` with
   - `proposal.md` — what & why; frontmatter `links` MUST reference the plan item it serves (P/M-IDs) or state plan-less explicitly. No fake plan items.
   - `tasks.md` — task breakdown as `- [ ]` checkboxes; sub-tasks nest one level. A task that changes behavior names the AC-IDs it serves. Tests never bind to tasks — tasks are transient and die at the digest; tests trace to ACs. A behavior-changing task with no AC to name means a criterion is missing in `criteria.md` — add it there first.
   - `open-questions.md` — blocking questions. When one appears, write it and **stop**; the human's answer becomes an ADR.

   Commit this folder before implementing: the branch is the proposal.
3. **Implement** against the corpus: capabilities get/update their folder (README / criteria / design), ACs are Gherkin with `@AC-xxx` tags, every AC gets a positive + negative test pair. An AC too complex to verify with a focused test pair is too big: split it in `criteria.md` before writing tests — granularity is flexible, the anchor is mandatory. ACs and their tests land in the same change; a capability whose ACs cannot be tested yet stays `draft`.
4. **Digest** (final commits before PR): update permanent `/docs` (including README index blocks), write ADRs for decisions made, update plan bookkeeping, then **delete `/changes/CH-xxx-slug/`**.
5. **PR + merge**: the PR is the review gate (human verifies meaning), the merge commit is the acceptance. Run `clue-verify` first. `main` must never contain `/changes/`.

Keep deltas small — git merges text, not meaning; small deltas plus a post-rebase consistency check are the mitigation.
