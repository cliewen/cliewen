# Remove the routing instruction from the agent hub in this throwaway fixture.
# The skills are left as they are, so this is a weakened form of the obligation, not its absence everywhere.
awk '!(/^Before editing, inspect the smallest relevant context/ || /^- \*\*Direct\*\*/ || /^- \*\*Tracked\*\*/ || /^Paths and diff size may warn/ || /^A route does not authorize a push/)' AGENTS.md > AGENTS.md.new
mv AGENTS.md.new AGENTS.md
git add -A
git commit -qm "variant: no-routing"
