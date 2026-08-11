---
id: CAP-008-design
type: design
status: active
links: [CAP-008]
title: Design for CAP-008 local verification
---

# Design — CAP-008 Local verification

CONTRIBUTING.md remains the sole repository-local command list. Its coverage commands pass each flag and its value as separate arguments, which preserves the report while keeping the block independent of how a shell splits `-flag=value` tokens — PowerShell's default native-argument mode splits a single-dash token containing `=` at the first dot in its value. The equals form stays valid Go flag syntax and is not deprecated; CI, which runs under Bash, keeps using it. Focused Go tests read that verification block and reject any single-dash `-flag=value` whose value contains a dot, so a documentation edit cannot silently reintroduce the defect on the coverage commands or on a command added later. The guard is the defect class rather than a list of approved command strings: naming only the two commands is what let the report be repaired while the profile line above it stayed broken.
