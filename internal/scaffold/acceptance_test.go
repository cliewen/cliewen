package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cliewen/cliewen/internal/acceptpolicy"
)

func acceptanceFile(t *testing.T, root, path, body string) {
	t.Helper()
	name := filepath.Join(root, path)
	if e := os.MkdirAll(filepath.Dir(name), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(name, []byte(body), 0644); e != nil {
		t.Fatal(e)
	}
}

func TestAC223_IntegrationPositive_FreshDefaultAndExplicitPR(t *testing.T) {
	for _, requested := range []string{"", "local", "pr"} {
		root := t.TempDir()
		acceptanceFile(t, root, "project.txt", "ordinary project")
		if _, e := RunWithAcceptance(root, requested); e != nil {
			t.Fatal(e)
		}
		want := "local"
		if requested == "pr" {
			want = "pr"
		}
		p, present, e := acceptpolicy.Load(root)
		if e != nil || !present || p != (acceptpolicy.Policy{Mode: want, Branch: "main"}) {
			t.Fatalf("policy %v %v %v", p, present, e)
		}
		before, _ := os.ReadFile(filepath.Join(root, acceptpolicy.Path))
		if _, e := Run(root); e != nil {
			t.Fatal(e)
		}
		after, _ := os.ReadFile(filepath.Join(root, acceptpolicy.Path))
		if string(after) != string(before) {
			t.Fatal("init changed policy")
		}
	}
}

func TestAC223_IntegrationNegative_ConflictingOptionsWriteNothing(t *testing.T) {
	for _, body := range []string{"# user choice\nmode: local\nbranch: trunk\n", "mode: invalid\nbranch: main\n"} {
		root := t.TempDir()
		acceptanceFile(t, root, acceptpolicy.Path, body)
		if _, e := RunWithAcceptance(root, "pr"); e == nil {
			t.Fatal("conflicting or invalid policy accepted")
		}
		if _, e := os.Stat(filepath.Join(root, "docs")); !os.IsNotExist(e) {
			t.Fatal("init wrote before policy resolution")
		}
		got, _ := os.ReadFile(filepath.Join(root, acceptpolicy.Path))
		if string(got) != body {
			t.Fatal("policy overwritten")
		}
	}
}

func TestAC227_IntegrationPositive_ExistingAdoptionDefaultsToLocal(t *testing.T) {
	for _, tc := range []struct{ path, body string }{{".clue/role.yaml", "role: adopter\n"}, {"docs/old.md", "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Legacy\n---\n"}} {
		root := t.TempDir()
		acceptanceFile(t, root, tc.path, tc.body)
		if _, e := Run(root); e != nil {
			t.Fatal(e)
		}
		p, _, e := acceptpolicy.Load(root)
		if e != nil || p.Mode != "local" {
			t.Fatalf("default: %v %v", p, e)
		}
	}
}

func TestAC227_IntegrationNegative_InitNeverOverwritesLocalChoice(t *testing.T) {
	root := t.TempDir()
	body := "# chosen before init\nmode: local\nbranch: trunk\n"
	acceptanceFile(t, root, acceptpolicy.Path, body)
	if _, e := Run(root); e != nil {
		t.Fatal(e)
	}
	got, _ := os.ReadFile(filepath.Join(root, acceptpolicy.Path))
	if string(got) != body {
		t.Fatal("choice changed")
	}
}
