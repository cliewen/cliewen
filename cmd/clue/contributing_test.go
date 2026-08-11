package main

import (
	"os"
	"os/exec"
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
	fenceStart := strings.Index(markdown[start:], "```text\n")
	if fenceStart < 0 {
		return false
	}
	fenceStart += start + len("```text\n")
	fenceEnd := strings.Index(markdown[fenceStart:], "```\n")
	if fenceEnd < 0 {
		return false
	}
	block := markdown[fenceStart : fenceStart+fenceEnd]
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
	profile := filepath.Join(t.TempDir(), "coverage.out")
	if err := os.WriteFile(profile, []byte("mode: set\ngithub.com/cliewen/cliewen/cmd/clue/main.go:42.48,43.31 1 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("go", "tool", "cover", "-func", profile).CombinedOutput(); err != nil {
		t.Fatalf("%q did not render a coverage profile: %v\n%s", portableCoverageReportCommand, err, output)
	}
}

func TestAC136_UnitNegative_EqualsFormIsNotAPortableCoverageCommand(t *testing.T) {
	const unsupported = "## Verify Locally\n\n```text\ngo tool cover -func=coverage.out\n```\n\nUse go tool cover -func coverage.out outside the verification block.\n"
	if localVerificationBlockUsesPortableCoverageCommand(unsupported) {
		t.Fatal("the deprecated equals form was accepted because the portable command appeared outside the verification block")
	}
}
