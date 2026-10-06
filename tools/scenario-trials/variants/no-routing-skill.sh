# Remove every carrier of the routing obligation from this throwaway fixture: the hub's routing
# paragraphs and the clue-delta skill that carries the routing reference. The adopter skills other than
# clue-delta are left as they are.
# A variant that finds nothing to remove fails the run, so scaffold drift cannot make it a silent no-op.
grep -q '^Before editing, inspect the smallest relevant context' AGENTS.md || { echo "variant removed nothing: the routing paragraph is not in AGENTS.md" >&2; exit 1; }
[ -d .agents/skills/clue-delta ] || { echo "variant removed nothing: no clue-delta skill directory" >&2; exit 1; }
awk '!(/^Before editing, inspect the smallest relevant context/ || /^- \*\*Direct\*\*/ || /^- \*\*Tracked\*\*/ || /^Paths and diff size may warn/ || /^A route does not authorize a push/)' AGENTS.md > AGENTS.md.new
mv AGENTS.md.new AGENTS.md
rm -rf .agents/skills/clue-delta .claude/skills/clue-delta
git add -A
git commit -qm "variant: no-routing-skill"
