import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..", "..");

function read(relativePath) {
  return fs.readFileSync(path.join(root, relativePath), "utf8");
}

test("Sanity: this repository recovers only unpublished clue releases", () => {
  const tagWorkflow = read(".github/workflows/tag-on-merge.yml");
  const releaseWorkflow = read(".github/workflows/release.yml");
  const gates = read(".github/scripts/release-gates.sh");
  const ci = read(".github/workflows/ci.yml");

  for (const required of [
    "internal/skills/source/shared/frontmatter.md.tmpl",
    "gh release view",
    "retag=true",
    "git push --force origin",
    "gh workflow run release.yml",
    "gh run watch",
    "group: tag-on-main",
    ".github/scripts/release-gates.sh",
  ]) {
    assert.ok(tagWorkflow.includes(required), `tag workflow lacks ${required}`);
  }
  assert.ok(ci.includes(".github/scripts/release-gates.sh"));
  assert.ok(gates.includes("CHANGELOG.md"));
  assert.ok(gates.includes("internal/migrate"));
  for (const required of [
    "Require the tag to name this run's commit",
    "Require no GitHub Release yet",
    "Require the tag not to have moved while building",
  ]) {
    assert.ok(releaseWorkflow.includes(required), `release workflow lacks ${required}`);
  }
});

test("Sanity: adopter output contains no release-process policy", () => {
  for (const relativePath of [
    "internal/skills/source/shared/change-tiers.md.tmpl",
    "internal/scaffold/templates/AGENTS.md",
    "internal/scaffold/templates/github/pull_request_template.md",
    "guide/operations.md",
  ]) {
    const content = read(relativePath);
    assert.ok(!content.includes("release route"), `${relativePath} defines a release route`);
    assert.ok(!content.includes("| Releases |"), `${relativePath} documents release policy`);
  }
});

test("Sanity: the clue CLI has no release command", () => {
  const main = read("cmd/clue/main.go");
  assert.ok(!main.includes('case "release":'));
});
