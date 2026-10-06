// This is this repository's scenario-trial harness, not part of the clue
// executable and not shipped to adopters. A trial starts only when the
// maintainer runs it; nothing here is wired into CI.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const usage = `usage:
  scenario-trials run -agent NAME [-scenario routing] [-runs 5] [-clue-commit REV] [-login PATH] [-out DIR] [-model NAME] [-effort WORD] [-variant NAME]
  scenario-trials check RUN_DIRECTORY   evaluate stored transcripts again with the current checks; runs no agent`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "run":
		runCommand(os.Args[2:])
	case "check":
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(2)
		}
		if err := Recheck(os.Args[2]); err != nil {
			fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
}

func runCommand(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	scenario := fs.String("scenario", "routing", "scenario: "+strings.Join(scenarioNames(), ", "))
	agent := fs.String("agent", "", "agent adapter: "+strings.Join(adapterNames(), ", "))
	runs := fs.Int("runs", 5, "number of runs")
	commit := fs.String("clue-commit", "HEAD", "revision the clue binary and skills are built from")
	login := fs.String("login", "", "where the agent's login lives, outside the repository; each adapter has a default under ~/.cliewen-trial")
	out := fs.String("out", "scenario-runs", "directory for run records")
	model := fs.String("model", "", "model to pass to the agent")
	effort := fs.String("effort", "", "reasoning effort, for an agent that takes one; the conditions record it")
	variant := fs.String("variant", baselineVariant, "method variant, a named removal applied inside the container only: "+strings.Join(variantNames(), ", "))
	_ = fs.Parse(args)
	ad, err := adapterFor(*agent)
	if err != nil {
		fatal(err)
	}
	if *login == "" {
		home, _ := os.UserHomeDir()
		*login = ad.DefaultLogin(home)
	}
	if err := Run(Options{Scenario: *scenario, Agent: *agent, Runs: *runs, Commit: *commit, Login: *login, OutDir: *out, Model: *model, Variant: *variant, Effort: *effort}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "scenario-trials:", err)
	os.Exit(1)
}
