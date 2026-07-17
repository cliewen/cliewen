package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/cliewen/cliewen/internal/scaffold"
)

// runInit materializes the Cliewen convention into a repository
// (CAP-001): docs taxonomy, AGENTS.md, skills, CI workflow template.
// Idempotent: existing files are reported and skipped, only README
// index blocks are regenerated.
func runInit(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	_ = fs.Parse(args)
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	rep, err := scaffold.Run(root)
	if err != nil {
		fmt.Fprintf(out, "clue init: %v\n", err)
		return 1
	}
	for _, p := range rep.Created {
		fmt.Fprintf(out, "created  %s\n", p)
	}
	for _, p := range rep.Skipped {
		fmt.Fprintf(out, "exists   %s (skipped — never overwritten)\n", p)
	}
	for _, p := range rep.Indexed {
		fmt.Fprintf(out, "indexed  %s\n", p)
	}
	fmt.Fprintf(out, "clue init: %d created, %d skipped, %d index block(s) regenerated\n", len(rep.Created), len(rep.Skipped), len(rep.Indexed))
	if len(rep.Created) > 0 {
		fmt.Fprintln(out, "next: run `clue validate` — it should be green; then read docs/README.md")
	}
	return 0
}
