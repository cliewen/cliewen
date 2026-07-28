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
	"runtime/debug"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
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

const usage = `clue — a verifiable thread from goal to test

Usage:
  clue init [path]
  clue scaffold [path]
  clue context <id> [path]
  clue validate [--forbid-changes] [--coverage] [--reality-gaps] [path]
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

  validate   Scan docs/ and changes/ under path (default ".") and check
             the frontmatter graph: core fields, unique IDs, link
             resolution, status vocabularies, folder READMEs, index
             integrity, and skill version drift.

             --forbid-changes  fail when /changes contains files — the
                               digest-before-merge gate used by CI.
             --coverage        print derived proof coverage by capability.
             --reality-gaps    print capabilities contradicted by incident
                               analyses after their corpus was green.

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
	case "validate":
		os.Exit(runValidate(os.Args[2:]))
	case "version", "--version":
		os.Exit(runVersion(os.Stdout))
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "clue: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	forbid := fs.Bool("forbid-changes", false, "fail when /changes contains files")
	coverage := fs.Bool("coverage", false, "print derived per-capability proof coverage; never a committed registry")
	realityGaps := fs.Bool("reality-gaps", false, "print capabilities contradicted by incident analyses; never a committed registry")
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
			fmt.Println(is)
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
			fmt.Printf("%s: %s\n", cc.Capability, cc.State)
		}
	}
	if *realityGaps {
		for _, gap := range corpus.RealityGaps(c) {
			fmt.Printf("%s: contradicted by %s\n", gap.Capability, strings.Join(gap.Analyses, ", "))
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
	fmt.Printf("clue validate: OK (%d artifacts%s)\n", len(c.Artifacts), notes)
	return 0
}

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
