// Package evidenceexport is a repository-owned compatibility example, not a
// harvester linked into clue. It exports the old Go/JVM/Cucumber conventions.
package evidenceexport

import (
	"fmt"
	"github.com/cliewen/cliewen/internal/evidence"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Issue struct{ Path, Msg string }

var (
	acIDRe           = regexp.MustCompile(`^([A-Z][A-Z0-9]*(-[A-Z][A-Z0-9]*)*)-([0-9]+)([a-z]*)$`)
	acCandidateTagRe = regexp.MustCompile(`@([A-Za-z][A-Za-z0-9_-]*[-_][A-Za-z0-9_-]*[0-9][A-Za-z0-9_-]*)`)
	fixedPurposeRe   = regexp.MustCompile(`^Test(Unit|Sanity|Arch)(_\w*)?$`)
	classifiedGoRe   = regexp.MustCompile(`^Test(.+?)_(Unit|Integration|E2E|Performance)(Positive|Negative)(_\w*)?$`)
	featureTagRe     = regexp.MustCompile(`@([A-Za-z][A-Za-z0-9_-]*)`)
	goOrdinalRe      = regexp.MustCompile(`^[0-9]+[a-z]*$`)
)

func goTestNames(text string) []string {
	f, _ := parser.ParseFile(token.NewFileSet(), "test.go", text, parser.AllErrors)
	if f == nil {
		return nil
	}
	var names []string
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") {
			names = append(names, fn.Name.Name)
		}
	}
	return names
}

func isScenarioHeader(line string) bool {
	return strings.HasPrefix(line, "Scenario:") || strings.HasPrefix(line, "Scenario Outline:") || strings.HasPrefix(line, "Scenario Template:")
}

func Prefixes(root string) map[string]bool {
	prefixes := map[string]bool{"AC": true}
	filepath.WalkDir(filepath.Join(root, "docs"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		pieces := strings.SplitN(strings.ReplaceAll(string(data), "\r\n", "\n"), "---", 3)
		if len(pieces) > 2 {
			var fm map[string]any
			if yaml.Unmarshal([]byte(pieces[1]), &fm) == nil {
				if v, ok := fm["ac-prefix"].(string); ok && fm["type"] == "criteria" {
					prefixes[v] = true
				}
			}
		}
		return nil
	})
	return prefixes
}

func Collect(root string, producerID string) (evidence.Producer, error) {
	prefixes := Prefixes(root)
	p := evidence.Producer{ID: producerID, Framework: "Go/JVM/Cucumber compatibility example", Include: []string{"**/*_test.go", "**/*Test.java", "**/*Tests.java", "**/*Test.kt", "**/*Tests.kt", "**/*.feature", "docs/**/*.md"}, Exclude: []string{".*/**", "**/.*/**", "node_modules/**", "**/node_modules/**", "vendor/**", "changes/**", "internal/scaffold/templates/**", "internal/migrate/testdata/**"}}
	pProducerExclude := p.Exclude
	var issues []Issue
	record := func(path, subject, ac, typ, direction string) {
		p.References = append(p.References, evidence.Reference{ID: ac, Path: path, Subject: subject, Type: typ, Direction: direction})
	}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Tests live in the code tree: skip the corpus, the transient
			// workspace and hidden directories.
			if rel, _ := filepath.Rel(root, p); rel != "." &&
				(strings.HasPrefix(d.Name(), ".") || d.Name() == "docs" || d.Name() == "changes" || d.Name() == "node_modules" || d.Name() == "vendor") {
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
		relCandidate, _ := filepath.Rel(root, p)
		for _, pattern := range pProducerExclude {
			if evidence.Match(pattern, filepath.ToSlash(relCandidate)) {
				return nil
			}
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		text := string(data)
		rel, _ := filepath.Rel(root, p)
		relSlash := filepath.ToSlash(rel)

		if isGo {
			if _, parseErr := parser.ParseFile(token.NewFileSet(), relSlash, text, parser.AllErrors); parseErr != nil {
				issues = append(issues, Issue{relSlash, "Go discovery failed: " + parseErr.Error()})
				return nil
			}
			for _, name := range goTestNames(text) {
				if name == "TestMain" {
					continue // the harness hook, not a test
				}
				if fixedPurposeRe.MatchString(name) {
					continue // Unit/Sanity/Arch need no AC
				}
				if ac, typ, direction, ambiguous, matched := parseGoClassifiedName(name, prefixes); matched {
					record(relSlash, "test "+name, ac, typ, direction)
					continue
				} else if ambiguous {
					issues = append(issues, Issue{relSlash, "test " + name + " has an ambiguous normalized criterion prefix (ADR-037)"})
					continue
				}
				if ac, ambiguous, matched := parseGoReference(name, prefixes); matched {
					record(relSlash, "test "+name, ac, "", "")
					continue
				} else if ambiguous {
					issues = append(issues, Issue{relSlash, "test " + name + " has an ambiguous normalized criterion prefix (ADR-037)"})
					continue
				} else {
					issues = append(issues, Issue{relSlash, "test " + name + " declares no purpose (ADR-006: normalized criterion prefix plus digits, Unit, Sanity or Arch)"})
				}
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
					ac, ok := normalizeCarrierACID(tag[1])
					if !ok {
						if inAnyDeclaredNamespace(tag[1], prefixes) {
							issues = append(issues, Issue{relSlash, "Cucumber tag " + tag[0] + " is not a supported canonical acceptance-criterion ID (ADR-037)"})
						}
						continue
					}
					acPrefix, _, _, valid := parseACID(ac)
					if !valid || !prefixes[acPrefix] {
						continue
					}
					if len(types) > 1 || len(directions) > 1 {
						issues = append(issues, Issue{relSlash, "Cucumber tag block for " + tag[1] + " must declare at most one test type and direction (ADR-032)"})
						record(relSlash, "scenario "+strings.TrimSpace(featureLines[i])+" at line "+fmt.Sprint(i+1), ac, "", "")
						continue
					}
					for typ := range types {
						for direction := range directions {
							record(relSlash, "scenario "+strings.TrimSpace(featureLines[i])+" at line "+fmt.Sprint(i+1), ac, typ, direction)
						}
					}
					if len(types) == 0 || len(directions) == 0 {
						record(relSlash, "scenario "+strings.TrimSpace(featureLines[i])+" at line "+fmt.Sprint(i+1), ac, "", "")
					}
				}
			}
			return nil
		}

		jvmEvidence, jvmIssues := harvestJVMEvidence(text, relSlash, prefixes)
		issues = append(issues, jvmIssues...)
		for _, evidence := range jvmEvidence {
			record(relSlash, evidence.subject, evidence.ac, evidence.testType, evidence.direction)
		}
		return nil
	})

	if err != nil {
		return p, err
	}
	for _, issue := range issues {
		p.Diagnostics = append(p.Diagnostics, evidence.Diagnostic{Path: issue.Path, Message: issue.Msg})
	}
	if len(issues) > 0 {
		p.References = nil
	}
	// Cucumber scenario identity includes its position; a scenario's identity
	// must not be a tag shared by distinct scenarios.
	p.Inputs, err = evidence.Snapshot(root, p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Fixture materializes an export for throwaway regression fixtures. Diagnostics
// are intentionally retained so consumer tests can assert the failure contract.
func Fixture(root string) error {
	p, err := Collect(root, "fixture")
	if err != nil {
		return err
	}
	data, err := evidence.Encode(evidence.Manifest{Version: 1, Producers: []evidence.Producer{p}})
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(root, ".clue"), 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, filepath.FromSlash(evidence.DefaultPath)), data, 0644)
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

func normalizeCarrierACID(raw string) (string, bool) {
	normalized := strings.ReplaceAll(raw, "_", "-")
	if _, _, _, ok := parseACID(normalized); !ok {
		return "", false
	}
	return normalized, true
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

func parseGoClassifiedName(name string, prefixes map[string]bool) (ac, testType, direction string, ambiguous, matched bool) {
	match := classifiedGoRe.FindStringSubmatch(name)
	if match == nil {
		return "", "", "", false, false
	}
	ac, ambiguous, matched = parseNormalizedACPart(match[1], prefixes)
	if matched {
		return ac, match[2], strings.ToLower(match[3]), false, true
	}
	return "", "", "", ambiguous, false
}

func parseGoReference(name string, prefixes map[string]bool) (ac string, ambiguous, matched bool) {
	if !strings.HasPrefix(name, "Test") {
		return "", false, false
	}
	part := strings.TrimPrefix(name, "Test")
	if separator := strings.IndexByte(part, '_'); separator >= 0 {
		part = part[:separator]
	}
	return parseNormalizedACPart(part, prefixes)
}

func parseNormalizedACPart(part string, prefixes map[string]bool) (ac string, ambiguous, matched bool) {
	if canonical := strings.ReplaceAll(part, "_", "-"); strings.Contains(part, "_") {
		if prefix, _, _, ok := parseACID(canonical); ok && prefixes[prefix] {
			return canonical, false, true
		}
	}

	if part == "" {
		return "", false, false
	}
	var candidates []string
	for prefix := range prefixes {
		normalizedPrefix := normalizeACPrefix(prefix)
		if !strings.HasPrefix(part, normalizedPrefix) {
			continue
		}
		ordinal := part[len(normalizedPrefix):]
		if !goOrdinalRe.MatchString(ordinal) {
			continue
		}
		candidate := prefix + "-" + ordinal
		if _, _, _, ok := parseACID(candidate); ok {
			candidates = append(candidates, candidate)
		}
	}
	sort.Strings(candidates)
	if len(candidates) == 1 {
		return candidates[0], false, true
	}
	if len(candidates) > 1 {
		return "", true, false
	}
	return "", false, false
}
