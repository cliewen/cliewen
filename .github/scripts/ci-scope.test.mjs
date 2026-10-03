import assert from "node:assert/strict";
import fs from "node:fs";
import test from "node:test";

import {
  classifyChange,
  hasCompleteDirectOverride,
  parseNulDelimited,
} from "./ci-scope.mjs";

test("AC-209 Unit positive: analysis checks corpus without becoming tracked", () => {
  assert.deepEqual(
    classifyChange([
      "docs/analysis/AN-024-example.md",
      "docs/analysis/README.md",
      ".clue/id-ledger.yaml",
    ]),
    { tracked: false, go: false, corpus: true, guide: false, release: false, override: false },
  );
});

test("AC-209 Unit positive: a tracked proposal in branch history selects tracked gates", () => {
  assert.deepEqual(
    classifyChange(
      ["cmd/clue/main.go", "docs/capabilities/CAP-002-validate/criteria.md"],
      ["changes/CH-200-example/proposal.md"],
    ),
    { tracked: true, go: true, corpus: true, guide: true, release: false, override: false },
  );
});

test("AC-211 Unit positive: a complete current-head override in either spelling suppresses tracked bookkeeping only", () => {
  for (const message of [
    `Fix the defect\n\nCliewen-Route: direct\nCliewen-Recommendation: tracked\nCliewen-Override: user chose direct; criterion risk accepted\n`,
    `Fix the defect\n\nCliewen-Route: simple\nCliewen-Recommendation: full\nCliewen-Override: user chose simple; criterion risk accepted\n`,
    `Fix the defect\n\nCliewen-Route: direct\nCliewen-Recommendation: full\nCliewen-Override: user chose direct; risk\n`,
  ]) {
    assert.deepEqual(
      classifyChange(["cmd/clue/main.go"], ["changes/CH-200-example/proposal.md"], message),
      { tracked: false, go: true, corpus: false, guide: false, release: false, override: true },
    );
  }
});

test("AC-211 Unit negative: incomplete or unknown trailers never override a tracked proposal", () => {
  for (const message of [
    "Cliewen-Route: direct\nCliewen-Recommendation: tracked\n",
    "Cliewen-Route: direct\nCliewen-Override: user chose direct; risk\n",
    "Cliewen-Recommendation: tracked\nCliewen-Override: user chose direct; risk\n",
    "Cliewen-Route: direct\nCliewen-Recommendation: tracked\nCliewen-Override: \n",
    "Cliewen-Route: simple\nCliewen-Recommendation: full\n",
    "Cliewen-Route: quick\nCliewen-Recommendation: tracked\nCliewen-Override: user chose quick; risk\n",
    "Cliewen-Route: direct\nCliewen-Recommendation: planned\nCliewen-Override: user chose direct; risk\n",
  ]) {
    assert.equal(hasCompleteDirectOverride(message), false);
    assert.equal(
      classifyChange(["cmd/clue/main.go"], ["changes/CH-200-example/proposal.md"], message).tracked,
      true,
    );
  }
});

test("AC-211 Unit positive: workflows read override trailers from the authored PR head", () => {
  for (const workflow of [".github/workflows/ci.yml", ".github/workflows/clue-validation.yml"]) {
    const source = fs.readFileSync(workflow, "utf8");
    assert.match(source, /HEAD_SHA: \$\{\{ github\.event\.pull_request\.head\.sha \|\| github\.sha \}\}/);
    assert.match(source, /git log -1 --format=%B "\$HEAD_SHA"/);
    assert.match(source, /git diff --name-only(?: -z)? "\$BASE_SHA" "\$GITHUB_SHA"/);
  }
});

test("Sanity: paths select relevant checks without deciding semantic route", () => {
  assert.deepEqual(classifyChange(["guide/what-is-cliewen.md"]), {
    tracked: false, go: false, corpus: false, guide: true, release: false, override: false,
  });
  assert.deepEqual(classifyChange(["cmd/clue/main.go"]), {
    tracked: false, go: true, corpus: false, guide: false, release: false, override: false,
  });
  assert.deepEqual(classifyChange([]), {
    tracked: false, go: false, corpus: false, guide: false, release: false, override: false,
  });
});

test("Sanity: this repository's exact release surface is a local direct-route specialization", () => {
  assert.deepEqual(
    classifyChange([
      "CHANGELOG.md",
      "internal/skills/source/shared/frontmatter.md.tmpl",
      "internal/migrate/migrate.go",
      ".agents/skills/clue-delta/skill.md",
      "internal/scaffold/templates/skills/clue-delta/skill.md",
    ]),
    { tracked: false, go: true, corpus: true, guide: true, release: true, override: false },
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

test("Unit: repository-owned exporter changes select Go checks", () => {
  assert.equal(classifyChange(["tools/export-evidence/main.go"]).go, true);
});
