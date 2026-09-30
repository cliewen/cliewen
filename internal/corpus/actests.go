package corpus

import (
	"github.com/cliewen/cliewen/internal/evidence"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Criteria declarations retain the classified/legacy/Human/draft contract.
// All executable references come from the common repository-owned export
// (ADR-071). The judge checks form and complete input freshness, never
// framework syntax, producer execution or the meaning of assertions.
var (
	acPrefixRe       = regexp.MustCompile(`^[A-Z][A-Z0-9]*(-[A-Z][A-Z0-9]*)*$`)
	acIDRe           = regexp.MustCompile(`^([A-Z][A-Z0-9]*(-[A-Z][A-Z0-9]*)*)-([0-9]+)([a-z]*)$`)
	acCandidateTagRe = regexp.MustCompile(`@([A-Za-z][A-Za-z0-9_-]*[-_][A-Za-z0-9_-]*[0-9][A-Za-z0-9_-]*)`)
	testTypeRe       = regexp.MustCompile(`^\s*Test-type:\s*(Unit|Integration|E2E|Performance|Human)(\s+\(single-direction\))?\s*$`)
)

type acDecl struct {
	path, status string
	retired      bool
	draft        bool // @draft (ADR-033): exempt from the active-file test requirement without retiring the AC
	testType     string
	single       bool
	invalidType  bool
	humanSingle  bool // Test-type: Human combined with (single-direction), which the Human class does not use (ADR-033)
}

// CriterionIdentity is a canonical acceptance-criterion declaration that the
// identity ledger must retain. Retired criteria stay declared as tombstones,
// so their ledger state is Retired rather than disappearing with their file.
type CriterionIdentity struct {
	ID      string
	Path    string
	Live    bool
	Retired bool
}

// Declaration is one criterion's evidence-contract classification, exported
// so a consumer outside the validator (migration parity, ADR-049) can derive
// a target manifest from the same declarations checkACTests enforces,
// without inventing a second reading of a criteria file's tag lines.
type Declaration struct {
	ID       string
	Path     string
	Status   string
	Retired  bool
	Draft    bool // @draft (ADR-033)
	TestType string
	Single   bool
}

// EvidenceRef is one source-qualified executable occurrence from a producer.
// Classification belongs to that executable, independent of its framework.
type EvidenceRef struct {
	Producer  string
	Path      string
	Subject   string
	Type      string
	Direction string
}

// AcceptanceEvidence returns every declared criterion's classification
// alongside each recorded evidence occurrence, in declaration and then
// evidence order. It shares harvestACs' single tree walk with checkACTests
// and Coverage so migration parity's target manifest reads the same
// declarations and evidence the validator already enforces.
func AcceptanceEvidence(c *Corpus) (map[string]Declaration, map[string][]EvidenceRef, []Issue) {
	declared, _, _, issues, locations := harvestACs(c)
	decls := make(map[string]Declaration, len(declared))
	for id, d := range declared {
		decls[id] = Declaration{ID: id, Path: d.path, Status: d.status, Retired: d.retired, Draft: d.draft, TestType: d.testType, Single: d.single}
	}
	return decls, locations, issues
}

// LedgerCriterionIdentities returns each live canonical declaration and every
// retained tombstone, ordered by ID. It shares harvestACs' namespace and
// tombstone reading with the evidence contract, so the ledger does not invent
// a second interpretation of an acceptance-criterion identity. A tombstone
// remains a retired identity even if its criteria artifact is no longer active.
func LedgerCriterionIdentities(c *Corpus) []CriterionIdentity {
	declared, _, _, _, _ := harvestACs(c)
	ids := make([]string, 0, len(declared))
	for id, d := range declared {
		if d.status == "active" || d.retired {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	out := make([]CriterionIdentity, 0, len(ids))
	for _, id := range ids {
		d := declared[id]
		out = append(out, CriterionIdentity{ID: id, Path: d.path, Live: d.status == "active" && !d.retired, Retired: d.retired})
	}
	return out
}

// harvestACs parses criterion declarations and imports validated evidence, shared by
// checkACTests (which enforces it) and Coverage (which derives a report
// from the same declarations without repeating the walk).
func harvestACs(c *Corpus) (declared map[string]acDecl, classified map[string]map[string]map[string]bool, tested map[string]bool, issues []Issue, locations map[string][]EvidenceRef) {
	declared = map[string]acDecl{}
	// The default namespace is always known, so an undeclared AC-xxx
	// reference fails as unknown, not as purpose-less.
	prefixes := map[string]bool{"AC": true}
	normalizedPrefixes := map[string]string{normalizeACPrefix("AC"): "AC"}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" {
			continue
		}
		prefix := "AC"
		if v, present := a.Fields["ac-prefix"]; present {
			s, _ := v.(string)
			if !acPrefixRe.MatchString(s) {
				issues = append(issues, Issue{a.Path, "ac-prefix must be uppercase alphanumeric segments joined by single hyphens, starting with a letter (ADR-037)"})
				continue
			}
			prefix = s
		}
		if previous, exists := normalizedPrefixes[normalizeACPrefix(prefix)]; exists && previous != prefix {
			issues = append(issues, Issue{a.Path, "ac-prefix " + prefix + " collides with " + previous + " after carrier normalization (ADR-037)"})
		} else {
			normalizedPrefixes[normalizeACPrefix(prefix)] = prefix
		}
		prefixes[prefix] = true
		// Tag lines are read per line: `@AC-012 @retired` on one line is
		// the tombstone form (ADR-007).
		lines := strings.Split(a.Body, "\n")
		for i, line := range lines {
			retired := strings.Contains(line, "@retired")
			draft := strings.Contains(line, "@draft")
			for _, m := range acCandidateTagRe.FindAllStringSubmatch(line, -1) {
				ac := m[1]
				acPrefix, _, _, valid := parseACID(ac)
				if !valid {
					// Only a near-miss inside this file's own namespace is a
					// malformed declaration; any other @token is prose.
					if inNamespace(ac, prefix) {
						issues = append(issues, Issue{a.Path, "tag @" + ac + " is not a canonical acceptance-criterion ID (ADR-037: use <PREFIX>-<digits><lowercase-suffix>)"})
					}
					continue
				}
				if acPrefix != prefix {
					issues = append(issues, Issue{a.Path, "tag @" + ac + " is outside this file's namespace " + prefix + " (ADR-009: fix the tag or move the AC to its capability)"})
					continue
				}
				if prev, dup := declared[ac]; dup {
					issues = append(issues, Issue{a.Path, "duplicate declaration of " + ac + " (already declared in " + prev.path + ")"})
					continue
				}
				d := acDecl{path: a.Path, status: a.Status, retired: retired, draft: draft}
				d.testType, d.single, d.invalidType = scenarioTestType(lines[i+1:])
				if d.testType == "Human" && d.single {
					d.testType, d.single, d.humanSingle = "", false, true
				}
				declared[ac] = d
			}
		}
	}

	tested = map[string]bool{}
	classified = map[string]map[string]map[string]bool{}
	locations = map[string][]EvidenceRef{}
	record := func(path, subject, ac, typ, direction string) {
		issues = append(issues, checkACRef(path, subject, ac, declared, tested)...)
		locations[ac] = append(locations[ac], EvidenceRef{Path: path, Subject: subject, Type: typ, Direction: direction})
		if typ == "" || direction == "" {
			return
		}
		if classified[ac] == nil {
			classified[ac] = map[string]map[string]bool{}
		}
		if classified[ac][typ] == nil {
			classified[ac][typ] = map[string]bool{}
		}
		classified[ac][typ][direction] = true
	}
	needed := false
	for _, d := range declared {
		if d.status == "active" && !d.retired && !d.draft && d.testType != "Human" && !d.humanSingle {
			needed = true
		}
	}
	if _, err := os.Lstat(filepath.Join(c.Root, filepath.FromSlash(evidence.DefaultPath))); err == nil || needed {
		m, err := evidence.Load(c.Root)
		if err != nil {
			issues = append(issues, Issue{evidence.DefaultPath, err.Error()})
		} else {
			for _, producer := range m.Producers {
				for _, ref := range producer.References {
					record(ref.Path, ref.Subject, ref.ID, ref.Type, ref.Direction)
					refs := locations[ref.ID]
					refs[len(refs)-1].Producer = producer.ID
				}
			}
		}
	}

	return declared, classified, tested, issues, locations
}

func parseACID(id string) (prefix, number, suffix string, ok bool) {
	match := acIDRe.FindStringSubmatch(id)
	if match == nil {
		return "", "", "", false
	}
	return match[1], match[3], match[4], true
}

// inNamespace reports whether a token that failed the canonical parse is still
// close enough to a declared namespace to be a malformed criterion reference
// rather than ordinary prose or runner metadata. Carrier underscores and case
// are folded in, because those are exactly the near-misses worth diagnosing:
// `SNAP_SQS_001` and `snap-sqs-001` are wrong spellings of a declared identity,
// while `java-17` belongs to no namespace at all and stays untouched (ADR-036:
// tags outside every declared AC namespace remain ordinary runner metadata).
func inNamespace(raw, prefix string) bool {
	candidate := strings.ToUpper(strings.ReplaceAll(raw, "_", "-"))
	return strings.HasPrefix(candidate, strings.ToUpper(prefix)+"-")
}

func inAnyDeclaredNamespace(raw string, prefixes map[string]bool) bool {
	for prefix := range prefixes {
		if inNamespace(raw, prefix) {
			return true
		}
	}
	return false
}

func normalizeACPrefix(prefix string) string {
	return strings.ReplaceAll(prefix, "-", "")
}

func canonicalACIDsInLine(line string) []string {
	var ids []string
	for _, match := range acCandidateTagRe.FindAllStringSubmatch(line, -1) {
		if _, _, _, ok := parseACID(match[1]); ok {
			ids = append(ids, match[1])
		}
	}
	return ids
}

func checkACTests(c *Corpus) []Issue {
	declared, classified, tested, issues, _ := harvestACs(c)

	for ac, d := range declared {
		// A Human-class criterion is satisfied by the acceptance brief, never
		// a code test; @draft exempts a not-yet-proven criterion from the
		// active-file test requirement without retiring it (ADR-033).
		exempt := d.draft || d.testType == "Human" || d.humanSingle
		live := d.status == "active" && !d.retired
		if live && d.humanSingle {
			issues = append(issues, Issue{d.path, ac + " declares Test-type: Human (single-direction), which the Human class does not use (ADR-033)"})
		}
		if live && !exempt && !tested[ac] {
			issues = append(issues, Issue{d.path, ac + " has no test (ADR-071: export an executable reference in .clue/evidence.yaml)"})
		}
		if d.status != "active" || d.retired || exempt || d.testType == "" {
			if live && !exempt && d.invalidType {
				issues = append(issues, Issue{d.path, ac + " has a Test-type that is not the first non-blank scenario-body line (ADR-032)"})
			}
			continue
		}
		coverage := classified[ac][d.testType]
		if d.single {
			if len(coverage) == 0 {
				issues = append(issues, Issue{d.path, ac + " has no " + d.testType + " evidence (ADR-032)"})
			}
			continue
		}
		for _, direction := range []string{"positive", "negative"} {
			if !coverage[direction] {
				issues = append(issues, Issue{d.path, ac + " has no " + d.testType + " " + direction + " evidence (ADR-032)"})
			}
		}
	}
	return issues
}

// scenarioTestType reads the Test-type declaration from the first non-blank
// scenario-body line. A Test-type elsewhere in the same scenario is malformed
// rather than an opt-in declaration (ADR-032).
func scenarioTestType(lines []string) (testType string, single, invalid bool) {
	inScenario := false
	firstBodyLine := true
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inScenario {
			if trimmed == "" {
				return "", false, false
			}
			if strings.HasPrefix(trimmed, "@") {
				continue
			}
			if !isScenarioHeader(trimmed) {
				return "", false, false
			}
			inScenario = true
			continue
		}
		if isScenarioHeader(trimmed) {
			break
		}
		if trimmed == "" {
			continue
		}
		if firstBodyLine {
			firstBodyLine = false
			if m := testTypeRe.FindStringSubmatch(line); m != nil {
				testType, single = m[1], m[2] != ""
				continue
			}
			if strings.HasPrefix(trimmed, "Test-type:") {
				return "", false, true
			}
			continue
		}
		if strings.HasPrefix(trimmed, "Test-type:") {
			return "", false, true
		}
	}
	return testType, single, false
}

func isScenarioHeader(line string) bool {
	return strings.HasPrefix(line, "Scenario:") || strings.HasPrefix(line, "Scenario Outline:") || strings.HasPrefix(line, "Scenario Template:")
}

// checkACRef records an AC reference from a test and reports it when it
// resolves to nothing or to a tombstone.
func checkACRef(path, subject, ac string, declared map[string]acDecl, tested map[string]bool) []Issue {
	tested[ac] = true
	d, ok := declared[ac]
	if !ok {
		return []Issue{{path, subject + " references " + ac + " which no criteria.md declares"}}
	}
	if d.retired {
		return []Issue{{path, subject + " references retired " + ac + " — remove the test or re-tag it (ADR-007)"}}
	}
	return nil
}
