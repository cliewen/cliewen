"""Conservative Python AST example: pytest native markers or custom decorators.

Accepts direct @pytest.mark.cliewen(ac=..., type=..., direction=...) or
@cliewen(ac=..., type=..., direction=...) with literal values. It does not run
tests or interpret comments/docstrings. Adapt discovery to the actual suite.
"""
import ast
import json
from pathlib import Path
import re
import sys


def collect(root, files):
    references, diagnostics = [], []

    def report(path, subject, message):
        diagnostics.append({"path": path, "subject": subject, "message": message})

    for path in files:
        if not path.endswith(".py") or not (Path(path).name.startswith("test_") or Path(path).stem.endswith("_test")):
            continue
        tree = ast.parse((root / path).read_text(encoding="utf-8"), filename=path)

        def walk(nodes, parents=()):
            for node in nodes:
                if not isinstance(node, (ast.ClassDef, ast.FunctionDef, ast.AsyncFunctionDef)):
                    continue
                subject = ".".join(parents + (node.name,))
                markers = []
                for decorator in node.decorator_list:
                    if not isinstance(decorator, ast.Call):
                        continue
                    name = ast.unparse(decorator.func)
                    if name in {"pytest.mark.cliewen", "cliewen"}:
                        try:
                            if decorator.args:
                                raise ValueError("use named literal metadata")
                            values = {arg.arg: ast.literal_eval(arg.value) for arg in decorator.keywords}
                            if set(values) != {"ac", "type", "direction"}:
                                raise ValueError("expected ac, type, direction")
                            markers.append(values)
                        except (ValueError, TypeError, SyntaxError):
                            report(path, subject, "unsupported or dynamic executable metadata")
                if isinstance(node, ast.ClassDef):
                    if markers:
                        report(path, subject, "container AC metadata cannot prove an executable")
                    bases = {ast.unparse(base) for base in node.bases}
                    unittest_class = bool(bases & {"unittest.TestCase", "TestCase", "unittest.IsolatedAsyncioTestCase", "IsolatedAsyncioTestCase"})
                    pytest_class = node.name.startswith("Test") and not any(isinstance(member, ast.FunctionDef) and member.name == "__init__" for member in node.body)
                    if unittest_class or pytest_class:
                        walk(node.body, parents + (node.name,))
                    elif any(isinstance(member, (ast.FunctionDef, ast.AsyncFunctionDef)) and member.decorator_list for member in node.body):
                        report(path, subject, "unsupported class discovery; adapt to the actual framework")
                    continue
                if not node.name.startswith("test" if parents else "test_"):
                    if markers:
                        report(path, subject, "AC metadata is not attached to a test executable")
                    continue
                fallback = re.fullmatch(r"test_([A-Z][A-Z0-9]*(?:_[A-Z][A-Z0-9]*)*_[0-9]+[a-z]*)_(Unit|Integration|E2E|Performance)(Positive|Negative)_.+", node.name)
                if fallback:
                    markers.append({"ac": fallback[1].replace("_", "-"), "type": fallback[2], "direction": fallback[3].lower()})
                if len(markers) > 1 and any(marker != markers[0] for marker in markers[1:]):
                    report(path, subject, "conflicting AC metadata")
                elif markers:
                    metadata = markers[0]
                    if not isinstance(metadata["ac"], str) or not re.fullmatch(r"[A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-[0-9]+[a-z]*", metadata["ac"]) or metadata["type"] not in {"Unit", "Integration", "E2E", "Performance"} or metadata["direction"] not in {"positive", "negative"}:
                        report(path, subject, "malformed AC metadata")
                    else:
                        references.append({"id": metadata["ac"], "path": path, "subject": subject,
                                           "type": metadata["type"], "direction": metadata["direction"]})
                # Nested function definitions are helpers, not discovered tests.
        walk(tree.body)
    return {"references": references, "diagnostics": diagnostics}


if __name__ == "__main__":
    request = json.load(sys.stdin)
    print(json.dumps(collect(Path(request["root"]), request["files"])))
