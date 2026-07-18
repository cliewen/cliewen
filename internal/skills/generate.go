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
	"clue-verify",
}

var outputRoots = []string{
	".agents/skills",
	"internal/scaffold/templates/skills",
}

type renderedFile struct {
	relativePath string
	content      []byte
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

	rendered := make([]renderedFile, 0, len(skillNames)+1)
	for _, name := range skillNames {
		var output bytes.Buffer
		if err := tmpl.ExecuteTemplate(&output, name+".md.tmpl", nil); err != nil {
			return nil, fmt.Errorf("render %s: %w", name, err)
		}
		rendered = append(rendered, renderedFile{
			relativePath: path.Join(name, "skill.md"),
			content:      normalize(output.Bytes()),
		})
	}

	mapping, err := sources.ReadFile("source/resources/clue-extract/mappings/openspec.md")
	if err != nil {
		return nil, fmt.Errorf("read OpenSpec mapping source: %w", err)
	}
	rendered = append(rendered, renderedFile{
		relativePath: "clue-extract/mappings/openspec.md",
		content:      normalize(mapping),
	})
	sort.Slice(rendered, func(i, j int) bool { return rendered[i].relativePath < rendered[j].relativePath })
	return rendered, nil
}

func normalize(content []byte) []byte {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	text = strings.TrimSpace(text) + "\n"
	return []byte(text)
}
