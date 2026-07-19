import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

export function parseChangedFiles(input) {
  return input
    .toString("utf8")
    .split("\0")
    .filter((file) => file.length > 0);
}

export function classifyChangedFiles(files) {
  const focusedGuide =
    files.length > 0 &&
    files.every(
      (file) =>
        file.startsWith("guide/") &&
        file.endsWith(".md"),
    );

  return {
    full: !focusedGuide,
    guide: focusedGuide || files.some((file) => file.startsWith("guide/")),
  };
}

function printGitHubOutputs(scope) {
  process.stdout.write(`full=${scope.full}\nguide=${scope.guide}\n`);
}

const invokedPath = process.argv[1] ? path.resolve(process.argv[1]) : "";
if (invokedPath === fileURLToPath(import.meta.url)) {
  printGitHubOutputs(classifyChangedFiles(parseChangedFiles(fs.readFileSync(0))));
}
