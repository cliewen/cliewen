package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cliewen/cliewen/internal/accept"
)

func runAccept(args []string, input io.Reader, out, errOut io.Writer, interactive bool) int {
	fs := flag.NewFlagSet("accept", flag.ContinueOnError)
	fs.SetOutput(errOut)
	base := fs.String("base", "", "full accepted base commit ID")
	// Kept only to refuse it plainly: an adopter running an older skill
	// against this binary learns why rather than seeing a usage error.
	brief := fs.String("brief", "", "no longer used; refused")
	check := fs.Bool("check", false, "preflight only; do not integrate")
	approval := fs.String("approval", "", "recorded human approval bound to candidate and base")
	// The documented candidate-first form also permits flags first.
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		args = append(append([]string{}, args[1:]...), args[0])
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *brief != "" {
		fmt.Fprintln(errOut, "clue accept: --brief is no longer used; the acceptance brief was removed. Accept with the candidate and base only, and record Human-class proof in the proposal's \"## Human evidence\" section.")
		return 2
	}
	if fs.NArg() != 1 || *base == "" {
		fmt.Fprintln(errOut, "usage: clue accept <candidate-sha> --base <base-sha> [--approval <file>] [--check]")
		return 2
	}
	p, err := accept.Check(accept.Request{Root: ".", Candidate: fs.Arg(0), Base: *base, Version: version})
	if err != nil {
		fmt.Fprintln(errOut, "clue accept:", err)
		return 1
	}
	if *approval != "" {
		if err = p.VerifyApproval(*approval); err != nil {
			fmt.Fprintln(errOut, "clue accept:", err)
			return 1
		}
	}
	if *check {
		fmt.Fprintln(out, "clue accept: preflight passed; no acceptance performed. Verification and review remain recorded claims.")
		return 0
	}
	var commit string
	if *approval != "" {
		commit, err = p.ConfirmRecorded(*approval, out)
	} else {
		commit, err = p.Confirm(input, out, interactive)
	}
	if err != nil {
		fmt.Fprintln(errOut, "clue accept:", err)
		return 1
	}
	fmt.Fprintln(out, "Accepted locally:", commit, "(not pushed)")
	return 0
}

func interactiveInput() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
