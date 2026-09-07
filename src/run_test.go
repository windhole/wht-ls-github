package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunHelpAndVersion(t *testing.T) {
	var stdout, stderr strings.Builder
	h := &host{
		stdout: &stdout,
		stderr: &stderr,
		args0:  "lsg",
	}
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"lsg", "--help"}
	if code := run(h); code != 0 {
		t.Fatalf("help code=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "使い方") {
		t.Fatalf("help: %s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	old := version
	version = "v9.9.9"
	t.Cleanup(func() { version = old })
	os.Args = []string{"lsg", "--version"}
	if code := run(h); code != 0 {
		t.Fatalf("version code=%d", code)
	}
	if stdout.String() != "lsg v9.9.9\n" {
		t.Fatalf("version: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	os.Args = []string{"lsg", "--nope"}
	if code := run(h); code != 2 {
		t.Fatalf("unknown code=%d", code)
	}
}

func TestRunListsRepos(t *testing.T) {
	var stdout, stderr strings.Builder
	newer := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	older := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	h := &host{
		lookPath: func(string) (string, error) { return "/bin/tool", nil },
		userHome: func() (string, error) { return "/Users/demo", nil },
		readDir: func(path string) ([]os.DirEntry, error) {
			if path != "/Users/demo/Documents/GitHub" {
				t.Fatalf("root %s", path)
			}
			return []os.DirEntry{
				memDirEntry{name: "older", dir: true},
				memDirEntry{name: "newer", dir: true},
			}, nil
		},
		stat: func(path string) (os.FileInfo, error) {
			switch path {
			case "/Users/demo/Documents/GitHub":
				return memFileInfo{name: "GitHub", dir: true}, nil
			case "/Users/demo/Documents/GitHub/older":
				return memFileInfo{name: "older", dir: true, modTime: older}, nil
			case "/Users/demo/Documents/GitHub/newer":
				return memFileInfo{name: "newer", dir: true, modTime: newer}, nil
			default:
				return nil, os.ErrNotExist
			}
		},
		capture: func(cwd, name string, args ...string) (string, error) {
			if name == "gh" {
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
		stdout: &stdout,
		stderr: &stderr,
		args0:  "lsg",
	}
	oldArgs := os.Args
	os.Args = []string{"lsg"}
	t.Cleanup(func() { os.Args = oldArgs })

	if code := run(h); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	out := stdout.String()
	newerIdx := strings.Index(out, "newer")
	olderIdx := strings.Index(out, "older")
	if newerIdx < 0 || olderIdx < 0 || newerIdx > olderIdx {
		t.Fatalf("order: %s", out)
	}
	if !strings.Contains(out, "Checking GitHub repositories") {
		t.Fatalf("missing progress: %s", out)
	}
}

func TestRunMissingScanDir(t *testing.T) {
	var stdout, stderr strings.Builder
	h := &host{
		lookPath: func(string) (string, error) { return "/bin/tool", nil },
		stat:     func(path string) (os.FileInfo, error) { return nil, os.ErrNotExist },
		stdout:   &stdout,
		stderr:   &stderr,
		args0:    "lsg",
	}
	oldArgs := os.Args
	os.Args = []string{"lsg", "--dir", "/missing"}
	t.Cleanup(func() { os.Args = oldArgs })

	if code := run(h); code != 1 {
		t.Fatalf("code=%d", code)
	}
	if !strings.Contains(stderr.String(), "/missing") {
		t.Fatalf("stderr: %s", stderr.String())
	}
}
