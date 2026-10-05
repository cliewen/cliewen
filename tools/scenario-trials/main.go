// This is this repository's scenario-trial harness, not part of the clue
// executable and not shipped to adopters. A trial starts only when the
// maintainer runs it; nothing here is wired into CI.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const usage = `usage:
  scenario-trials run [-scenario routing] [-agent claude] [-runs 5] [-clue-commit REV] [-token-file FILE] [-out DIR] [-model NAME]
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
	home, _ := os.UserHomeDir()
	scenario := fs.String("scenario", "routing", "scenario name")
	agent := fs.String("agent", "claude", "agent adapter")
	runs := fs.Int("runs", 5, "number of runs")
	commit := fs.String("clue-commit", "HEAD", "revision the clue binary and skills are built from")
	token := fs.String("token-file", filepath.Join(home, ".cliewen-trial", "token"), "file holding the agent login token, outside the repository")
	out := fs.String("out", "scenario-runs", "directory for run records")
	model := fs.String("model", "", "model to pass to the agent")
	variant := fs.String("variant", "baseline", "method variant; only baseline exists before M-107")
	_ = fs.Parse(args)
	if *variant != "baseline" {
		fatal(fmt.Errorf("method variant %q does not exist yet", *variant))
	}
	if err := Run(Options{Scenario: *scenario, Agent: *agent, Runs: *runs, Commit: *commit, TokenFile: *token, OutDir: *out, Model: *model, Variant: *variant}); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "scenario-trials:", err)
	os.Exit(1)
}
