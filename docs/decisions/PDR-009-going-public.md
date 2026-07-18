---
id: PDR-009
type: decision
status: verified
links: [G-003]
title: Going public is a campaign — readiness first, one release, then the flip
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, planning conversation — visibility goal, guide form, and release sequencing chosen explicitly)
---

# PDR-009 — Going public is a campaign

## Context and problem statement

The methodology, its distribution story, and its foreign-soil validation are complete, and the repository's purpose — a methodology strangers can adopt — cannot be served while no stranger can reach it. Flipping a repository public is effectively irreversible: clones, module proxies, and search caches persist whatever was visible. How does Cliewen go public without shipping a stale install story, an insider-only corpus as the sole documentation, and no front door for contributors?

## Decision outcome

**The repository goes public at the end of a dedicated readiness campaign, in a fixed order: readiness changes land first, release v0.4.0 is cut from them, then visibility flips and the guide goes live in the same act.**

- **Readiness precedes the flip.** The corpus must be stranger-safe, the community front door (contribution, conduct, security-reporting paths) must exist, and a human-readable guide must be **deployment-ready**: built green with dead links failing the build, its deploy pipeline in place but gated, because the host cannot serve a Pages site for a private repository. The campaign's milestones are the checklist.
- **The flip and the guide going live are one act.** Flipping visibility, ungating the deploy, and verifying the guide live are a single final step; the repository is not announced-by-visibility with the guide dark, and the guide cannot precede the flip.
- **The goes-public release is v0.4.0, cut last.** It is the release that strangers first see: it carries the readiness work and an install story free of private-repo caveats, and the flip follows it. Earlier releases keep their private-era text as historical record.
- **The guide is a handwritten site, not a rendered corpus.** A static documentation site (VitePress on GitHub Pages) explains the methodology, the directory taxonomy, the change loop, and the skills in narrative form, and links to the corpus on GitHub as the living example. The corpus is written for practitioners and agents inside the loop; rendering it verbatim would make every corpus convention a site-breakage risk and would not serve a newcomer.
- **Citation policy: the corpus never cites a private repository's artifacts as normative references.** A reader must be able to resolve every reference the corpus relies on. Historical adopter evidence naming private repos (extraction targets, adopter-CI proofs) stays — evidence is not rewritten — with the analysis index noting that some cited adopter repositories are private.

**Carrier:** the campaign plan's milestones and exit criteria; CONTRIBUTING.md once it exists. Repo-local — nothing ships to adopters.

### Rejected: staying private with invited collaborators

Keeps adoption evidence self-selected — the foreign-soil trials already showed that the interesting failures come from ground the maintainer did not prepare.

### Rejected: flipping public without a guide

The corpus map orients insiders; the root README is one page. A stranger who bounces off the corpus on day one is evidence lost, and first impressions cannot be re-run.

### Rejected: rendering the corpus as the public site

Couples the site's build to corpus frontmatter and conventions, publishes insider status machinery as if it were documentation, and still would not answer a newcomer's first questions in order.
