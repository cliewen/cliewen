import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

export function parseNulDelimited(input) {
  return input
    .toString("utf8")
    .split("\0")
    .filter((value) => value.length > 0);
}

export function hasCompleteSimpleOverride(message) {
  const lines = message.toString("utf8").split(/\r?\n/);
  const route = lines.some((line) => /^Cliewen-Route:\s*simple\s*$/i.test(line));
  const recommendation = lines.some((line) => /^Cliewen-Recommendation:\s*full\s*$/i.test(line));
  const risk = lines.some((line) => /^Cliewen-Override:\s*\S.+$/i.test(line));
  return route && recommendation && risk;
}

export function classifyChange(files, historyFiles = [], headMessage = "") {
  const releaseFiles = new Set([
    "CHANGELOG.md",
    "internal/migrate/migrate.go",
    "internal/skills/source/shared/frontmatter.md.tmpl",
  ]);
  const generatedSkill = /^(?:\.agents\/skills|internal\/scaffold\/templates\/skills)\/(?:clue-analysis|clue-delta|clue-extract|clue-plan|clue-upgrade|clue-verify)\/skill\.md$/;
  const release =
    files.includes("internal/skills/source/shared/frontmatter.md.tmpl") &&
    files.every((file) => releaseFiles.has(file) || generatedSkill.test(file));

  const proposal = historyFiles.some((file) => /^changes\/CH-[^/]+\/proposal\.md$/.test(file));
  const override = hasCompleteSimpleOverride(headMessage);
  const full = proposal && !override;

  const guide = full || release || files.some((file) => file.startsWith("guide/"));
  const corpus =
    full ||
    release ||
    files.some((file) =>
      /^(?:docs\/|changes\/|\.clue\/|\.agents\/|AGENTS\.md$)/.test(file),
    );
  const go =
    full ||
    release ||
    files.some((file) =>
      /^(?:cmd\/|internal\/|\.github\/|go\.mod$|go\.sum$)/.test(file),
    );

  return { full, go, corpus, guide, release, override };
}

function printGitHubOutputs(scope) {
  for (const [name, value] of Object.entries(scope)) {
    process.stdout.write(`${name}=${value}\n`);
  }
}

const invokedPath = process.argv[1] ? path.resolve(process.argv[1]) : "";
if (invokedPath === fileURLToPath(import.meta.url)) {
  const historyPath = process.argv[2];
  const messagePath = process.argv[3];
  const history = historyPath ? parseNulDelimited(fs.readFileSync(historyPath)) : [];
  const message = messagePath ? fs.readFileSync(messagePath, "utf8") : "";
  printGitHubOutputs(classifyChange(parseNulDelimited(fs.readFileSync(0)), history, message));
}
