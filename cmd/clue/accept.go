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
	brief := fs.String("brief", "", "path to the completed acceptance brief")
	check := fs.Bool("check", false, "preflight only; do not integrate")
	approval := fs.String("approval", "", "recorded human approval bound to candidate, base and brief hash")
	// The documented candidate-first form also permits flags first.
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		args = append(append([]string{}, args[1:]...), args[0])
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 || *base == "" || *brief == "" {
		fmt.Fprintln(errOut, "usage: clue accept <candidate-sha> --base <base-sha> --brief <file> [--approval <file>] [--check]")
		return 2
	}
	p, err := accept.Check(accept.Request{Root: ".", Candidate: fs.Arg(0), Base: *base, BriefPath: *brief, Version: version})
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
