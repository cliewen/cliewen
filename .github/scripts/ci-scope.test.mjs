import assert from "node:assert/strict";
import test from "node:test";

import {
  classifyChangedFiles,
  parseChangedFiles,
} from "./ci-scope.mjs";

test("Sanity: a guide-Markdown-only diff uses focused checks", () => {
  assert.deepEqual(
    classifyChangedFiles([
      "guide/what-is-cliewen.md",
      "guide/reference/nested.md",
    ]),
    { full: false, guide: true, release: false },
  );
});

test("Sanity: a mixed guide and corpus diff fails closed", () => {
  assert.deepEqual(
    classifyChangedFiles([
      "guide/methodology.md",
      "docs/decisions/PDR-011-plain-changes-bypass-cliewen.md",
    ]),
    { full: true, guide: true, release: false },
  );
});

test("Sanity: guide configuration is not editorial prose", () => {
  assert.deepEqual(
    classifyChangedFiles(["guide/.vitepress/config.mts"]),
    { full: true, guide: true, release: false },
  );
});

test("Sanity: code and empty input fail closed", () => {
  assert.deepEqual(classifyChangedFiles(["cmd/clue/main.go"]), {
    full: true,
    guide: false,
    release: false,
  });
  assert.deepEqual(classifyChangedFiles([]), {
    full: true,
    guide: false,
    release: false,
  });
});

test("Sanity: an exact release surface uses the short release route", () => {
  assert.deepEqual(
    classifyChangedFiles([
      "CHANGELOG.md",
      "internal/skills/source/shared/frontmatter.md.tmpl",
      "internal/migrate/migrate.go",
      ".agents/skills/clue-delta/skill.md",
      "internal/scaffold/templates/skills/clue-delta/skill.md",
    ]),
    { full: false, guide: false, release: true },
  );
});

test("Sanity: a release-looking change fails closed when it changes code", () => {
  assert.deepEqual(
    classifyChangedFiles([
      "CHANGELOG.md",
      "internal/skills/source/shared/frontmatter.md.tmpl",
      "cmd/clue/main.go",
    ]),
    { full: true, guide: false, release: false },
  );
});

test("Unit: changed paths are read as NUL-delimited data", () => {
  assert.deepEqual(
    parseChangedFiles(
      Buffer.from("guide/file with spaces.md\0guide/other.md\0"),
    ),
    ["guide/file with spaces.md", "guide/other.md"],
  );
});
