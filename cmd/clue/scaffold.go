package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/cliewen/cliewen/internal/scaffold"
)

// runScaffold regenerates the taxonomy README index blocks (CAP-005):
// the standalone exposure of the engine clue init runs. It materializes
// nothing — a missing docs tree is an error, missing folder READMEs are
// reported, never invented.
func runScaffold(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("scaffold", flag.ExitOnError)
	_ = fs.Parse(args)
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	rep, err := scaffold.Regen(root)
	if err != nil {
		fmt.Fprintf(errOut, "clue scaffold: %v\n", err)
		return 1
	}
	for _, p := range rep.Indexed {
		fmt.Fprintf(out, "indexed  %s\n", p)
	}
	for _, p := range rep.MissingReadmes {
		fmt.Fprintf(out, "missing  %s (folder has no README — `clue validate` requires one; scaffold does not invent it)\n", p)
	}
	fmt.Fprintf(out, "clue scaffold: %d index block(s) regenerated\n", len(rep.Indexed))
	return 0
}
