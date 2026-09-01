package role

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, root, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, DefaultPath), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAC153_UnitPositive_DeclaredRolesLoad(t *testing.T) {
	for _, want := range []Role{Source, Adopter} {
		root := t.TempDir()
		write(t, root, "role: "+string(want)+"\n")
		got, declared, err := Load(root)
		if err != nil {
			t.Fatalf("%s: %v", want, err)
		}
		if !declared {
			t.Fatalf("%s: declared = false, want true", want)
		}
		if got != want {
			t.Fatalf("role = %q, want %q", got, want)
		}
		if !Exists(root) {
			t.Fatalf("%s: Exists = false, want true", want)
		}
	}
}

// An undeclared repository is an adopter. Every repository adopted before
// the marker existed is one, and none of them can be expected to carry a
// file that did not exist when they were onboarded.
func TestAC153_UnitPositive_UndeclaredRepositoryIsAnAdopter(t *testing.T) {
	root := t.TempDir()
	got, declared, err := Load(root)
	if err != nil {
		t.Fatalf("a missing marker must not be an error: %v", err)
	}
	if declared {
		t.Fatal("declared = true, want false")
	}
	if got != Adopter {
		t.Fatalf("role = %q, want %q", got, Adopter)
	}
	if Exists(root) {
		t.Fatal("Exists = true for a repository with no marker")
	}
}

func TestAC153_UnitNegative_UnknownOrMalformedRoleFails(t *testing.T) {
	for name, body := range map[string]string{
		"unknown value": "role: upstream\n",
		"empty value":   "role:\n",
		"not yaml":      "role: [source\n",
	} {
		root := t.TempDir()
		write(t, root, body)
		if _, _, err := Load(root); err == nil {
			t.Fatalf("%s: Load succeeded, want an error", name)
		}
	}
}

func TestAC153_UnitPositive_WriteRoundTripsThroughLoad(t *testing.T) {
	root := t.TempDir()
	if err := Write(root, Source); err != nil {
		t.Fatal(err)
	}
	got, declared, err := Load(root)
	if err != nil || !declared || got != Source {
		t.Fatalf("Load after Write = (%q, %v, %v), want (source, true, nil)", got, declared, err)
	}
	data, err := os.ReadFile(filepath.Join(root, DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	want, err := Bytes(Source)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(want) {
		t.Fatalf("written bytes = %q, want %q", data, want)
	}
}

func TestAC153_UnitNegative_AnInvalidRoleIsNeverRendered(t *testing.T) {
	if _, err := Bytes(Role("upstream")); err == nil {
		t.Fatal("Bytes rendered an invalid role")
	}
	if err := Write(t.TempDir(), Role("")); err == nil {
		t.Fatal("Write accepted an invalid role")
	}
}
