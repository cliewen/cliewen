# Tasks

- [x] Re-derive AN-006's plain-change overhead evidence at head: identify whether a comparable plain change has occurred since PDR-011 shipped the plain tier, measure it the same way (or state why no observable instance exists), and update AN-006's evidence boundary and finding.
- [ ] Re-derive AN-010's adopter-change overhead evidence at head: extend or re-check the pinned adopter history against current `origin/main`, and update AN-010's observed facts and inferences with what changed and what did not.
- [ ] Re-derive AN-012's adopter-configuration cost evidence at head: re-run `clue validate` against the same pinned adopter corpus with the current release, and update AN-012's finding on the version-upgrade-path and CI-wall-fork costs with what changed and what did not.
- [x] Ask the human, via `open-questions.md`, whether P-013 designates a successor plan and whether that plan carries AN-008 pattern C's named door (widening `supersedes:`). Answered: yes — P-014 opens with M-069 as the door.
- [ ] Write the M-066 digest: record what each re-derivation changed, did not change, and declined in P-013's milestone table; set P-013 `status: completed`; designate the successor per the human's answer (or record none); update `docs/plans/README.md` and any other generated indexes; add a CHANGELOG entry; delete this change workspace.
- [ ] Run focused and full local verification and the automatic agentic review loop, then prepare the reviewed PR handoff.
