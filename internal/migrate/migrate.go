// Package migrate upgrades mechanical Cliewen corpus and managed-carrier
// changes without treating init as an updater. It plans every write first,
// refuses ambiguous or locally modified inputs, and applies the complete plan
// only after all preconditions hold.
package migrate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cliewen/cliewen/internal/corpus"
	"github.com/cliewen/cliewen/internal/scaffold"
	"gopkg.in/yaml.v3"
)

const (
	// MigrationReversalCost adds the required ADR-035 field when the adopter
	// explicitly supplies the semantic routing choice.
	MigrationReversalCost = "MIG-001"
	// MigrationStatusLifecycle maps the old architecture/analysis status to
	// the current default lifecycle.
	MigrationStatusLifecycle = "MIG-002"
	// MigrationManagedCarriers refreshes generated skills and the thin caller.
	MigrationManagedCarriers = "MIG-003"
	// MigrationQualifiedReferences reports bare forge references. It never
	// repairs one: nothing in the file says which repository was meant, and
	// defaulting to the adopter's own slug would convert a confidently wrong
	// reference into a confidently wrong qualified one that no later check
	// can question.
	MigrationQualifiedReferences = "MIG-004"
	// MigrationClaudeEntryPoint reports a Claude Code entry point that never
	// reaches AGENTS.md. It never repairs one: the missing case is already
	// solved by re-running the non-destructive init, and an entry point the
	// adopter wrote themselves is their prose (PDR-022).
	MigrationClaudeEntryPoint = "MIG-005"
)

// Options controls planning. Preview is the default; applying a plan is a
// separate command choice so a caller cannot accidentally turn a report into
// writes.
type Options struct {
	ReversalCost string
}

// MigrationDefinition is one ordered, append-only mechanical migration in
// the supported registry.
type MigrationDefinition struct {
	ID          string
	Description string
}

var orderedMigrations = []MigrationDefinition{
	{ID: MigrationReversalCost, Description: "add explicit inferred-meaning reversal routing"},
	{ID: MigrationStatusLifecycle, Description: "map historical architecture and analysis status to the default lifecycle"},
	{ID: MigrationManagedCarriers, Description: "refresh generated skills, mirrors, and the thin CI caller"},
	{ID: MigrationQualifiedReferences, Description: "report external references that name no repository"},
	{ID: MigrationClaudeEntryPoint, Description: "report a Claude Code entry point that never reaches the routing hub"},
}

// Registry returns the migration order without exposing mutable package state.
func Registry() []MigrationDefinition {
	return append([]MigrationDefinition(nil), orderedMigrations...)
}

// Change is one file replacement in a complete migration plan.
type Change struct {
	Path        string
	Migration   string
	Description string
	Existed     bool
	Before      []byte
	After       []byte
}

// Finding is a blocking condition. A plan with any finding is never applied.
type Finding struct {
	Path      string
	Migration string
	Message   string
}

// Notice is an explicitly non-blocking boundary, such as a missing Claude
// mirror that init never materialized or a symlinked mirror outside the repo.
type Notice struct {
	Path      string
	Migration string
	Message   string
}

// MigrationPlan is the complete preflight result. Changes contain exact before/after
// bytes; callers must not write a subset of them.
type MigrationPlan struct {
	Changes  []Change
	Findings []Finding
	Notices  []Notice
	Target   string
}

var (
	fieldLineRe     = regexp.MustCompile(`(?m)^([ \t]*)([A-Za-z][A-Za-z0-9_-]*):([ \t]*)([^\r\n]*)(\r?\n|$)`)
	callerUsesRe    = regexp.MustCompile(`(?m)^([ \t]*)uses:([ \t]*)(cliewen/cliewen/\.github/workflows/clue-validation\.yml)@([^\s#]+)([ \t]*(?:#.*)?)(\r?\n|$)`)
	callerVersionRe = regexp.MustCompile(`(?m)^([ \t]*)clue-version:([ \t]*)([^\s#]+)([ \t]*(?:#.*)?)(\r?\n|$)`)
)

// releaseManifest is one published release's generated-carrier digests.
type releaseManifest struct {
	Version string
	Files   map[string]string
}

// legacyDigests is the committed release manifest for generated carriers.
// A carrier is replaceable only when its bytes match the release that wrote
// it. This is what distinguishes an old generated file from a locally edited
// one without a hidden state file or a destructive overwrite.
//
// It is an ordered, append-only list rather than a map: releases are reported
// oldest-first, and ordering them by comparison would need semver parsing that
// a lexicographic sort gets wrong the moment a minor reaches two digits
// (0.10.0 before 0.9.0). Append each new release at the end.
var legacyDigests = []releaseManifest{
	{Version: "0.5.1", Files: map[string]string{
		"clue-analysis/skill.md":            "9dc9cbf5482dfce9ba78bcc47eae205ab7f65a29061c5614580fbcfb704d430f",
		"clue-delta/skill.md":               "a6ede82283531f71edb53773a7cad477ceae1a7cd97957d74be0f12dec4a9ddc",
		"clue-extract/mappings/openspec.md": "a6d40a7d9b213d58c935b1d343dcb27e3efacde7b1e8a80745c60e724bbc68e4",
		"clue-extract/skill.md":             "98730f6edc2f2170ce1179d3f907ab127db9baab4b1ddedc68718deac1411bd1",
		"clue-plan/skill.md":                "b23af904b43543cfe2616fc7273b9f630105230983a844de415c393024acb048",
		"clue-verify/skill.md":              "788a1d6219248f8f8cd169d2858f20617aae036a7f56c06099f199c8e0776b76",
	}},
	{Version: "0.6.0", Files: map[string]string{
		"clue-analysis/skill.md":            "fb97a73b58a466b45d8d47f9cb3a7b3e54d73c9c93877e23637e827fe7bdb6e7",
		"clue-delta/skill.md":               "322c9d3bbce9580ba2b0091a2bbf4743eea3c283f5c679af7c4eb7a122ed364f",
		"clue-extract/mappings/openspec.md": "3673fe7bb66237d8d73c42158581c9ed33c5c53cb3aa89bdbde6312151dd37ea",
		"clue-extract/skill.md":             "d197867a2a33a5adf8d827c9823313a111ff5d19660e9eb7f6d237cdf9cfc1f3",
		"clue-plan/skill.md":                "deb6221dd6b612eac16c8d36acdc8f93807bc02c2bf1c3827f5a1e45e6dba2d7",
		"clue-verify/skill.md":              "e2e078899834b9e5719e300aedb80dec470d2675c6188927f77c69af57e94a34",
	}},
	{Version: "0.7.0", Files: map[string]string{
		"clue-analysis/skill.md":            "6194cb804f08133ef3b3e87c39a7616bf9cb2115fe759598e796ba66df967dd4",
		"clue-delta/skill.md":               "e336c2df69fc1be9d5aeebe81b3989dae65187703172ff283120c5aff76e73bc",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "97d0aaa79644bea5bd6bf6000354d7e66d0ca005eb7fe936b672d961001cfc38",
		"clue-extract/skill.md":             "960a77ba250367bca319e708272b78fdeed01683adef5c83c55d6d0f8b4d584e",
		"clue-plan/skill.md":                "0a2b4e3a7c8dbe63247087e08942b3211c42e7d6693de4751fa842a5ebc87497",
		"clue-verify/skill.md":              "faa1b050e5a440d00ee7cff110aca4e1516507f7c8b035ac2e4b36d51f068bed",
	}},
	{Version: "0.8.0", Files: map[string]string{
		"clue-analysis/skill.md":            "8a4967df0c74e6fe2b3ef6f857a925d8798ed16152d3ebf503bc64abaab34c5b",
		"clue-delta/skill.md":               "2fd3516de8b8ab9bc7f1fe24e7fca42df5c33b6fbed8c3ce08fb1317e87c347c",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "97d0aaa79644bea5bd6bf6000354d7e66d0ca005eb7fe936b672d961001cfc38",
		"clue-extract/skill.md":             "070b5bf6ff4fd6159166d645d01693de0404e42124cbbcc5dcfd14139ba3af64",
		"clue-plan/skill.md":                "af8dee537456644e99607a702da136467f7fd8967bea950fe8b5a3e5bf6f8661",
		"clue-verify/skill.md":              "44fc206073addf90763def2ca682a86ddbf10395e0dfa72332bda56da094bfec",
	}},
	{Version: "0.9.0", Files: map[string]string{
		"clue-analysis/skill.md":            "c472b39311b43cd2c2c8bba1c7de9e1361bf2465de3dc3f6cf3879c5fcf0377f",
		"clue-delta/skill.md":               "a229f48814099930a7df5781b1302bbc97bb602f2179e56f5942218a8b518145",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "d58fc65df6f240197e291cd48b250f11676e6d9992147934cd87800ccd58905d",
		"clue-extract/skill.md":             "ef0bc657288a1d5adec5d8980a8c705aff4c9f964d9eb27fd62987bd9cabadb0",
		"clue-plan/skill.md":                "bf6a4e4cd1fa5b6ed8b318687f5da51c0837415e66ae9e48c1c34ff4b2c7e324",
		"clue-verify/skill.md":              "23a409a732255da58412785f7b06fe52923966fe16bcd23250369faf09c920b6",
	}},
	{Version: "0.10.0", Files: map[string]string{
		"clue-analysis/skill.md":            "ba8034351ddd23cd798d3bd024f8bb96dc04ee144cd58c9f0ad0b16ec832292d",
		"clue-delta/skill.md":               "3a73f16e42caf9f232abe435374acc87bf28fb9f0039849e2dd61f1fb30503de",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "d58fc65df6f240197e291cd48b250f11676e6d9992147934cd87800ccd58905d",
		"clue-extract/skill.md":             "814ca9d06a9d10630593603f938ff6c0ee631ccafb1366ddd9b024e731c0be45",
		"clue-plan/skill.md":                "2b498fba3343f1fe336291b56776fa3467b397a15e6704608a9c5946d52dc7f3",
		"clue-verify/skill.md":              "26328abab11cc271b1dc88bf05660d4ad523ae1dddc2d746874794902f9575b6",
	}},
	{Version: "0.11.0", Files: map[string]string{
		"clue-analysis/skill.md":            "53f2eb7297fb250a70299a851998cd3177965a7d5596f135ea648130e01de3c1",
		"clue-delta/skill.md":               "95b34bfd7c4b80e145fa63a5522688264b02be4849dd9fb5550a825f510f6b2e",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "1144b08bad8e752d8286efd6f6a0085f179729c728bc7ec56e19b0f5694578e9",
		"clue-extract/skill.md":             "b9b4698b7b143c8f386ce61712c712df71a5e299a906be1831a7cba6504d6e27",
		"clue-plan/skill.md":                "2130fc3367272b608644a454b3d1347f6efe8e437b072fa5a46e4ec9c52a5e40",
		"clue-verify/skill.md":              "2e03bd7a0ba62137d08a4c5251dc6bc113c2e52e1202b5f6405f26063a32f758",
	}},
	// 0.11.1 is deliberately absent. Its stamp merged, but the tagging gate
	// rejected the release for missing migration guidance, so no tag and no
	// artifacts were ever published. No repository can be carrying 0.11.1
	// carriers, and listing digests for a release that does not exist would
	// invite a later reader to trust a version they cannot install.
	{Version: "0.11.2", Files: map[string]string{
		"clue-analysis/skill.md":            "d1683922661b8ca9618110184f7bd7722f1af34a2e1e034fd74214c20033f9e1",
		"clue-delta/skill.md":               "fecad3c155534b07ca22d16cfcf1a6c25fa9b0aacf9d29545d80f705254f76b7",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "1144b08bad8e752d8286efd6f6a0085f179729c728bc7ec56e19b0f5694578e9",
		"clue-extract/skill.md":             "d65bef335a7bd4bc40e86931c2ed757ddf908166c80616613f51f1a01d32c7b0",
		"clue-plan/skill.md":                "c10aa23feeb26f96bcc3d8698a6780e10ecd29e4e83d775b44b341133168ce46",
		"clue-verify/skill.md":              "0d5e3b4e627ec5891b552b5602af5731b0d3a582bf19fc828ea32c0f543c9543",
	}},
}

// Plan scans the target and returns the complete deterministic migration.
// It performs no writes.
func Plan(root string, opts Options) (MigrationPlan, error) {
	target, err := scaffold.PairVersion()
	if err != nil {
		return MigrationPlan{}, err
	}
	result := MigrationPlan{Target: target}
	if opts.ReversalCost != "" && opts.ReversalCost != "low" && opts.ReversalCost != "high" {
		result.Findings = append(result.Findings, Finding{Migration: MigrationReversalCost, Message: "reversal-cost must be low or high"})
		return result, nil
	}
	if err := planCorpus(root, opts, &result); err != nil {
		return MigrationPlan{}, err
	}
	carriers, err := scaffold.ManagedCarrierFiles()
	if err != nil {
		return MigrationPlan{}, err
	}
	planCarriers(root, target, carriers, &result)
	planQualifiedReferences(root, &result)
	planClaudeEntryPoint(root, &result)
	sortPlan(&result)
	return result, nil
}

// Apply writes a complete, previously planned migration. It rechecks every
// source byte immediately before writing so a changed target cannot be
// overwritten after preview. Any preflight finding is a hard stop.
func Apply(root string, plan MigrationPlan) error {
	if len(plan.Findings) > 0 {
		return errors.New("migration has unresolved findings; no files were changed")
	}
	for _, change := range plan.Changes {
		full := filepath.Join(root, filepath.FromSlash(change.Path))
		if hasLinkBoundary(root, change.Path) {
			return fmt.Errorf("%s: path is symlinked; no files were changed", change.Path)
		}
		got, err := os.ReadFile(full)
		if !change.Existed && os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("%s: read before apply: %w; no files were changed", change.Path, err)
		}
		if !change.Existed {
			return fmt.Errorf("%s: appeared after planning; rerun preview; no files were changed", change.Path)
		}
		if !bytes.Equal(got, change.Before) {
			return fmt.Errorf("%s: changed after planning; rerun preview; no files were changed", change.Path)
		}
	}
	var applied []appliedFile
	for _, change := range plan.Changes {
		path := filepath.Join(root, filepath.FromSlash(change.Path))
		entry, err := applyFile(path, change.After, change.Existed)
		if err != nil {
			rollbackErrors := rollback(applied)
			if len(rollbackErrors) > 0 {
				return fmt.Errorf("%s: apply failed: %w; rollback also failed: %s", change.Path, err, strings.Join(rollbackErrors, "; "))
			}
			return fmt.Errorf("%s: apply failed: %w; earlier writes were rolled back", change.Path, err)
		}
		applied = append(applied, entry)
	}
	for _, entry := range applied {
		if entry.backup != "" {
			_ = os.Remove(entry.backup)
		}
	}
	return nil
}

func planCorpus(root string, opts Options, result *MigrationPlan) error {
	docs := filepath.Join(root, "docs")
	info, err := os.Stat(docs)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("no docs tree under %q — migration applies to an existing Cliewen corpus", root)
	}
	return filepath.WalkDir(docs, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if hasLinkBoundary(root, rel) {
			// Classify before reporting, so a prose README behind a link does
			// not block the plan. An unreadable one is exactly what the
			// finding is for, so a failed read reports rather than skips.
			data, readErr := os.ReadFile(filePath)
			if readErr != nil || isArtifactFile(entry.Name(), data) {
				result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "corpus artifact is behind a symlink; resolve the repository-owned path before migrating"})
			}
			return nil
		}
		before, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		if !isArtifactFile(entry.Name(), before) {
			return nil
		}
		after, changes, findings := migrateArtifact(rel, before, opts.ReversalCost)
		result.Findings = append(result.Findings, findings...)
		if !bytes.Equal(before, after) {
			result.Changes = append(result.Changes, Change{Path: rel, Migration: migrationForArtifact(changes), Description: strings.Join(changes, "; "), Existed: true, Before: before, After: after})
		}
		return nil
	})
}

// isArtifactFile reports whether a walked docs file participates in the
// migration. Only READMEs need the test: a folder README is prose, but a
// capability's README carries frontmatter and is as much an artifact as any
// other file the validator judges. corpus.Scan draws the line at a complete
// frontmatter block — an unclosed one counts as prose there — so drawing it
// the same way here keeps the migration from reporting a file that
// clue validate accepts, and from silently skipping one it rejects.
func isArtifactFile(name string, data []byte) bool {
	if name != "README.md" {
		return true
	}
	_, _, ok, err := splitFrontmatter(string(data))
	return ok && err == nil
}

func migrateArtifact(rel string, before []byte, reversalCost string) ([]byte, []string, []Finding) {
	text := string(before)
	end, fields, ok, err := readFrontmatter(text)
	if err != nil {
		return before, nil, []Finding{{Path: rel, Migration: MigrationReversalCost, Message: "frontmatter is ambiguous; resolve it before migration: " + err.Error()}}
	}
	if !ok {
		return before, nil, nil
	}
	var changes []string
	var findings []Finding
	front := text[:end]
	if fields.Type == "decision" && fields.HasReversalCost {
		findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "decisions route cost by record type and must not carry reversal-cost; resolve it by hand before resuming"})
	} else if fields.HasReversalCost && fields.ReversalCost != "low" && fields.ReversalCost != "high" {
		findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "reversal-cost is not low or high; resolve its semantic classification by hand before resuming"})
	}
	if fields.Type == "decision" && fields.HasProvenance {
		findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "decisions carry provenance in status, not in a provenance field; resolve it by hand before resuming"})
	}
	if fields.Type != "decision" && fields.HasProvenance && fields.Provenance != "inferred" && fields.Provenance != "verified" {
		findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "provenance is not inferred or verified; resolve it by hand before resuming"})
	}
	if fields.Type != "decision" && fields.HasReversalCost && !fields.HasProvenance {
		findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "reversal-cost requires a provenance field; resolve it by hand before resuming"})
	}
	if fields.Provenance == "inferred" && fields.Type != "decision" && !fields.HasReversalCost {
		if reversalCost == "" {
			findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "inferred meaning has no reversal-cost; choose --reversal-cost=low or --reversal-cost=high"})
		} else if !fields.ProvenanceSimple {
			findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "provenance syntax is not a plain scalar; add reversal-cost by hand before resuming"})
		} else {
			updated, inserted := insertAfterField(front, "provenance", "reversal-cost: "+reversalCost)
			if !inserted {
				findings = append(findings, Finding{Path: rel, Migration: MigrationReversalCost, Message: "provenance field is not a safe top-level line; add reversal-cost by hand before resuming"})
			} else {
				front = updated
				changes = append(changes, "add reversal-cost: "+reversalCost)
			}
		}
	}
	if fields.Status == "verified" && fields.Type != "decision" {
		if fields.Type == "architecture" || fields.Type == "analysis" {
			if !fields.StatusSimple {
				findings = append(findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "status syntax is not a safe top-level line; change verified to active by hand before resuming"})
			} else {
				updated, replaced := replaceField(front, "status", "active")
				if !replaced {
					findings = append(findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "status syntax is not a safe top-level line; change verified to active by hand before resuming"})
				} else {
					front = updated
					changes = append(changes, "status verified -> active")
				}
			}
		} else {
			findings = append(findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "status verified on a non-architecture/non-analysis type is semantic; choose its current lifecycle by hand"})
		}
	}
	if fields.Status == "retired" && fields.Type != "goal" && fields.Type != "plan" && fields.Type != "decision" && fields.Type != "change" && fields.Type != "tasks" && fields.Type != "open-questions" {
		findings = append(findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "retired is no longer a resting status; delete or supersede this artifact by hand before resuming"})
	}
	if len(findings) > 0 {
		return before, nil, findings
	}
	if len(changes) == 0 {
		return before, nil, nil
	}
	return []byte(front + text[end:]), changes, nil
}

type artifactFields struct {
	Type             string
	Status           string
	StatusSimple     bool
	Provenance       string
	ProvenanceSimple bool
	HasProvenance    bool
	ReversalCost     string
	HasReversalCost  bool
}

// openingFence reports the byte length of the leading `---` line, accepting
// either line ending so a CRLF artifact is read as-is rather than normalized.
func openingFence(text string) (int, bool) {
	switch {
	case strings.HasPrefix(text, "---\r\n"):
		return len("---\r\n"), true
	case strings.HasPrefix(text, "---\n"):
		return len("---\n"), true
	}
	return 0, false
}

// splitFrontmatter locates the frontmatter block by scanning lines in the
// original text. It never rewrites line endings: every terminator the file
// already had — including a mixed set — is left exactly where it was, and
// edits below reuse the terminator of the line they touch.
func splitFrontmatter(text string) (int, string, bool, error) {
	openLength, ok := openingFence(text)
	if !ok {
		return 0, "", false, nil
	}
	for pos := openLength; pos <= len(text); {
		line, next := text[pos:], len(text)
		if lineEnd := strings.IndexByte(line, '\n'); lineEnd >= 0 {
			line, next = line[:lineEnd], pos+lineEnd+1
		}
		if strings.TrimSuffix(line, "\r") == "---" {
			return next, text[openLength:pos], true, nil
		}
		if next == pos {
			break
		}
		pos = next
	}
	return 0, "", false, errors.New("frontmatter closing fence is missing")
}

func readFrontmatter(text string) (int, artifactFields, bool, error) {
	end, inner, ok, err := splitFrontmatter(text)
	if err != nil || !ok {
		return 0, artifactFields{}, false, err
	}
	front := text[:end]
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(inner), &raw); err != nil {
		return 0, artifactFields{}, false, err
	}
	field := func(name string) (string, bool) {
		value, exists := raw[name]
		if !exists {
			return "", false
		}
		str, ok := value.(string)
		return str, ok
	}
	typ, _ := field("type")
	status, statusOK := field("status")
	provenance, provenanceOK := field("provenance")
	reversalCost, _ := field("reversal-cost")
	_, hasProvenance := raw["provenance"]
	_, hasCostField := raw["reversal-cost"]
	return end, artifactFields{Type: typ, Status: status, StatusSimple: statusOK && isSimpleScalar(front, "status", status), Provenance: provenance, ProvenanceSimple: provenanceOK && isSimpleScalar(front, "provenance", provenance), HasProvenance: hasProvenance, ReversalCost: reversalCost, HasReversalCost: hasCostField}, true, nil
}

// topLevelField returns the submatch indexes of the one unindented line
// declaring key. An indented match belongs to a nested mapping or to block
// scalar content, and editing it would silently rewrite something other than
// the field being migrated; several top-level matches are equally ambiguous.
// Both cases return false so the caller raises a finding instead of writing.
func topLevelField(front, key string) ([]int, bool) {
	var found [][]int
	for _, match := range fieldLineRe.FindAllStringSubmatchIndex(front, -1) {
		if match[3] != match[2] {
			continue
		}
		if front[match[4]:match[5]] != key {
			continue
		}
		found = append(found, match)
	}
	if len(found) != 1 {
		return nil, false
	}
	return found[0], true
}

// trailingComment returns the ` # …` remainder of a field value, which a
// replacement must carry over. A `#` not preceded by whitespace is part of
// the scalar, not a comment.
func trailingComment(rawValue string) string {
	hash := strings.Index(rawValue, "#")
	if hash <= 0 || (rawValue[hash-1] != ' ' && rawValue[hash-1] != '\t') {
		return ""
	}
	return " " + rawValue[hash:]
}

func isSimpleScalar(front, key, value string) bool {
	match, ok := topLevelField(front, key)
	if !ok {
		return false
	}
	raw := strings.TrimSpace(front[match[8]:match[9]])
	if raw == value {
		return true
	}
	if comment := trailingComment(raw); comment != "" {
		return strings.TrimSpace(strings.TrimSuffix(raw, strings.TrimPrefix(comment, " "))) == value
	}
	return false
}

func insertAfterField(front, key, line string) (string, bool) {
	match, ok := topLevelField(front, key)
	if !ok {
		return front, false
	}
	// Reuse the anchor line's own terminator so an inserted line matches the
	// file it lands in instead of imposing one ending on the whole artifact.
	eol := front[match[10]:match[11]]
	if eol == "" {
		return front, false
	}
	end := match[1]
	return front[:end] + line + eol + front[end:], true
}

func replaceField(front, key, value string) (string, bool) {
	match, ok := topLevelField(front, key)
	if !ok {
		return front, false
	}
	lineStart, lineEnd := match[0], match[1]
	prefix := front[lineStart:match[8]]
	suffix := trailingComment(front[match[8]:match[9]])
	eol := front[match[10]:match[11]]
	return front[:lineStart] + prefix + value + suffix + eol + front[lineEnd:], true
}

func migrationForArtifact(changes []string) string {
	for _, change := range changes {
		if strings.HasPrefix(change, "status ") {
			return MigrationStatusLifecycle
		}
	}
	return MigrationReversalCost
}

func planCarriers(root, target string, expected map[string][]byte, result *MigrationPlan) {
	for rel, want := range expected {
		if strings.HasPrefix(rel, ".github/workflows/") {
			planCaller(root, rel, want, target, result)
			continue
		}
		planManagedFile(root, rel, want, result)
	}
}

func planManagedFile(root, rel string, want []byte, result *MigrationPlan) {
	if hasLinkBoundary(root, rel) {
		result.Notices = append(result.Notices, Notice{Path: rel, Migration: MigrationManagedCarriers, Message: "managed carrier is behind a symlink; migration does not write through that boundary"})
		return
	}
	full := filepath.Join(root, filepath.FromSlash(rel))
	got, err := os.ReadFile(full)
	if os.IsNotExist(err) {
		if strings.HasPrefix(rel, ".claude/") {
			result.Notices = append(result.Notices, Notice{Path: rel, Migration: MigrationManagedCarriers, Message: "managed mirror is not present; init or its symlink boundary owns materialization"})
			return
		}
		legacyRel := strings.TrimPrefix(rel, ".agents/skills/")
		if legacyRel != rel {
			if oldVersion := skillVersion(root, legacyRel); oldVersion != "" && releaseDigest(oldVersion, legacyRel) == "" {
				result.Changes = append(result.Changes, Change{Path: rel, Migration: MigrationManagedCarriers, Description: fmt.Sprintf("add generated carrier from release %s content", oldVersion), After: want})
				return
			}
		}
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "managed carrier is missing; run clue init or restore it before migrating"})
		return
	}
	if err != nil {
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "managed carrier cannot be read: " + err.Error()})
		return
	}
	if bytes.Equal(got, want) {
		return
	}
	legacyRel := legacyCarrierRel(rel)
	if oldVersion := legacyVersion(legacyRel, got); oldVersion != "" {
		result.Changes = append(result.Changes, Change{Path: rel, Migration: MigrationManagedCarriers, Description: fmt.Sprintf("replace generated carrier %s with release %s content", legacyRel, oldVersion), Existed: true, Before: got, After: want})
		return
	}
	result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "managed carrier differs from every supported generated release; local edits are never overwritten"})
}

func legacyCarrierRel(rel string) string {
	legacyRel := strings.TrimPrefix(rel, ".agents/skills/")
	legacyRel = strings.TrimPrefix(legacyRel, ".claude/skills/")
	if strings.HasSuffix(legacyRel, "/SKILL.md") {
		legacyRel = strings.TrimSuffix(legacyRel, "/SKILL.md") + "/skill.md"
	}
	return legacyRel
}

func planCaller(root, rel string, want []byte, target string, result *MigrationPlan) {
	if hasLinkBoundary(root, rel) {
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "thin CI caller is behind a symlink; resolve the repository-owned path before migrating"})
		return
	}
	full := filepath.Join(root, filepath.FromSlash(rel))
	got, err := os.ReadFile(full)
	if os.IsNotExist(err) {
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "thin CI caller is missing; run clue init before migrating"})
		return
	}
	if err != nil {
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: "thin CI caller cannot be read: " + err.Error()})
		return
	}
	updated, description, ok, message := updateCaller(got, want, target)
	if !ok {
		result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationManagedCarriers, Message: message})
		return
	}
	if !bytes.Equal(got, updated) {
		result.Changes = append(result.Changes, Change{Path: rel, Migration: MigrationManagedCarriers, Description: description, Existed: true, Before: got, After: updated})
	}
}

func updateCaller(got, want []byte, target string) ([]byte, string, bool, string) {
	text := string(got)
	if regexp.MustCompile(`(?m)^\s*steps:`).MatchString(text) {
		return got, "", false, "caller contains steps; copied validation logic is semantic and must be replaced by hand"
	}
	uses := callerUsesRe.FindAllStringSubmatchIndex(text, -1)
	versions := callerVersionRe.FindAllStringSubmatchIndex(text, -1)
	if len(uses) != 1 || len(versions) != 1 {
		return got, "", false, "caller must contain exactly one upstream uses line and one clue-version line"
	}
	wantText := string(want)
	wantUses := callerUsesRe.FindStringSubmatch(wantText)
	wantVersion := callerVersionRe.FindStringSubmatch(wantText)
	if len(wantUses) == 0 || len(wantVersion) == 0 {
		return got, "", false, "embedded thin caller template is missing its upstream reference or version"
	}
	// The target argument is carried separately so callers cannot accidentally
	// update to an arbitrary value if a future template changes shape.
	if target == "" || wantVersion[3] != target {
		return got, "", false, "embedded carrier target version is inconsistent"
	}
	updated := text
	updated = replaceMatch(updated, uses[0], wantUses[4], 8, 9)
	// The first replacement changes indexes, so locate the version again.
	versions = callerVersionRe.FindAllStringSubmatchIndex(updated, -1)
	updated = replaceMatch(updated, versions[0], wantVersion[3], 6, 7)
	description := fmt.Sprintf("update upstream validation reference and clue-version to %s", target)
	return []byte(updated), description, true, ""
}

func replaceMatch(text string, match []int, value string, start, end int) string {
	return text[:match[start]] + value + text[match[end]:]
}

// releaseDigest returns the digest a named release published for rel, or ""
// when that release did not ship the file.
func releaseDigest(version, rel string) string {
	for _, release := range legacyDigests {
		if release.Version == version {
			return release.Files[rel]
		}
	}
	return ""
}

func skillVersion(root, rel string) string {
	skillDir := path.Dir(rel)
	if path.Base(skillDir) == "mappings" {
		skillDir = path.Dir(skillDir)
	}
	skillRel := path.Join(skillDir, "skill.md")
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(".agents/skills/"+skillRel)))
	if err != nil {
		return ""
	}
	return legacyVersion(skillRel, data)
}

// legacyVersion names the earliest release that published these exact bytes
// for rel, or "" when no supported release did. Identical bytes can span
// several releases; the manifest's own order decides, so the answer does not
// depend on how release strings compare.
func legacyVersion(rel string, data []byte) string {
	digest := sha256.Sum256(data)
	hexDigest := hex.EncodeToString(digest[:])
	for _, release := range legacyDigests {
		if release.Files[rel] == hexDigest {
			return release.Version
		}
	}
	return ""
}

func linkedAncestor(root, rel string) bool {
	dir := path.Dir(rel)
	prefix := ""
	for _, segment := range strings.Split(dir, "/") {
		prefix = path.Join(prefix, segment)
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(prefix)))
		if err != nil {
			if os.IsNotExist(err) {
				return false
			}
			return false
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func hasLinkBoundary(root, rel string) bool {
	if linkedAncestor(root, rel) {
		return true
	}
	info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil && info.Mode()&fs.ModeSymlink != 0
}

type appliedFile struct {
	path    string
	backup  string
	created bool
}

func applyFile(path string, data []byte, existed bool) (appliedFile, error) {
	entry := appliedFile{path: path}
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		if !existed {
			return entry, errors.New("file appeared after preflight")
		}
		mode = info.Mode().Perm()
	} else if existed {
		return entry, err
	} else if !os.IsNotExist(err) {
		return entry, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return entry, err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".clue-migrate-*")
	if err != nil {
		return entry, err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return entry, err
	}
	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		return entry, err
	}
	if err := temp.Close(); err != nil {
		return entry, err
	}
	if !existed {
		if err := os.Rename(tempName, path); err != nil {
			return entry, err
		}
		entry.created = true
		return entry, nil
	}
	backup, err := os.CreateTemp(filepath.Dir(path), ".clue-migrate-backup-*")
	if err != nil {
		return entry, err
	}
	backupName := backup.Name()
	if err := backup.Close(); err != nil {
		_ = os.Remove(backupName)
		return entry, err
	}
	if err := os.Remove(backupName); err != nil {
		return entry, err
	}
	if err := os.Rename(path, backupName); err != nil {
		return entry, err
	}
	if err := os.Rename(tempName, path); err != nil {
		_ = os.Rename(backupName, path)
		return entry, err
	}
	entry.backup = backupName
	return entry, nil
}

func rollback(applied []appliedFile) []string {
	var failures []string
	for i := len(applied) - 1; i >= 0; i-- {
		entry := applied[i]
		if entry.created {
			if err := os.Remove(entry.path); err != nil && !os.IsNotExist(err) {
				failures = append(failures, entry.path+": "+err.Error())
			}
			continue
		}
		if entry.backup == "" {
			continue
		}
		if err := os.Remove(entry.path); err != nil && !os.IsNotExist(err) {
			failures = append(failures, entry.path+": "+err.Error())
			continue
		}
		if err := os.Rename(entry.backup, entry.path); err != nil {
			failures = append(failures, entry.path+": "+err.Error())
		}
	}
	return failures
}

func sortPlan(plan *MigrationPlan) {
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Path < plan.Changes[j].Path })
	sort.Slice(plan.Findings, func(i, j int) bool {
		if plan.Findings[i].Path != plan.Findings[j].Path {
			return plan.Findings[i].Path < plan.Findings[j].Path
		}
		if plan.Findings[i].Migration != plan.Findings[j].Migration {
			return plan.Findings[i].Migration < plan.Findings[j].Migration
		}
		return plan.Findings[i].Message < plan.Findings[j].Message
	})
	// Notices tiebreak like findings do. Path alone leaves two notices about
	// the same file in whatever order planning happened to append them, and
	// sort.Slice is not stable — a preview that reordered itself between runs
	// would look like the target changed when nothing did.
	sort.Slice(plan.Notices, func(i, j int) bool {
		if plan.Notices[i].Path != plan.Notices[j].Path {
			return plan.Notices[i].Path < plan.Notices[j].Path
		}
		if plan.Notices[i].Migration != plan.Notices[j].Migration {
			return plan.Notices[i].Migration < plan.Notices[j].Migration
		}
		return plan.Notices[i].Message < plan.Notices[j].Message
	})
}

// planQualifiedReferences reports every external reference that names no
// repository, in the same shape the registry uses for semantic cases.
//
// This migration is deliberately report-only. A bare forge number cannot be
// qualified mechanically: the file does not say which repository was meant,
// and the one reference already known to be wrong in a Cliewen corpus resolved
// perfectly in the wrong namespace. Guessing the adopter's own slug would
// cement exactly that mistake behind a form no later check can question. The
// adopter resolves each one in a reviewed change.
func planQualifiedReferences(root string, result *MigrationPlan) {
	// Scan's second return is the parse-level issue list, not a fatal: an
	// adopter's pre-upgrade corpus is exactly where a stray file without
	// frontmatter is likely, and treating that as a reason to report nothing
	// would silence the migration for the population it exists to serve. The
	// artifacts it did parse are still every artifact it could read.
	scanned, _ := corpus.Scan(root)
	if scanned == nil {
		return
	}
	for _, issue := range corpus.BareReferenceIssues(scanned) {
		result.Findings = append(result.Findings, Finding{
			Path:      issue.Path,
			Migration: MigrationQualifiedReferences,
			Message:   issue.Msg + "; this cannot be repaired mechanically, because nothing here says which repository was meant",
		})
	}
}

// entryPointLocation is a place Claude Code loads a project entry point from,
// with the import that reaches the routing hub from there. An import path is
// relative to the file holding it, so the two locations do not accept the same
// line: under `.claude/`, a bare `@AGENTS.md` names `.claude/AGENTS.md`, which
// no scaffold emits and no session ever loads.
type entryPointLocation struct {
	rel      string
	spelling string
	importRe *regexp.Regexp
}

// entryPointLocations are the places Claude Code loads a project entry point
// from. Both are loaded when both exist, so the hub is reachable as soon as
// any one of them imports it.
var entryPointLocations = []entryPointLocation{
	{rel: "CLAUDE.md", spelling: "@AGENTS.md", importRe: hubImportRe("")},
	{rel: ".claude/CLAUDE.md", spelling: "@../AGENTS.md", importRe: hubImportRe("../")},
}

// hubImportRe matches an import of the routing hub that resolves to the root
// AGENTS.md from a file at the given directory offset. An import is a path
// token, not a line: Claude Code reads `@path` anywhere in the prose, so a
// sentence like "See @AGENTS.md for the rules" routes the session exactly as a
// bare line does. The token must stand alone, because a path run together with
// the text around it is a different path, and one Claude Code would fail to
// load. Code spans and fences are excluded before matching, not here.
func hubImportRe(offset string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)(?:^|[ \t])@(?:\./)?` + regexp.QuoteMeta(offset) + `AGENTS\.md(?:[ \t]|$)`)
}

// planClaudeEntryPoint reports an adopted repository whose Claude Code
// sessions never receive the routing hub, and repairs nothing (PDR-022).
// Claude Code reads CLAUDE.md and not AGENTS.md, so without a bridge the
// adopter's agent lists the lifecycle skills and is never told when to invoke
// them — it holds the manuals with no instruction to open them.
//
// Neither case is a Finding. A missing pointer must not block a carrier
// upgrade: refusing to refresh an adopter's skills until they add a vendor
// file would be a hard stop out of all proportion to a notice.
func planClaudeEntryPoint(root string, result *MigrationPlan) {
	hub := filepath.Join(root, "AGENTS.md")
	var present []entryPointLocation
	for _, loc := range entryPointLocations {
		full := filepath.Join(root, filepath.FromSlash(loc.rel))
		info, err := os.Lstat(full)
		if err != nil {
			continue
		}
		// A symlink to the hub is the vendor's other documented bridge, and
		// reading through it would find AGENTS.md, which of course imports no
		// copy of itself. Only that target routes: a link to anything else is
		// judged on the content it resolves to, like any other file.
		if info.Mode()&fs.ModeSymlink != 0 && linksToHub(full, hub) {
			return
		}
		data, err := os.ReadFile(full)
		if err != nil {
			// The file is there and nothing loads from it — a dangling link,
			// or one the operator cannot read. Reporting it as absent would
			// send them to `clue init`, which skips what already exists.
			present = append(present, loc)
			continue
		}
		if loc.importRe.MatchString(outsideCode(string(data))) {
			return
		}
		present = append(present, loc)
	}
	if len(present) > 0 {
		result.Notices = append(result.Notices, Notice{
			Path:      present[0].rel,
			Migration: MigrationClaudeEntryPoint,
			Message:   fmt.Sprintf("exists but never imports AGENTS.md, so Claude Code reads no routing; add a line containing just `%s` — migration does not edit a file you wrote", present[0].spelling),
		})
		return
	}
	result.Notices = append(result.Notices, Notice{
		Path:      entryPointLocations[0].rel,
		Migration: MigrationClaudeEntryPoint,
		Message:   "absent, so Claude Code reads no routing; run `clue init` to materialize the pointer, which never overwrites an existing file",
	})
}

// linksToHub reports whether an entry point is a symlink resolving to the
// repository's own AGENTS.md. Both sides are resolved, so a checkout reached
// through a link — a worktree, a shared tree — compares equal to itself.
func linksToHub(entry, hub string) bool {
	target, err := filepath.EvalSymlinks(entry)
	if err != nil {
		return false
	}
	resolved, err := filepath.EvalSymlinks(hub)
	if err != nil {
		return false
	}
	return target == resolved
}

// outsideCode blanks the two places Claude Code's import parser does not look
// — fenced code blocks and inline code spans — leaving the prose it does read.
// An import shown as an example must not be mistaken for one that loads:
// reporting a routed repository is noise, while reading an example as the real
// thing leaves the gap unreported, which is the failure this migration exists
// to catch. Line structure is preserved so the caller still matches per line.
func outsideCode(doc string) string {
	lines := strings.Split(strings.ReplaceAll(doc, "\r\n", "\n"), "\n")
	var fence string // the open fence's marker run, empty outside a block
	for i, line := range lines {
		marker := fenceMarker(line)
		if fence == "" {
			if marker != "" {
				fence = marker
				lines[i] = ""
				continue
			}
			lines[i] = blankCodeSpans(line)
			continue
		}
		// Only a fence of the same character and at least the opening
		// length closes the block, so a markdown example wrapped in ````
		// keeps the ``` inside it fenced rather than flipping the state.
		// A closing fence also carries no info string: a ```go line inside
		// a fenced example is content, and reading it as the end of the
		// block would hand the rest of the example back to the matcher as
		// prose — the direction that reports a repository routed when
		// nothing loads.
		if marker != "" && marker[0] == fence[0] && len(marker) >= len(fence) && closesFence(line, marker) {
			fence = ""
		}
		lines[i] = ""
	}
	return strings.Join(lines, "\n")
}

// fenceMarker returns the leading backtick or tilde run of a fence line, or ""
// when the line opens or closes nothing.
func fenceMarker(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" || (trimmed[0] != '`' && trimmed[0] != '~') {
		return ""
	}
	run := len(trimmed) - len(strings.TrimLeft(trimmed, string(trimmed[0])))
	if run < 3 {
		return ""
	}
	return trimmed[:run]
}

// closesFence reports whether a fence line is bare — nothing after its marker
// run but whitespace. Only a bare fence ends a block; a run followed by an
// info string opens one, so inside an open block it is content.
func closesFence(line, marker string) bool {
	rest := strings.TrimLeft(line, " \t")
	return strings.TrimSpace(strings.TrimPrefix(rest, marker)) == ""
}

// blankCodeSpans removes backtick-delimited spans from one line of prose. A
// span closes on a run of the same length, and an unclosed run opens nothing —
// the same reading that makes `@AGENTS.md` in backticks a mention rather than
// an import.
func blankCodeSpans(line string) string {
	var out strings.Builder
	for i := 0; i < len(line); {
		if line[i] != '`' {
			out.WriteByte(line[i])
			i++
			continue
		}
		run := len(line[i:]) - len(strings.TrimLeft(line[i:], "`"))
		if end := strings.Index(line[i+run:], strings.Repeat("`", run)); end >= 0 {
			// A space, not nothing: removing the span outright would join
			// the text on either side into a token neither side wrote.
			out.WriteByte(' ')
			i += run + end + run
			continue
		}
		out.WriteString(line[i : i+run])
		i += run
	}
	return out.String()
}
