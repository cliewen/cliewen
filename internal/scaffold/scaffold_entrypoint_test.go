package scaffold

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// AC-071 positive: init emits a Claude Code entry point, and it only points.
// Claude Code loads CLAUDE.md and not AGENTS.md, so an adopter without this
// file gets the lifecycle skills listed and never the routing that says when
// to invoke them (PDR-022).
func TestAC071_UnitPositive_InitEmitsAnEntryPointThatOnlyPoints(t *testing.T) {
	root, rep := runInto(t)
	if !slices.Contains(rep.Created, "CLAUDE.md") {
		t.Fatalf("the report does not name CLAUDE.md as created: %v", rep.Created)
	}
	entry, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(entry)
	// The import must stand alone: Claude Code skips imports inside code
	// spans and fences, so a mention in backticks would silently load nothing.
	var imported bool
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == "@AGENTS.md" {
			imported = true
		}
	}
	if !imported {
		t.Errorf("CLAUDE.md carries no bare @AGENTS.md import, so the hub never loads:\n%s", content)
	}
	if !strings.Contains(content, "Rules go in `AGENTS.md`") {
		t.Errorf("CLAUDE.md does not send rules back to the hub:\n%s", content)
	}

	// The pointer contract, mechanically: nothing the hub says is restated
	// here. Duplication is the failure PDR-022 forbids, because a second copy
	// drifts and only one assistant ever sees the divergence.
	hub, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	sentences := longSentences(string(hub))
	if len(sentences) < 10 {
		t.Fatalf("only %d comparable sentences in AGENTS.md — the duplication guard would pass vacuously", len(sentences))
	}
	for _, sentence := range sentences {
		if strings.Contains(content, sentence) {
			t.Errorf("CLAUDE.md duplicates a sentence of AGENTS.md rather than pointing at it: %q", sentence)
		}
	}
}

// Every emitted text file ends with a newline. An adopter's repository is
// linted, reviewed, and diffed by tools that all treat a missing final
// newline as damage, and the scaffold must not be the thing that introduces
// it on the first commit.
func TestSanity_EmittedTextFilesEndWithANewline(t *testing.T) {
	root, rep := runInto(t)
	for _, rel := range rep.Created {
		if !strings.HasSuffix(rel, ".md") && !strings.HasSuffix(rel, ".yml") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 || data[len(data)-1] != '\n' {
			t.Errorf("emitted %s does not end with a newline", rel)
		}
	}
}

// AC-071 negative: an adopter who already keeps a CLAUDE.md keeps every byte
// of it. The file is theirs on arrival — that is what makes emitting a
// vendor's flagship file honest rather than a claim on their repository.
func TestAC071_UnitNegative_AnExistingEntryPointIsNeverTouched(t *testing.T) {
	root := t.TempDir()
	own := "# Ours\n\nRun the integration suite with `make it`.\n"
	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte(own), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != own {
		t.Fatalf("an existing CLAUDE.md was rewritten:\n%s", got)
	}
	if !slices.Contains(rep.Skipped, "CLAUDE.md") {
		t.Errorf("the report does not name CLAUDE.md as skipped: %v", rep.Skipped)
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); err != nil {
		t.Errorf("the rest of the convention was not created: %v", err)
	}
}

// longSentences returns the prose sentences worth comparing between two
// carriers. Short ones ("The skills are generated.") recur innocently; a
// shared long one is copied text.
func longSentences(doc string) []string {
	var out []string
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "|") {
			continue
		}
		for _, sentence := range strings.Split(line, ". ") {
			if sentence = strings.TrimSpace(sentence); len(sentence) >= 60 {
				out = append(out, sentence)
			}
		}
	}
	return out
}
