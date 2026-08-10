package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
)

// defaultContextDepth is the root and what it links to directly. A corpus
// densifies as campaigns close — a hub artifact acquires enough inbound and
// outbound edges that breadth-first search pours through it into everything —
// so an unbounded slice stops being a slice and becomes the repository. The
// bound is a default rather than a limit: --depth widens it, and the frontier
// below names every artifact the bound held back.
const defaultContextDepth = 1

func runContext(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("context", flag.ContinueOnError)
	fs.SetOutput(errOut)
	depth := fs.String("depth", "", `link hops beyond the root to print: a number, or "all" (default 1)`)
	stats := fs.Bool("stats", false, "print the slice's artifact and byte counts")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 || fs.NArg() > 2 {
		fmt.Fprintln(errOut, `usage: clue context [--depth=<n>|all] [--stats] <id> [path]`)
		return 2
	}
	id := fs.Arg(0)
	root := "."
	if fs.NArg() == 2 {
		root = fs.Arg(1)
	}

	opts := corpus.ContextOptions{Depth: defaultContextDepth}
	switch *depth {
	case "":
	case "all":
		opts.Depth = corpus.DepthAll
	default:
		n, err := strconv.Atoi(*depth)
		if err != nil || n < 0 {
			fmt.Fprintf(errOut, "clue context: --depth takes a non-negative number or %q, got %q\n", "all", *depth)
			return 2
		}
		opts.Depth = n
	}

	c, issues := corpus.Scan(root)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Fprintln(errOut, issue)
		}
		fmt.Fprintf(errOut, "clue context: %d scan issue(s)\n", len(issues))
		return 1
	}
	artifacts, frontier, unfollowed, err := corpus.Context(c, id, opts)
	if err != nil {
		fmt.Fprintln(errOut, "clue context:", err)
		return 1
	}
	for _, issue := range unfollowed {
		fmt.Fprintln(errOut, "clue context:", issue)
	}

	bytes := 0
	for i, artifact := range artifacts {
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "===== %s | %s =====\n", artifact.ID, artifact.Path)
		content := c.Contents[artifact.Path]
		bytes += len(content)
		fmt.Fprint(out, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprintln(out)
		}
	}

	printFrontier(out, frontier, opts.Depth)
	if *stats {
		// The figure is artifact content, which is what a slice costs a reader
		// to load; the boundary lines and this report are not counted.
		fmt.Fprintf(out, "\n----- %d artifact(s), %d content byte(s); %d beyond the slice -----\n", len(artifacts), bytes, len(frontier))
	}
	return 0
}

// printFrontier names what the bound held back. Only the next hop is listed by
// name: those are the artifacts one widening away, and they are the ones a
// reader can act on. Listing the whole remainder would grow with the corpus and
// recreate the cost the bound exists to remove, so the rest is a count.
func printFrontier(out io.Writer, frontier []corpus.Frontier, depth int) {
	if len(frontier) == 0 {
		return
	}
	beyond := 0
	for _, f := range frontier {
		if f.Hops != depth+1 {
			beyond++
		}
	}
	fmt.Fprintf(out, "\n----- frontier: %d artifact(s) reached but not printed -----\n", len(frontier))
	for _, f := range frontier {
		if f.Hops != depth+1 {
			continue
		}
		fmt.Fprintf(out, "%s | %s\n", f.Artifact.ID, f.Artifact.Title)
	}
	if beyond > 0 {
		fmt.Fprintf(out, "... and %d further artifact(s) more than %d hop(s) out\n", beyond, depth+1)
	}
	fmt.Fprintf(out, "widen with --depth=%d, or --depth=all\n", depth+1)
}
