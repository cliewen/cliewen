# Remove every carrier of the routing obligation from this throwaway fixture: the hub's routing
# paragraphs and the clue-delta skill that carries the routing reference. The adopter skills other than
# clue-delta are left as they are.
awk '!(/^Before editing, inspect the smallest relevant context/ || /^- \*\*Direct\*\*/ || /^- \*\*Tracked\*\*/ || /^Paths and diff size may warn/ || /^A route does not authorize a push/)' AGENTS.md > AGENTS.md.new
mv AGENTS.md.new AGENTS.md
rm -rf .agents/skills/clue-delta .claude/skills/clue-delta
git add -A
git commit -qm "variant: no-routing-skill"
