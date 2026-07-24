//go:build unix

package corpus

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// AC-037: a named pipe is not a manifest. Opening a FIFO with no writer blocks
// until one arrives, so a resolver that accepted any non-directory entry would
// hang the judge here rather than fail it — a validate that never returns is
// worse than one that disagrees (ADR-028). The verdict is taken off a
// goroutine so a regression is a bounded failure instead of a stuck suite.
func TestAC037_ANamedPipeIsNotAManifest(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".agents", "skills", "clue-delta")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, "skill.md"), 0o644); err != nil {
		t.Skip("this host cannot create a FIFO: ", err)
	}

	done := make(chan []Issue, 1)
	go func() { done <- checkSkillVersions(&Corpus{Root: root}, "0.1.0") }()
	select {
	case issues := <-done:
		if len(issues) != 0 {
			t.Fatalf("a named pipe is not a manifest, got %v", issues)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("checkSkillVersions blocked on a FIFO named skill.md instead of skipping it")
	}
}
