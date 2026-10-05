set -e
mkdir -p /home/node/work && cd /home/node/work
git init -q
cat > tool.js <<'JS'
const fs = require("fs");

function report(path) {
  const lines = fs.readFileSync(path, "utf8").trim().split("\n").slice(1);
  const amounts = lines.map((l) => parseFloat(l.split(",")[1]));
  return { rows: amounts.length, total: amounts.reduce((a, b) => a + b, 0) };
}

if (require.main === module) {
  const r = report(process.argv[2]);
  console.log(`rows: ${r.rows}  total: ${r.total.toFixed(2)}`);
}

module.exports = { report };
JS
printf 'item,amount\napple,1.50\npear,2.25\n' > sales.csv
printf '# Sales report\n\nRun `node tool.js sales.csv` to print a report of rows and total.\n' > README.md
clue init . >/dev/null
git add -A
git commit -qm init
