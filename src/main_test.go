package main

import (
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	args, err := parseArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	if args.help || args.version || args.dir != "" {
		t.Fatalf("empty: %+v", args)
	}

	args, err = parseArgs([]string{"--dir", "/tmp/repos"})
	if err != nil {
		t.Fatal(err)
	}
	if args.dir != "/tmp/repos" {
		t.Fatalf("dir: %+v", args)
	}

	args, err = parseArgs([]string{"--version"})
	if err != nil || !args.version {
		t.Fatalf("version: %+v %v", args, err)
	}

	args, err = parseArgs([]string{"--help", "--version"})
	if err != nil || !args.help || args.version {
		t.Fatalf("help wins: %+v %v", args, err)
	}

	if _, err := parseArgs([]string{"--unknown"}); err == nil {
		t.Fatal("unknown flag should fail")
	}
	if _, err := parseArgs([]string{"extra"}); err == nil {
		t.Fatal("positional should fail")
	}
}

func TestResolveScanDir(t *testing.T) {
	t.Parallel()

	h := &host{
		userHome: func() (string, error) { return "/Users/demo", nil },
	}
	got, err := resolveScanDir(h, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "/Users/demo/Documents/GitHub" {
		t.Fatalf("default: %s", got)
	}
	got, err = resolveScanDir(h, "/custom")
	if err != nil || got != "/custom" {
		t.Fatalf("specified: %s %v", got, err)
	}

	h.userHome = func() (string, error) { return "", nil }
	if _, err := resolveScanDir(h, ""); err == nil {
		t.Fatal("empty home should fail")
	}
}

func TestPrintVersion(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	old := version
	version = "v0.1.0"
	t.Cleanup(func() { version = old })
	printVersion(&b)
	if b.String() != "lsg v0.1.0\n" {
		t.Fatalf("got %q", b.String())
	}
}
