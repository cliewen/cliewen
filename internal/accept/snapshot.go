package accept

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type snapshotNode struct {
	mode, oid, target string
	children          []string
}

func safeSnapshotPath(name string) bool {
	if !filepath.IsLocal(filepath.FromSlash(name)) || strings.Contains(name, "\\") {
		return false
	}
	for _, part := range strings.Split(name, "/") {
		if part == "" || part == "." || part == ".." || strings.EqualFold(part, ".git") {
			return false
		}
	}
	return true
}

// Materialization reads committed blobs only. Internal link targets are resolved
// in this revision's virtual graph and copied as ordinary snapshot content;
// OS symlinks, checkout filters and live working-tree targets are never used.
func materialize(root, revision, destination string) error {
	entries, err := command(root, "ls-tree", "-rz", revision).Output()
	if err != nil {
		return err
	}
	nodes := map[string]*snapshotNode{"": {mode: "tree"}}
	var names []string
	for _, entry := range bytes.Split(entries, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		meta, name, ok := strings.Cut(string(entry), "\t")
		fields := strings.Fields(meta)
		if !ok || len(fields) != 3 || fields[1] != "blob" || (fields[0] != "100644" && fields[0] != "100755" && fields[0] != "120000") {
			return fmt.Errorf("local acceptance does not support submodules or entry %q", name)
		}
		if !safeSnapshotPath(name) {
			return fmt.Errorf("unsafe tracked path %q", name)
		}
		if _, exists := nodes[name]; exists {
			return fmt.Errorf("duplicate tracked path %q", name)
		}
		nodes[name] = &snapshotNode{mode: fields[0], oid: fields[2]}
		names = append(names, name)
		parent := path.Dir(name)
		if parent == "." {
			parent = ""
		}
		for parent != "" {
			if existing := nodes[parent]; existing != nil && existing.mode != "tree" {
				return fmt.Errorf("tracked path traverses non-directory %q", parent)
			}
			if nodes[parent] == nil {
				nodes[parent] = &snapshotNode{mode: "tree"}
			}
			parent = path.Dir(parent)
			if parent == "." {
				parent = ""
			}
		}
	}
	for name := range nodes {
		if name == "" {
			continue
		}
		parent := path.Dir(name)
		if parent == "." {
			parent = ""
		}
		if nodes[parent] == nil || nodes[parent].mode != "tree" {
			return fmt.Errorf("tracked path traverses non-directory %q", parent)
		}
		nodes[parent].children = append(nodes[parent].children, name)
	}
	for _, node := range nodes {
		sort.Strings(node.children)
	}
	c := command(root, "cat-file", "--batch")
	input, err := c.StdinPipe()
	if err != nil {
		return err
	}
	output, err := c.StdoutPipe()
	if err != nil {
		input.Close()
		return err
	}
	if err = c.Start(); err != nil {
		input.Close()
		return err
	}
	blobsRead := false
	defer func() {
		input.Close()
		if !blobsRead && c.Process != nil {
			_ = c.Process.Kill()
		}
		_ = c.Wait()
	}()
	reader := bufio.NewReader(output)
	for _, name := range names {
		node := nodes[name]
		if _, err = io.WriteString(input, node.oid+"\n"); err != nil {
			return err
		}
		header, e := reader.ReadString('\n')
		if e != nil {
			return e
		}
		info := strings.Fields(header)
		if len(info) != 3 || info[0] != node.oid || info[1] != "blob" {
			return fmt.Errorf("missing candidate blob %s", node.oid)
		}
		size, e := strconv.ParseInt(info[2], 10, 64)
		if e != nil || size < 0 {
			return fmt.Errorf("invalid candidate blob size")
		}
		if node.mode == "120000" {
			if size == 0 || size > 65536 {
				return fmt.Errorf("invalid internal link target at %q", name)
			}
			var target bytes.Buffer
			if _, err = io.CopyN(&target, reader, size); err != nil {
				return err
			}
			node.target = target.String()
		} else {
			filename := filepath.Join(destination, filepath.FromSlash(name))
			if err = os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
				return err
			}
			f, e := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if e != nil {
				return e
			}
			_, err = io.CopyN(f, reader, size)
			closeErr := f.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
		}
		if b, e := reader.ReadByte(); e != nil || b != '\n' {
			return fmt.Errorf("invalid candidate blob terminator")
		}
	}
	blobsRead = true
	var resolve func(string, map[string]bool) (string, error)
	resolve = func(name string, chain map[string]bool) (string, error) {
		if name == "" {
			return "", nil
		}
		parts := strings.Split(name, "/")
		for i := range parts {
			prefix := strings.Join(parts[:i+1], "/")
			node := nodes[prefix]
			if node == nil {
				return "", fmt.Errorf("missing internal link target %q", name)
			}
			if node.mode == "120000" {
				if chain[prefix] {
					return "", fmt.Errorf("internal link cycle at %q", prefix)
				}
				target := node.target
				if strings.ContainsAny(target, "\\\x00") || path.IsAbs(target) || (len(target) > 1 && target[1] == ':') {
					return "", fmt.Errorf("external or unsafe link target at %q", prefix)
				}
				target = path.Clean(path.Join(path.Dir(prefix), target))
				if target == "." {
					target = ""
				}
				if target != "" && !safeSnapshotPath(target) {
					return "", fmt.Errorf("external or unsafe link target at %q", prefix)
				}
				if i+1 < len(parts) {
					target = path.Join(target, strings.Join(parts[i+1:], "/"))
				}
				chain[prefix] = true
				resolved, e := resolve(target, chain)
				delete(chain, prefix)
				return resolved, e
			}
			if i+1 < len(parts) && node.mode != "tree" {
				return "", fmt.Errorf("link target traverses non-directory %q", prefix)
			}
		}
		return name, nil
	}
	remaining := (len(nodes) + 1) * 32
	var copyNode func(string, string, map[string]bool) error
	copyNode = func(source, target string, stack map[string]bool) error {
		remaining--
		if remaining < 0 {
			return fmt.Errorf("internal link expansion exceeds snapshot budget")
		}
		canonical, e := resolve(source, map[string]bool{})
		if e != nil {
			return e
		}
		if stack[canonical] {
			return fmt.Errorf("internal directory link cycle at %q", source)
		}
		stack[canonical] = true
		defer delete(stack, canonical)
		node := nodes[canonical]
		filename := filepath.Join(destination, filepath.FromSlash(target))
		if node.mode == "tree" {
			if e = os.MkdirAll(filename, 0700); e != nil {
				return e
			}
			for _, child := range node.children {
				if e = copyNode(child, path.Join(target, path.Base(child)), stack); e != nil {
					return e
				}
			}
			return nil
		}
		if e = os.MkdirAll(filepath.Dir(filename), 0700); e != nil {
			return e
		}
		from, e := os.Open(filepath.Join(destination, filepath.FromSlash(canonical)))
		if e != nil {
			return e
		}
		defer from.Close()
		to, e := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, e = io.Copy(to, from)
		closeErr := to.Close()
		if e != nil {
			return e
		}
		return closeErr
	}
	for _, name := range names {
		if nodes[name].mode == "120000" {
			if err = copyNode(name, name, map[string]bool{}); err != nil {
				return err
			}
		}
	}
	return nil
}
