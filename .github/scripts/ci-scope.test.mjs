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
    { full: false, guide: true },
  );
});

test("Sanity: a mixed guide and corpus diff fails closed", () => {
  assert.deepEqual(
    classifyChangedFiles([
      "guide/methodology.md",
      "docs/decisions/PDR-011-plain-changes-bypass-cliewen.md",
    ]),
    { full: true, guide: true },
  );
});

test("Sanity: guide configuration is not editorial prose", () => {
  assert.deepEqual(
    classifyChangedFiles(["guide/.vitepress/config.mts"]),
    { full: true, guide: true },
  );
});

test("Sanity: code and empty input fail closed", () => {
  assert.deepEqual(classifyChangedFiles(["cmd/clue/main.go"]), {
    full: true,
    guide: false,
  });
  assert.deepEqual(classifyChangedFiles([]), {
    full: true,
    guide: false,
  });
});

test("Unit: changed paths are read as NUL-delimited data", () => {
  assert.deepEqual(
    parseChangedFiles(
      Buffer.from("guide/file with spaces.md\0guide/other.md\0"),
    ),
    ["guide/file with spaces.md", "guide/other.md"],
  );
});
