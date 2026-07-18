---
id: CH-024
type: change
status: open
links: [P-002, M-007, CAP-004]
title: Generate standalone skills from shared canonical sources
---

# CH-024 — Generate standalone skills from shared canonical sources

## What

Replace independently authored `clue-*` skill files with deterministic generated outputs composed from skill-specific templates and shared instruction fragments. Keep every installed skill self-contained and preserve the five existing skill names, while making shared methodology rules single-source at authoring time.

Add a repository generator that writes both `.agents/skills/` and the embedded `clue init` skill templates, plus a test that fails when either generated tree differs from the canonical sources.

## Why

The skills intentionally repeat cross-cutting rules so each installed skill can stand alone, but independent copies can drift when a methodology rule changes. Generated standalone outputs retain portability and explicit skill routing while moving consistency from editorial discipline to a machine-enforced build step.

This change prepares P-002/M-007 by making the skill set internally coherent before foreign-soil trials. It does not count as a trial or as the methodology adjustment that M-007 requires from a qualifying external finding.

## Decision

Record the canonical-source and generated-output boundary as an architectural decision because it changes the durable structure and maintenance path of the versioned skill artifacts.
