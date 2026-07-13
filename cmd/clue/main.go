// clue is the deterministic judge for a Cliewen corpus: stateless, no
// AI, no orchestration (Foundation §13). Commands are added only when a
// linter rule or skill needs them.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
)

// version is the release stamp, injected at build time via
// `-ldflags "-X main.version=<semver>"` (see .github/workflows/release.yml).
// When no stamp is injected, main falls back to the module version Go
// embeds in `go install module@vX.Y.Z` builds; checkout and commit builds
// report "dev" and are exempt from skill-drift checks (ADR-011).
var version = "dev"

// pseudoVersion matches Go pseudo-versions (…-yyyymmddhhmmss-abcdefabcdef):
// a commit, not a release.
var pseudoVersion = regexp.MustCompile(`[.-][0-9]{14}-[0-9a-f]{12}$`)

// releaseFromModuleVersion maps the module version embedded by
// `go install module@vX.Y.Z` to a bare-semver release stamp, or "" when it
// names no release: "(devel)", a pseudo-version (branch/commit install, or
// the VCS-derived version a checkout build embeds since Go 1.24), or a
// build with local modifications ("+dirty").
func releaseFromModuleVersion(v string) string {
	if v == "" || v == "(devel)" {
		return ""
	}
	base, meta, _ := strings.Cut(v, "+")
	if meta == "dirty" || pseudoVersion.MatchString(base) {
		return ""
	}
	return strings.TrimPrefix(base, "v")
}

const usage = `clue — a verifiable thread from goal to test

Usage:
  clue validate [--forbid-changes] [path]
  clue version

Commands:
  validate   Scan docs/ and changes/ under path (default ".") and check
             the frontmatter graph: core fields, unique IDs, link
             resolution, status vocabularies, folder READMEs, index
             integrity, and skill version drift.

             --forbid-changes  fail when /changes contains files — the
                               digest-before-merge gate used by CI.

  version    Print the release version this clue was built from.

Exit codes: 0 corpus valid · 1 issues found · 2 usage error
`

func main() {
	if version == "dev" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			if v := releaseFromModuleVersion(bi.Main.Version); v != "" {
				version = v
			}
		}
	}
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "validate":
		os.Exit(runValidate(os.Args[2:]))
	case "version", "--version":
		os.Exit(runVersion(os.Stdout))
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "clue: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	forbid := fs.Bool("forbid-changes", false, "fail when /changes contains files")
	_ = fs.Parse(args)
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}

	c, issues := corpus.Scan(root)
	issues = append(issues, corpus.Validate(c, corpus.Options{ForbidChanges: *forbid, Version: version})...)
	if len(issues) > 0 {
		for _, is := range issues {
			fmt.Println(is)
		}
		fmt.Fprintf(os.Stderr, "clue validate: %d issue(s)\n", len(issues))
		return 1
	}
	if n := inferredCount(c); n > 0 {
		fmt.Printf("clue validate: OK (%d artifacts, %d born inferred awaiting verification)\n", len(c.Artifacts), n)
	} else {
		fmt.Printf("clue validate: OK (%d artifacts)\n", len(c.Artifacts))
	}
	return 0
}

// runVersion prints the release stamp (AC-019). "dev" for source builds.
func runVersion(w io.Writer) int {
	fmt.Fprintf(w, "clue %s\n", version)
	return 0
}

// inferredCount is the visible to-do list of unverified meaning (ADR-010).
func inferredCount(c *corpus.Corpus) int {
	n := 0
	for _, a := range c.Artifacts {
		if p, _ := a.Fields["provenance"].(string); p == "inferred" {
			n++
		}
	}
	return n
}
