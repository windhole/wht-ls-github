package main

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxWorkers = 8

type dirInfo struct {
	name  string
	mtime time.Time
}

func defaultScanDir(home string) string {
	return filepath.Join(home, "Documents", "GitHub")
}

func listRepoDirs(h *host, root string) ([]dirInfo, error) {
	entries, err := h.readDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]dirInfo, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		full := filepath.Join(root, e.Name())
		info, err := h.stat(full)
		if err != nil {
			continue
		}
		out = append(out, dirInfo{name: e.Name(), mtime: info.ModTime()})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].mtime.After(out[j].mtime)
	})
	return out, nil
}

func inspectAll(h *host, root string, dirs []dirInfo) []repoRow {
	rows := make([]*repoRow, len(dirs))
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	for i, d := range dirs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, d dirInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			rows[i] = inspectRepo(h, filepath.Join(root, d.name), d.name)
		}(i, d)
	}
	wg.Wait()

	out := make([]repoRow, 0, len(dirs))
	for _, r := range rows {
		if r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func inspectRepo(h *host, fullPath, dirName string) *repoRow {
	var (
		ghRaw      string
		ghErr      error
		statusText string
		statusErr  error
		devExists  bool
		fetchDone  = make(chan struct{})
	)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		ghRaw, ghErr = h.capture(fullPath, "gh", "repo", "view", "--json", "isPrivate")
	}()
	go func() {
		defer wg.Done()
		_, _ = h.capture(fullPath, "git", "fetch")
		close(fetchDone)
	}()
	go func() {
		defer wg.Done()
		statusText, statusErr = h.capture(fullPath, "git", "status", "-s")
	}()
	go func() {
		defer wg.Done()
		_, err := h.stat(filepath.Join(fullPath, ".devcontainer"))
		devExists = err == nil
	}()

	<-fetchDone
	syncRaw, err := h.capture(fullPath, "git", "rev-list", "--left-right", "--count", "HEAD...@{u}")
	if err != nil {
		syncRaw = "0\t0"
	}
	wg.Wait()

	if ghErr != nil || ghRaw == "" {
		return nil
	}
	isPrivate, ok := parseGhPrivate(ghRaw)
	if !ok {
		return nil
	}
	if statusErr != nil {
		return nil
	}

	ahead, behind := parseAheadBehind(syncRaw)
	return &repoRow{
		Repo:       dirName,
		Visibility: visibilityLabel(isPrivate),
		Sync:       syncLabel(ahead, behind),
		Clean:      cleanLabel(statusText),
		Dev:        devLabel(devExists),
	}
}
