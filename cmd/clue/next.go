package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/cliewen/cliewen/internal/corpus"
)

// runNext reports the next actionable milestone without changing any corpus
// state. The command is deliberately an orientation aid, not an orchestration
// or claim mechanism: a human still decides which candidate to start.
func runNext(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("next", flag.ContinueOnError)
	fs.SetOutput(errOut)
	all := fs.Bool("all", false, "list every unfinished milestone in active plans")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(errOut, "usage: clue next [--all] [path]")
		return 2
	}
	root := "."
	if fs.NArg() == 1 {
		root = fs.Arg(0)
	}
	c, scanIssues := corpus.Scan(root)
	if len(scanIssues) > 0 {
		for _, issue := range scanIssues {
			fmt.Fprintln(errOut, issue)
		}
		fmt.Fprintf(errOut, "clue next: %d scan issue(s)\n", len(scanIssues))
		return 1
	}

	actionable := corpus.UnfinishedMilestones(c)
	if *all {
		if len(actionable) == 0 {
			fmt.Fprintln(out, "No unfinished milestone in active plans.")
		} else {
			fmt.Fprintln(out, "Unfinished milestones in active plans:")
			for _, milestone := range actionable {
				printMilestone(out, milestone)
			}
		}
	} else if len(actionable) == 0 {
		fmt.Fprintln(out, "No actionable milestone in active plans.")
	} else {
		fmt.Fprintln(out, "Next unfinished milestone (human selection required):")
		printMilestone(out, actionable[0])
		if len(actionable) > 1 {
			fmt.Fprintf(out, "Other active alternatives: %d (use --all to list them)\n", len(actionable)-1)
		}
	}

	draft := corpus.DraftUnfinishedMilestones(c)
	if len(draft) > 0 {
		fmt.Fprintln(out, "Proposed unfinished milestones in draft plans (not actionable):")
		for _, milestone := range draft {
			printMilestone(out, milestone)
		}
	}
	return 0
}

func printMilestone(out io.Writer, milestone corpus.Milestone) {
	fmt.Fprintf(out, "%s/%s | %s | %s | %s\n", milestone.Plan.ID, milestone.ID, milestone.Status, milestone.Name, milestone.Plan.Path)
	if milestone.ExitCriterion != "" && milestone.ExitCriterion != milestone.Name {
		fmt.Fprintf(out, "  Exit criterion: %s\n", milestone.ExitCriterion)
	}
}
