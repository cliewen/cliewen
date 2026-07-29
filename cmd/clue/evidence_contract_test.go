package main

import (
	"os"
	"strings"
	"testing"
)

func TestSanity_LiveEvidenceModelCarriersAgree(t *testing.T) {
	requiredByFile := map[string][]string{
		"README.md":                                    {"acceptance criterion → acceptance evidence", "Human proof"},
		"AGENTS.md":                                    {"criterion → acceptance evidence", "Human proof in the acceptance brief"},
		"docs/README.md":                               {"acceptance evidence", "classified Go/JVM/Cucumber test reference", "Human acceptance brief"},
		"docs/goals/G-001-verifiable-thread.md":        {"goal → capability → acceptance criterion → acceptance evidence", "Human-class criteria"},
		"docs/architecture/architecture.md":            {"acceptance evidence", "Human acceptance brief"},
		"docs/architecture/core.md":                    {"acceptance criterion → acceptance evidence", "classified executable evidence", "does not execute tests"},
		"docs/capabilities/CAP-002-validate/README.md": {"AC-009", "positive and negative unit tests", "Human-class criteria", "@draft", "does not execute tests"},
		"internal/corpus/actests.go":                   {"positive/negative direction", "Human needs no code reference", "@draft exempts only", "Cucumber feature tags", "does not execute any test runner"},
		"cmd/clue/main.go":                             {"a verifiable thread from goal to acceptance evidence"},
		".claude-plugin/marketplace.json":              {"a verifiable thread from goal to acceptance evidence"},
		"guide/.vitepress/config.mts":                  {"A verifiable thread from goal to acceptance evidence"},
		"guide/index.md":                               {"acceptance evidence connected in Git", "classified test references or genuine Human proof", "without executing tests"},
	}
	for rel, required := range requiredByFile {
		content, err := os.ReadFile(repoPath(rel))
		if err != nil {
			t.Fatal(err)
		}
		text := string(content)
		for _, want := range required {
			if !strings.Contains(text, want) {
				t.Errorf("%s is missing current evidence-model carrier %q", rel, want)
			}
		}
		for _, stale := range []string{
			"a verifiable thread from goal to test",
			"goal → capability → acceptance criterion → test",
			"goal → plan → change → capability → criterion → test",
			"mechanical AC↔test link is the remaining half",
		} {
			if strings.Contains(text, stale) {
				t.Errorf("%s still carries stale evidence-model claim %q", rel, stale)
			}
		}
	}
}
