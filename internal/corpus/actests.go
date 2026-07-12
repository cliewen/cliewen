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
// exactly one purpose (ADR-005, ADR-006). Go carries the declaration in the
// function name: TestAC004_…, TestUnit_…, TestSanity_…, TestArch_….
var (
	acTagRe       = regexp.MustCompile(`@AC-(\d+)\b`)
	testFuncRe    = regexp.MustCompile(`(?m)^func (Test\w*)\s*\(`)
	testPurposeRe = regexp.MustCompile(`^Test(AC(\d+)|Unit|Sanity|Arch)(_\w*)?$`)
)

func checkACTests(c *Corpus) []Issue {
	type decl struct{ path, status string }
	declared := map[string]decl{}
	for _, a := range c.Artifacts {
		if a.Type != "criteria" {
			continue
		}
		for _, m := range acTagRe.FindAllStringSubmatch(a.Body, -1) {
			declared["AC-"+m[1]] = decl{a.Path, a.Status}
		}
	}

	tested := map[string]bool{}
	var issues []Issue
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
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return nil
		}
		rel, _ := filepath.Rel(c.Root, p)
		relSlash := filepath.ToSlash(rel)
		for _, m := range testFuncRe.FindAllStringSubmatch(string(data), -1) {
			name := m[1]
			if name == "TestMain" {
				continue // the harness hook, not a test
			}
			pm := testPurposeRe.FindStringSubmatch(name)
			if pm == nil {
				issues = append(issues, Issue{relSlash, "test " + name + " declares no purpose (ADR-006: prefix AC<digits>, Unit, Sanity or Arch)"})
				continue
			}
			if pm[2] != "" { // the AC<digits> purpose
				ac := "AC-" + pm[2]
				tested[ac] = true
				if _, ok := declared[ac]; !ok {
					issues = append(issues, Issue{relSlash, "test " + name + " references " + ac + " which no criteria.md declares"})
				}
			}
		}
		return nil
	})

	for ac, d := range declared {
		if d.status == "active" && !tested[ac] {
			issues = append(issues, Issue{d.path, ac + " has no test (convention per ADR-005: a test named TestAC" + strings.TrimPrefix(ac, "AC-") + "_…)"})
		}
	}
	return issues
}
