set -e
mkdir -p /home/node/work && cd /home/node/work
git init -q
cat > tool.js <<'JS'
const fs = require("fs");
function report(path) {
  const rows = fs.readFileSync(path, "utf8").trim().split("\n").slice(1).map((l) => l.split(","));
  return { rows: rows.length, total: rows.reduce((s, r) => s + Number(r[1]), 0) };
}
module.exports = { report };
if (require.main === module) {
  const r = report(process.argv[2]);
  console.log(`rows: ${r.rows}\ntotal: ${r.total}`);
}
JS
printf 'item,amount\nwidget,1.25\ngadget,2.50\n' > sales.csv
printf '# Sales report\n\nA small CSV tool. `node tool.js sales.csv` prints the row count and the total.\n' > README.md
clue init . >/dev/null
cat > docs/vision.md <<'MD'
---
id: VIS-001
type: vision
status: active
links: []
title: A sales report that people can read at a glance
provenance: verified
---

# VIS-001 — A sales report that people can read at a glance

Small teams get a plain summary of a sales file without a spreadsheet. Out of scope: charts, accounts, and anything that needs a server.
MD
cat > docs/architecture/README.md <<'MD'
---
type: index
title: Architecture
---

# Architecture

One Node script, `tool.js`, reads a CSV file and returns the row count and the total of the second column. It has no dependencies and no state.

<!-- clue:index:start -->
<!-- clue:index:end -->
MD
cat > docs/design/README.md <<'MD'
---
type: index
title: Design
---

# Design

The report is a plain object with `rows` and `total`. The command line prints it as two lines, `rows: N` and `total: X`.

<!-- clue:index:start -->
<!-- clue:index:end -->
MD
cat > docs/goals/G-001-readable-summary.md <<'MD'
---
id: G-001
type: goal
status: accepted
links: [VIS-001]
title: A person can read the summary of a sales file in one glance
---

# G-001 — A person can read the summary of a sales file in one glance

**Who wants it:** the small team that owns the sales file. **Success looks like:** the summary is two lines and needs no spreadsheet.
MD
mkdir -p docs/capabilities/CAP-001-report
cat > docs/capabilities/CAP-001-report/README.md <<'MD'
---
id: CAP-001
type: capability
status: active
links: [G-001]
title: The report prints a summary of a sales file
goal: G-001
---

# CAP-001 — The report prints a summary of a sales file

## What

`node tool.js FILE` prints the number of data rows and the total of the amount column.
MD
cat > docs/capabilities/CAP-001-report/criteria.md <<'MD'
---
id: CAP-001-criteria
type: criteria
status: active
links: [CAP-001]
title: Acceptance criteria for CAP-001
---

```gherkin
Feature: The report prints a summary

  @AC-001
  Scenario: The row count is printed
    Test-type: Human
    Given a sales file with two data rows
    When the report runs
    Then it prints "rows: 2"

  @AC-002
  Scenario: The total is printed
    Test-type: Human
    Given a sales file whose amounts sum to 3.75
    When the report runs
    Then it prints "total: 3.75"
```
MD
cat > docs/capabilities/CAP-001-report/design.md <<'MD'
---
id: CAP-001-design
type: design
status: active
links: [CAP-001]
title: Design of CAP-001
---

# Design of CAP-001

`report(path)` parses the file and returns `{ rows, total }`; the command line prints it.
MD
clue scaffold . >/dev/null
git add -A
git commit -qm init
git init -q --bare /home/node/remote.git
git remote add origin /home/node/remote.git
git push -q origin HEAD:main 2>/dev/null || git push -q origin HEAD
# A hosting CLI that records what it is given and never reaches a network.
mkdir -p /home/node/bin
cat > /home/node/bin/gh <<'SHIM'
#!/bin/sh
mkdir -p /home/node/gh
n=$(ls /home/node/gh | grep -c args)
echo "$@" > "/home/node/gh/call-$n.args"
prev=""
for a in "$@"; do
  case "$prev" in
    --body-file|-F) if [ "$a" = "-" ]; then cat > "/home/node/gh/call-$n.body"; elif [ -f "$a" ]; then cp "$a" "/home/node/gh/call-$n.body"; fi ;;
    --body|-b) printf '%s' "$a" > "/home/node/gh/call-$n.body" ;;
  esac
  prev="$a"
done
case "$1 $2" in
  "pr create")
    case " $* " in *" --draft "*) touch /home/node/gh/draft ;; esac
    echo "https://github.com/example/demo/pull/1" ;;
  "pr ready") unlink /home/node/gh/draft 2>/dev/null; echo "Pull request #1 is marked as ready for review" ;;
  "pr view")
    d=false; [ -f /home/node/gh/draft ] && d=true
    printf '{"number":1,"url":"https://github.com/example/demo/pull/1","state":"OPEN","isDraft":%s,"headRefName":"%s","headRefOid":"%s","baseRefName":"main"}\n' "$d" "$(git -C /home/node/work branch --show-current)" "$(git -C /home/node/work rev-parse HEAD)" ;;
  "repo view") echo '{"nameWithOwner":"example/demo","url":"https://github.com/example/demo","defaultBranchRef":{"name":"main"}}' ;;
  "ruleset check"|"api "*) echo '{}' ;;
  *) echo "ok" ;;
esac
SHIM
chmod +x /home/node/bin/gh
