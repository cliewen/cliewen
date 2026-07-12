package corpus

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The AC↔test contract (P-001/M-002) and the test-purpose taxonomy:
// every AC declared in an active criteria.md has at least one test, every
// test that names an AC points at one that exists, and every test declares
// exactly one purpose (ADR-005, ADR-006). AC IDs are namespaced per
// criteria file via the optional `ac-prefix` frontmatter field, default
// `AC` (ADR-009). Go carries the declaration in the function name
// (TestAC004_…, TestMG010_…, TestUnit_…); JVM test files carry it as
// @Tag("MG_010"), harvested at file level — per-test purpose enforcement
// there is the adopting repo's ArchUnit rule (ADR-009).
var (
	acTagRe        = regexp.MustCompile(`@([A-Z][A-Z0-9]*)-(\d+)\b`)
	acPrefixRe     = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)
	testFuncRe     = regexp.MustCompile(`(?m)^func (Test\w*)\s*\(`)
	fixedPurposeRe = regexp.MustCompile(`^Test(Unit|Sanity|Arch)(_\w*)?$`)
	acNameRe       = regexp.MustCompile(`^Test([A-Z][A-Z0-9]*?)(\d+)(_\w*)?$`)
	jvmTagRe       = regexp.MustCompile(`@Tag\(\s*"([^"]+)"\s*\)`)
	jvmACRe        = regexp.MustCompile(`^([A-Z][A-Z0-9]*)-(\d+)$`)
)

type acDecl struct {
	path, status string
	retired      bool
}

func checkACTests(c *Corpus) []Issue {
	declared := map[string]acDecl{}
	// The default namespace is always known, so an undeclared AC-xxx
	// reference fails as unknown, not as purpose-less.
	prefixes := map[string]bool{"AC": true}
	var issues []Issue
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
		for _, line := range strings.Split(a.Body, "\n") {
			retired := strings.Contains(line, "@retired")
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
				declared[ac] = acDecl{a.Path, a.Status, retired}
			}
		}
	}

	tested := map[string]bool{}
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
		if !isGo && !isJVM {
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
				am := acNameRe.FindStringSubmatch(name)
				if am == nil || !prefixes[am[1]] {
					issues = append(issues, Issue{relSlash, "test " + name + " declares no purpose (ADR-006: prefix <ACPREFIX><digits>, Unit, Sanity or Arch)"})
					continue
				}
				issues = append(issues, checkACRef(relSlash, "test "+name, am[1]+"-"+am[2], declared, tested)...)
			}
			return nil
		}

		// JVM files: harvest @Tag values at file level (ADR-009). Tags in
		// a known AC namespace are references; everything else is runner
		// metadata clue ignores (ADR-006).
		for _, m := range jvmTagRe.FindAllStringSubmatch(text, -1) {
			norm := strings.ReplaceAll(m[1], "_", "-")
			am := jvmACRe.FindStringSubmatch(norm)
			if am == nil || !prefixes[am[1]] {
				continue
			}
			issues = append(issues, checkACRef(relSlash, `tag "`+m[1]+`"`, norm, declared, tested)...)
		}
		return nil
	})

	for ac, d := range declared {
		if d.status == "active" && !d.retired && !tested[ac] {
			issues = append(issues, Issue{d.path, ac + " has no test (convention per ADR-005/ADR-009: a Go test named Test" + strings.ReplaceAll(ac, "-", "") + "_… or a framework tag \"" + strings.ReplaceAll(ac, "-", "_") + "\")"})
		}
	}
	return issues
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
