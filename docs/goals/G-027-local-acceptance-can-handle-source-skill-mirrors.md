---
id: G-027
type: goal
status: proposed
links: [VIS-001]
title: Local acceptance can handle the source repository's managed skill mirrors
---

# Accept source work locally with managed mirrors

The source repository selects local acceptance, but its committed Claude skill entry points and reference directories use Git symlinks. The current local command refuses all tracked symlinks before validating a snapshot. Consequently a source change cannot receive native local preflight merely because source role is now permitted.

The command should eventually support a reviewed, complete and safe source snapshot or provide an explicitly accepted alternative for this repository. This is a proposed follow-up from the readable-reference change, not authority to weaken snapshot checks or silently accept work. The observed source tree is the accepted release commit `5be6182506f1d2f58da3ea87ee359fe259d0578e`; `git ls-tree -r main` names the managed symlinks, and the local command's materializer rejects their mode.
