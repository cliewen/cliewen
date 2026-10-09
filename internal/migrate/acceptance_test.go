package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cliewen/cliewen/internal/acceptpolicy"
)

func policyFile(t *testing.T, root, path, body string) {
	t.Helper()
	name := filepath.Join(root, path)
	if e := os.MkdirAll(filepath.Dir(name), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(name, []byte(body), 0644); e != nil {
		t.Fatal(e)
	}
}

func TestAC224_IntegrationPositive_MigrationMakesLegacyPRExplicit(t *testing.T) {
	root := t.TempDir()
	policyFile(t, root, "docs/g.md", "---\nid: G-001\ntype: goal\nstatus: accepted\nlinks: []\ntitle: Legacy\n---\n")
	var plan MigrationPlan
	planAcceptancePolicy(root, &plan)
	if len(plan.Changes) != 1 || plan.Changes[0].Path != acceptpolicy.Path || len(plan.Findings) != 0 {
		t.Fatalf("plan %+v", plan)
	}
	if _, e := os.Stat(filepath.Join(root, acceptpolicy.Path)); !os.IsNotExist(e) {
		t.Fatal("preview wrote policy")
	}
	if e := Apply(root, plan); e != nil {
		t.Fatal(e)
	}
	p, present, e := acceptpolicy.Load(root)
	if e != nil || !present || p.Mode != "pr" {
		t.Fatalf("policy %v %v %v", p, present, e)
	}
	for _, mode := range []string{"pr", "local"} {
		body := "# preserved\nmode: " + mode + "\nbranch: main\n"
		policyFile(t, root, acceptpolicy.Path, body)
		plan = MigrationPlan{}
		planAcceptancePolicy(root, &plan)
		if len(plan.Changes) != 0 || len(plan.Findings) != 0 {
			t.Fatalf("changed existing %+v", plan)
		}
		if e := Apply(root, plan); e != nil {
			t.Fatal(e)
		}
		got, _ := os.ReadFile(filepath.Join(root, acceptpolicy.Path))
		if string(got) != body {
			t.Fatal("migration rewrote existing policy")
		}
	}
}

func TestAC224_IntegrationNegative_MalformedPolicyBlocksAllWrites(t *testing.T) {
	root := t.TempDir()
	policyFile(t, root, acceptpolicy.Path, "mode: wrong\nbranch: main\n")
	plan := MigrationPlan{Changes: []Change{{Path: "other.txt", After: []byte("must not write")}}}
	planAcceptancePolicy(root, &plan)
	if len(plan.Findings) != 1 {
		t.Fatalf("findings %+v", plan.Findings)
	}
	if e := Apply(root, plan); e == nil {
		t.Fatal("malformed policy allowed apply")
	}
	if _, e := os.Stat(filepath.Join(root, "other.txt")); !os.IsNotExist(e) {
		t.Fatal("partial write")
	}
}

func TestUnit_AcceptanceMigrationLeavesFreshRootsAndUserHubsAlone(t *testing.T) {
	root := t.TempDir()
	var plan MigrationPlan
	planAcceptancePolicy(root, &plan)
	if len(plan.Changes) != 0 {
		t.Fatal("migration assumed prior adoption")
	}
	policyFile(t, root, acceptpolicy.Path, "mode: local\nbranch: main\n")
	hub := "PR acceptance is the default.\n"
	policyFile(t, root, "AGENTS.md", hub)
	plan = MigrationPlan{}
	planAcceptancePolicy(root, &plan)
	if len(plan.Notices) != 1 || len(plan.Changes) != 0 {
		t.Fatalf("hub plan %+v", plan)
	}
	got, _ := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if string(got) != hub {
		t.Fatal("hub rewritten")
	}
}

func TestSanity_AcceptanceMigrationRefusesSymlinkBoundary(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	policyFile(t, outside, "acceptance.yaml", "mode: local\nbranch: main\n")
	if e := os.Symlink(outside, filepath.Join(root, ".clue")); e != nil {
		t.Skipf("symlinks unavailable: %v", e)
	}
	var plan MigrationPlan
	planAcceptancePolicy(root, &plan)
	if len(plan.Findings) != 1 || len(plan.Changes) != 0 {
		t.Fatalf("symlink plan %+v", plan)
	}
	if e := Apply(root, plan); e == nil {
		t.Fatal("symlink policy allowed migration")
	}
	got, _ := os.ReadFile(filepath.Join(outside, "acceptance.yaml"))
	if string(got) != "mode: local\nbranch: main\n" {
		t.Fatal("outside policy changed")
	}
}
