import assert from "node:assert/strict";
import fs from "node:fs";
import test from "node:test";

import {
  classifyChange,
  hasCompleteSimpleOverride,
  parseNulDelimited,
} from "./ci-scope.mjs";

test("AC-139 Unit positive: analysis checks corpus without becoming full", () => {
  assert.deepEqual(
    classifyChange([
      "docs/analysis/AN-024-example.md",
      "docs/analysis/README.md",
      ".clue/id-ledger.yaml",
    ]),
    { full: false, go: false, corpus: true, guide: false, release: false, override: false },
  );
});

test("AC-139 Unit positive: a full proposal in branch history selects full gates", () => {
  assert.deepEqual(
    classifyChange(
      ["cmd/clue/main.go", "docs/capabilities/CAP-002-validate/criteria.md"],
      ["changes/CH-200-example/proposal.md"],
    ),
    { full: true, go: true, corpus: true, guide: true, release: false, override: false },
  );
});

test("AC-139 Unit positive: a complete current-head override suppresses full bookkeeping only", () => {
  const message = `Fix the defect\n\nCliewen-Route: simple\nCliewen-Recommendation: full\nCliewen-Override: user chose simple; criterion risk accepted\n`;
  assert.deepEqual(
    classifyChange(["cmd/clue/main.go"], ["changes/CH-200-example/proposal.md"], message),
    { full: false, go: true, corpus: false, guide: false, release: false, override: true },
  );
});

test("AC-139 Unit negative: incomplete trailers never override a full proposal", () => {
  for (const message of [
    "Cliewen-Route: simple\nCliewen-Recommendation: full\n",
    "Cliewen-Route: simple\nCliewen-Override: user chose simple; risk\n",
    "Cliewen-Recommendation: full\nCliewen-Override: user chose simple; risk\n",
    "Cliewen-Route: simple\nCliewen-Recommendation: full\nCliewen-Override: \n",
  ]) {
    assert.equal(hasCompleteSimpleOverride(message), false);
    assert.equal(
      classifyChange(["cmd/clue/main.go"], ["changes/CH-200-example/proposal.md"], message).full,
      true,
    );
  }
});

test("AC-139 Unit positive: workflows read override trailers from the authored PR head", () => {
  for (const workflow of [".github/workflows/ci.yml", ".github/workflows/clue-validation.yml"]) {
    const source = fs.readFileSync(workflow, "utf8");
    assert.match(source, /HEAD_SHA: \$\{\{ github\.event\.pull_request\.head\.sha \|\| github\.sha \}\}/);
    assert.match(source, /git log -1 --format=%B "\$HEAD_SHA"/);
    assert.match(source, /git diff --name-only(?: -z)? "\$BASE_SHA" "\$GITHUB_SHA"/);
  }
});

test("Sanity: paths select relevant checks without deciding semantic route", () => {
  assert.deepEqual(classifyChange(["guide/what-is-cliewen.md"]), {
    full: false, go: false, corpus: false, guide: true, release: false, override: false,
  });
  assert.deepEqual(classifyChange(["cmd/clue/main.go"]), {
    full: false, go: true, corpus: false, guide: false, release: false, override: false,
  });
  assert.deepEqual(classifyChange([]), {
    full: false, go: false, corpus: false, guide: false, release: false, override: false,
  });
});

test("Sanity: this repository's exact release surface is a local simple specialization", () => {
  assert.deepEqual(
    classifyChange([
      "CHANGELOG.md",
      "internal/skills/source/shared/frontmatter.md.tmpl",
      "internal/migrate/migrate.go",
      ".agents/skills/clue-delta/skill.md",
      "internal/scaffold/templates/skills/clue-delta/skill.md",
    ]),
    { full: false, go: true, corpus: true, guide: true, release: true, override: false },
  );
});

test("Sanity: release-looking work with implementation is ordinary checked work", () => {
  const scope = classifyChange([
    "CHANGELOG.md",
    "internal/skills/source/shared/frontmatter.md.tmpl",
    "cmd/clue/main.go",
  ]);
  assert.equal(scope.release, false);
  assert.equal(scope.go, true);
});

test("Unit: NUL-delimited paths preserve spaces", () => {
  assert.deepEqual(
    parseNulDelimited(Buffer.from(" guide/file with spaces.md \0guide/other.md\0")),
    [" guide/file with spaces.md ", "guide/other.md"],
  );
});
