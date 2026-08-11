package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const portableCoverageReportCommand = "go tool cover -func coverage.out"

func localVerificationBlockUsesPortableCoverageCommand(markdown string) bool {
	start := strings.Index(markdown, "## Verify Locally\n")
	if start < 0 {
		return false
	}
	block := markdown[start:]
	return strings.Contains(block, "\n"+portableCoverageReportCommand+"\n") && !strings.Contains(block, "go tool cover -func=coverage.out")
}

func readContributing(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "CONTRIBUTING.md"))
	if err != nil {
		t.Fatalf("CONTRIBUTING.md not found: %v", err)
	}
	return string(data)
}

func TestAC136_UnitPositive_ContributingDocumentsPortableCoverageCommand(t *testing.T) {
	if !localVerificationBlockUsesPortableCoverageCommand(readContributing(t)) {
		t.Fatalf("CONTRIBUTING.md must document %q without the deprecated equals form", portableCoverageReportCommand)
	}
}

func TestAC136_UnitNegative_EqualsFormIsNotAPortableCoverageCommand(t *testing.T) {
	const unsupported = "## Verify Locally\n\n```text\ngo tool cover -func=coverage.out\n```\n"
	if localVerificationBlockUsesPortableCoverageCommand(unsupported) {
		t.Fatal("the deprecated equals form was accepted as a portable coverage command")
	}
}
