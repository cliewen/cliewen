// This is this repository's evidence producer, not part of the clue executable.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cliewen/cliewen/internal/evidence"
	"github.com/cliewen/cliewen/internal/evidenceexport"
)

func main() {
	root := flag.String("root", ".", "repository root")
	flag.Parse()
	p, err := evidenceexport.Collect(*root, "cliewen-tests")
	if err == nil {
		p.Include = append(p.Include, "internal/evidenceexport/**/*.go", "tools/export-evidence/**/*.go")
		p.Inputs, err = evidence.Snapshot(*root, p)
	}
	if err == nil {
		if len(p.Diagnostics) > 0 {
			for _, d := range p.Diagnostics {
				fmt.Fprintf(os.Stderr, "%s: %s\n", d.Path, d.Message)
			}
		}
		err = evidence.Write(*root, evidence.Manifest{Version: evidence.Version, Producers: []evidence.Producer{p}})
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
