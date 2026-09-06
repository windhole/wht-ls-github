package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type memFileInfo struct {
	name    string
	modTime time.Time
	dir     bool
}

func (m memFileInfo) Name() string { return m.name }
func (m memFileInfo) Size() int64  { return 0 }
func (m memFileInfo) Mode() os.FileMode {
	if m.dir {
		return os.ModeDir
	}
	return 0
}
func (m memFileInfo) ModTime() time.Time { return m.modTime }
func (m memFileInfo) IsDir() bool        { return m.dir }
func (m memFileInfo) Sys() any           { return nil }

type memDirEntry struct {
	name string
	dir  bool
}

func (m memDirEntry) Name() string { return m.name }
func (m memDirEntry) IsDir() bool  { return m.dir }
func (m memDirEntry) Type() os.FileMode {
	if m.dir {
		return os.ModeDir
	}
	return 0
}
func (m memDirEntry) Info() (os.FileInfo, error) {
	return memFileInfo{name: m.name, dir: m.dir}, nil
}

func TestListRepoDirsSortsByMtimeAndSkipsHidden(t *testing.T) {
	t.Parallel()

	newer := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	older := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	h := &host{
		readDir: func(path string) ([]os.DirEntry, error) {
			if path != "/repos" {
				t.Fatalf("root: %s", path)
			}
			return []os.DirEntry{
				memDirEntry{name: "old", dir: true},
				memDirEntry{name: "new", dir: true},
				memDirEntry{name: ".hidden", dir: true},
				memDirEntry{name: "file.txt", dir: false},
			}, nil
		},
		stat: func(path string) (os.FileInfo, error) {
			switch path {
			case "/repos/old":
				return memFileInfo{name: "old", dir: true, modTime: older}, nil
			case "/repos/new":
				return memFileInfo{name: "new", dir: true, modTime: newer}, nil
			default:
				return nil, os.ErrNotExist
			}
		},
	}

	dirs, err := listRepoDirs(h, "/repos")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 2 {
		t.Fatalf("len=%d", len(dirs))
	}
	if dirs[0].name != "new" || dirs[1].name != "old" {
		t.Fatalf("order: %#v", dirs)
	}
}

func TestInspectRepoSkipsWhenGhFails(t *testing.T) {
	t.Parallel()

	h := &host{
		capture: func(cwd, name string, args ...string) (string, error) {
			if name == "gh" {
				return "", fmt.Errorf("not a github repo")
			}
			if name == "git" && len(args) > 0 && args[0] == "status" {
				return "", nil
			}
			if name == "git" && len(args) > 0 && args[0] == "rev-list" {
				return "0\t0", nil
			}
			return "", nil
		},
		stat: func(path string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}
	if row := inspectRepo(h, "/repos/foo", "foo"); row != nil {
		t.Fatalf("expected skip, got %+v", row)
	}
}

func TestInspectRepoBuildsRowAfterFetch(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var sawRevList bool
	fetchReleased := make(chan struct{})
	h := &host{
		capture: func(cwd, name string, args ...string) (string, error) {
			switch {
			case name == "gh":
				return `{"isPrivate":true}`, nil
			case name == "git" && len(args) > 0 && args[0] == "fetch":
				<-fetchReleased
				return "", nil
			case name == "git" && len(args) > 0 && args[0] == "status":
				return " M README.md", nil
			case name == "git" && len(args) > 0 && args[0] == "rev-list":
				mu.Lock()
				sawRevList = true
				mu.Unlock()
				return "2\t1", nil
			default:
				return "", fmt.Errorf("unexpected %s %v", name, args)
			}
		},
		stat: func(path string) (os.FileInfo, error) {
			if strings.HasSuffix(path, ".devcontainer") {
				return memFileInfo{name: ".devcontainer", dir: true}, nil
			}
			return nil, os.ErrNotExist
		},
	}

	done := make(chan *repoRow, 1)
	go func() {
		done <- inspectRepo(h, "/repos/foo", "foo")
	}()

	select {
	case <-done:
		t.Fatal("rev-list ran before fetch finished")
	case <-time.After(30 * time.Millisecond):
	}
	close(fetchReleased)
	row := <-done
	if row == nil {
		t.Fatal("row is nil")
	}
	if row.Repo != "foo" || row.Visibility != "🔒 Priv" || row.Clean != "⚠️ Mod" || row.Dev != "📦" {
		t.Fatalf("row: %+v", row)
	}
	if row.Sync != "🔄 Diverged (↑2 ↓1)" {
		t.Fatalf("sync: %s", row.Sync)
	}
	mu.Lock()
	defer mu.Unlock()
	if !sawRevList {
		t.Fatal("rev-list was not called")
	}
}

func TestInspectAllPreservesOrderAndDropsSkipped(t *testing.T) {
	t.Parallel()

	h := &host{
		capture: func(cwd, name string, args ...string) (string, error) {
			base := filepath.Base(cwd)
			if name == "gh" {
				if base == "skip" {
					return "", fmt.Errorf("missing")
				}
				return `{"isPrivate":false}`, nil
			}
			if name == "git" && len(args) > 0 && args[0] == "status" {
				return "", nil
			}
			if name == "git" && len(args) > 0 && args[0] == "rev-list" {
				return "0\t0", nil
			}
			return "", nil
		},
		stat: func(path string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}

	dirs := []dirInfo{
		{name: "first"},
		{name: "skip"},
		{name: "third"},
	}
	rows := inspectAll(h, "/repos", dirs)
	if len(rows) != 2 {
		t.Fatalf("len=%d %#v", len(rows), rows)
	}
	if rows[0].Repo != "first" || rows[1].Repo != "third" {
		t.Fatalf("order: %#v", rows)
	}
}

func TestRunMissingTools(t *testing.T) {
	var stdout, stderr strings.Builder
	h := &host{
		lookPath: func(string) (string, error) { return "", os.ErrNotExist },
		stdout:   &stdout,
		stderr:   &stderr,
		args0:    "wht-ls-github",
	}
	oldArgs := os.Args
	os.Args = []string{"wht-ls-github"}
	t.Cleanup(func() { os.Args = oldArgs })

	code := run(h)
	if code != 1 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "git") {
		t.Fatalf("stderr: %s", stderr.String())
	}
}
