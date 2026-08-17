package corpus

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cliewen/cliewen/internal/importedchange"
	"github.com/cliewen/cliewen/internal/ledger"
)

// Options control which gates Validate applies.
type Options struct {
	// ForbidChanges fails when /changes contains any file: the
	// digest-before-merge gate. CI uses it; local runs during a change
	// loop do not.
	ForbidChanges bool
	// Version is the running clue binary's release stamp, used to detect
	// skill drift (ADR-011). Empty or "dev" exempts the build from the
	// drift comparison; skill presence and mutual consistency are still
	// checked.
	Version string
}

// defaultLifecycle is the status vocabulary for every artifact type that has
// no semantic reason to differ — including adopter-defined types the validator
// does not recognize (ADR-025, ADR-026). Types that need a different shape are
// listed in statusVocabExceptions. The docs/README.md status table mirrors
// both; together they keep the `status` field honest (Foundation §7: every
// field must have a consumer). It has no `retired` value: retirement means
// deleting the file, not holding a status (ADR-034) — a value no committed
// file may ever carry is not a value the vocabulary should offer.
var defaultLifecycle = []string{"draft", "active"}

// statusVocabExceptions holds the types whose lifecycle is not the default,
// each for the reason named in ADR-025.
var statusVocabExceptions = map[string][]string{
	"goal":            {"proposed", "accepted"},         // proposed goals are the inbox (ADR-002)
	"plan":            {"draft", "active", "completed"}, // completed is immutable (C-008)
	"decision":        {"inferred", "verified"},         // provenance lives in status (ADR-010)
	"change":          {"open"},                         // transient workspace artifact
	"tasks":           {"open"},                         // transient workspace artifact
	"open-questions":  {"open", "resolved"},             // transient workspace artifact
	"imported-change": {"in-progress", "complete"},      // durable record, never retired (ADR-050)
}

// statusVocabFor returns the allowed status set for a type: its exception if it
// has one, otherwise the default lifecycle. Unknown types resolve to the
// default, which is what lets adopters add their own types (ADR-026).
func statusVocabFor(typ string) []string {
	if v, ok := statusVocabExceptions[typ]; ok {
		return v
	}
	return defaultLifecycle
}

var (
	indexStart  = "<!-- clue:index:start -->"
	indexEnd    = "<!-- clue:index:end -->"
	milestoneRe = regexp.MustCompile(`\bM-\d+\b`)
	mdLinkRe    = regexp.MustCompile(`\]\(([^)#\s]+)\)`)
	externalRe  = regexp.MustCompile(`^[a-z][a-z0-9+.-]*:`) // http:, https:, mailto:, ...
)

// Validate applies the graph and layout rules to a scanned corpus.
func Validate(c *Corpus, opts Options) []Issue {
	var issues []Issue
	issues = append(issues, checkCoreFields(c)...)
	issues = append(issues, checkDuplicateIDs(c)...)
	issues = append(issues, checkStatusVocab(c)...)
	issues = append(issues, checkFrontmatterHygiene(c)...)
	issues = append(issues, checkLinks(c)...)
	issues = append(issues, checkSupersedes(c)...)
	issues = append(issues, checkFolderReadmes(c)...)
	issues = append(issues, checkIndexes(c)...)
	issues = append(issues, checkExternalReferences(c)...)
	issues = append(issues, checkForeignPointers(c)...)
	issues = append(issues, checkACTests(c)...)
	issues = append(issues, checkProvenance(c)...)
	issues = append(issues, checkReality(c)...)
	issues = append(issues, checkConstraints(c)...)
	issues = append(issues, checkTypeFields(c)...)
	issues = append(issues, checkDecisionTaxonomy(c)...)
	issues = append(issues, checkProseLayout(c)...)
	issues = append(issues, checkSkippedTasks(c)...)
	issues = append(issues, checkProposalPlanItem(c)...)
	issues = append(issues, checkMilestoneStatus(c)...)
	issues = append(issues, checkSkillVersions(c, opts.Version)...)
	issues = append(issues, checkLedger(c)...)
	issues = append(issues, checkImportedChanges(c)...)
	if opts.ForbidChanges && c.HasChanges {
		issues = append(issues, Issue{"changes", "transient workspace present — digest before merge (main must never contain /changes)"})
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Msg < issues[j].Msg
	})
	return issues
}

var decisionIDRe = regexp.MustCompile(`^(ADR|PDR|IDR)-[0-9]+$`)

// checkDecisionTaxonomy enforces AC-143's subject-typed decision boundary.
// The generic `decision` artifact type keeps provenance and lifecycle uniform;
// the ID and filename carry the subject type a reader chooses from the index.
func checkDecisionTaxonomy(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type == "log" {
			issues = append(issues, Issue{a.Path, "legacy decision logs are not supported — classify future-shaping rows into ADR, PDR, or IDR records in a reviewed full change"})
			continue
		}
		if a.Type != "decision" {
			continue
		}
		base := path.Base(a.Path)
		validName := decisionIDRe.MatchString(a.ID) && path.Dir(a.Path) == "docs/decisions" && strings.HasPrefix(base, a.ID+"-") && strings.HasSuffix(base, ".md")
		if !validName {
			issues = append(issues, Issue{a.Path, "decision records must live in docs/decisions and use an ADR-<number>-, PDR-<number>-, or IDR-<number>- filename matching their ID"})
		}
	}
	return issues
}

// checkProvenance validates the provenance and reversal-cost fields and
// prevents high-cost inferred meaning from backing an active capability
// (ADR-010, ADR-035).
func checkProvenance(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		v, present := a.Fields["provenance"]
		rawCost, hasCost := a.Fields["reversal-cost"]
		if a.Type == "decision" && hasCost {
			issues = append(issues, Issue{a.Path, "decisions route cost by record type and must not carry reversal-cost (ADR-035)"})
		}
		if !present {
			// A decision already got the "must not carry reversal-cost"
			// issue above; telling it to add provenance too would name a
			// field decisions are equally forbidden to carry (ADR-010).
			if hasCost && a.Type != "decision" {
				issues = append(issues, Issue{a.Path, "reversal-cost requires a provenance field (ADR-035)"})
			}
			continue // absent means human-authored
		}
		if a.Type == "decision" {
			issues = append(issues, Issue{a.Path, "decisions carry provenance in status, not in a provenance field (ADR-010)"})
			continue
		}
		s, valid := v.(string)
		if !valid || (s != "inferred" && s != "verified") {
			issues = append(issues, Issue{a.Path, "provenance must be inferred or verified (ADR-010)"})
			continue
		}
		cost, costIsString := rawCost.(string)
		if hasCost && (!costIsString || (cost != "low" && cost != "high")) {
			issues = append(issues, Issue{a.Path, "reversal-cost must be low or high (ADR-035)"})
		}
		if s == "inferred" && !hasCost {
			issues = append(issues, Issue{a.Path, "provenance inferred requires reversal-cost low or high (ADR-035)"})
		}
	}
	for _, b := range ProvenanceBacklog(c).Blockers {
		issues = append(issues, Issue{b.Artifact.Path, "high-cost inferred artifact " + b.Artifact.ID + " blocks active capability " + b.Capability.ID + " — verify it or classify it low only when reversal is cheap (ADR-035)"})
	}
	return issues
}

// checkFrontmatterHygiene guards the file shapes that hide or duplicate
// frontmatter (AC-034, AC-035): a UTF-8 BOM anywhere in a corpus markdown
// file conceals a fence from the parser, and a complete frontmatter block
// opening an artifact body is source frontmatter an extraction failed to
// replace.
func checkFrontmatterHygiene(c *Corpus) []Issue {
	const bom = "\uFEFF"
	var issues []Issue
	for _, rel := range c.MDFiles {
		if strings.Contains(c.Contents[rel], bom) {
			issues = append(issues, Issue{rel, "file contains a UTF-8 byte-order mark — strip the byte-order mark (it hides frontmatter fences from the parser)"})
		}
	}
	for _, a := range c.Artifacts {
		// Reuse the scanner's own fence semantics: only a closed block
		// whose content parses as a non-empty YAML mapping is leftover
		// frontmatter. Prose between two thematic breaks is not.
		body := strings.TrimLeft(a.Body, "\n"+bom)
		if fields, _, ok, _ := parseFrontmatter(body); ok && len(fields) > 0 {
			issues = append(issues, Issue{a.Path, "body opens with a complete second frontmatter block — leftover source frontmatter; an artifact carries exactly one frontmatter block"})
		}
	}
	return issues
}

// checkConstraints enforces the convention register's fields (AC-089):
// every constraint names its source (the doc or catalog that states the
// rule) and an enforcement class. `agent` marks the promotion backlog —
// rules awaiting a machine check; the CLI reports their count.
//
// `partial` and `human` are the classes ADR-045 gave declarations: a rule
// no longer leaves the backlog by relabelling, only by naming the machine
// that holds its subset and pricing the judgment that stays. A `machine`
// rule owes neither, because there is no residual to state, and an `agent`
// rule keeps its promotion trigger.
func checkConstraints(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type != "constraint" {
			continue
		}
		if s, ok := a.Fields["source"].(string); !ok || s == "" {
			issues = append(issues, Issue{a.Path, "constraint missing or empty source field"})
		} else {
			issues = append(issues, checkConstraintSource(c, a, s)...)
		}
		e, _ := a.Fields["enforcement"].(string)
		switch e {
		case "machine", "partial", "agent", "human":
		case "":
			issues = append(issues, Issue{a.Path, "constraint missing or empty enforcement field"})
		default:
			issues = append(issues, Issue{a.Path, "enforcement " + e + " not allowed (allowed: machine, partial, agent, human)"})
		}
		// Every declared class prices its residual. Only `partial` must also
		// name a machine: a `human` rule may have a fragment worth naming and
		// usually does not, and demanding the heading anyway would invite
		// "**Checked by:** nothing" — a sentence that reads as a check.
		var required []string
		switch e {
		case "partial":
			required = []string{"Checked by", "Residual"}
		case "human":
			required = []string{"Residual"}
		}
		for _, d := range required {
			if !strings.Contains(a.Body, "**"+d+":**") {
				issues = append(issues, Issue{a.Path, "enforcement " + e + " needs a **" + d + ":** declaration (ADR-045)"})
			}
		}
	}
	return issues
}

// constraintSourceIDRe matches a corpus decision, constraint, architecture,
// goal, capability, analysis, plan, imported-change, or log ID inside
// free-form source: prose.
// M- (milestone) and AC- (criterion) identities are deliberately excluded:
// neither is a top-level c.ByID entry, so a milestone or criterion mentioned
// in passing (as this very check's own commentary does) is not a claim the
// source resolves to.
var constraintSourceIDRe = regexp.MustCompile(`\b(?:ADR|PDR|ARCH|CAP|AN|LOG|C|G|P|IC)-\d+\b`)

// constraintSourcePathRe matches a bare or relative markdown file path
// inside free-form source: prose, such as "AGENTS.md" or
// "docs/decisions/log.md".
var constraintSourcePathRe = regexp.MustCompile(`[\w./-]+\.md`)

// constraintSourceFragments and constraintSourceSkills name the shared
// skill-source fragments and managed skills a source: may point at by name
// rather than by file path — real files under the repository the check can
// still resolve against.
var constraintSourceFragments = []string{"review-boundary", "durable-work", "decision-records", "local-conventions", "change-tiers", "frontmatter"}

// checkConstraintSource enforces that a constraint's source: resolves to
// something live rather than merely being non-empty (M-067, log.md
// 2026-08-08): a corpus-ID-shaped token must resolve to a live artifact, and
// a token naming a markdown path, a shared skill-source fragment, or a
// managed skill must resolve to a real file. Locating prose that names
// neither — "step 2 (Propose)", a quoted sentence — asserts nothing this
// check can verify and is left to human review: deciding that a paragraph
// still says what its source claims is a judgment about English, which is
// this rule's own stated residual.
func checkConstraintSource(c *Corpus, a *Artifact, source string) []Issue {
	var issues []Issue
	seen := map[string]bool{}
	for _, id := range constraintSourceIDRe.FindAllString(source, -1) {
		if seen[id] {
			continue
		}
		seen[id] = true
		if _, ok := c.ByID[id]; !ok {
			issues = append(issues, Issue{a.Path, "constraint source names " + id + ", which does not resolve to a live artifact"})
		}
	}
	for _, p := range constraintSourcePathRe.FindAllString(source, -1) {
		if seen[p] {
			continue
		}
		seen[p] = true
		if _, err := os.Stat(filepath.Join(c.Root, filepath.FromSlash(p))); err != nil {
			issues = append(issues, Issue{a.Path, "constraint source names " + p + ", which does not resolve to a live file"})
		}
	}
	for _, frag := range constraintSourceFragments {
		if !strings.Contains(source, frag) || seen[frag] {
			continue
		}
		seen[frag] = true
		p := filepath.Join(c.Root, "internal", "skills", "source", "shared", frag+".md.tmpl")
		if _, err := os.Stat(p); err != nil {
			issues = append(issues, Issue{a.Path, "constraint source names the " + frag + " fragment, which does not resolve to a live file"})
		}
	}
	for name := range legacyCliewenSkillNames {
		if !strings.Contains(source, name) || seen[name] {
			continue
		}
		seen[name] = true
		p := filepath.Join(c.Root, "internal", "skills", "source", "skills", name+".md.tmpl")
		if _, err := os.Stat(p); err != nil {
			issues = append(issues, Issue{a.Path, "constraint source names the " + name + " skill, which does not resolve to a live file"})
		}
	}
	return issues
}

func checkCoreFields(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		var missing []string
		for _, f := range CoreFields {
			v, present := a.Fields[f]
			if !present {
				missing = append(missing, f)
				continue
			}
			// links may legitimately be an empty list; the rest must be
			// non-empty strings.
			if f != "links" {
				if s, ok := v.(string); !ok || s == "" {
					missing = append(missing, f)
				}
			}
		}
		if len(missing) > 0 {
			issues = append(issues, Issue{a.Path, "missing or empty core field(s): " + strings.Join(missing, ", ")})
		}
	}
	return issues
}

func checkDuplicateIDs(c *Corpus) []Issue {
	var issues []Issue
	for id, as := range c.ByID {
		if len(as) > 1 {
			var paths []string
			for _, a := range as {
				paths = append(paths, a.Path)
			}
			sort.Strings(paths)
			issues = append(issues, Issue{paths[0], "duplicate id " + id + " (also in " + strings.Join(paths[1:], ", ") + ")"})
		}
	}
	return issues
}

// checkLedger cross-checks the corpus against .clue/id-ledger.yaml
// (ADR-048). A corpus without a ledger file yet is unaffected — the gate
// that keeps this rule from firing before a repository has run the
// backfill migration. Where a ledger exists, it rejects a live artifact
// missing from the ledger or whose ID the ledger marks reserved or retired,
// and a malformed entry: a numeric-kind entry with no valid decimal
// component, or an opaque-kind entry carrying one.
func checkLedger(c *Corpus) []Issue {
	if !ledger.Exists(c.Root) {
		return nil
	}
	l, err := ledger.Load(c.Root)
	if err != nil {
		return []Issue{{ledger.DefaultPath, "ledger: " + err.Error()}}
	}
	var issues []Issue
	for _, e := range l.Entries() {
		switch e.State {
		case ledger.StateReserved, ledger.StateLive, ledger.StateRetired:
		default:
			issues = append(issues, Issue{ledger.DefaultPath, "entry " + e.ID + " has invalid state " + string(e.State)})
		}
		switch e.Kind {
		case ledger.KindNumeric:
			if !ledger.ValidNumericEntry(e) {
				issues = append(issues, Issue{ledger.DefaultPath, "entry " + e.ID + " is numeric-kind but its ID, prefix, and component do not agree"})
			}
		case ledger.KindOpaque:
			if e.Component != nil || e.Prefix != "" {
				issues = append(issues, Issue{ledger.DefaultPath, "entry " + e.ID + " is opaque-kind but carries numeric fields"})
			}
		default:
			issues = append(issues, Issue{ledger.DefaultPath, "entry " + e.ID + " has invalid kind " + string(e.Kind)})
		}
	}
	for id, as := range c.ByID {
		entry, ok := l.Lookup(id)
		if !ok {
			for _, a := range as {
				issues = append(issues, Issue{a.Path, "id " + id + " is missing from " + ledger.DefaultPath})
			}
			continue
		}
		if entry.State == ledger.StateLive {
			continue
		}
		for _, a := range as {
			issues = append(issues, Issue{a.Path, "id " + id + " is marked " + string(entry.State) + " in " + ledger.DefaultPath + " and cannot be used by a live artifact"})
		}
	}
	for _, criterion := range LedgerCriterionIdentities(c) {
		entry, ok := l.Lookup(criterion.ID)
		if !ok {
			issues = append(issues, Issue{criterion.Path, "criterion " + criterion.ID + " is missing from " + ledger.DefaultPath})
			continue
		}
		want := ledger.StateLive
		if criterion.Retired {
			want = ledger.StateRetired
		} else if !criterion.Live {
			continue
		}
		if entry.State != want {
			issues = append(issues, Issue{criterion.Path, "criterion " + criterion.ID + " is marked " + string(entry.State) + " in " + ledger.DefaultPath + " but its declaration is " + string(want)})
		}
	}
	return issues
}

// checkImportedChanges enforces ADR-050's imported-change record contract:
// every record names its pinned origin, and a record whose status claims
// complete may not name a proof link the corpus cannot back — a criterion
// that does not exist, is @draft, or is retired. An in-progress record is
// exempt from that last check by design: it is allowed to name work not yet
// proven, because it is still declaring what remains to be done. Dependency
// links (the ordinary links: field) resolve like any other artifact link
// and need no rule here — checkLinks already enforces that.
func checkImportedChanges(c *Corpus) []Issue {
	var issues []Issue
	var declared map[string]Declaration
	for _, a := range c.Artifacts {
		if a.Type != "imported-change" {
			continue
		}
		sourceRevision, _ := a.Fields["source-revision"].(string)
		sourceLocation, _ := a.Fields["source-location"].(string)
		if sourceRevision == "" || sourceLocation == "" {
			issues = append(issues, Issue{a.Path, "imported-change missing source-revision or source-location (ADR-050)"})
		}
		links := importedchange.ParseProofLinks(a.Body)
		if len(links) == 0 {
			issues = append(issues, Issue{a.Path, "imported-change has no proof-links table under a \"## Proof links\" heading with Task and Criterion columns (ADR-050)"})
			continue
		}
		for _, l := range links {
			if l.Task == "" || l.Criterion == "" {
				issues = append(issues, Issue{a.Path, "proof-links row missing task or criterion"})
			}
		}
		if a.Status != "complete" {
			continue
		}
		if declared == nil {
			declared, _ = AcceptanceEvidence(c)
		}
		provable := make(map[string]bool, len(declared))
		for id, d := range declared {
			provable[id] = !d.Draft && !d.Retired
		}
		for _, id := range importedchange.UnprovenLinks(links, provable) {
			issues = append(issues, Issue{a.Path, "imported-change is complete but proof-linked criterion " + id + " does not exist, is @draft, or is retired (ADR-050)"})
		}
	}
	return issues
}

func checkStatusVocab(c *Corpus) []Issue {
	var issues []Issue
	for _, a := range c.Artifacts {
		if a.Type == "" {
			continue // reported by checkCoreFields
		}
		if a.Status == "" {
			continue // reported by checkCoreFields
		}
		allowed := statusVocabFor(a.Type)
		ok := false
		for _, s := range allowed {
			if a.Status == s {
				ok = true
				break
			}
		}
		if !ok {
			issues = append(issues, Issue{a.Path, "status " + a.Status + " not allowed for type " + a.Type + " (allowed: " + strings.Join(allowed, ", ") + ")"})
		}
	}
	return issues
}

func checkLinks(c *Corpus) []Issue {
	// Milestones (M-xxx) live inside their plan file, not as separate
	// artifacts: harvest them from plan bodies.
	milestones := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.Type == "plan" {
			for _, m := range milestoneRe.FindAllString(a.Body, -1) {
				milestones[m] = true
			}
		}
	}
	supersededBy := supersededByIndex(c)
	acOwners := acceptanceCriterionOwners(c)
	var issues []Issue
	for _, a := range c.Artifacts {
		for _, l := range a.Links {
			if milestoneRe.MatchString(l) {
				if !milestones[l] {
					issues = append(issues, Issue{a.Path, "link " + l + " not found in any plan"})
				}
				continue
			}
			if _, ok := acOwners[l]; ok {
				continue
			}
			if _, ok := c.ByID[l]; !ok {
				if successor, retired := supersededBy[l]; retired {
					issues = append(issues, Issue{a.Path, "link " + l + " was retired — repoint to its successor " + successor + ", which names " + l + " in its supersedes field"})
				} else {
					issues = append(issues, Issue{a.Path, "link " + l + " resolves to no artifact"})
				}
			}
		}
	}
	return issues
}

// supersededByIndex maps each retired ID named in some artifact's
// supersedes field to the ID of the artifact that named it (ADR-034).
func supersededByIndex(c *Corpus) map[string]string {
	idx := map[string]string{}
	for _, a := range c.Artifacts {
		for _, s := range a.Supersedes {
			idx[s] = a.ID
		}
	}
	return idx
}

// checkSupersedes enforces ADR-034: retiring an artifact means deleting
// its file in the same change. A supersedes entry naming an ID that
// still resolves to a live artifact means the retirement was declared
// but not actually carried out. A retired ID claimed by more than one
// successor is an unresolved ambiguity, not a fact the validator can
// silently pick a winner for.
func checkSupersedes(c *Corpus) []Issue {
	var issues []Issue
	claimants := map[string][]*Artifact{} // superseded ID -> claiming artifacts
	for _, a := range c.Artifacts {
		claimed := map[string]bool{} // dedupe a repeated entry within one artifact's own list
		for _, s := range a.Supersedes {
			if s == a.ID {
				issues = append(issues, Issue{a.Path, "supersedes its own id " + s + " — an artifact cannot retire itself (ADR-034)"})
				continue
			}
			if _, ok := c.ByID[s]; ok {
				issues = append(issues, Issue{a.Path, "supersedes " + s + " but that artifact still exists in the corpus — delete it in this change (ADR-034)"})
			}
			if claimed[s] {
				continue
			}
			claimed[s] = true
			claimants[s] = append(claimants[s], a)
		}
	}
	for s, by := range claimants {
		if len(by) < 2 {
			continue
		}
		var ids []string
		for _, a := range by {
			ids = append(ids, a.ID)
		}
		sort.Strings(ids)
		for _, a := range by {
			issues = append(issues, Issue{a.Path, s + " is claimed as superseded by more than one artifact: " + strings.Join(ids, ", ") + " (ADR-034)"})
		}
	}
	return issues
}

func checkFolderReadmes(c *Corpus) []Issue {
	var issues []Issue
	docs := filepath.Join(c.Root, "docs")
	if info, err := os.Stat(docs); err != nil || !info.IsDir() {
		return nil
	}
	_ = filepath.WalkDir(docs, func(p string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(p, "README.md")); err != nil {
			rel, _ := filepath.Rel(c.Root, p)
			issues = append(issues, Issue{filepath.ToSlash(rel), "folder has no README.md"})
		}
		return nil
	})
	return issues
}

// isTaxonomyReadme reports whether rel is one of the READMEs carrying a
// generated index block: docs/README.md and each docs/<folder>/README.md.
// The judge and the filler-row count read the same set.
func isTaxonomyReadme(rel string) bool {
	parts := strings.Split(rel, "/")
	return rel == "docs/README.md" ||
		(len(parts) == 3 && parts[0] == "docs" && parts[2] == "README.md")
}

// checkIndexes enforces the generated-index contract on the taxonomy
// READMEs (docs/README.md and each docs/<folder>/README.md): index
// markers must exist, every link inside the block must resolve to a
// file, and every sibling artifact or artifact-bearing subfolder must
// be referenced in the block.
func checkIndexes(c *Corpus) []Issue {
	var issues []Issue
	for _, rel := range c.MDFiles {
		if !isTaxonomyReadme(rel) {
			continue
		}
		text := c.Contents[rel]
		start := strings.Index(text, indexStart)
		end := strings.Index(text, indexEnd)
		if start < 0 || end < 0 || end < start {
			issues = append(issues, Issue{rel, "index markers missing or malformed (" + indexStart + " … " + indexEnd + ")"})
			continue
		}
		block := text[start:end]
		dir := path.Dir(rel)

		targets := map[string]bool{}
		for _, m := range mdLinkRe.FindAllStringSubmatch(block, -1) {
			t := m[1]
			if externalRe.MatchString(t) {
				continue
			}
			t = path.Clean(t)
			targets[t] = true
			if _, err := os.Stat(filepath.Join(c.Root, filepath.FromSlash(path.Join(dir, t)))); err != nil {
				issues = append(issues, Issue{rel, "index references missing file " + t})
			}
		}

		entries, err := os.ReadDir(filepath.Join(c.Root, filepath.FromSlash(dir)))
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				if !dirHasMarkdown(filepath.Join(c.Root, filepath.FromSlash(dir), name)) {
					continue
				}
				covered := false
				for t := range targets {
					if t == name || strings.HasPrefix(t, name+"/") {
						covered = true
						break
					}
				}
				if !covered {
					issues = append(issues, Issue{rel, "index does not reference subfolder " + name + "/"})
				}
			} else if strings.HasSuffix(name, ".md") && name != "README.md" && !targets[name] {
				issues = append(issues, Issue{rel, "index does not reference sibling " + name})
			}
		}
	}
	return issues
}

func dirHasMarkdown(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return found
}
