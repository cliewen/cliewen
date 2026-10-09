package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/cliewen/cliewen/internal/acceptpolicy"
)

func TestAC223_IntegrationPositive_InitCLISelectsAcceptance(t *testing.T) {
	for _, mode := range []string{"local", "pr"} {
		root := t.TempDir()
		var out, errs bytes.Buffer
		if code := runInit([]string{"--acceptance=" + mode, root}, &out, &errs); code != 0 {
			t.Fatalf("init: %d %s", code, errs.String())
		}
		p, present, e := acceptpolicy.Load(root)
		if e != nil || !present || p.Mode != mode {
			t.Fatalf("CLI policy %v %v %v", p, present, e)
		}
	}
}

func TestAC223_IntegrationNegative_InitCLIRefusesInvalidOrConflictingChoice(t *testing.T) {
	root := t.TempDir()
	var out, errs bytes.Buffer
	if code := runInit([]string{"--acceptance=invalid", root}, &out, &errs); code != 1 {
		t.Fatalf("invalid returned %d", code)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("invalid flag wrote files")
	}
	if e := os.MkdirAll(filepath.Join(root, ".clue"), 0755); e != nil {
		t.Fatal(e)
	}
	body := "mode: local\nbranch: main\n"
	if e := os.WriteFile(filepath.Join(root, acceptpolicy.Path), []byte(body), 0644); e != nil {
		t.Fatal(e)
	}
	if code := runInit([]string{"--acceptance=pr", root}, &out, &errs); code != 1 {
		t.Fatalf("conflict returned %d", code)
	}
	if _, e := os.Stat(filepath.Join(root, "docs")); !os.IsNotExist(e) {
		t.Fatal("conflict wrote docs")
	}
}
