package idalloc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cliewen/cliewen/internal/ledger"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func coordinatedRepository(t *testing.T) (string, string) {
	t.Helper()
	base := t.TempDir()
	remote := filepath.Join(base, "remote.git")
	root := filepath.Join(base, "seed")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, base, "init", "--bare", remote)
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	l.MarkLive("CH-170")
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "seed")
	git(t, root, "remote", "add", "origin", remote)
	git(t, root, "push", "-u", "origin", "main")
	git(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	if err := Coordinate(root, "origin", 30*time.Second); err != nil {
		t.Fatalf("Coordinate: %v", err)
	}
	git(t, root, "add", ".clue/id-ledger.yaml")
	git(t, root, "commit", "-m", "coordinate")
	git(t, root, "push", "origin", "main")
	return root, remote
}

func TestAC171_IntegrationPositive_ConcurrentClonesReceiveUniqueSequentialIDs(t *testing.T) {
	_, remote := coordinatedRepository(t)
	const workers = 10
	roots := make([]string, workers)
	for i := range workers {
		roots[i] = filepath.Join(t.TempDir(), "clone")
		git(t, filepath.Dir(roots[i]), "clone", "--quiet", remote, roots[i])
	}
	var wg sync.WaitGroup
	ids := make(chan string, workers)
	errs := make(chan error, workers)
	for _, root := range roots {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := Allocate(root, "CH", "", 1, 30*time.Second)
			if err != nil {
				errs <- err
				return
			}
			ids <- got[0]
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for err := range errs {
		t.Errorf("Allocate: %v", err)
	}
	var got []string
	for id := range ids {
		got = append(got, id)
	}
	sort.Strings(got)
	want := []string{"CH-171", "CH-172", "CH-173", "CH-174", "CH-175", "CH-176", "CH-177", "CH-178", "CH-179", "CH-180"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("allocated IDs = %v, want %v\nremote claims:\n%s\nremote history:\n%s", got, want, git(t, remote, "show", Ref+":id-allocations.yaml"), git(t, remote, "log", "--oneline", "--decorate", Ref))
	}
}

func TestAC172_IntegrationPositive_ConcurrentWorktreesShareTheRemoteCAS(t *testing.T) {
	root, _ := coordinatedRepository(t)
	worktrees := []string{filepath.Join(t.TempDir(), "one"), filepath.Join(t.TempDir(), "two")}
	for i, worktree := range worktrees {
		git(t, root, "worktree", "add", "--quiet", "-b", fmt.Sprintf("worktree-%d", i), worktree, "main")
	}
	var wg sync.WaitGroup
	results := make(chan string, len(worktrees))
	errs := make(chan error, len(worktrees))
	for _, worktree := range worktrees {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids, err := Allocate(worktree, "CH", "", 1, 30*time.Second)
			if err != nil {
				errs <- err
				return
			}
			results <- ids[0]
		}()
	}
	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		t.Errorf("Allocate: %v", err)
	}
	var got []string
	for id := range results {
		got = append(got, id)
	}
	sort.Strings(got)
	if strings.Join(got, ",") != "CH-171,CH-172" {
		t.Fatalf("worktree allocations = %v", got)
	}
}

func TestSanity_GitUnionDriverCombinesConcurrentLedgerEvents(t *testing.T) {
	root := t.TempDir()
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	base := "version: 2\ncoordination:\n    mode: local\nevents:\n    - {id: CH-001, kind: numeric, state: live, prefix: CH, component: \"1\"}\n"
	ledgerPath := filepath.Join(root, ledger.DefaultPath)
	if err := os.WriteFile(ledgerPath, []byte(base), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "base ledger")
	for _, branch := range []string{"one", "two"} {
		git(t, root, "checkout", "-b", branch, "main")
		file, err := os.OpenFile(ledgerPath, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			t.Fatal(err)
		}
		component := "2"
		if branch == "two" {
			component = "3"
		}
		_, writeErr := fmt.Fprintf(file, "    - {id: CH-00%s, kind: numeric, state: reserved, prefix: CH, component: \"%s\"}\n", component, component)
		closeErr := file.Close()
		if writeErr != nil || closeErr != nil {
			t.Fatalf("append event: write=%v close=%v", writeErr, closeErr)
		}
		git(t, root, "add", ledger.DefaultPath)
		git(t, root, "commit", "-m", "event "+branch)
	}
	git(t, root, "checkout", "main")
	git(t, root, "merge", "--no-edit", "one")
	git(t, root, "merge", "--no-edit", "two")
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatalf("merged ledger: %v", err)
	}
	for _, id := range []string{"CH-002", "CH-003"} {
		if _, ok := l.Lookup(id); !ok {
			t.Fatalf("union merge lost %s", id)
		}
	}
}

func TestAC173_IntegrationPositive_BatchAndReadOnlySync(t *testing.T) {
	root, remote := coordinatedRepository(t)
	ids, err := Allocate(root, "CH", "", 3, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "CH-171,CH-172,CH-173" {
		t.Fatalf("batch = %v", ids)
	}
	reader := filepath.Join(t.TempDir(), "reader")
	git(t, filepath.Dir(reader), "clone", "--quiet", remote, reader)
	if err := Sync(reader, "", 30*time.Second); err != nil {
		t.Fatal(err)
	}
	l, err := ledger.Load(reader)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		if e, ok := l.Lookup(id); !ok || e.State != ledger.StateReserved {
			t.Fatalf("synced %s = %+v, ok=%v", id, e, ok)
		}
	}
}

func TestAC173_IntegrationNegative_SyncRefusesAnUncoordinatedLedger(t *testing.T) {
	root := t.TempDir()
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	if err := Sync(root, "", time.Second); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("Sync error = %v", err)
	}
}

func TestAC175_IntegrationNegative_MissingEstablishedRefFailsClosed(t *testing.T) {
	root, remote := coordinatedRepository(t)
	git(t, remote, "update-ref", "-d", Ref)
	before, err := os.ReadFile(filepath.Join(root, ledger.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Allocate(root, "CH", "", 1, 5*time.Second); err == nil || !strings.Contains(err.Error(), "restore it rather than recreating") {
		t.Fatalf("Allocate error = %v", err)
	}
	after, err := os.ReadFile(filepath.Join(root, ledger.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("failed coordinated allocation changed the local ledger")
	}
}

func TestAC175_IntegrationNegative_MalformedRemoteFailsClosed(t *testing.T) {
	root, remote := coordinatedRepository(t)
	git(t, remote, "update-ref", "-d", Ref)
	bad := filepath.Join(t.TempDir(), "bad")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, bad, "init", "-b", "main")
	git(t, bad, "config", "user.name", "Test User")
	git(t, bad, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(bad, "id-allocations.yaml"), []byte("version: 99\nclaims: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, bad, "add", ".")
	git(t, bad, "commit", "-m", "malformed allocator")
	badOID := git(t, bad, "rev-parse", "HEAD")
	git(t, bad, "push", remote, badOID+":"+Ref)
	assertAllocationFailsWithoutLocalWrite(t, root, "malformed")
}

func TestAC175_IntegrationNegative_UnreachableRemoteFailsClosed(t *testing.T) {
	root, _ := coordinatedRepository(t)
	assertAllocationFailsWithoutLocalWriteAt(t, root, "missing-remote", "")
}

func TestAC177_IntegrationPositive_SyncRecoversClaimAfterLocalSaveFailure(t *testing.T) {
	root, _ := coordinatedRepository(t)
	ids, err := allocate(root, "CH", "", 1, 30*time.Second, func(*ledger.Ledger) error {
		return fmt.Errorf("injected local save failure")
	})
	if err == nil || strings.Join(ids, ",") != "CH-171" {
		t.Fatalf("Allocate ids=%v error=%v", ids, err)
	}
	before, loadErr := ledger.Load(root)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if before.IsUsed("CH-171") {
		t.Fatal("failed local save nevertheless persisted CH-171")
	}
	if err := Sync(root, "", 30*time.Second); err != nil {
		t.Fatal(err)
	}
	after, loadErr := ledger.Load(root)
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if e, ok := after.Lookup("CH-171"); !ok || e.State != ledger.StateReserved {
		t.Fatalf("recovered claim = %+v, ok=%v", e, ok)
	}
}

func assertAllocationFailsWithoutLocalWrite(t *testing.T, root, want string) {
	t.Helper()
	assertAllocationFailsWithoutLocalWriteAt(t, root, "", want)
}

func assertAllocationFailsWithoutLocalWriteAt(t *testing.T, root, remote, want string) {
	t.Helper()
	before, err := os.ReadFile(filepath.Join(root, ledger.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = Allocate(root, "CH", remote, 1, 5*time.Second)
	if err == nil || (want != "" && !strings.Contains(err.Error(), want)) {
		t.Fatalf("Allocate error = %v, want text %q", err, want)
	}
	after, readErr := os.ReadFile(filepath.Join(root, ledger.DefaultPath))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(after) != string(before) {
		t.Fatal("failed coordinated allocation changed the local ledger")
	}
}
