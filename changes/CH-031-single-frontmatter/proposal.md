---
id: CH-031
type: change
status: open
links: []
title: One frontmatter per artifact — BOM-safe extraction and a duplicate-frontmatter lint
---

# CH-031 — One frontmatter per artifact

## What

Make `clue validate` reject the two file shapes that let a converted artifact carry two frontmatter blocks: a UTF-8 byte-order mark anywhere in a corpus file, and an artifact whose body opens with a second frontmatter fence. Extend the `clue-extract` target contract so a conversion must replace any existing frontmatter — even when an invisible byte prefix hides its opening fence — and leave exactly one block per file.

## Why

This change is plan-less: it fixes a defect found in the field, not a plan item. During the Tank Royale extraction (its CH-001, PR robocode-dev/tank-royale#218), 37 of 41 MADR ADRs began with a UTF-8 BOM. The extraction did not recognize their existing frontmatter behind the BOM, prepended the Cliewen block instead of replacing it, and `clue validate` passed the result: the leftover source frontmatter — acceptance dates included — sat unnoticed in the body. The judge must catch what the extractor got wrong, and the extractor's contract must say what right is.

## Decision boundary

Both lints are cheap to reverse — a decision-log row, not an ADR. No capability meaning changes; CAP-002 gains two criteria (AC-034, AC-035) for the new checks. The `clue-extract` contract gains one item; no mapping file changes.
