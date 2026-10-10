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
	resolutionRemaining := (len(nodes) + 1) * 32
	var walk func(string, []string, map[string]bool) (string, error)
	walk = func(current string, parts []string, chain map[string]bool) (string, error) {
		for _, part := range parts {
			resolutionRemaining--
			if resolutionRemaining < 0 {
				return "", fmt.Errorf("internal link resolution exceeds snapshot budget")
			}
			if nodes[current] == nil || nodes[current].mode != "tree" {
				return "", fmt.Errorf("link target traverses non-directory %q", current)
			}
			switch part {
			case "", ".":
				continue
			case "..":
				if current == "" {
					return "", fmt.Errorf("external link target escapes snapshot")
				}
				current = path.Dir(current)
				if current == "." {
					current = ""
				}
				continue
			}
			if strings.EqualFold(part, ".git") || !safeSnapshotPath(part) {
				return "", fmt.Errorf("unsafe internal link component %q", part)
			}
			prefix := path.Join(current, part)
			node := nodes[prefix]
			if node == nil {
				return "", fmt.Errorf("missing internal link target %q", prefix)
			}
			if node.mode != "120000" {
				current = prefix
				continue
			}
			if chain[prefix] {
				return "", fmt.Errorf("internal link cycle at %q", prefix)
			}
			target := node.target
			if strings.ContainsAny(target, "\\\x00") || path.IsAbs(target) || (len(target) > 1 && target[1] == ':') {
				return "", fmt.Errorf("external or unsafe link target at %q", prefix)
			}
			chain[prefix] = true
			resolved, e := walk(current, strings.Split(target, "/"), chain)
			delete(chain, prefix)
			if e != nil {
				return "", e
			}
			current = resolved
		}
		return current, nil
	}
	resolve := func(name string, chain map[string]bool) (string, error) {
		return walk("", strings.Split(name, "/"), chain)
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
