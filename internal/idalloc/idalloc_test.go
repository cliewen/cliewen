package idalloc

import (
	"context"
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
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
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

func TestSanity_CoordinateRequiresLedgerUnionMergeRule(t *testing.T) {
	root := t.TempDir()
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	if err := Coordinate(root, "origin", time.Second); err == nil || !strings.Contains(err.Error(), "merge=union") {
		t.Fatalf("Coordinate error = %v", err)
	}
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
		if err := os.WriteFile(filepath.Join(worktree, "worktree-marker"), []byte(fmt.Sprintf("marker-%d", i)), 0o644); err != nil {
			t.Fatal(err)
		}
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
	for i, worktree := range worktrees {
		marker, err := os.ReadFile(filepath.Join(worktree, "worktree-marker"))
		if err != nil || string(marker) != fmt.Sprintf("marker-%d", i) {
			t.Fatalf("worktree %d marker changed: %q, %v", i, marker, err)
		}
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
	remoteBefore := git(t, remote, "rev-parse", Ref)
	if err := Sync(reader, "", 30*time.Second); err != nil {
		t.Fatal(err)
	}
	if remoteAfter := git(t, remote, "rev-parse", Ref); remoteAfter != remoteBefore {
		t.Fatalf("sync changed allocator head from %s to %s", remoteBefore, remoteAfter)
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

func TestAC176_IntegrationPositive_MigratedHighWaterSurvivesCoordination(t *testing.T) {
	base := t.TempDir()
	remote := filepath.Join(base, "remote.git")
	root := filepath.Join(base, "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, base, "init", "--bare", remote)
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "counters:\n    CH: \"9\"\nentries:\n    - {id: CH-007, kind: numeric, state: retired, prefix: CH, component: \"7\"}\n"
	if err := os.WriteFile(filepath.Join(root, ledger.DefaultPath), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "legacy ledger")
	git(t, root, "remote", "add", "origin", remote)
	git(t, root, "push", "-u", "origin", "main")
	git(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	l.ConvertV2()
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ledger.DefaultPath)
	git(t, root, "commit", "-m", "migrate ledger")
	git(t, root, "push", "origin", "main")
	if err := Coordinate(root, "origin", 30*time.Second); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ledger.DefaultPath)
	git(t, root, "commit", "-m", "coordinate")
	git(t, root, "push", "origin", "main")
	ids, err := Allocate(root, "CH", "", 1, 30*time.Second)
	if err != nil || strings.Join(ids, ",") != "CH-010" {
		t.Fatalf("coordinated allocation after migration = %v, %v", ids, err)
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
	}, push)
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

func TestAC177_IntegrationPositive_AmbiguousPushSuccessDoesNotAllocateAgain(t *testing.T) {
	root, remote := coordinatedRepository(t)
	pushes := 0
	lostResponse := func(ctx context.Context, root, remote, oid string) error {
		pushes++
		if err := push(ctx, root, remote, oid); err != nil {
			return err
		}
		return fmt.Errorf("injected lost push response")
	}
	ids, err := allocate(root, "CH", "", 1, 30*time.Second, func(l *ledger.Ledger) error { return l.Save() }, lostResponse)
	if err != nil || strings.Join(ids, ",") != "CH-171" {
		t.Fatalf("Allocate ids=%v error=%v", ids, err)
	}
	if pushes != 1 {
		t.Fatalf("push attempts = %d, want 1", pushes)
	}
	claims := git(t, remote, "show", Ref+":id-allocations.yaml")
	if strings.Count(claims, "id: CH-171") != 1 || strings.Contains(claims, "id: CH-172") {
		t.Fatalf("ambiguous success created another claim:\n%s", claims)
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

// The allocate path imported remote claims but not the remote high-water mark,
// so a boundary carried by no claim — a migrated legacy counter whose artifact
// was deleted in another clone — never reached the checked-in ledger. Nothing
// visibly broke while allocation stayed coordinated, because the remote is
// re-read every time; the damage showed only when something fell back to local
// allocation and reissued a number the allocator had already burned.
func TestAC170_IntegrationPositive_AllocateImportsRemoteHighWater(t *testing.T) {
	base := t.TempDir()
	remote := filepath.Join(base, "remote.git")
	root := filepath.Join(base, "repo")
	if err := os.MkdirAll(filepath.Join(root, ".clue"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, base, "init", "--bare", remote)
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")

	// A legacy counter for AC stands at 50 with no surviving AC artifact, so
	// only the high-water record carries that boundary.
	legacy := "counters:\n    AC: \"50\"\n    CH: \"170\"\nentries:\n    - {id: CH-170, kind: numeric, state: live, prefix: CH, component: \"170\"}\n"
	if err := os.WriteFile(filepath.Join(root, ledger.DefaultPath), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "legacy ledger")
	git(t, root, "remote", "add", "origin", remote)
	git(t, root, "push", "-u", "origin", "main")
	git(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	l.ConvertV2()
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	if err := Coordinate(root, "origin", 30*time.Second); err != nil {
		t.Fatal(err)
	}

	// Another contributor's checked-in ledger never carried the AC boundary.
	withoutHighWater := "version: 2\ncoordination:\n    mode: git\n    remote: origin\nevents:\n    - {id: CH-170, kind: numeric, state: live, prefix: CH, component: \"170\"}\n"
	if err := os.WriteFile(filepath.Join(root, ledger.DefaultPath), []byte(withoutHighWater), 0o644); err != nil {
		t.Fatal(err)
	}

	ids, err := Allocate(root, "CH", "", 1, 30*time.Second)
	if err != nil || strings.Join(ids, ",") != "CH-171" {
		t.Fatalf("Allocate = %v, %v; want CH-171", ids, err)
	}
	reloaded, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	next, err := reloaded.NextNumeric("AC")
	if err != nil {
		t.Fatal(err)
	}
	if next != "AC-051" {
		t.Fatalf("local allocation after coordinated allocation = %s, want AC-051; the remote high-water boundary was not imported", next)
	}
}

// Allocation across worktrees is proven elsewhere; what was never proven is
// that the two branches then merge. Each worktree commits its own ledger and
// both land on main through a real three-way merge with the union driver in
// effect, which is the whole point of the checked-in event log.
func TestAC172_IntegrationPositive_ParallelWorktreeAllocationsMergeCleanly(t *testing.T) {
	root, _ := coordinatedRepository(t)
	base := t.TempDir()
	var allocated []string
	for i, name := range []string{"one", "two"} {
		worktree := filepath.Join(base, name)
		git(t, root, "worktree", "add", "--quiet", "-b", "branch-"+name, worktree, "main")
		ids, err := Allocate(worktree, "CH", "", 1, 30*time.Second)
		if err != nil {
			t.Fatalf("worktree %d: %v", i, err)
		}
		allocated = append(allocated, ids[0])
		git(t, worktree, "add", ledger.DefaultPath)
		git(t, worktree, "commit", "-m", "allocate in "+name)
	}
	sort.Strings(allocated)
	if strings.Join(allocated, ",") != "CH-171,CH-172" {
		t.Fatalf("worktree allocations = %v", allocated)
	}

	git(t, root, "merge", "--no-edit", "branch-one")
	git(t, root, "merge", "--no-edit", "branch-two")

	merged, err := ledger.Load(root)
	if err != nil {
		t.Fatalf("merged ledger: %v", err)
	}
	if merged.Damage() != "" {
		t.Fatalf("merging two parallel allocations damaged the ledger: %s", merged.Damage())
	}
	for _, id := range allocated {
		if e, ok := merged.Lookup(id); !ok || e.State != ledger.StateReserved {
			t.Fatalf("%s = %+v, ok=%v after merge", id, e, ok)
		}
	}
	next, err := merged.NextNumeric("CH")
	if err != nil {
		t.Fatal(err)
	}
	if next != "CH-173" {
		t.Fatalf("next after merging both allocations = %s, want CH-173", next)
	}
	for _, name := range []string{"one", "two"} {
		git(t, root, "worktree", "remove", "--force", filepath.Join(base, name))
	}
}

// The hazard `clue id next` warns about, demonstrated rather than assumed: in
// local mode two worktrees allocate the same number. The ledger itself
// survives — the identical events fold to one entry — so the collision has to
// surface as two artifacts holding one ID, which is the corpus rule's job.
func TestSanity_LocalModeWorktreesCollideYetTheLedgerStillMerges(t *testing.T) {
	root := t.TempDir()
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	l.MarkLive("CH-001")
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "base")

	base := t.TempDir()
	var allocated []string
	for _, name := range []string{"one", "two"} {
		worktree := filepath.Join(base, name)
		git(t, root, "worktree", "add", "--quiet", "-b", "branch-"+name, worktree, "main")
		wl, loadErr := ledger.Load(worktree)
		if loadErr != nil {
			t.Fatal(loadErr)
		}
		id, nextErr := wl.NextNumeric("CH")
		if nextErr != nil {
			t.Fatal(nextErr)
		}
		if saveErr := wl.Save(); saveErr != nil {
			t.Fatal(saveErr)
		}
		allocated = append(allocated, id)
		git(t, worktree, "add", ledger.DefaultPath)
		git(t, worktree, "commit", "-m", "allocate in "+name)
	}
	if allocated[0] != allocated[1] {
		t.Fatalf("local-mode worktrees allocated %v; the collision this warns about did not happen", allocated)
	}

	git(t, root, "merge", "--no-edit", "branch-one")
	git(t, root, "merge", "--no-edit", "branch-two")
	merged, err := ledger.Load(root)
	if err != nil {
		t.Fatalf("the ledger itself must survive a local-mode collision: %v", err)
	}
	if _, ok := merged.Lookup(allocated[0]); !ok {
		t.Fatalf("merged ledger lost %s", allocated[0])
	}
	for _, name := range []string{"one", "two"} {
		git(t, root, "worktree", "remove", "--force", filepath.Join(base, name))
	}
}

// unionRepository is a repository with the union rule installed, a version-two
// ledger, and however many bare remotes the caller asks for.
func unionRepository(t *testing.T, remotes ...string) string {
	t.Helper()
	base := t.TempDir()
	root := filepath.Join(base, "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-b", "main")
	git(t, root, "config", "user.name", "Test User")
	git(t, root, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(root, ".gitattributes"), []byte(ledger.UnionAttribute+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	l, err := ledger.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	l.MarkLive("CH-001")
	if err := l.Save(); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-m", "base")
	for i, name := range remotes {
		bare := filepath.Join(base, name+".git")
		git(t, base, "init", "--bare", bare)
		git(t, root, "remote", "add", name, bare)
		if i == 0 {
			git(t, root, "push", "-u", name, "main")
		} else {
			git(t, root, "push", name, "main")
		}
		git(t, bare, "symbolic-ref", "HEAD", "refs/heads/main")
	}
	return root
}

// coordinateOnBranch enables coordination against remote on a new branch cut
// from main, and commits the ledger change.
func coordinateOnBranch(t *testing.T, root, branch, remote string) {
	t.Helper()
	git(t, root, "checkout", "-q", "-b", branch, "main")
	if err := Coordinate(root, remote, 30*time.Second); err != nil {
		t.Fatalf("Coordinate(%s): %v", remote, err)
	}
	git(t, root, "add", ledger.DefaultPath)
	git(t, root, "commit", "-m", "coordinate through "+remote)
}

// Two branches enabling coordination against the same remote write the same
// lines, so there is nothing for the union driver to combine.
func TestSanity_IdenticalCoordinationOnTwoBranchesMergesCleanly(t *testing.T) {
	root := unionRepository(t, "origin")
	coordinateOnBranch(t, root, "one", "origin")
	coordinateOnBranch(t, root, "two", "origin")
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "--no-edit", "one")
	git(t, root, "merge", "--no-edit", "two")

	l, err := ledger.Load(root)
	if err != nil {
		t.Fatalf("merged ledger: %v", err)
	}
	if l.Damage() != "" {
		t.Fatalf("identical coordination reported damage: %s", l.Damage())
	}
	if c := l.Coordination(); c.Mode != "git" || c.Remote != "origin" {
		t.Fatalf("coordination = %+v, want git through origin", c)
	}
}

// The damage a real merge actually produces. Union merge combines line by line,
// so `version` and `mode` are identical on both sides and merge as context: the
// result is not a repeated header but a repeated `remote` inside coordination.
// Two branches each coordinating to their own remote is a disagreement no
// reading can settle, so the ledger fails closed and names the choice.
func TestSanity_DisagreeingRemotesMergeIntoADecidableFailure(t *testing.T) {
	root := unionRepository(t, "origin", "other")
	coordinateOnBranch(t, root, "one", "origin")
	coordinateOnBranch(t, root, "two", "other")
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "--no-edit", "one")
	git(t, root, "merge", "--no-edit", "two")

	data, err := os.ReadFile(filepath.Join(root, ledger.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	// The shape this asserts is Git's, not ours: if a future Git combines the
	// file differently, this is the test that says so.
	if got := strings.Count(string(data), "remote:"); got != 2 {
		t.Fatalf("union merge produced %d remote lines, want 2:\n%s", got, data)
	}
	if got := strings.Count(string(data), "version: 2"); got != 1 {
		t.Fatalf("union merge repeated the header %d times; the damage shape has changed:\n%s", got, data)
	}

	_, err = ledger.Load(root)
	if err == nil {
		t.Fatal("a ledger naming two different allocator remotes was accepted")
	}
	for _, want := range []string{"two different allocator remotes", "origin", "other", "only a person can decide"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want it to mention %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "unmarshal errors") || strings.Contains(err.Error(), "already defined") {
		t.Fatalf("error = %q, want a decidable message rather than a parser one", err)
	}
}

// A whole-file conversion on one branch against an append on another still
// merges cleanly, because Git matches the lines both sides share. This is the
// case the deterministic backfill ordering protects: were the converted file
// written in a different order each run, the whole file would become one
// conflicting hunk and the union driver would concatenate both copies.
func TestSanity_WholeFileConversionMergesWithAConcurrentAppend(t *testing.T) {
	root := unionRepository(t)
	ledgerPath := filepath.Join(root, ledger.DefaultPath)
	base, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}

	git(t, root, "checkout", "-q", "-b", "converted", "main")
	converted := strings.Replace(string(base), "mode: local", "mode: local\n    # rewritten in place", 1)
	if err := os.WriteFile(ledgerPath, []byte(converted), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "commit", "-am", "rewrite the ledger")

	git(t, root, "checkout", "-q", "-b", "appended", "main")
	appended := string(base) + "    - {id: CH-002, kind: numeric, state: reserved, prefix: CH, component: \"2\"}\n"
	if err := os.WriteFile(ledgerPath, []byte(appended), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, root, "commit", "-am", "append an event")

	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "--no-edit", "converted")
	git(t, root, "merge", "--no-edit", "appended")

	l, err := ledger.Load(root)
	if err != nil {
		t.Fatalf("merged ledger: %v", err)
	}
	if l.Damage() != "" {
		t.Fatalf("a rewrite against an append damaged the ledger: %s", l.Damage())
	}
	for _, id := range []string{"CH-001", "CH-002"} {
		if _, ok := l.Lookup(id); !ok {
			t.Fatalf("merge lost %s", id)
		}
	}
}
