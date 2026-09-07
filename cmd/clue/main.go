// clue is the deterministic judge for a Cliewen corpus: stateless, no
// AI, no orchestration (Foundation §13). Commands are added only when a
// linter rule, a skill, or an onboarding criterion needs them.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/cliewen/cliewen/internal/carriers"
	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/ledger"
	"github.com/cliewen/cliewen/internal/migrate"
	"github.com/cliewen/cliewen/internal/parity"
	"github.com/cliewen/cliewen/internal/refs"
	"github.com/cliewen/cliewen/internal/release"
)

// version is the release stamp, injected at build time via
// `-ldflags "-X main.version=<semver>"` (see .github/workflows/release.yml).
// When no stamp is injected, main falls back to the module version Go
// embeds in `go install module@vX.Y.Z` builds; checkout and commit builds
// report "dev" and are exempt from skill-drift checks (ADR-011).
var version string

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

// resolvedVersion preserves every explicit linker stamp, including "dev".
// Only an absent stamp may take a Go module-version fallback.
func resolvedVersion(stamp, moduleVersion string) string {
	if stamp != "" {
		return stamp
	}
	if v := releaseFromModuleVersion(moduleVersion); v != "" {
		return v
	}
	return "dev"
}

const usage = `clue — a verifiable thread from goal to acceptance evidence

Usage:
  clue init [path]
  clue scaffold [path]
  clue context [--depth=<n>|all] [--stats] <id> [path]
  clue id next <prefix> [path]
  clue id live <id> [path]
  clue refs [--apply] [--timeout=<duration>] [path]
  clue carriers <inventory> [path]
  clue validate [--forbid-changes] [--coverage] [--reality-gaps] [--index-rows] [--read-cost] [--intent] [path]
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

             The architecture and design overview templates and the vision
             are marked bootstrap placeholders. Replace them with concise,
             repository-specific truth before clue validate can pass. The
             use-case folder is materialized empty and may stay that way:
             use cases are optional and absence is never a gap.

  scaffold   Regenerate the taxonomy README index blocks under path
             (default ".") from folder contents: entries whose targets
             survive keep their hand-written lines, missing entries are
             appended, prose outside the clue:index markers is never
             touched. Materializes nothing — missing folder READMEs are
             reported, and a path without a docs/ tree is an error.

  context    Print one artifact and its outgoing links, with complete markdown
             content in deterministic order. Acceptance-criterion and milestone
             IDs resolve to the artifact that declares them. Path defaults
             to ".". The slice is bounded so that a corpus which has grown
             dense still answers with a slice rather than with itself; what
             the bound held back is named, never silently dropped.
             --depth=<n>       follow n link hops beyond the root (default 1);
                               --depth=all follows every edge to exhaustion.
             --stats           print the slice's artifact and byte counts.
             Use cases whose links name the root are listed afterwards by
             identity, title, and path alone. Nothing is followed and no
             content is added, so the slice and its bound are unchanged;
             this is how a capability reaches the journey that governs it
             without the edge being written in two files.

  migrate    Preview a versioned corpus and managed-carrier migration; use
             --apply to write the complete safe plan and
             --reversal-cost=low|high to resolve missing inferred-meaning routing.
             Existing prose and locally modified generated files are never
             overwritten. Path defaults to ".".

  id next    Allocate the next numeric ID for a prefix through the
             persisted identity ledger (.clue/id-ledger.yaml): an O(1)
             counter increment, never a corpus scan. Prints the new
             ID and reserves it in the ledger. Path defaults to ".".

  id live    Mark a previously reserved ID as live after its artifact
             has been created. Refuses an ID that is not reserved.
             Path defaults to ".".

  refs       Resolve the external addresses docs/ and changes/ point at,
             classifying each: reachable, restricted (it exists, this
             runner may not read it), redirected, gone, or unreachable.
             Only "gone" is an error; an outage elsewhere never
             condemns a corpus. Use
             --apply to rewrite redirected addresses in place. A clue:
             identity is never followed. Never make this a required
             check: another host's uptime must not gate a merge.
             Path defaults to ".".

  parity     Compare a pinned source manifest against the target manifest
             derived from path's corpus and ledger (default "."), reporting
             every missing criterion, orphaned tag, changed direction or
             evidence location, stale source fingerprint, unjustified
             @draft/Human/retirement disposition, and a disposition whose
             plan door is absent from the target corpus or already claimed
             by another deferred criterion. The report is derived and never
             edited by hand; a clean run is the only passing result.

             --out=<path>  also write the report to path, for a migration
                           workflow to upload as a CI artifact.

  report     Render the derived region of the extraction report at the
             given path (required) from the pinned source manifest the
             region's opening marker names, leaving every byte outside the
             region alone. A report states its criterion counts and
             mapping table only there, and validate re-renders the region
             and fails on any difference, so a typed figure is not a
             supported way to write one. An optional second path is the
             repository the manifest is named relative to, default ".".

  carriers   Reconcile a pinned operational-carrier inventory against
             path's corpus (default "."), reporting a stale reference to a
             path the inventory names as deleted, a lost fingerprint on a
             mapped target that drifted after the rehearsal recorded it as
             current, and a missing mapped target. The report is derived
             and never edited by hand; a clean run is the only passing
             result. A blocked entry has no target and is not reconciled —
             its presence is the record of a known gap.

             --out=<path>  also write the report to path, for a migration
                           workflow to upload as a CI artifact.

  validate   Scan docs/ and changes/ under path (default ".") and check
             the frontmatter graph: core fields, unique IDs, link
             resolution, status vocabularies, folder READMEs, index
             integrity, activated architecture and design overviews, skill
             version drift, and every derived extraction
             report region against the manifest it names.

             --forbid-changes  fail when /changes contains files — the
                               digest-before-merge gate used by CI.
             --coverage        print derived proof coverage by capability,
                               then any pointer to proof in another
                               repository, listed apart as named but
                               locally unproven.
             --reality-gaps    print capabilities contradicted by incident
                               analyses after their corpus was green.
             --index-rows      print index rows that only restate their own
                               link, say nothing about the artifact, or give
                               a constraint a badge other than enforcement.
             --read-cost       print multi-document artifacts and identities
                               whose default context slice exceeds the budget.
             --intent          print the direction this corpus states, its
                               provenance, and the use cases it has written
                               down. A state report: a missing vision is a
                               state and not an issue, and no coverage
                               figure is computed, because use cases are
                               optional and a ratio over an optional
                               artifact reads as a target.

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

Release notice: init, scaffold, context, migrate, id, refs, and report print one line to
standard error when a newer release exists. Standard output and the exit code
are identical with it and without it, and the answer is cached for a day. Never
from validate or version, never when CI carries a value, and never when
CLUE_NO_UPDATE_NOTIFIER is set at all, the empty string included.

Exit codes: 0 corpus valid · 1 issues found · 2 usage error
`

func main() {
	moduleVersion := ""
	if bi, ok := debug.ReadBuildInfo(); ok {
		moduleVersion = bi.Main.Version
	}
	version = resolvedVersion(version, moduleVersion)
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	command := os.Args[1]
	code := run(command, os.Args[2:])
	emitNotice(command, os.Stderr, os.LookupEnv, runtimeCheck())
	os.Exit(code)
}

// run dispatches one command and returns its exit code. It exists as its own
// function so that main has a place to stand after the command has finished —
// the ambient release notice is printed there, and a command that exits from
// inside its own case leaves no such place.
func run(command string, args []string) int {
	switch command {
	case "init":
		return runInit(args, os.Stdout, os.Stderr)
	case "scaffold":
		return runScaffold(args, os.Stdout, os.Stderr)
	case "context":
		return runContext(args, os.Stdout, os.Stderr)
	case "migrate":
		return runMigrate(args, os.Stdout, os.Stderr)
	case "id":
		return runID(args, os.Stdout, os.Stderr)
	case "refs":
		return runRefs(args, os.Stdout, os.Stderr)
	case "parity":
		return runParity(args, os.Stdout, os.Stderr)
	case "report":
		return runReport(args, os.Stdout, os.Stderr)
	case "carriers":
		return runCarriers(args, os.Stdout, os.Stderr)
	case "validate":
		return runValidate(args, os.Stdout)
	case "latest":
		return runLatest(args, os.Stdout, os.Stderr, runtimeCheck())
	case "version", "--version":
		return runVersion(os.Stdout)
	case "help", "--help", "-h":
		fmt.Print(usage)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "clue: unknown command %q\n\n%s", command, usage)
		return 2
	}
}

// notifierCommands are the workflow commands that may carry an ambient
// release notice (PDR-023). The exclusions are the point of the list.
//
// `validate`, `parity`, and `carriers` are absent because all three are
// deterministic judges: their verdict, exit code, and output stay a
// statement about the inputs alone, and a line that depends on another
// system's present state is not that. `version` is absent because it is the
// one command guaranteed to answer instantly, offline, and identically
// forever (ADR-042). `latest` is absent because it is the check, and would
// otherwise report twice.
var notifierCommands = map[string]bool{
	"init":     true,
	"scaffold": true,
	"context":  true,
	"migrate":  true,
	"id":       true,
	"refs":     true,
	"report":   true,
}

// emitNotice writes the ambient release notice, when one is allowed and there
// is one to write (PDR-023). It runs after the command it rode in on, so no
// command's output is interleaved with it, and it writes only to the error
// stream — the exit code and standard output of every command are identical
// with the notice and without it.
//
// Everything it consults is a parameter, so the whole path from a command name
// to what lands on the stream is provable without a particular machine's
// environment.
func emitNotice(command string, errOut io.Writer, env func(string) (string, bool), opts release.Options) {
	if !notifierAllowed(command, env) {
		return
	}
	if line := release.Notice(opts); line != "" {
		fmt.Fprintln(errOut, line)
	}
}

// notifierAllowed reports whether an ambient release notice may be printed.
// Every input is a parameter so the gate is provable without a clock and
// without a particular machine's environment.
//
// There is deliberately no terminal condition. The first implementation had
// one, and it excluded the audience the notice exists for: a coding agent
// captures a command's output through a pipe rather than a terminal — reading
// the output is the point — so the notice was switched off in exactly the
// sessions PDR-023 built it to reach, leaving only the prose fallback that had
// already been observed failing. `CI` carries the ephemeral-runner rejection on
// its own, and what the terminal check uniquely bought is narrower than it
// looked: a script consumes standard output and an exit code, and both are
// byte-identical with the notice and without it.
func notifierAllowed(command string, env func(string) (string, bool)) bool {
	if !notifierCommands[command] {
		return false
	}
	// Set by every CI system worth naming. An ephemeral runner never reaches
	// the release list, so no adopter's push spends a request on a fact no one
	// is there to read (ADR-042). Its value decides, because "CI=" is the
	// conventional way a shell says "not CI" and reading that as a runner
	// would silence a developer's own machine.
	if v, _ := env("CI"); v != "" {
		return false
	}
	// The documented opt-out, for anyone who wants the tool to say nothing it
	// was not asked to say. Presence decides, not value: the documentation
	// offers it as a switch you set, and "CLUE_NO_UPDATE_NOTIFIER=" is exactly
	// what a wrapper script or a templated CI matrix produces. A switch that
	// fails in the spelling its own documentation invites is not a switch.
	_, optedOut := env("CLUE_NO_UPDATE_NOTIFIER")
	return !optedOut
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

// runID dispatches identity-ledger subcommands. The subcommand shape leaves
// room for future "id retire" or "id lookup" operations without a new
// top-level command (ADR-048).
func runID(args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errOut, "clue id: expected a subcommand (next or live)")
		return 2
	}
	switch args[0] {
	case "next":
		return runIDNext(args[1:], out, errOut)
	case "live":
		return runIDLive(args[1:], out, errOut)
	default:
		fmt.Fprintf(errOut, "clue id: unknown subcommand %q\n", args[0])
		return 2
	}
}

// runIDLive promotes an allocator-reserved ID once its artifact has been
// created, preserving the ledger's reserved -> live state transition.
func runIDLive(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("id live", flag.ContinueOnError)
	fs.SetOutput(errOut)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(errOut, "clue id live: expected an ID")
		return 2
	}
	id := fs.Arg(0)
	root := "."
	if fs.NArg() > 1 {
		root = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		fmt.Fprintln(errOut, "clue id live: expected an ID and at most one repository path")
		return 2
	}

	l, err := ledger.Load(root)
	if err != nil {
		fmt.Fprintf(errOut, "clue id live: %v\n", err)
		return 2
	}
	if err := l.PromoteReserved(id); err != nil {
		fmt.Fprintf(errOut, "clue id live: %v\n", err)
		return 2
	}
	if err := l.Save(); err != nil {
		fmt.Fprintf(errOut, "clue id live: %v\n", err)
		return 2
	}
	return 0
}

// runIDNext allocates the next numeric ID for a prefix through the ledger:
// an O(1) counter increment, never a corpus scan (ADR-048, AC-101).
func runIDNext(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("id next", flag.ContinueOnError)
	fs.SetOutput(errOut)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(errOut, "clue id next: expected a prefix")
		return 2
	}
	prefix := fs.Arg(0)
	root := "."
	if fs.NArg() > 1 {
		root = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		fmt.Fprintln(errOut, "clue id next: expected a prefix and at most one repository path")
		return 2
	}

	l, err := ledger.Load(root)
	if err != nil {
		fmt.Fprintf(errOut, "clue id next: %v\n", err)
		return 2
	}
	if !ledger.Exists(root) {
		fmt.Fprintln(errOut, "clue id next: identity ledger is missing; run `clue migrate --apply` first")
		return 2
	}
	id, err := l.NextNumeric(prefix)
	if err != nil {
		fmt.Fprintf(errOut, "clue id next: %v\n", err)
		return 2
	}
	if err := l.Save(); err != nil {
		fmt.Fprintf(errOut, "clue id next: %v\n", err)
		return 2
	}
	fmt.Fprintln(out, id)
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
	indexRows := fs.Bool("index-rows", false, "print index rows that only restate their own link, say nothing about the artifact, or mismatch constraint enforcement; never a committed registry")
	readCost := fs.Bool("read-cost", false, "print multi-document artifacts and bounded context slices over the read-cost budget; never a committed registry")
	intent := fs.Bool("intent", false, "print the direction this corpus states and the use cases it has written down; a state report, never a coverage figure")
	_ = fs.Parse(args)
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}

	c, issues := corpus.Scan(root)
	issues = append(issues, corpus.Validate(c, corpus.Options{ForbidChanges: *forbid, Version: version})...)
	// A derived extraction report is checked against the manifest it names,
	// not against the corpus graph, so the check lives beside the parity
	// contract it renders (ADR-054) rather than inside the graph rules.
	issues = append(issues, parity.CheckReports(root)...)
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
	undescribed := corpus.IndexDescriptionBacklog(c)
	constraintBadges := corpus.IndexRowConstraintBadgeBacklog(c)
	multiDocument := corpus.MultiDocumentBacklog(c)
	overBudget := corpus.ContextSliceBudgetBacklog(c)
	if *indexRows {
		for _, row := range fillerRows {
			fmt.Fprintf(out, "%s: %s states only its own link\n", row.Readme, row.Target)
		}
		for _, row := range undescribed {
			fmt.Fprintf(out, "%s: %s states its record but not what it is about\n", row.Readme, row.Target)
		}
		for _, row := range constraintBadges {
			fmt.Fprintf(out, "%s: %s has index badge %q, not enforcement %q\n", row.Readme, row.Target, row.Badge, row.Enforcement)
		}
	}
	// A state report, never a scorecard. It prints no ratio of goals or
	// capabilities with use cases, because a percentage over an optional
	// artifact reads as a target and the only way to move it is to write
	// artifacts nobody needs (PDR-054).
	if *intent {
		state := corpus.Intent(c)
		if state.Vision.Present {
			meaning := "human-authored"
			if state.Vision.Inferred {
				meaning = "inferred — no human has confirmed it"
			}
			fmt.Fprintf(out, "vision: %s %s (%s, %s) at %s\n", state.Vision.ID, state.Vision.Title, state.Vision.Status, meaning, state.Vision.Path)
		} else {
			fmt.Fprintf(out, "vision: none — this corpus states no direction (%s is absent)\n", corpus.VisionPath)
		}
		for _, use := range state.UseCases {
			crosses := "no capability"
			if len(use.Capabilities) > 0 {
				crosses = strings.Join(use.Capabilities, ", ")
			}
			fmt.Fprintf(out, "use case: %s %s (%s) crosses %s\n", use.ID, use.Title, use.Status, crosses)
		}
		if len(state.UseCases) == 0 {
			fmt.Fprintln(out, "use cases: none — they are optional, and absence is not a gap")
		}
	}
	if *readCost {
		for _, item := range multiDocument {
			fmt.Fprintf(out, "%s: %d rendered document(s)\n", item.Artifact.Path, item.Documents)
		}
		for _, item := range overBudget {
			fmt.Fprintf(out, "%s: default context slice prints %d artifact(s) (budget %d)\n", item.Identity, item.Artifacts, corpus.DefaultContextSliceBudget)
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
	if n := len(undescribed); n > 0 {
		notes += fmt.Sprintf(", %d index row(s) not saying what the artifact is about", n)
	}
	if n := len(constraintBadges); n > 0 {
		notes += fmt.Sprintf(", %d constraint index badge(s) not stating enforcement", n)
	}
	if n := len(multiDocument); n > 0 {
		notes += fmt.Sprintf(", %d multi-document artifact(s)", n)
	}
	if n := len(overBudget); n > 0 {
		notes += fmt.Sprintf(", %d identity slice(s) over the %d-artifact budget", n, corpus.DefaultContextSliceBudget)
	}
	fmt.Fprintf(out, "clue validate: OK (%d artifacts%s)\n", len(c.Artifacts), notes)
	return 0
}

// runParity compares a pinned source manifest against the target manifest
// clue derives from path's corpus and ledger (ADR-049), printing every
// finding and exiting non-zero on any required parity failure class.
func runParity(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("parity", flag.ContinueOnError)
	fs.SetOutput(errOut)
	outPath := fs.String("out", "", "also write the report to this path, for a CI artifact")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(errOut, "clue parity: expected a source manifest path")
		return 2
	}
	root := "."
	if fs.NArg() > 1 {
		root = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		fmt.Fprintln(errOut, "clue parity: expected a source manifest path and at most one repository path")
		return 2
	}

	source, err := parity.LoadSourceManifest(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(errOut, "clue parity: %v\n", err)
		return 2
	}
	target, err := parity.DeriveTargetManifest(root)
	if err != nil {
		fmt.Fprintf(errOut, "clue parity: %v\n", err)
		return 2
	}
	report := parity.Compare(source, target)

	var buf strings.Builder
	for _, f := range report.Findings {
		fmt.Fprintf(&buf, "%s %s: %s\n", f.Class, f.ID, f.Detail)
	}
	if !report.Failed() {
		fmt.Fprintf(&buf, "clue parity: OK (%d source entries; %d deferred criterion(s))\n", len(source.Entries), report.Deferred)
	} else {
		fmt.Fprintf(&buf, "clue parity: %d finding(s); %d deferred criterion(s)\n", len(report.Findings), report.Deferred)
	}
	fmt.Fprint(out, buf.String())
	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(buf.String()), 0o644); err != nil {
			fmt.Fprintf(errOut, "clue parity: %v\n", err)
			return 2
		}
	}
	if report.Failed() {
		return 1
	}
	return 0
}

// runReport renders a durable extraction report's derived region from the
// pinned source manifest the region names (ADR-054). It writes only inside
// the region: the report's prose is the author's, its figures are generated.
func runReport(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(errOut)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(errOut, "clue report: expected an extraction report path")
		return 2
	}
	root := "."
	if fs.NArg() > 1 {
		root = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		fmt.Fprintln(errOut, "clue report: expected a report path and at most one repository path")
		return 2
	}
	if err := parity.WriteReport(root, fs.Arg(0)); err != nil {
		fmt.Fprintf(errOut, "clue report: %v\n", err)
		return 2
	}
	fmt.Fprintf(out, "clue report: %s derived from its manifest\n", filepath.ToSlash(fs.Arg(0)))
	return 0
}

// runCarriers reconciles a pinned operational-carrier inventory against the
// current corpus (ADR-051), printing every finding and exiting non-zero on
// any of the three required failure classes.
func runCarriers(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("carriers", flag.ContinueOnError)
	fs.SetOutput(errOut)
	outPath := fs.String("out", "", "also write the report to this path, for a CI artifact")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(errOut, "clue carriers: expected a carrier inventory path")
		return 2
	}
	root := "."
	if fs.NArg() > 1 {
		root = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		fmt.Fprintln(errOut, "clue carriers: expected a carrier inventory path and at most one repository path")
		return 2
	}

	inv, err := carriers.LoadInventory(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(errOut, "clue carriers: %v\n", err)
		return 2
	}
	report, err := carriers.Reconcile(inv, root)
	if err != nil {
		fmt.Fprintf(errOut, "clue carriers: %v\n", err)
		return 2
	}

	var buf strings.Builder
	for _, f := range report.Findings {
		fmt.Fprintf(&buf, "%s %s: %s\n", f.Class, f.ID, f.Detail)
	}
	if !report.Failed() {
		fmt.Fprintf(&buf, "clue carriers: OK (%d entries)\n", len(inv.Entries))
	} else {
		fmt.Fprintf(&buf, "clue carriers: %d finding(s)\n", len(report.Findings))
	}
	fmt.Fprint(out, buf.String())
	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(buf.String()), 0o644); err != nil {
			fmt.Fprintf(errOut, "clue carriers: %v\n", err)
			return 2
		}
	}
	if report.Failed() {
		return 1
	}
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
			fmt.Fprintln(out, release.QuietLine(report))
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
	case report.SourceBuild:
		fmt.Fprintf(out, "clue %s — the newest release is %s; a source build has no release to compare\n", report.Current, report.Latest)
	case !report.Comparable:
		// The host answered and the running stamp could not be read as a
		// release. Saying "you are current" here would be the one thing this
		// command exists to stop: a repository unable to tell, being reassured.
		fmt.Fprintf(out, "clue %s — the newest release is %s; this build's version cannot be compared with it\n", report.Current, report.Latest)
	case report.Ahead:
		// An unreleased build: a release candidate, or a local tag ahead of the
		// list. Nothing to catch up with, but it is not the newest release
		// either, and saying it is would state a fact the list never gave.
		fmt.Fprintf(out, "clue %s — nothing newer is published; the newest release is %s\n", report.Current, report.Latest)
	default:
		fmt.Fprintf(out, "clue %s — this is the newest release\n", report.Current)
	}
	return 0
}

// runVersion prints the release stamp (AC-019). "dev" for source builds.
func runVersion(w io.Writer) int {
	fmt.Fprintf(w, "clue %s\n", version)
	return 0
}

// agentConstraintCount is the visible promotion backlog of the convention
// register (AC-089): rules an agent must hold until a machine can. It counts
// `agent` and nothing else — `partial` and `human` are declared rather than
// queued, and a rule leaves the backlog by gaining a check or by stating what
// cannot be mechanized, never by relabelling (ADR-045).
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
