package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type host struct {
	lookPath func(string) (string, error)
	capture  func(cwd, name string, args ...string) (string, error)
	readDir  func(path string) ([]os.DirEntry, error)
	stat     func(path string) (os.FileInfo, error)
	userHome func() (string, error)
	stdout   io.Writer
	stderr   io.Writer
	args0    string
}

func defaultHost() *host {
	return &host{
		lookPath: exec.LookPath,
		capture:  captureCmd,
		readDir:  os.ReadDir,
		stat:     os.Stat,
		userHome: os.UserHomeDir,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		args0:    os.Args[0],
	}
}

func captureCmd(cwd, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return out, fmt.Errorf("%w: %s", err, msg)
		}
		return out, err
	}
	return out, nil
}

func (h *host) hasCommand(name string) bool {
	_, err := h.lookPath(name)
	return err == nil
}

func (h *host) invocation() string {
	return filepath.Base(h.args0)
}

func (h *host) printf(format string, args ...any) {
	fmt.Fprintf(h.stdout, format, args...)
}

func (h *host) errorf(format string, args ...any) {
	fmt.Fprintf(h.stderr, format, args...)
}
