package corpus

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The AC↔acceptance-evidence contract and test-purpose taxonomy:
// every unannotated legacy AC in active criteria has one supported reference;
// a machine proof type additionally requires references classified by that
// type and positive/negative direction, unless it declares single-direction;
// Human needs no code reference because the acceptance brief is its proof;
// and @draft exempts only that not-yet-proven criterion. Every code reference
// names a live AC and every test declares exactly one purpose (ADR-005,
// ADR-006, ADR-032, ADR-033). AC IDs are namespaced by criteria file via
// optional `ac-prefix`, default `AC` (ADR-009). Go carries purpose, proof type,
// and direction in the function name; JVM test files use JUnit tags harvested
// at file level; Cucumber feature tags are harvested at scenario level.
// The judge validates those references but does not execute any test runner.
var (
	acTagRe        = regexp.MustCompile(`@([A-Z][A-Z0-9]*)-(\d+)\b`)
	acPrefixRe     = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)
	testFuncRe     = regexp.MustCompile(`(?m)^func (Test\w*)\s*\(`)
	fixedPurposeRe = regexp.MustCompile(`^Test(Unit|Sanity|Arch)(_\w*)?$`)
	acNameRe       = regexp.MustCompile(`^Test([A-Z][A-Z0-9]*?)(\d+)(_\w*)?$`)
	classifiedGoRe = regexp.MustCompile(`^Test([A-Z][A-Z0-9]*?)(\d+)_(Unit|Integration|E2E|Performance)(Positive|Negative)(_\w*)?$`)
	jvmTagRe       = regexp.MustCompile(`@Tag\(\s*"([^"]+)"\s*\)`)
	jvmACRe        = regexp.MustCompile(`^([A-Z][A-Z0-9]*)-(\d+)$`)
	testTypeRe     = regexp.MustCompile(`^\s*Test-type:\s*(Unit|Integration|E2E|Performance|Human)(\s+\(single-direction\))?\s*$`)
	featureTagRe   = regexp.MustCompile(`@([A-Za-z][A-Za-z0-9-]*)`)
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

// harvestACs parses every criteria.md tag-line declaration and walks the
// tree for classified Go, JVM, and Cucumber evidence, shared by
// checkACTests (which enforces it) and Coverage (which derives a report
// from the same declarations without repeating the walk).
func harvestACs(c *Corpus) (declared map[string]acDecl, classified map[string]map[string]map[string]bool, tested map[string]bool, issues []Issue) {
	declared = map[string]acDecl{}
	// The default namespace is always known, so an undeclared AC-xxx
	// reference fails as unknown, not as purpose-less.
	prefixes := map[string]bool{"AC": true}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" {
			continue
		}
		prefix := "AC"
		if v, present := a.Fields["ac-prefix"]; present {
			s, _ := v.(string)
			if !acPrefixRe.MatchString(s) {
				issues = append(issues, Issue{a.Path, "ac-prefix must be uppercase letters/digits starting with a letter (ADR-009)"})
				continue
			}
			prefix = s
		}
		prefixes[prefix] = true
		// Tag lines are read per line: `@AC-012 @retired` on one line is
		// the tombstone form (ADR-007).
		lines := strings.Split(a.Body, "\n")
		for i, line := range lines {
			retired := strings.Contains(line, "@retired")
			draft := strings.Contains(line, "@draft")
			for _, m := range acTagRe.FindAllStringSubmatch(line, -1) {
				ac := m[1] + "-" + m[2]
				if m[1] != prefix {
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
	record := func(path, subject, ac, typ, direction string) {
		issues = append(issues, checkACRef(path, subject, ac, declared, tested)...)
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
	_ = filepath.WalkDir(c.Root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			// Tests live in the code tree: skip the corpus, the transient
			// workspace and hidden directories.
			if rel, _ := filepath.Rel(c.Root, p); rel != "." &&
				(strings.HasPrefix(d.Name(), ".") || d.Name() == "docs" || d.Name() == "changes") {
				return fs.SkipDir
			}
			return nil
		}
		isGo := strings.HasSuffix(d.Name(), "_test.go")
		isJVM := strings.HasSuffix(d.Name(), "Test.kt") || strings.HasSuffix(d.Name(), "Test.java") ||
			strings.HasSuffix(d.Name(), "Tests.kt") || strings.HasSuffix(d.Name(), "Tests.java")
		isFeature := strings.HasSuffix(d.Name(), ".feature")
		if !isGo && !isJVM && !isFeature {
			return nil
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return nil
		}
		text := string(data)
		rel, _ := filepath.Rel(c.Root, p)
		relSlash := filepath.ToSlash(rel)

		if isGo {
			for _, m := range testFuncRe.FindAllStringSubmatch(text, -1) {
				name := m[1]
				if name == "TestMain" {
					continue // the harness hook, not a test
				}
				if fixedPurposeRe.MatchString(name) {
					continue // Unit/Sanity/Arch need no AC
				}
				if cm := classifiedGoRe.FindStringSubmatch(name); cm != nil && prefixes[cm[1]] {
					record(relSlash, "test "+name, cm[1]+"-"+cm[2], cm[3], strings.ToLower(cm[4]))
					continue
				}
				am := acNameRe.FindStringSubmatch(name)
				if am == nil || !prefixes[am[1]] {
					issues = append(issues, Issue{relSlash, "test " + name + " declares no purpose (ADR-006: prefix <ACPREFIX><digits>, Unit, Sanity or Arch)"})
					continue
				}
				record(relSlash, "test "+name, am[1]+"-"+am[2], "", "")
			}
			return nil
		}

		if isFeature {
			featureLines := strings.Split(text, "\n")
			for i := 0; i < len(featureLines); {
				if !strings.HasPrefix(strings.TrimSpace(featureLines[i]), "@") {
					i++
					continue
				}
				var tags [][]string
				for i < len(featureLines) && strings.HasPrefix(strings.TrimSpace(featureLines[i]), "@") {
					tags = append(tags, featureTagRe.FindAllStringSubmatch(featureLines[i], -1)...)
					i++
				}
				if i == len(featureLines) || !isScenarioHeader(strings.TrimSpace(featureLines[i])) {
					continue
				}
				types, directions := map[string]bool{}, map[string]bool{}
				for _, tag := range tags {
					switch strings.ToLower(tag[1]) {
					case "unit":
						types["Unit"] = true
					case "integration":
						types["Integration"] = true
					case "e2e":
						types["E2E"] = true
					case "performance":
						types["Performance"] = true
					case "positive", "negative":
						directions[strings.ToLower(tag[1])] = true
					}
				}
				for _, tag := range tags {
					am := jvmACRe.FindStringSubmatch(tag[1])
					if am == nil || !prefixes[am[1]] {
						continue
					}
					if len(types) > 1 || len(directions) > 1 {
						issues = append(issues, Issue{relSlash, "Cucumber tag block for " + tag[1] + " must declare at most one test type and direction (ADR-032)"})
						record(relSlash, "tag "+tag[0], tag[1], "", "")
						continue
					}
					for typ := range types {
						for direction := range directions {
							record(relSlash, "tag "+tag[0], tag[1], typ, direction)
						}
					}
					if len(types) == 0 || len(directions) == 0 {
						record(relSlash, "tag "+tag[0], tag[1], "", "")
					}
				}
			}
			return nil
		}

		// JVM files: harvest @Tag values at file level (ADR-009). Tags in
		// a known AC namespace are references; everything else is runner
		// metadata clue ignores (ADR-006).
		jvmTags := jvmTagRe.FindAllStringSubmatch(text, -1)
		jvmTypes, jvmDirections := []string{}, []string{}
		for _, other := range jvmTags {
			switch strings.ToLower(other[1]) {
			case "unit":
				jvmTypes = append(jvmTypes, "Unit")
			case "integration":
				jvmTypes = append(jvmTypes, "Integration")
			case "e2e":
				jvmTypes = append(jvmTypes, "E2E")
			case "performance":
				jvmTypes = append(jvmTypes, "Performance")
			case "positive", "negative":
				jvmDirections = append(jvmDirections, strings.ToLower(other[1]))
			}
		}
		for _, m := range jvmTags {
			norm := strings.ReplaceAll(m[1], "_", "-")
			am := jvmACRe.FindStringSubmatch(norm)
			if am == nil || !prefixes[am[1]] {
				continue
			}
			if len(jvmTypes) == 0 || len(jvmDirections) == 0 {
				record(relSlash, `tag "`+m[1]+`"`, norm, "", "")
				continue
			}
			for _, typ := range jvmTypes {
				for _, direction := range jvmDirections {
					record(relSlash, `tag "`+m[1]+`"`, norm, typ, direction)
				}
			}
		}
		return nil
	})

	return declared, classified, tested, issues
}

func checkACTests(c *Corpus) []Issue {
	declared, classified, tested, issues := harvestACs(c)

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
			issues = append(issues, Issue{d.path, ac + " has no test (convention per ADR-005/ADR-009: a Go test named Test" + strings.ReplaceAll(ac, "-", "") + "_… or a framework tag \"" + strings.ReplaceAll(ac, "-", "_") + "\")"})
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
