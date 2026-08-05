# Decisions

Three records live here ([PDR-006](PDR-006-decision-records-are-typed.md), [PDR-003](PDR-003-decision-log.md)), routed by two questions. Is reversing the decision cheap and local? Then it is a row in the **[decision log](log.md)** — one dated line per decision. Otherwise, is it about architecture or about how the project works? **ADRs** (Architectural Decision Records, in Nygard's strict sense) record decisions about the structure of the software and the corpus format; **PDRs** (Project/Process Decision Records) record decisions about how the project works — change tiers, acceptance mechanics, validation strategy. A decision that adopts a well-established practice cites it by name and records only the local why and deviations.

ADRs and PDRs share the MADR format and two-tier provenance: `inferred` (binding once merged, but no human has signed it) → `verified` (at least one human has explicitly approved — the act that makes provenance auditable). Every record carries `author: agent|human` and `accepted-by:`. **Rejected alternatives are half of "why does the system look like this"** — rejected records live here too. **Decisions are never deleted** — retention applies to the decision, not the file format: a record demoted under the litmus test survives as a dated log row, a record filed under the wrong type is renamed into the right series, and git history keeps the full provenance.

**Retirement, generalized to every artifact type, is deletion** ([ADR-034](ADR-034-retirement-is-deletion.md)): the file goes, a `supersedes:` frontmatter field on the successor (or, for a demoted decision with no successor record, on `log.md` itself) names the dead ID, and git history is the archive. `clue validate` rejects a `supersedes:` entry whose ID still exists on disk and any surviving link to an ID that no longer resolves. Criteria tombstones and completed plans are the named exceptions — see ADR-034.

**Merge binds, approval signs** ([PDR-004](PDR-004-merge-binds-approval-signs.md), superseding [PDR-001](PDR-001-pr-approval-promotes-adrs.md)): merging a PR makes its `inferred` decisions binding — in force immediately, no signature required — but does not touch their status. Only an explicit human approval (review approval, review comment, or a stated "approved") flips a decision to `verified`: each approver signs `accepted-by:`, approvals accumulate, and the acceptance date is the first approval. An explicit objection keeps a decision `inferred` regardless of other approvals and becomes an open question; unresolved reviewer disagreement is an objection. The agent performs the clerical signing, citing approver, date, and venue.

The CLI reports inferred ADRs and PDRs as decisions awaiting verification, separately from high-cost inferred non-decision meaning that blocks capability activation. It never turns the missing signature itself into a validation failure: merge binds and approval signs remain separate acts.

**`accepted-by:` records only approval given under Cliewen's merge boundary** ([ADR-029](ADR-029-accepted-by-is-cliewen-approval-only.md)), never acceptance a record already carried before it entered this corpus. A record extracted from a source with its own acceptance history — a MADR record's decision-makers, for instance — preserves that history as body prose with its original names, roles, and dates, and keeps `accepted-by: []`, exactly the shape any unsigned record already has.

**Records are timeless.** Context states the problem, not the episode: a motivating incident earns at most one sentence, and the change history lives in git log and the plans. Concrete mechanisms appear as decision content — the chosen option, the rejected options, the carrier — never as narrative.

**Same-change carrier rule for method decisions:** a decision that changes a methodology contract inventories every live carrier that states the affected contract and updates that complete inventory in the same change. Live carriers include current corpus truth, canonical and generated skills, templates, public or contributor guidance, implementation explanations, CLI text, and distribution metadata. Historical analyses, completed plans, and changelog entries remain pinned history. Focused guards hold stable repaired claims, but no current mechanism derives an arbitrary contract's complete carrier set, so the general obligation remains agent-enforced. The foundation new projects receive has exactly one authoritative form: the output of `clue init` plus the rules of the `clue` binary — and CAP-001's criteria hold that output to account.

<!-- clue:index:start -->
- [ADR-001 — Implementation language for the clue CLI: Go](ADR-001-implementation-language.md) · `verified`
- [ADR-002 — The inbox is goals with status: proposed](ADR-002-inbox-is-proposed-goals.md) · `verified`
- [ADR-005 — Tests reference ACs via framework-native tags; function names where no tags exist](ADR-005-test-reference-convention.md) · `verified`
- [ADR-006 — Every test declares its purpose from a small taxonomy](ADR-006-test-purpose-taxonomy.md) · `verified`
- [ADR-007 — AC lifecycle — meaning-immutable IDs, retirement by tombstone](ADR-007-ac-lifecycle.md) · `verified`
- [ADR-008 — Brownfield extraction is one generic skill with per-source mappings](ADR-008-extraction-is-a-skill.md) · `verified`
- [ADR-009 — AC IDs are namespaced — criteria declare an ac-prefix](ADR-009-ac-id-namespaces.md) · `verified`
- [ADR-010 — Extracted artifacts carry a provenance field, born inferred](ADR-010-provenance-field.md) · `verified`
- [ADR-011 — clue and the skills are versioned — tag-stamped binary, per-skill markers, drift is a failure](ADR-011-version-stamping.md) · `verified` · location-only enrollment superseded by ADR-022
- [ADR-012 — Release notes are user-facing and come from CHANGELOG.md — extracted verbatim, missing section fails the release](ADR-012-release-notes-from-changelog.md) · `verified`
- [ADR-013 — What ships to adopters is generic; AGENTS.md is the repo-local layer](ADR-013-ships-generic-vs-repo-local.md) · `verified` · index-block clause superseded by ADR-019
- [ADR-017 — Prose conventions register as constraint artifacts with enforcement classes](ADR-017-conventions-are-constraints.md) · `verified`
- [ADR-018 — The init scaffolding is embedded in the clue binary](ADR-018-init-templates-embedded.md) · `verified`
- [ADR-019 — Index regeneration runs in clue init; ADR-013's emits-empty clause is superseded](ADR-019-init-regenerates-indexes.md) · `verified`
- [ADR-020 — The scaffolded register seeds only conventions without a versioned carrier](ADR-020-scaffolded-register-scope.md) · `verified`
- [ADR-021 — Skills are generated as standalone artifacts from shared canonical sources](ADR-021-generated-standalone-skills.md) · `verified`
- [PDR-001 — PR approval is ADR acceptance — the agent performs the clerical promotion](PDR-001-pr-approval-promotes-adrs.md) · `verified` · superseded by PDR-004
- [PDR-002 — A light change tier — the PR description is the proposal](PDR-002-light-change-tier.md) · `verified` · the light/full split now applies after the plain-route boundary in PDR-011
- [PDR-003 — ADRs are for expensive-to-reverse decisions; the rest is a decision log](PDR-003-decision-log.md) · `verified` · superseded by PDR-006
- [PDR-004 — Merge makes a decision binding; approval verifies it — approvers sign, first signature dates it](PDR-004-merge-binds-approval-signs.md) · `verified`
- [PDR-005 — Validation requires foreign soil — trials on external repos, as findings not adoptions](PDR-005-foreign-soil-trials.md) · `verified`
- [PDR-006 — Decision records are typed — ADRs for architecture, PDRs for project/process, log rows for the cheap](PDR-006-decision-records-are-typed.md) · `verified`
- [PDR-007 — The PR is the authorization boundary — changes root at main and humans merge](PDR-007-review-boundary.md) · `verified` · PDR-011 narrows the one-in-flight slot to Cliewen changes; PDR-016 refines the initiating-author scope and supersedes mandatory rebase for published PRs
- [PDR-008 — A declared plan revision may ride with its implementing change](PDR-008-plan-revisions-may-ride.md) · `verified`
- [LOG-001 — Decision log](log.md) · `active` · dated rows for the cheap-to-reverse (ADR-003 and ADR-004 demoted here)
- [PDR-009 — Going public is a campaign — readiness first, one release, then the flip](PDR-009-going-public.md) · `verified`
- [ADR-022 — Cliewen skills declare ownership in frontmatter](ADR-022-skill-ownership-marker.md) · `verified`
- [PDR-010 — Community participation enters through structured intake, private safety channels, and human review](PDR-010-community-participation.md) · `verified` · PDR-011 removes Cliewen fields from plain PRs
- [ADR-023 — The public guide is an isolated VitePress site with a visibility-gated Pages deployment](ADR-023-public-guide-architecture.md) · `verified`
- [PDR-011 — Plain changes stay outside Cliewen while retaining human merge](PDR-011-plain-changes-bypass-cliewen.md) · `verified`
- [PDR-012 — Every Cliewen change receives an automatic agentic review before publication](PDR-012-agentic-review-before-publication.md) · `verified`
- [ADR-024 — The public guide is canonical at the cliewen.dev root](ADR-024-custom-domain-root.md) · `verified`
- [PDR-013 — Cliewen has an explicit core behind a red line; everything else is periphery adopters may extend](PDR-013-explicit-core-red-line.md) · `verified`
- [ADR-025 — One default status lifecycle — draft → active → retired — plus named exceptions](ADR-025-one-status-lifecycle.md) · `verified`
- [ADR-026 — Unknown artifact types are adopter extensions, validated against the default lifecycle](ADR-026-adopter-types-default-lifecycle.md) · `verified`
- [ADR-027 — Quality scenarios are constraints — the quality type folds into the register](ADR-027-quality-scenarios-are-constraints.md) · `verified`
- [ADR-028 — A skill's manifest is resolved by case-folded name, so the judge reaches one verdict on every filesystem](ADR-028-deterministic-skill-manifest.md) · `verified`
- [ADR-029 — accepted-by records only approval given under Cliewen's merge boundary](ADR-029-accepted-by-is-cliewen-approval-only.md) · `verified`
- [PDR-014 — Installation distribution reopens on the active campaign, not a successor plan](PDR-014-distribution-reopens-on-the-active-campaign.md) · `verified`
- [ADR-030 — Installation is a checksum-verifying script, and release asset names are an append-only contract](ADR-030-verified-install-scripts.md) · `verified`
- [ADR-031 — The Claude Code plugin ships a bootstrap, and the managed skills never ride in it](ADR-031-plugin-ships-a-bootstrap-only.md) · `verified`
- [ADR-032 — Acceptance criteria declare classified proof and paired directions](ADR-032-classified-ac-evidence.md) · `verified`
- [PDR-015 — Merging the release pull request cuts the release; the tag is derived, not performed](PDR-015-merging-cuts-the-release.md) · `verified`
- [PDR-016 — Hosted PR state carries review findings and updater handoffs across agents](PDR-016-pr-state-carries-agent-handoffs.md) · `verified` · refines PDR-007 clause 2 and supersedes clause 6
- [PDR-017 — The merge gate carries an acceptance brief, and the review loop adds an advisory scenario-resolution step](PDR-017-merge-gate-has-content.md) · `verified`
- [ADR-033 — Human proof class, a per-criterion draft exemption, and derived coverage](ADR-033-human-proof-and-draft-criteria.md) · `verified`
- [ADR-034 — Retirement is deletion; supersedes carries the pointer forward](ADR-034-retirement-is-deletion.md) · `verified`
- [ADR-035 — Cost bounds inferred provenance and incident analyses return an edge from reality](ADR-035-bounded-provenance-and-reality-edges.md) · `verified`
- [PDR-018 — Behavior changes remain full until adopter evidence supports a narrower loop](PDR-018-behavior-changes-remain-full.md) · `verified`
- [PDR-019 — Methodology contract changes update every live carrier in the same change](PDR-019-methodology-contract-carriers-move-together.md) · `verified`
- [PDR-020 — Extraction rehearses before it mutates](PDR-020-extraction-rehearsal-before-mutation.md) · `verified`
- [ADR-036 — JVM evidence is credited only from one statically attributable executable](ADR-036-jvm-evidence-per-executable.md) · `verified`
- [ADR-037 — Brownfield criterion IDs preserve segmented prefixes and letter suffixes](ADR-037-brownfield-ac-id-grammar.md) · `verified`
- [PDR-021 — Full Cliewen changes are accepted with a merge commit that preserves their branch history](PDR-021-supported-merge-commit-history.md) · `verified`
- [ADR-038 — The CI wall is an upstream reusable workflow with a thin caller](ADR-038-upstream-validation-workflow.md) · `verified`
- [ADR-039 — Versioned corpus migrations plan safe mechanical upgrades](ADR-039-versioned-corpus-migrations.md) · `verified`
- [ADR-040 — External references name their target, and resolving them stays outside the judge](ADR-040-qualified-external-references.md) · `verified`
- [PDR-022 — A scaffolded vendor entry point may exist, and it may only point at AGENTS.md](PDR-022-vendor-entry-points-only-point.md) · `verified`
- [ADR-041 — Generated index rows state their record, and rows that state only their link are counted](ADR-041-index-rows-state-their-record.md) · `inferred`
- [ADR-042 — A release check reaches the network, reports, and writes nothing](ADR-042-release-check-outside-the-judge.md) · `verified`
- [ADR-043 — The managed skill set includes a human-authorized upgrade entry point](ADR-043-upgrade-skill-is-a-managed-carrier.md) · `inferred`
- [PDR-023 — The tool carries the notice and the hub carries the instruction, and no vendor configuration is ever emitted](PDR-023-tool-notice-and-hub-instruction.md) · `inferred`
- [ADR-044 — The judge judges a repository state, never a transition](ADR-044-judge-reads-state-not-transitions.md) · `inferred`
- [ADR-045 — Every constraint names the machine that holds it or the judgment that remains](ADR-045-register-names-the-machine.md) · `inferred`
- [ADR-046 — An index row says what its artifact is about](ADR-046-index-rows-say-what-the-artifact-is-about.md) · `inferred` — An appended row seeds a description from the artifact's own body; the sentence is curated thereafter, and a row still lacking one is counted rather than failed on.
- [ADR-047 — Diagram representation preserves links and assets](ADR-047-diagram-representation-preserves-links.md) · `inferred` — Diagram syntax alone does not establish whether a diagram is readable.
<!-- clue:index:end -->
