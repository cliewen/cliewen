package acceptpolicy

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func put(t *testing.T, root, path, body string) {
	t.Helper()
	name := filepath.Join(root, path)
	if e := os.MkdirAll(filepath.Dir(name), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(name, []byte(body), 0644); e != nil {
		t.Fatal(e)
	}
}

func TestAC231_UnitPositive_ExplicitAndDefaultPolicy(t *testing.T) {
	root := t.TempDir()
	p, present, e := Load(root)
	if e != nil || present || p != (Policy{Local, "main"}) {
		t.Fatalf("default %v %v %v", p, present, e)
	}
	for _, mode := range []string{Local, PR} {
		put(t, root, Path, "# preserve this comment\nmode: "+mode+"\nbranch: release/stable\n")
		p, present, e = Load(root)
		if e != nil || !present || p != (Policy{mode, "release/stable"}) {
			t.Fatalf("explicit %v %v %v", p, present, e)
		}
		put(t, root, Path, string(p.Bytes()))
		if got, _, e := Load(root); e != nil || got != p {
			t.Fatalf("roundtrip %v %v", got, e)
		}
	}
	put(t, root, ".clue/role.yaml", "role: source\n")
	if p, _, e = Load(root); e != nil || p.Mode != PR {
		t.Fatalf("source PR %v %v", p, e)
	}
}

func TestAC231_UnitNegative_InvalidPolicyIsNotGuessed(t *testing.T) {
	root := t.TempDir()
	for _, body := range []string{"", "mode: unknown\nbranch: main\n", "mode: local\n", "mode: local\nbranch: bad..name\n", "mode: local\nbranch: main\nextra: ignored\n", "mode: local\nbranch: main\n---\nmode: pr\nbranch: main\n", "mode: local\nmode: pr\nbranch: main\n"} {
		put(t, root, Path, body)
		if _, _, e := Load(root); e == nil {
			t.Fatalf("accepted malformed %q", body)
		}
	}
	put(t, root, Path, "mode: local\nbranch: main\n")
	put(t, root, ".clue/role.yaml", "role: unknown\n")
	if _, _, e := Load(root); e == nil {
		t.Fatal("malformed role accepted")
	}
	if _, e := SelectInit(t.TempDir(), "invalid"); e == nil {
		t.Fatal("invalid option accepted")
	}
}

func TestSanity_PolicyBranchValidationAgreesWithGit(t *testing.T) {
	for _, branch := range []string{"main", "trunk", "release/stable", "#tag", "@", "-dash", "unicøde", "", "/main", "main/", "a//b", "a..b", "a@{b", ".hidden", "a/.hidden", "a.lock", "a.lock/b", "a.", "a\\b", "a b", "a\t", "a~b", "a:b", "a?b", "a*b", "a[b", "a^b"} {
		gitOK := exec.Command("git", "check-ref-format", "refs/heads/"+branch).Run() == nil
		if got := Validate(Policy{Local, branch}) == nil; got != gitOK {
			t.Fatalf("branch %q: policy %v, git %v", branch, got, gitOK)
		}
	}
}

func TestAC231_UnitPositive_SourceAndAdopterSharePolicySelection(t *testing.T) {
	for _, roleName := range []string{"source", "adopter"} {
		for _, requested := range []string{"", Local, PR} {
			root := t.TempDir()
			put(t, root, ".clue/role.yaml", "role: "+roleName+"\n")
			want := Local
			if requested == PR {
				want = PR
			}
			if p, e := SelectInit(root, requested); e != nil || p != (Policy{want, "main"}) {
				t.Fatalf("%s/%s: %v %v", roleName, requested, p, e)
			}
			put(t, root, Path, string((Policy{want, "main"}).Bytes()))
			if p, _, e := Load(root); e != nil || p != (Policy{want, "main"}) {
				t.Fatalf("read %s/%s: %v %v", roleName, requested, p, e)
			}
		}
	}
}
