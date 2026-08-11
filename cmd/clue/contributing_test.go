package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	portableCoverageProfileCommand = "go test ./... -coverprofile coverage.out"
	portableCoverageReportCommand  = "go tool cover -func coverage.out"
)

// shellSplitRe matches a single-dash flag written as -flag=value whose value
// contains a dot. PowerShell 7's default Windows native-argument mode splits
// such a token at that dot, so the program receives two arguments instead of
// one; the space-separated form and the double-dash form are unaffected. The
// pattern is the defect class rather than the two commands that hit it, so a
// command added to the block later is covered without editing this test.
var shellSplitRe = regexp.MustCompile(`(^|\s)-[A-Za-z][\w-]*=\S*\.\S*`)

// verifyLocallyBlock returns the trimmed command lines of the fenced block
// under CONTRIBUTING.md's "Verify Locally" heading. A missing heading, a
// missing fence, or an unterminated fence is a failure to read the block, not
// an empty block: reporting it as absent is what keeps a later section's fence
// from being validated in its place.
func verifyLocallyBlock(markdown string) ([]string, bool) {
	lines := strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n")
	i := 0
	for ; i < len(lines) && lines[i] != "## Verify Locally"; i++ {
	}
	for i++; ; i++ {
		if i >= len(lines) || strings.HasPrefix(lines[i], "## ") {
			return nil, false
		}
		if strings.HasPrefix(lines[i], "```") {
			break
		}
	}
	var block []string
	for i++; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "```") {
			return block, true
		}
		block = append(block, strings.TrimSpace(lines[i]))
	}
	return nil, false
}

func localVerificationBlockIsShellPortable(markdown string) bool {
	block, ok := verifyLocallyBlock(markdown)
	if !ok {
		return false
	}
	documented := map[string]bool{}
	for _, line := range block {
		if shellSplitRe.MatchString(line) {
			return false
		}
		documented[line] = true
	}
	return documented[portableCoverageProfileCommand] && documented[portableCoverageReportCommand]
}

func readContributing(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "CONTRIBUTING.md"))
	if err != nil {
		t.Fatalf("CONTRIBUTING.md not found: %v", err)
	}
	return string(data)
}

func TestAC136_UnitPositive_ContributingDocumentsPortableCoverageCommands(t *testing.T) {
	if !localVerificationBlockIsShellPortable(readContributing(t)) {
		t.Fatalf("CONTRIBUTING.md must document %q and %q, and pass no -flag=value with a dot in the value", portableCoverageProfileCommand, portableCoverageReportCommand)
	}
	profile := filepath.Join(t.TempDir(), "coverage.out")
	if err := os.WriteFile(profile, []byte("mode: set\ngithub.com/cliewen/cliewen/cmd/clue/main.go:42.48,43.31 1 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("go", "tool", "cover", "-func", profile).CombinedOutput(); err != nil {
		t.Fatalf("%q did not render a coverage profile: %v\n%s", portableCoverageReportCommand, err, output)
	}
}

// The block holds both portable commands, so only the equals-form line can
// fail this fixture — which is what makes it evidence that the check reads
// every command in the block rather than the two it names.
func TestAC136_UnitNegative_EqualsFormInsideTheBlockIsRejected(t *testing.T) {
	unsupported := "## Verify Locally\n\n```text\n" + portableCoverageProfileCommand + "\n" + portableCoverageReportCommand + "\ngo tool cover -html=coverage.out\n```\n"
	if localVerificationBlockIsShellPortable(unsupported) {
		t.Fatal("a -flag=value command with a dot in the value was accepted inside the verification block")
	}
}

// No equals form appears anywhere here, so only the block scoping can fail the
// fixture: a portable command documented in prose does not put it in the block
// a contributor actually runs.
func TestAC136_UnitNegative_PortableCommandsOutsideTheBlockDoNotCount(t *testing.T) {
	unsupported := "## Verify Locally\n\n```text\ngo build ./...\n```\n\nThen run " + portableCoverageProfileCommand + " and " + portableCoverageReportCommand + ".\n"
	if localVerificationBlockIsShellPortable(unsupported) {
		t.Fatal("the coverage commands were accepted from prose outside the verification block")
	}
}
