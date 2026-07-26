---
id: CH-066-open-questions
type: open-questions
status: resolved
links: [CH-066]
title: Open questions for CH-066
---

# Open questions

**Resolved — how should Linux be served, once Homebrew cannot serve it?** The change was first built around a Homebrew formula covering macOS and Linux. That channel turned out to be deprecated upstream, and its replacement is macOS-only, so the question became real rather than theoretical: accept a deprecated dependency, drop Linux to the manual download, or find another channel. Resolved by the maintainer in favour of using nothing deprecated, which left a verified install script as the only route giving all three platforms one command without new infrastructure. Recorded in [ADR-030](../../docs/decisions/ADR-030-verified-install-scripts.md).

Two implementation unknowns were settled empirically rather than by decision, and are noted so a reviewer can see they were not assumed:

1. `formats: [binary]` does carry the `.exe` extension into the Windows asset name, and the generated `SHA256SUMS` lists **upload** names rather than on-disk paths — so an adopter's `sha256sum -c --ignore-missing SHA256SUMS` passes unchanged. Verified against a local dry run and a simulated adopter wall.
2. Both install scripts were run end to end against the live v0.7.0 release: correct release resolution, download, checksum verification, install, and a working `clue version`. The failure path was exercised too — a tampered binary is refused and nothing is installed.
