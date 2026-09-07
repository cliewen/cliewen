// Package idalloc coordinates numeric identity claims through an ordinary
// fast-forward-only Git branch. The checked-in ledger remains the lifecycle
// register; this package's remote journal contains permanent claims only.
package idalloc

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/cliewen/cliewen/internal/ledger"
	"gopkg.in/yaml.v3"
)

const Ref = "refs/heads/clue/id-allocator"

type claimFile struct {
	Version int            `yaml:"version"`
	Claims  []ledger.Entry `yaml:"claims"`
}

type state struct {
	byID map[string]ledger.Entry
}

func newState() *state { return &state{byID: map[string]ledger.Entry{}} }

func (s *state) merge(claims []ledger.Entry) error {
	for _, e := range claims {
		e.State = ledger.StateReserved
		if e.Kind != ledger.KindNumeric || !ledger.ValidNumericEntry(e) {
			return fmt.Errorf("allocator claim %s is not a valid numeric identity", e.ID)
		}
		if old, ok := s.byID[e.ID]; ok {
			if old.Prefix != e.Prefix || old.Component.Cmp(e.Component) != 0 {
				return fmt.Errorf("allocator has conflicting metadata for %s", e.ID)
			}
			continue
		}
		s.byID[e.ID] = e
	}
	return nil
}

func (s *state) claims() []ledger.Entry {
	ids := make([]string, 0, len(s.byID))
	for id := range s.byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]ledger.Entry, 0, len(ids))
	for _, id := range ids {
		out = append(out, s.byID[id])
	}
	return out
}

func (s *state) next(prefix string, count int) ([]string, error) {
	if !ledger.ValidNumericPrefix(prefix) {
		return nil, fmt.Errorf("prefix %q is not a canonical numeric prefix", prefix)
	}
	if count < 1 {
		return nil, fmt.Errorf("count must be positive")
	}
	n := new(big.Int)
	for _, e := range s.byID {
		if e.Prefix == prefix && e.Component.Cmp(n) > 0 {
			n.Set(e.Component)
		}
	}
	ids := make([]string, 0, count)
	for len(ids) < count {
		n.Add(n, big.NewInt(1))
		id := fmt.Sprintf("%s-%03d", prefix, n)
		if _, used := s.byID[id]; used {
			continue
		}
		e := ledger.Entry{ID: id, Kind: ledger.KindNumeric, State: ledger.StateReserved, Prefix: prefix, Component: new(big.Int).Set(n)}
		s.byID[id] = e
		ids = append(ids, id)
	}
	return ids, nil
}

func marshalState(s *state) ([]byte, error) {
	return yaml.Marshal(claimFile{Version: 1, Claims: s.claims()})
}

func parseState(data []byte) (*state, error) {
	var f claimFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	if f.Version != 1 {
		return nil, fmt.Errorf("allocator protocol version must be 1, not %d", f.Version)
	}
	s := newState()
	if err := s.merge(f.Claims); err != nil {
		return nil, err
	}
	return s, nil
}

type gitError struct {
	args []string
	out  string
	err  error
}

func (e *gitError) Error() string {
	if e.out != "" {
		return fmt.Sprintf("git %s: %s", strings.Join(e.args, " "), e.out)
	}
	return fmt.Sprintf("git %s: %v", strings.Join(e.args, " "), e.err)
}

func runGit(ctx context.Context, root string, input []byte, env []string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{"-C", root}, args...)
	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Env = append(os.Environ(), env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		out := strings.TrimSpace(stderr.String())
		if out == "" {
			out = strings.TrimSpace(stdout.String())
		}
		return nil, &gitError{args: args, out: out, err: err}
	}
	return stdout.Bytes(), nil
}

func token() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("create allocator transaction identity: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func remoteOID(ctx context.Context, root, remote string) (string, error) {
	out, err := runGit(ctx, root, nil, nil, "ls-remote", "--refs", remote, Ref)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "", nil
	}
	if len(fields) != 2 || fields[1] != Ref {
		return "", fmt.Errorf("unexpected allocator ref response from %s", remote)
	}
	return fields[0], nil
}

func fetch(ctx context.Context, root, remote string) (string, *state, error) {
	nonce, err := token()
	if err != nil {
		return "", nil, err
	}
	tmpRef := "refs/clue/tmp/id-allocator-" + nonce
	defer func() { _, _ = runGit(context.Background(), root, nil, nil, "update-ref", "-d", tmpRef) }()
	if _, err := runGit(ctx, root, nil, nil, "fetch", "--quiet", "--no-tags", remote, "+"+Ref+":"+tmpRef); err != nil {
		return "", nil, err
	}
	oidBytes, err := runGit(ctx, root, nil, nil, "rev-parse", tmpRef)
	if err != nil {
		return "", nil, err
	}
	oid := strings.TrimSpace(string(oidBytes))
	data, err := runGit(ctx, root, nil, nil, "show", oid+":id-allocations.yaml")
	if err != nil {
		return "", nil, err
	}
	s, err := parseState(data)
	if err != nil {
		return "", nil, fmt.Errorf("remote allocator is malformed: %w", err)
	}
	return oid, s, nil
}

func commit(ctx context.Context, root, parent, message string, s *state) (string, error) {
	nonce, err := token()
	if err != nil {
		return "", err
	}
	data, err := marshalState(s)
	if err != nil {
		return "", err
	}
	blob, err := runGit(ctx, root, data, nil, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	treeLine := fmt.Sprintf("100644 blob %s\tid-allocations.yaml\n", strings.TrimSpace(string(blob)))
	tree, err := runGit(ctx, root, []byte(treeLine), nil, "mktree")
	if err != nil {
		return "", err
	}
	args := []string{"commit-tree", strings.TrimSpace(string(tree))}
	if parent != "" {
		args = append(args, "-p", parent)
	}
	env := []string{
		"GIT_AUTHOR_NAME=Cliewen identity allocator",
		"GIT_AUTHOR_EMAIL=allocator@cliewen.local",
		"GIT_COMMITTER_NAME=Cliewen identity allocator",
		"GIT_COMMITTER_EMAIL=allocator@cliewen.local",
	}
	oid, err := runGit(ctx, root, []byte(message+"\n\ntransaction: "+nonce+"\n"), env, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(oid)), nil
}

func push(ctx context.Context, root, remote, oid string) error {
	_, err := runGit(ctx, root, nil, nil, "push", "--porcelain", remote, oid+":"+Ref)
	return err
}

func withTimeout(timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := fn(ctx); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("allocator operation timed out after %s", timeout)
		}
		return err
	}
	return nil
}

// Coordinate initializes the remote journal and enables Git coordination in
// the checked-in ledger only after the remote claim is durable.
func Coordinate(root, remote string, timeout time.Duration) error {
	return withTimeout(timeout, func(ctx context.Context) error {
		l, err := ledger.Load(root)
		if err != nil {
			return err
		}
		if l.Version() != 2 {
			return fmt.Errorf("identity ledger must be migrated to version 2 before coordination")
		}
		for {
			oid, err := remoteOID(ctx, root, remote)
			if err != nil {
				return err
			}
			s := newState()
			if oid != "" {
				var fetched string
				fetched, s, err = fetch(ctx, root, remote)
				if err != nil {
					return err
				}
				oid = fetched
			}
			before := len(s.byID)
			if err := s.merge(l.Claims()); err != nil {
				return err
			}
			if oid == "" || len(s.byID) != before {
				newOID, err := commit(ctx, root, oid, "clue: initialize identity allocator", s)
				if err != nil {
					return err
				}
				if err := push(ctx, root, remote, newOID); err != nil {
					now, probeErr := remoteOID(ctx, root, remote)
					if probeErr == nil && now != oid {
						continue
					}
					return err
				}
			}
			if err := l.MergeClaims(s.claims()); err != nil {
				return err
			}
			if err := l.SetGitCoordination(remote); err != nil {
				return err
			}
			return l.Save()
		}
	})
}

// Sync imports permanent remote claims into the checked-in ledger without
// changing any existing lifecycle state.
func Sync(root, remoteOverride string, timeout time.Duration) error {
	return withTimeout(timeout, func(ctx context.Context) error {
		l, err := ledger.Load(root)
		if err != nil {
			return err
		}
		coord := l.Coordination()
		if coord.Mode != "git" {
			return fmt.Errorf("identity ledger is not configured for Git coordination")
		}
		remote := coord.Remote
		if remoteOverride != "" {
			remote = remoteOverride
		}
		oid, err := remoteOID(ctx, root, remote)
		if err != nil {
			return err
		}
		if oid == "" {
			return fmt.Errorf("coordinated allocator ref %s is missing; restore it rather than recreating it", Ref)
		}
		_, s, err := fetch(ctx, root, remote)
		if err != nil {
			return err
		}
		if err := l.MergeClaims(s.claims()); err != nil {
			return err
		}
		return l.Save()
	})
}

// Allocate atomically claims count IDs through the configured remote. It
// returns the remotely durable IDs even when the subsequent local save fails,
// allowing the caller to report the recovery boundary precisely.
func Allocate(root, prefix, remoteOverride string, count int, timeout time.Duration) (ids []string, err error) {
	return allocate(root, prefix, remoteOverride, count, timeout, func(l *ledger.Ledger) error { return l.Save() })
}

func allocate(root, prefix, remoteOverride string, count int, timeout time.Duration, save func(*ledger.Ledger) error) (ids []string, err error) {
	err = withTimeout(timeout, func(ctx context.Context) error {
		l, loadErr := ledger.Load(root)
		if loadErr != nil {
			return loadErr
		}
		coord := l.Coordination()
		if coord.Mode != "git" {
			return fmt.Errorf("identity ledger is not configured for Git coordination")
		}
		remote := coord.Remote
		if remoteOverride != "" {
			remote = remoteOverride
		}
		for {
			oid, remoteErr := remoteOID(ctx, root, remote)
			if remoteErr != nil {
				return remoteErr
			}
			if oid == "" {
				return fmt.Errorf("coordinated allocator ref %s is missing; restore it rather than recreating it", Ref)
			}
			fetched, s, remoteErr := fetch(ctx, root, remote)
			if remoteErr != nil {
				return remoteErr
			}
			oid = fetched
			if remoteErr = s.merge(l.Claims()); remoteErr != nil {
				return remoteErr
			}
			ids, remoteErr = s.next(prefix, count)
			if remoteErr != nil {
				return remoteErr
			}
			message := "clue: reserve " + strings.Join(ids, ", ")
			newOID, remoteErr := commit(ctx, root, oid, message, s)
			if remoteErr != nil {
				return remoteErr
			}
			if remoteErr = push(ctx, root, remote, newOID); remoteErr != nil {
				now, probeErr := remoteOID(ctx, root, remote)
				if probeErr == nil && now != oid {
					ids = nil
					continue
				}
				return remoteErr
			}
			if remoteErr = l.MergeClaims(s.claims()); remoteErr != nil {
				return remoteErr
			}
			return save(l)
		}
	})
	return ids, err
}
