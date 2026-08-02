// clue is the deterministic judge for a Cliewen corpus: stateless, no
// AI, no orchestration (Foundation §13). Commands are added only when a
// linter rule, a skill, or an onboarding criterion needs them.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/migrate"
	"github.com/cliewen/cliewen/internal/refs"
	"github.com/cliewen/cliewen/internal/release"
)

// version is the release stamp, injected at build time via
// `-ldflags "-X main.version=<semver>"` (see .github/workflows/release.yml).
// When no stamp is injected, main falls back to the module version Go
// embeds in `go install module@vX.Y.Z` builds; checkout and commit builds
// report "dev" and are exempt from skill-drift checks (ADR-011).
var version = "dev"

// pseudoVersion matches Go pseudo-versions (…-yyyymmddhhmmss-abcdefabcdef):
// a commit, not a release.
var pseudoVersion = regexp.MustCompile(`[.-][0-9]{14}-[0-9a-f]{12}$`)

// releaseFromModuleVersion maps the module version embedded by
// `go install module@vX.Y.Z` to a bare-semver release stamp, or "" when it
// names no release: "(devel)", a pseudo-version (branch/commit install, or
// the VCS-derived version a checkout build embeds since Go 1.24), or a
// build with local modifications ("+dirty").
func releaseFromModuleVersion(v string) string {
	if v == "" || v == "(devel)" {
		return ""
	}
	base, meta, _ := strings.Cut(v, "+")
	if meta == "dirty" || pseudoVersion.MatchString(base) {
		return ""
	}
	return strings.TrimPrefix(base, "v")
}

const usage = `clue — a verifiable thread from goal to acceptance evidence

Usage:
  clue init [path]
  clue scaffold [path]
  clue context <id> [path]
  clue refs [--apply] [--timeout=<duration>] [path]
  clue validate [--forbid-changes] [--coverage] [--reality-gaps] [path]
  clue latest [--quiet] [--timeout=<duration>]
  clue version

Commands:
  init       Materialize the Cliewen convention under path (default "."):
             the docs/ taxonomy, AGENTS.md routing hub, agent skills
             (.agents/skills + .claude/skills mirror), and a CI workflow
             template. Idempotent: existing files are never overwritten
             (they are reported and skipped); README index blocks between
             the clue:index markers are regenerated. A target directory
             that is a symlink — a skills tree shared across checkouts —
             is left untouched: nothing is written through it, and it is
             reported as linked.

  scaffold   Regenerate the taxonomy README index blocks under path
             (default ".") from folder contents: entries whose targets
             survive keep their hand-written lines, missing entries are
             appended, prose outside the clue:index markers is never
             touched. Materializes nothing — missing folder READMEs are
             reported, and a path without a docs/ tree is an error.

  context    Print one artifact and the transitive closure of its outgoing
             links, with complete markdown content in deterministic order.
             Acceptance-criterion and milestone IDs resolve to the artifact
             that declares them. Path defaults to ".".

  migrate    Preview a versioned corpus and managed-carrier migration; use
             --apply to write the complete safe plan and
             --reversal-cost=low|high to resolve missing inferred-meaning routing.
             Existing prose and locally modified generated files are never
             overwritten. Path defaults to ".".

  refs       Resolve the external addresses docs/ and changes/ point at,
             classifying each: reachable, restricted (it exists, this
             runner may not read it), redirected, gone, or unreachable.
             Only "gone" is an error; an outage elsewhere never
             condemns a corpus. Use
             --apply to rewrite redirected addresses in place. A clue:
             identity is never followed. Never make this a required
             check: another host's uptime must not gate a merge.
             Path defaults to ".".

  validate   Scan docs/ and changes/ under path (default ".") and check
             the frontmatter graph: core fields, unique IDs, link
             resolution, status vocabularies, folder READMEs, index
             integrity, and skill version drift.

             --forbid-changes  fail when /changes contains files — the
                               digest-before-merge gate used by CI.
             --coverage        print derived proof coverage by capability,
                               then any pointer to proof in another
                               repository, listed apart as named but
                               locally unproven.
             --reality-gaps    print capabilities contradicted by incident
                               analyses after their corpus was green.

  latest     Report whether a newer clue release exists, and how to install
             it on the machine this is running on — the PowerShell script on
             Windows, the shell script on macOS and Linux, go install where
             no prebuilt asset exists — followed by the clue migrate sequence
             that moves the repository with it. Reaches the network, so it is
             never part of a validation verdict and must never be a required
             check. It writes no file in the repository and never replaces
             the binary. Being unable to tell — offline, a timeout, a rate
             limit, an unrecognized answer — is reported calmly and exits 0.
             The answer is cached in the user's cache directory for a day.

             --quiet           print one line when behind and nothing when
                               current, for a session start where output is
                               context.
             --timeout=<dur>   request budget (default 3s).

  version    Print the release version this clue was built from.

Exit codes: 0 corpus valid · 1 issues found · 2 usage error
`

func main() {
	if version == "dev" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			if v := releaseFromModuleVersion(bi.Main.Version); v != "" {
				version = v
			}
		}
	}
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "init":
		os.Exit(runInit(os.Args[2:], os.Stdout, os.Stderr))
	case "scaffold":
		os.Exit(runScaffold(os.Args[2:], os.Stdout, os.Stderr))
	case "context":
		os.Exit(runContext(os.Args[2:], os.Stdout, os.Stderr))
	case "migrate":
		os.Exit(runMigrate(os.Args[2:], os.Stdout, os.Stderr))
	case "refs":
		os.Exit(runRefs(os.Args[2:], os.Stdout, os.Stderr))
	case "validate":
		os.Exit(runValidate(os.Args[2:], os.Stdout))
	case "latest":
		os.Exit(runLatest(os.Args[2:], os.Stdout, os.Stderr, runtimeCheck()))
	case "version", "--version":
		os.Exit(runVersion(os.Stdout))
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "clue: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runMigrate(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	fs.SetOutput(errOut)
	apply := fs.Bool("apply", false, "write the complete migration plan after preflight")
	reversalCost := fs.String("reversal-cost", "", "explicitly classify missing inferred-artifact routing as low or high")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(errOut, "clue migrate: expected at most one repository path")
		return 2
	}

	plan, err := migrate.Plan(root, migrate.Options{ReversalCost: *reversalCost})
	if err != nil {
		fmt.Fprintf(errOut, "clue migrate: %v\n", err)
		return 2
	}
	mode := "preview"
	if *apply {
		mode = "apply"
	}
	fmt.Fprintf(out, "clue migrate: %s for target pair %s\n", mode, plan.Target)
	for _, notice := range plan.Notices {
		fmt.Fprintf(out, "notice %s %s: %s\n", notice.Migration, notice.Path, notice.Message)
	}
	for _, change := range plan.Changes {
		fmt.Fprintf(out, "%s %s: %s\n", change.Migration, change.Path, change.Description)
	}
	for _, finding := range plan.Findings {
		fmt.Fprintf(out, "finding %s %s: %s\n", finding.Migration, finding.Path, finding.Message)
	}
	if len(plan.Findings) > 0 {
		fmt.Fprintf(out, "clue migrate: %d finding(s); no files changed\n", len(plan.Findings))
		return 1
	}
	if !*apply {
		if len(plan.Changes) == 0 {
			fmt.Fprintln(out, "clue migrate: no changes needed")
		} else {
			fmt.Fprintf(out, "clue migrate: %d file(s) would change; preview only\n", len(plan.Changes))
		}
		return 0
	}
	if err := migrate.Apply(root, plan); err != nil {
		fmt.Fprintf(errOut, "clue migrate: %v\n", err)
		return 1
	}
	if len(plan.Changes) == 0 {
		fmt.Fprintln(out, "clue migrate: no changes needed")
	} else {
		fmt.Fprintf(out, "clue migrate: applied %d file(s)\n", len(plan.Changes))
	}
	return 0
}

// runValidate takes its writer so the command's own output is observable: a
// criterion that describes what a user sees is only proven when a test reads
// what the command actually printed.
func runValidate(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	forbid := fs.Bool("forbid-changes", false, "fail when /changes contains files")
	coverage := fs.Bool("coverage", false, "print derived per-capability proof coverage; never a committed registry")
	realityGaps := fs.Bool("reality-gaps", false, "print capabilities contradicted by incident analyses; never a committed registry")
	indexRows := fs.Bool("index-rows", false, "print index rows that only restate their own link; never a committed registry")
	_ = fs.Parse(args)
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}

	c, issues := corpus.Scan(root)
	issues = append(issues, corpus.Validate(c, corpus.Options{ForbidChanges: *forbid, Version: version})...)
	provenance := corpus.ProvenanceBacklog(c)
	if len(issues) > 0 {
		for _, is := range issues {
			fmt.Fprintln(out, is)
		}
		fmt.Fprintf(os.Stderr, "clue validate: %d issue(s)", len(issues))
		if len(provenance.BlockerArtifacts) > 0 {
			fmt.Fprintf(os.Stderr, ", %d high-cost inferred activation blocker(s)", len(provenance.BlockerArtifacts))
		}
		if len(provenance.Decisions) > 0 {
			fmt.Fprintf(os.Stderr, ", %d inferred decision(s) awaiting verification", len(provenance.Decisions))
		}
		fmt.Fprintln(os.Stderr)
		return 1
	}
	if *coverage {
		for _, cc := range corpus.Coverage(c) {
			fmt.Fprintf(out, "%s: %s\n", cc.Capability, cc.State)
		}
		// A pointer to proof in another repository is listed apart from
		// coverage, never inside it. Naming it says a human can go and look;
		// counting it would be importing a verdict this judge cannot see.
		for _, p := range corpus.ForeignPointers(c) {
			fmt.Fprintf(out, "%s: named but locally unproven\n", p)
		}
	}
	if *realityGaps {
		for _, gap := range corpus.RealityGaps(c) {
			fmt.Fprintf(out, "%s: contradicted by %s\n", gap.Capability, strings.Join(gap.Analyses, ", "))
		}
	}
	// No command clears this count: index regeneration preserves any row
	// whose target still exists, so every repair is by hand. Listing the rows
	// is what makes the number actionable (ADR-041).
	fillerRows := corpus.IndexRowBacklog(c)
	if *indexRows {
		for _, row := range fillerRows {
			fmt.Fprintf(out, "%s: %s states only its own link\n", row.Readme, row.Target)
		}
	}
	notes := ""
	if n := len(provenance.BlockerArtifacts); n > 0 {
		notes += fmt.Sprintf(", %d high-cost inferred activation blocker(s)", n)
	}
	if n := len(provenance.Decisions); n > 0 {
		notes += fmt.Sprintf(", %d inferred decision(s) awaiting verification", n)
	}
	if n := agentConstraintCount(c); n > 0 {
		notes += fmt.Sprintf(", %d agent-enforced constraint(s) awaiting machine checks", n)
	}
	if n := len(fillerRows); n > 0 {
		notes += fmt.Sprintf(", %d index row(s) stating only their own link", n)
	}
	fmt.Fprintf(out, "clue validate: OK (%d artifacts%s)\n", len(c.Artifacts), notes)
	return 0
}

// runtimeCheck is the real machine's options: the running stamp and the host
// this binary was built for. It is a named function rather than a literal at
// the call site so a test can assert that the platform a user's run reports on
// is the one they are actually running (AC-075).
func runtimeCheck() release.Options {
	return release.Options{
		Current:  version,
		Platform: release.Platform{OS: runtime.GOOS, Arch: runtime.GOARCH},
	}
}

// runLatest reports whether a newer release exists (AC-075, AC-076, AC-077).
//
// It takes the check's options from its caller rather than building them, so a
// test can supply the platform and the network the run sees; main supplies the
// real ones. It never returns non-zero for a network answer: only a usage
// error fails, because nothing this command learns is a defect in the
// repository (ADR-042).
func runLatest(args []string, out, errOut io.Writer, opts release.Options) int {
	fs := flag.NewFlagSet("latest", flag.ContinueOnError)
	fs.SetOutput(errOut)
	quiet := fs.Bool("quiet", false, "print one line when behind and nothing when current")
	timeout := fs.Duration("timeout", 0, "request budget (default 3s)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		// The check reads no repository, so a path would be a silent no-op
		// rather than the scoping the other commands give it.
		fmt.Fprintln(errOut, "clue latest: takes no path — it reports on the installed clue, not on a repository")
		return 2
	}
	if *timeout != 0 {
		opts.Timeout = *timeout
	}

	report := release.Check(opts)

	switch {
	case report.Known && report.Behind:
		if *quiet {
			fmt.Fprintf(out, "clue %s is behind %s — run \"clue latest\" for the upgrade recipe\n", report.Current, report.Latest)
			return 0
		}
		fmt.Fprintf(out, "clue %s — %s is available\n\n", report.Current, report.Latest)
		fmt.Fprintln(out, "install it on this machine:")
		for _, line := range report.Recipe {
			fmt.Fprintf(out, "  %s\n", line)
		}
		fmt.Fprintln(out, "\nthen move this repository with it, in a reviewed change:")
		fmt.Fprintln(out, "  clue migrate          # preview the corpus and carrier upgrade")
		fmt.Fprintln(out, "  clue migrate --apply  # write it")
	case *quiet:
		// Current, or unable to tell. Neither is worth a line where the output
		// is a coding agent's context.
	case !report.Known:
		fmt.Fprintf(out, "clue %s — could not reach the release list; nothing to report\n", report.Current)
	case !isRelease(report.Current):
		fmt.Fprintf(out, "clue %s — the newest release is %s; a source build has no release to compare\n", report.Current, report.Latest)
	case !report.Comparable:
		// The host answered and the running stamp could not be read as a
		// release. Saying "you are current" here would be the one thing this
		// command exists to stop: a repository unable to tell, being reassured.
		fmt.Fprintf(out, "clue %s — the newest release is %s; this build's version cannot be compared with it\n", report.Current, report.Latest)
	default:
		fmt.Fprintf(out, "clue %s — this is the newest release\n", report.Current)
	}
	return 0
}

// isRelease reports whether a version stamp names a release at all. The same
// exemption the drift rule makes: a source build has nothing to compare.
func isRelease(v string) bool { return v != "" && v != "dev" }

// runVersion prints the release stamp (AC-019). "dev" for source builds.
func runVersion(w io.Writer) int {
	fmt.Fprintf(w, "clue %s\n", version)
	return 0
}

// agentConstraintCount is the visible promotion backlog of the convention
// register (AC-023): rules an agent must hold until clue can.
func agentConstraintCount(c *corpus.Corpus) int {
	n := 0
	for _, a := range c.Artifacts {
		if e, _ := a.Fields["enforcement"].(string); a.Type == "constraint" && e == "agent" {
			n++
		}
	}
	return n
}

// runRefs resolves the corpus's external addresses.
//
// ADR-040 keeps this outside the judge on purpose: it needs the network, so
// its answer can differ between two runs over the same revision. Only a gone
// reference is an error — restricted and unreachable say nothing about whether
// the corpus is correct, and treating them as failures would let someone
// else's outage, or the absence of a credential, block unrelated work.
func runRefs(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("refs", flag.ContinueOnError)
	fs.SetOutput(errOut)
	apply := fs.Bool("apply", false, "rewrite redirected addresses in place")
	timeout := fs.Duration("timeout", 0, "per-request budget (default 10s)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(errOut, "clue refs: expected at most one repository path")
		return 2
	}

	report, err := refs.Resolve(root, refs.Options{Apply: *apply, Timeout: *timeout})

	// The findings are printed before the error is handled. A rewrite that
	// fails part-way has already resolved every address, and discarding that
	// would throw away the whole network run — including any gone reference —
	// leaving the user with an error and no idea which files were touched.
	for _, r := range report.References {
		switch r.Result {
		case refs.Reachable:
			continue // the ordinary case is not worth a line
		case refs.Redirected:
			note := ""
			if r.Frozen {
				note = "  [pinned history — --apply leaves it alone; decide by hand]"
			}
			fmt.Fprintf(out, "%s:%d: redirected %s -> %s%s\n", r.Path, r.Line, r.URL, r.Target, note)
		case refs.Gone:
			fmt.Fprintf(out, "%s:%d: gone %s\n", r.Path, r.Line, r.URL)
		case refs.Restricted:
			fmt.Fprintf(out, "%s:%d: restricted %s\n", r.Path, r.Line, r.URL)
		case refs.Unreachable:
			detail := r.Detail
			if detail != "" {
				detail = " (" + detail + ")"
			}
			fmt.Fprintf(out, "%s:%d: unreachable %s%s\n", r.Path, r.Line, r.URL, detail)
		}
	}

	if err != nil {
		fmt.Fprintf(errOut, "clue refs: %v\n", err)
		return 2
	}

	// Redirects are split by what --apply can actually do with them. A redirect
	// inside pinned history is reported and left alone, so counting it towards
	// the advice would tell the user to run a command that will not touch it.
	frozen, rewritable := 0, 0
	for _, r := range report.References {
		if r.Result == refs.Redirected {
			if r.Frozen {
				frozen++
			} else {
				rewritable++
			}
		}
	}

	counts := report.Counts()
	mode := "preview"
	if report.Applied {
		mode = "applied"
	}
	fmt.Fprintf(out, "clue refs: %s — %d reachable, %d restricted, %d redirected, %d gone, %d unreachable\n",
		mode, counts[refs.Reachable], counts[refs.Restricted], counts[refs.Redirected],
		counts[refs.Gone], counts[refs.Unreachable])
	if frozen > 0 {
		fmt.Fprintf(out, "clue refs: %d redirected address(es) sit in pinned history and are never rewritten; decide those by hand\n", frozen)
	}
	if rewritable > 0 && !report.Applied {
		fmt.Fprintln(out, "clue refs: rerun with --apply to rewrite the redirected addresses")
	}
	if report.HasErrors() {
		return 1
	}
	return 0
}
