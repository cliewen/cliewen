# clue-verify

Pre-merge checklist. Run before opening or updating any PR. When the `clue` CLI exists, `clue validate` performs the mechanical half; until then, check by hand. **Never fix a failure by weakening the check.**

- [ ] Every artifact touched has frontmatter: `id`, `type`, `status`, `links`, `title` (+ type extensions: ADRs `author`/`accepted-by`; constraints `source`/`enforcement`; capabilities `goal`).
- [ ] Every `links` entry resolves to an existing ID.
- [ ] The proposal references a real plan item, or is declared plan-less.
- [ ] Plan bookkeeping updated (milestone statuses reflect this merge).
- [ ] No completed (`status: completed`) plan was modified.
- [ ] Every AC has a test tag with a positive + negative pair — or the capability is honestly `status: draft` with the gap stated.
- [ ] Every `/docs/**` folder has README.md; index blocks list every sibling artifact and reference no deleted file.
- [ ] Change assessed against every constraint in `/docs/constraints/` and every quality scenario in `/docs/quality/`.
- [ ] Diagrams are inline Mermaid and readable when rendered.
- [ ] `/changes/CH-xxx-slug/` is deleted in the digest commit; after merge, `main` contains no `/changes/`.
- [ ] Decisions made during the change are ADRs (`author: agent` starts `inferred`; only a human promotes to `verified`).
