"""Bounded Java/Kotlin source example, not a universal Gatling parser.

Supports literal scenario names assigned to a variable, then registered by
setUp(variable.inject...()).assertions(...), in the same simulation source.
Unsupported AC shapes are diagnosed. Review whether the assertions prove the AC.
"""
import json
from pathlib import Path
import re
import sys

TOKEN = re.compile(r'//[^\n]*|/\*[\s\S]*?\*/|"(?:\\.|[^"\\])*"|\b[A-Za-z_$][\w$]*\b|[^\s]')
TITLE = re.compile(r'^\[([A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-\d+[a-z]*) (Unit|Integration|E2E|Performance) (Positive|Negative)\]')


def collect(root, files):
    references, diagnostics = [], []
    for source in files:
        if not source.endswith(("Simulation.java", "Simulation.kt")):
            continue
        text = (root / source).read_text(encoding="utf-8")
        tokens = [m[0] for m in TOKEN.finditer(text) if not m[0].startswith(("//", "/*"))]
        classes = re.findall(r'\bclass\s+\w+', " ".join(tokens))
        simulation = re.search(r'\bclass\s+\w+\s+(?:extends|:)\s+(?:[\w]+\s*\.\s*)*Simulation\b', " ".join(tokens))
        registered = set()
        for i, token in enumerate(tokens):
            if token != "setUp" or tokens[i + 1:i + 2] != ["("]:
                continue
            depth, j, names = 1, i + 2, set()
            while j < len(tokens) and depth:
                if tokens[j] == "(": depth += 1
                if tokens[j] == ")": depth -= 1
                if re.fullmatch(r"[A-Za-z_$][\w$]*", tokens[j]) and tokens[j + 1:j + 2] == ["."] and tokens[j + 2:j + 3] and tokens[j + 2].startswith("inject"):
                    names.add(tokens[j])
                j += 1
            if tokens[j:j + 3] == [".", "assertions", "("]:
                if tokens[j + 3:j + 4] != [")"]:
                    registered.update(names)
        for i, token in enumerate(tokens):
            if token != "scenario" or tokens[i + 1:i + 2] != ["("]:
                continue
            literal = tokens[i + 2] if i + 2 < len(tokens) else ""
            if not literal.startswith('"'):
                diagnostics.append({"path": source, "message": "dynamic scenario name needs an attributable export"})
                continue
            title = json.loads(literal)
            metadata = TITLE.match(title)
            if not metadata:
                if title.startswith("["):
                    diagnostics.append({"path": source, "message": "malformed scenario evidence title"})
                continue
            variable = tokens[i - 2] if i >= 2 and tokens[i - 1] == "=" else ""
            if len(classes) != 1 or not simulation or not variable or variable not in registered or tokens[i + 3:i + 4] != [")"]:
                diagnostics.append({"path": source, "subject": title, "message": "AC scenario needs literal attribution, registered simulation and nonempty assertions"})
                continue
            references.append({"id": metadata[1], "path": source, "subject": title,
                               "type": metadata[2], "direction": metadata[3].lower()})
    return {"references": references, "diagnostics": diagnostics}


if __name__ == "__main__":
    request = json.load(sys.stdin)
    print(json.dumps(collect(Path(request["root"]), request["files"])))
