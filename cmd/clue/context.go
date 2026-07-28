package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
)

func runContext(args []string, out, errOut io.Writer) int {
	if len(args) < 1 || len(args) > 2 {
		fmt.Fprintln(errOut, "usage: clue context <id> [path]")
		return 2
	}
	id := args[0]
	root := "."
	if len(args) == 2 {
		root = args[1]
	}

	c, issues := corpus.Scan(root)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Fprintln(errOut, issue)
		}
		fmt.Fprintf(errOut, "clue context: %d scan issue(s)\n", len(issues))
		return 1
	}
	artifacts, err := corpus.Context(c, id)
	if err != nil {
		fmt.Fprintln(errOut, "clue context:", err)
		return 1
	}

	for i, artifact := range artifacts {
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "===== %s | %s =====\n", artifact.ID, artifact.Path)
		content := c.Contents[artifact.Path]
		fmt.Fprint(out, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprintln(out)
		}
	}
	return 0
}
