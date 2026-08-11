---
id: CAP-008-design
type: design
status: active
links: [CAP-008]
title: Design for CAP-008 local verification
---

# Design — CAP-008 Local verification

CONTRIBUTING.md remains the sole repository-local command list. Its coverage commands pass each flag and its value as separate arguments, which preserves the report while keeping the block independent of how a shell splits `-flag=value` tokens — PowerShell's default native-argument mode splits a single-dash token containing `=` at the first dot in its value. The equals form stays valid Go flag syntax and is not deprecated; CI, which runs under Bash, keeps using it. Focused Go tests read that verification block and reject the equals form of the coverage report, so a documentation edit cannot silently reintroduce it.
