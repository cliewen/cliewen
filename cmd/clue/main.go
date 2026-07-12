// clue is the deterministic judge for a Cliewen corpus: stateless, no
// AI, no orchestration (Foundation §13). Commands are added only when a
// linter rule or skill needs them.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cliewen/cliewen/internal/corpus"
)

const usage = `clue — a verifiable thread from goal to test

Usage:
  clue validate [--forbid-changes] [path]

Commands:
  validate   Scan docs/ and changes/ under path (default ".") and check
             the frontmatter graph: core fields, unique IDs, link
             resolution, status vocabularies, folder READMEs, index
             integrity.

             --forbid-changes  fail when /changes contains files — the
                               digest-before-merge gate used by CI.

Exit codes: 0 corpus valid · 1 issues found · 2 usage error
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "validate":
		os.Exit(runValidate(os.Args[2:]))
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
	issues = append(issues, corpus.Validate(c, corpus.Options{ForbidChanges: *forbid})...)
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
