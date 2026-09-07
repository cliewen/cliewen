// Package skills renders Cliewen's independently installable agent skills from
// skill-specific templates and shared instruction fragments (ADR-021).
package skills

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:generate go run ./cmd/generate

//go:embed source
var sources embed.FS

var skillNames = []string{
	"clue-analysis",
	"clue-delta",
	"clue-extract",
	"clue-plan",
	"clue-upgrade",
	"clue-verify",
}

var outputRoots = []string{
	".agents/skills",
	"internal/scaffold/templates/skills",
}

// OutputRoots names the repository-relative directories the generator owns:
// it deletes and rewrites each managed skill beneath them, and reports any
// other file there as drift. Hand-authored skill trees must stay outside.
func OutputRoots() []string {
	return append([]string(nil), outputRoots...)
}

type renderedFile struct {
	relativePath string
	content      []byte
}

type skillRoute struct {
	heading   string
	file      string
	condition string
}

type skillDefinition struct {
	description string
	routes      []skillRoute
}

var skillDefinitions = map[string]skillDefinition{
	"clue-analysis": {
		description: "Investigate an unclear risk or unknown before planning or implementation, and leave durable findings for the next workflow.",
		routes: []skillRoute{
			{heading: "Analysis workflow", file: "analysis-workflow.md", condition: "Before starting the investigation"},
			{heading: "Decision records", file: "decision-records.md", condition: "If a finding chooses or rejects a consequential course"},
		},
	},
	"clue-plan": {
		description: "Create or revise a verifiable campaign plan through the reviewed Cliewen change loop.",
		routes: []skillRoute{
			{heading: "Planning workflow", file: "planning-workflow.md", condition: "Before creating or revising a plan"},
			{heading: "Intent model", file: "intent-model.md", condition: "Before deciding which goals a plan serves, and whether any journey needs a use case"},
			{heading: "Intent discovery", file: "intent-discovery.md", condition: "When the repository states no usable vision"},
			{heading: "Decision records", file: "decision-records.md", condition: "When plan meaning or another consequential choice is recorded"},
		},
	},
	"clue-delta": {
		description: "Run a chosen full Cliewen change from proposal through implementation, digest, verification, and human-controlled merge.",
		routes: []skillRoute{
			{heading: "Change routing", file: "change-scope-and-tiers.md", condition: "Before recommending a route or starting a full change"},
			{heading: "Review boundary", file: "review-boundary.md", condition: "Before branching, publishing, updating a hosted PR, or handing work to a human"},
			{heading: "Change loop", file: "change-loop.md", condition: "After the user chooses the recommended full loop"},
			{heading: "Intent model", file: "intent-model.md", condition: "When the change touches what the product means, or the acceptance brief must state the vision it proceeds under"},
			{heading: "Decision records", file: "decision-records.md", condition: "When the change makes, rejects, or carries a decision"},
			{heading: "Repository-local conventions", file: "repository-local-conventions.md", condition: "Before applying repository-specific implementation or digest rules"},
			{heading: "Durable work state", file: "durable-work-state.md", condition: "When a change starts or resumes, a suggestion arrives, or a merge is reported"},
		},
	},
	"clue-extract": {
		description: "Transform one brownfield specification corpus into Cliewen through a report-only rehearsal and a human-authorized mutation.",
		routes: []skillRoute{
			{heading: "Review boundary", file: "review-boundary.md", condition: "Before branching, publishing, updating a hosted PR, or handing work to a human"},
			{heading: "Boundaries", file: "boundaries.md", condition: "Before beginning an extraction"},
			{heading: "Rehearsal before mutation", file: "rehearsal-before-mutation.md", condition: "After proposal and before changing the target corpus, tests, routing, or hosted state"},
			{heading: "Target contract", file: "target-contract.md", condition: "After the human authorizes mutation and while constructing the target corpus"},
			{heading: "Intent model", file: "intent-model.md", condition: "Before proposing a vision or any use case for the target corpus"},
			{heading: "Intent discovery", file: "intent-discovery.md", condition: "When the source repository states no usable vision"},
			{heading: "Source mappings", file: "source-mappings.md", condition: "When the source uses a supported format or needs a new mapping"},
			{heading: "Decision records", file: "decision-records.md", condition: "When extraction classifies or records a consequential choice"},
			{heading: "Repository-local conventions", file: "repository-local-conventions.md", condition: "Before reconciling source instructions with repository-specific rules"},
			{heading: "Durable work state", file: "durable-work-state.md", condition: "When extraction work starts or resumes, a suggestion arrives, or a merge is reported"},
		},
	},
	"clue-upgrade": {
		description: "Check for a newer Cliewen release and, only with explicit human authorization, carry out its coordinated repository upgrade.",
		routes: []skillRoute{
			{heading: "Upgrade workflow", file: "upgrade-workflow.md", condition: "Before checking or acting on an available release, and before recommending the upgrade's route"},
			{heading: "Change routing", file: "change-scope-and-tiers.md", condition: "When the upgrade escalates a decision of this repository's own to the full loop"},
			{heading: "Review boundary", file: "review-boundary.md", condition: "If an upgrade change begins and before branching, publishing, or handing it off"},
			{heading: "Decision records", file: "decision-records.md", condition: "When the upgrade requires a consequential local choice"},
			{heading: "Repository-local conventions", file: "repository-local-conventions.md", condition: "Before applying repository-specific upgrade or verification rules"},
			{heading: "Durable work state", file: "durable-work-state.md", condition: "When an upgrade change starts or resumes, a suggestion arrives, or a merge is reported"},
		},
	},
	"clue-verify": {
		description: "Verify a chosen full Cliewen change and run its bounded adversarial review before claiming the hosted pull request is ready.",
		routes: []skillRoute{
			{heading: "Change routing", file: "change-scope-and-tiers.md", condition: "Before confirming that full-loop verification applies"},
			{heading: "Review boundary", file: "review-boundary.md", condition: "Before inspecting or updating hosted pull-request state and before the readiness handoff"},
			{heading: "Verification checklist", file: "verification-checklist.md", condition: "Before running readiness verification"},
			{heading: "Agentic review loop", file: "agentic-review-loop.md", condition: "After the complete candidate is committed and its applicable local checks pass"},
			{heading: "Decision records", file: "decision-records.md", condition: "When verification encounters or evaluates a consequential choice"},
			{heading: "Repository-local conventions", file: "repository-local-conventions.md", condition: "Before selecting and running repository-specific checks"},
			{heading: "Durable work state", file: "durable-work-state.md", condition: "When review work starts or resumes, a suggestion arrives, or a merge is reported"},
		},
	},
}

// Drift describes one committed output that differs from the canonical render.
type Drift struct {
	Path   string
	Reason string
}

func (d Drift) String() string {
	return d.Path + ": " + d.Reason
}

// Write replaces the generator-owned skill directories under root.
func Write(root string) error {
	rendered, err := render()
	if err != nil {
		return err
	}
	for _, outputRoot := range outputRoots {
		for _, name := range skillNames {
			owned := filepath.Join(root, filepath.FromSlash(outputRoot), name)
			if err := os.RemoveAll(owned); err != nil {
				return fmt.Errorf("remove generated skill directory %s: %w", owned, err)
			}
		}
		for _, file := range rendered {
			target := filepath.Join(root, filepath.FromSlash(outputRoot), filepath.FromSlash(file.relativePath))
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("create generated skill directory for %s: %w", target, err)
			}
			if err := os.WriteFile(target, file.content, 0o644); err != nil {
				return fmt.Errorf("write generated skill %s: %w", target, err)
			}
		}
	}
	return nil
}

// Check compares both committed output trees with the canonical render.
func Check(root string) ([]Drift, error) {
	rendered, err := render()
	if err != nil {
		return nil, err
	}
	expected := make(map[string][]byte, len(rendered))
	for _, file := range rendered {
		expected[file.relativePath] = file.content
	}

	var drifts []Drift
	for _, outputRoot := range outputRoots {
		for rel, want := range expected {
			full := filepath.Join(root, filepath.FromSlash(outputRoot), filepath.FromSlash(rel))
			got, readErr := os.ReadFile(full)
			switch {
			case os.IsNotExist(readErr):
				drifts = append(drifts, Drift{Path: filepath.ToSlash(filepath.Join(outputRoot, rel)), Reason: "generated file is missing"})
			case readErr != nil:
				return nil, fmt.Errorf("read generated skill %s: %w", full, readErr)
			case !bytes.Equal(got, want):
				drifts = append(drifts, Drift{Path: filepath.ToSlash(filepath.Join(outputRoot, rel)), Reason: "content differs from canonical skill sources"})
			}
		}

		for _, name := range skillNames {
			owned := filepath.Join(root, filepath.FromSlash(outputRoot), name)
			walkErr := filepath.WalkDir(owned, func(filePath string, entry fs.DirEntry, walkErr error) error {
				if os.IsNotExist(walkErr) {
					return nil
				}
				if walkErr != nil {
					return walkErr
				}
				if entry.IsDir() {
					return nil
				}
				rel, relErr := filepath.Rel(filepath.Join(root, filepath.FromSlash(outputRoot)), filePath)
				if relErr != nil {
					return relErr
				}
				slashRel := filepath.ToSlash(rel)
				if _, ok := expected[slashRel]; !ok {
					drifts = append(drifts, Drift{Path: filepath.ToSlash(filepath.Join(outputRoot, rel)), Reason: "unexpected file in generator-owned skill directory"})
				}
				return nil
			})
			if walkErr != nil && !os.IsNotExist(walkErr) {
				return nil, fmt.Errorf("inspect generated skill directory %s: %w", owned, walkErr)
			}
		}
	}
	sort.Slice(drifts, func(i, j int) bool {
		if drifts[i].Path == drifts[j].Path {
			return drifts[i].Reason < drifts[j].Reason
		}
		return drifts[i].Path < drifts[j].Path
	})
	return drifts, nil
}

func render() ([]renderedFile, error) {
	tmpl, err := template.New("skills").Option("missingkey=error").ParseFS(sources, "source/shared/*.md.tmpl", "source/skills/*.md.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse skill sources: %w", err)
	}

	var rendered []renderedFile
	for _, name := range skillNames {
		var output bytes.Buffer
		if err := tmpl.ExecuteTemplate(&output, name+".md.tmpl", nil); err != nil {
			return nil, fmt.Errorf("render %s: %w", name, err)
		}
		files, splitErr := splitSkill(name, normalize(output.Bytes()))
		if splitErr != nil {
			return nil, splitErr
		}
		rendered = append(rendered, files...)
	}

	const mappingsDir = "source/resources/clue-extract/mappings"
	walkErr := fs.WalkDir(sources, mappingsDir, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		content, readErr := sources.ReadFile(filePath)
		if readErr != nil {
			return readErr
		}
		rel := strings.TrimPrefix(filePath, mappingsDir+"/")
		rendered = append(rendered, renderedFile{
			relativePath: path.Join("clue-extract/mappings", rel),
			content:      normalize(content),
		})
		return nil
	})
	if walkErr != nil {
		return nil, fmt.Errorf("read clue-extract mapping sources: %w", walkErr)
	}
	sort.Slice(rendered, func(i, j int) bool { return rendered[i].relativePath < rendered[j].relativePath })
	return rendered, nil
}

func splitSkill(name string, complete []byte) ([]renderedFile, error) {
	definition, ok := skillDefinitions[name]
	if !ok {
		return nil, fmt.Errorf("split %s: no skill definition", name)
	}
	text := string(complete)
	if !strings.HasPrefix(text, "---\n") {
		return nil, fmt.Errorf("split %s: missing frontmatter", name)
	}
	frontmatterEnd := strings.Index(text[4:], "\n---\n")
	if frontmatterEnd < 0 {
		return nil, fmt.Errorf("split %s: unterminated frontmatter", name)
	}
	frontmatterEnd += 9
	frontmatter := text[:frontmatterEnd]
	body := strings.TrimLeft(text[frontmatterEnd:], "\n")
	title := "# " + name
	titleIndex := strings.Index(body, title+"\n")
	if titleIndex < 0 {
		return nil, fmt.Errorf("split %s: missing title %q", name, title)
	}
	preamble := strings.TrimSpace(body[:titleIndex])
	body = strings.TrimPrefix(body[titleIndex:], title+"\n")
	body = strings.TrimLeft(body, "\n")

	sections := map[string]string{}
	var heading string
	var section strings.Builder
	flush := func() {
		if heading != "" {
			sections[heading] = strings.TrimSpace(section.String()) + "\n"
		}
		section.Reset()
	}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "## ") {
			flush()
			heading = strings.TrimPrefix(line, "## ")
			section.WriteString(line)
			section.WriteByte('\n')
			continue
		}
		if heading == "" {
			if strings.TrimSpace(line) != "" {
				return nil, fmt.Errorf("split %s: unrouted content before first section: %q", name, line)
			}
			continue
		}
		section.WriteString(line)
		section.WriteByte('\n')
	}
	flush()

	var router strings.Builder
	router.WriteString(frontmatter)
	if preamble != "" {
		router.WriteString("\n")
		router.WriteString(preamble)
		router.WriteString("\n\n")
	}
	router.WriteString(title)
	router.WriteString("\n\n")
	router.WriteString(definition.description)
	router.WriteString("\n\n## Routing\n\nRead each reference when its condition is reached, before taking action governed by it. The references are required instructions, not optional background.\n\n")
	files := []renderedFile{{relativePath: path.Join(name, "skill.md")}}
	for _, route := range definition.routes {
		content, found := sections[route.heading]
		if !found {
			return nil, fmt.Errorf("split %s: routed section %q is missing", name, route.heading)
		}
		delete(sections, route.heading)
		router.WriteString("- ")
		router.WriteString(route.condition)
		router.WriteString(", read [")
		router.WriteString(route.heading)
		router.WriteString("](references/")
		router.WriteString(route.file)
		router.WriteString(").\n")
		files = append(files, renderedFile{
			relativePath: path.Join(name, "references", route.file),
			content:      normalize([]byte(content)),
		})
	}
	if len(sections) != 0 {
		remaining := make([]string, 0, len(sections))
		for name := range sections {
			remaining = append(remaining, name)
		}
		sort.Strings(remaining)
		return nil, fmt.Errorf("split %s: sections have no route: %s", name, strings.Join(remaining, ", "))
	}
	files[0].content = normalize([]byte(router.String()))
	return files, nil
}

func normalize(content []byte) []byte {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	text = strings.TrimSpace(text) + "\n"
	return []byte(text)
}
