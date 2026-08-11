---
id: CAP-008-design
type: design
status: active
links: [CAP-008]
title: Design for CAP-008 local verification
---

# Design — CAP-008 Local verification

CONTRIBUTING.md remains the sole repository-local command list. Its coverage report passes the `-func` flag and the coverage profile as separate arguments, which preserves the report while avoiding the Windows Go invocation failure. Focused Go tests read that verification block and reject the deprecated equals form, so a documentation edit cannot silently reintroduce it.
