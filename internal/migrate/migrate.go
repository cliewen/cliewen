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
	Path    string
	Message string
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
	frontmatterStart = "---\n"
	fieldLineRe      = regexp.MustCompile(`(?m)^([ \t]*)([A-Za-z][A-Za-z0-9_-]*):([ \t]*)([^\r\n]*)(\r?\n|$)`)
	callerUsesRe     = regexp.MustCompile(`(?m)^([ \t]*)uses:([ \t]*)(cliewen/cliewen/\.github/workflows/clue-validation\.yml)@([^\s#]+)([ \t]*(?:#.*)?)(\r?\n|$)`)
	callerVersionRe  = regexp.MustCompile(`(?m)^([ \t]*)clue-version:([ \t]*)([^\s#]+)([ \t]*(?:#.*)?)(\r?\n|$)`)
)

// legacyDigests is the committed release manifest for generated carriers.
// A carrier is replaceable only when its bytes match the release that wrote
// it. This is what distinguishes an old generated file from a locally edited
// one without a hidden state file or a destructive overwrite.
var legacyDigests = map[string]map[string]string{
	"0.5.1": {
		"clue-analysis/skill.md":            "9dc9cbf5482dfce9ba78bcc47eae205ab7f65a29061c5614580fbcfb704d430f",
		"clue-delta/skill.md":               "a6ede82283531f71edb53773a7cad477ceae1a7cd97957d74be0f12dec4a9ddc",
		"clue-extract/mappings/openspec.md": "a6d40a7d9b213d58c935b1d343dcb27e3efacde7b1e8a80745c60e724bbc68e4",
		"clue-extract/skill.md":             "98730f6edc2f2170ce1179d3f907ab127db9baab4b1ddedc68718deac1411bd1",
		"clue-plan/skill.md":                "b23af904b43543cfe2616fc7273b9f630105230983a844de415c393024acb048",
		"clue-verify/skill.md":              "788a1d6219248f8f8cd169d2858f20617aae036a7f56c06099f199c8e0776b76",
	},
	"0.6.0": {
		"clue-analysis/skill.md":            "fb97a73b58a466b45d8d47f9cb3a7b3e54d73c9c93877e23637e827fe7bdb6e7",
		"clue-delta/skill.md":               "322c9d3bbce9580ba2b0091a2bbf4743eea3c283f5c679af7c4eb7a122ed364f",
		"clue-extract/mappings/openspec.md": "3673fe7bb66237d8d73c42158581c9ed33c5c53cb3aa89bdbde6312151dd37ea",
		"clue-extract/skill.md":             "d197867a2a33a5adf8d827c9823313a111ff5d19660e9eb7f6d237cdf9cfc1f3",
		"clue-plan/skill.md":                "deb6221dd6b612eac16c8d36acdc8f93807bc02c2bf1c3827f5a1e45e6dba2d7",
		"clue-verify/skill.md":              "e2e078899834b9e5719e300aedb80dec470d2675c6188927f77c69af57e94a34",
	},
	"0.7.0": {
		"clue-analysis/skill.md":            "6194cb804f08133ef3b3e87c39a7616bf9cb2115fe759598e796ba66df967dd4",
		"clue-delta/skill.md":               "e336c2df69fc1be9d5aeebe81b3989dae65187703172ff283120c5aff76e73bc",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "97d0aaa79644bea5bd6bf6000354d7e66d0ca005eb7fe936b672d961001cfc38",
		"clue-extract/skill.md":             "960a77ba250367bca319e708272b78fdeed01683adef5c83c55d6d0f8b4d584e",
		"clue-plan/skill.md":                "0a2b4e3a7c8dbe63247087e08942b3211c42e7d6693de4751fa842a5ebc87497",
		"clue-verify/skill.md":              "faa1b050e5a440d00ee7cff110aca4e1516507f7c8b035ac2e4b36d51f068bed",
	},
	"0.8.0": {
		"clue-analysis/skill.md":            "8a4967df0c74e6fe2b3ef6f857a925d8798ed16152d3ebf503bc64abaab34c5b",
		"clue-delta/skill.md":               "2fd3516de8b8ab9bc7f1fe24e7fca42df5c33b6fbed8c3ce08fb1317e87c347c",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "97d0aaa79644bea5bd6bf6000354d7e66d0ca005eb7fe936b672d961001cfc38",
		"clue-extract/skill.md":             "070b5bf6ff4fd6159166d645d01693de0404e42124cbbcc5dcfd14139ba3af64",
		"clue-plan/skill.md":                "af8dee537456644e99607a702da136467f7fd8967bea950fe8b5a3e5bf6f8661",
		"clue-verify/skill.md":              "44fc206073addf90763def2ca682a86ddbf10395e0dfa72332bda56da094bfec",
	},
	"0.9.0": {
		"clue-analysis/skill.md":            "c472b39311b43cd2c2c8bba1c7de9e1361bf2465de3dc3f6cf3879c5fcf0377f",
		"clue-delta/skill.md":               "a229f48814099930a7df5781b1302bbc97bb602f2179e56f5942218a8b518145",
		"clue-extract/mappings/madr.md":     "96dfb83213387ab16b4927870b21533646eeb8fadae164555cdeccd51fe503e5",
		"clue-extract/mappings/openspec.md": "d58fc65df6f240197e291cd48b250f11676e6d9992147934cd87800ccd58905d",
		"clue-extract/skill.md":             "ef0bc657288a1d5adec5d8980a8c705aff4c9f964d9eb27fd62987bd9cabadb0",
		"clue-plan/skill.md":                "bf6a4e4cd1fa5b6ed8b318687f5da51c0837415e66ae9e48c1c34ff4b2c7e324",
		"clue-verify/skill.md":              "23a409a732255da58412785f7b06fe52923966fe16bcd23250369faf09c920b6",
	},
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "README.md" {
			return nil
		}
		rel, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if hasLinkBoundary(root, rel) {
			result.Findings = append(result.Findings, Finding{Path: rel, Migration: MigrationStatusLifecycle, Message: "corpus artifact is behind a symlink; resolve the repository-owned path before migrating"})
			return nil
		}
		before, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		after, changes, findings := migrateArtifact(rel, before, opts.ReversalCost)
		result.Findings = append(result.Findings, findings...)
		if !bytes.Equal(before, after) {
			result.Changes = append(result.Changes, Change{Path: rel, Migration: migrationForArtifact(changes), Description: strings.Join(changes, "; "), Existed: true, Before: before, After: after})
		}
		return nil
	})
}

func migrateArtifact(rel string, before []byte, reversalCost string) ([]byte, []string, []Finding) {
	text := strings.ReplaceAll(string(before), "\r\n", "\n")
	start, end, fields, _, eol, ok, err := readFrontmatter(text)
	if err != nil {
		return before, nil, []Finding{{Path: rel, Migration: MigrationReversalCost, Message: "frontmatter is ambiguous; resolve it before migration: " + err.Error()}}
	}
	if !ok {
		return before, nil, nil
	}
	var changes []string
	var findings []Finding
	front := text[start:end]
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
			updated, inserted := insertAfterField(front, "provenance", "reversal-cost: "+reversalCost, eol)
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
				updated, replaced := replaceField(front, "status", "active", eol)
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
	updated := text[:start] + front + text[end:]
	if strings.Contains(string(before), "\r\n") {
		updated = strings.ReplaceAll(updated, "\n", "\r\n")
	}
	return []byte(updated), changes, nil
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

func readFrontmatter(text string) (int, int, artifactFields, string, string, bool, error) {
	eol := "\n"
	if strings.Contains(text, "\r\n") {
		eol = "\r\n"
		text = strings.ReplaceAll(text, "\r\n", "\n")
	}
	if !strings.HasPrefix(text, frontmatterStart) {
		return 0, 0, artifactFields{}, text, eol, false, nil
	}
	rest := text[len(frontmatterStart):]
	endOffset := strings.Index(rest, "\n---\n")
	closingLength := len("\n---\n")
	if endOffset < 0 {
		if strings.HasSuffix(rest, "\n---") {
			endOffset = len(rest) - len("\n---")
			closingLength = len("\n---")
		} else {
			return 0, 0, artifactFields{}, text, eol, false, errors.New("frontmatter closing fence is missing")
		}
	}
	end := len(frontmatterStart) + endOffset + closingLength
	front := text[:end]
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(text[len(frontmatterStart):len(frontmatterStart)+endOffset]), &raw); err != nil {
		return 0, 0, artifactFields{}, text, eol, false, err
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
	return 0, end, artifactFields{Type: typ, Status: status, StatusSimple: statusOK && isSimpleScalar(front, "status", status), Provenance: provenance, ProvenanceSimple: provenanceOK && isSimpleScalar(front, "provenance", provenance), HasProvenance: hasProvenance, ReversalCost: reversalCost, HasReversalCost: hasCostField}, text[end:], eol, true, nil
}

func isSimpleScalar(front, key, value string) bool {
	for _, match := range fieldLineRe.FindAllStringSubmatch(front, -1) {
		if match[2] == key {
			raw := strings.TrimSpace(match[4])
			if raw == value {
				return true
			}
			if hash := strings.Index(raw, "#"); hash > 0 && (raw[hash-1] == ' ' || raw[hash-1] == '\t') && strings.TrimSpace(raw[:hash]) == value {
				return true
			}
			return false
		}
	}
	return false
}

func insertAfterField(front, key, line, eol string) (string, bool) {
	for _, match := range fieldLineRe.FindAllStringSubmatchIndex(front, -1) {
		if front[match[4]:match[5]] != key {
			continue
		}
		end := match[1]
		return front[:end] + line + eol + front[end:], true
	}
	return front, false
}

func replaceField(front, key, value, eol string) (string, bool) {
	for _, match := range fieldLineRe.FindAllStringSubmatchIndex(front, -1) {
		if front[match[4]:match[5]] != key {
			continue
		}
		lineStart, lineEnd := match[0], match[1]
		prefix := front[lineStart:match[8]]
		rawValue := front[match[8]:match[9]]
		suffix := ""
		if hash := strings.Index(rawValue, "#"); hash >= 0 {
			suffix = rawValue[hash:]
		}
		if suffix != "" {
			suffix = " " + suffix
		}
		return front[:lineStart] + prefix + value + suffix + eol + front[lineEnd:], true
	}
	return front, false
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
		result.Notices = append(result.Notices, Notice{Path: rel, Message: "managed carrier is behind a symlink; migration does not write through that boundary"})
		return
	}
	full := filepath.Join(root, filepath.FromSlash(rel))
	got, err := os.ReadFile(full)
	if os.IsNotExist(err) {
		if strings.HasPrefix(rel, ".claude/") {
			result.Notices = append(result.Notices, Notice{Path: rel, Message: "managed mirror is not present; init or its symlink boundary owns materialization"})
			return
		}
		legacyRel := strings.TrimPrefix(rel, ".agents/skills/")
		if legacyRel != rel {
			if oldVersion := skillVersion(root, legacyRel); oldVersion != "" && legacyDigests[oldVersion][legacyRel] == "" {
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
	if digestMatchesLegacy(legacyRel, got) {
		oldVersion := legacyVersion(legacyRel, got)
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

func digestMatchesLegacy(rel string, data []byte) bool {
	digest := sha256.Sum256(data)
	hexDigest := hex.EncodeToString(digest[:])
	for _, files := range legacyDigests {
		if files[rel] == hexDigest {
			return true
		}
	}
	return false
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

func legacyVersion(rel string, data []byte) string {
	digest := sha256.Sum256(data)
	hexDigest := hex.EncodeToString(digest[:])
	versions := make([]string, 0, len(legacyDigests))
	for version, files := range legacyDigests {
		if files[rel] == hexDigest {
			versions = append(versions, version)
		}
	}
	sort.Strings(versions)
	if len(versions) == 0 {
		return "unknown"
	}
	return versions[0]
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
	sort.Slice(plan.Notices, func(i, j int) bool { return plan.Notices[i].Path < plan.Notices[j].Path })
}
