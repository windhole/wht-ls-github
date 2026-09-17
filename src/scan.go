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

type inspectOpts struct {
	showArchived bool
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

func filterDirs(dirs []dirInfo, proj, wh bool) []dirInfo {
	if !proj && !wh {
		return dirs
	}
	out := make([]dirInfo, 0, len(dirs))
	for _, d := range dirs {
		if matchDirName(d.name, proj, wh) {
			out = append(out, d)
		}
	}
	return out
}

func inspectAll(h *host, root string, dirs []dirInfo, opts inspectOpts) []repoRow {
	rows := make([]*repoRow, len(dirs))
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	for i, d := range dirs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, d dirInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			rows[i] = inspectRepo(h, filepath.Join(root, d.name), d.name, opts)
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

func inspectRepo(h *host, fullPath, dirName string, opts inspectOpts) *repoRow {
	revOut, revErr := h.capture(fullPath, "git", "rev-parse", "--is-inside-work-tree")
	devExists := false
	if _, err := h.stat(filepath.Join(fullPath, ".devcontainer")); err == nil {
		devExists = true
	}

	if !isGitWorkTreeOutput(revOut, revErr) {
		return &repoRow{
			Repo:       dirName,
			Visibility: labelNA,
			Sync:       labelNoGit,
			Clean:      labelNA,
			Dev:        devLabel(devExists),
		}
	}

	var (
		ghRaw      string
		ghErr      error
		statusText string
		statusErr  error
		fetchDone  = make(chan struct{})
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		ghRaw, ghErr = h.capture(fullPath, "gh", "repo", "view", "--json", "isPrivate,isArchived")
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

	<-fetchDone
	wg.Wait()

	if ghErr != nil || ghRaw == "" {
		clean := labelNA
		if statusErr == nil {
			clean = cleanLabel(statusText)
		}
		return &repoRow{
			Repo:       dirName,
			Visibility: labelNA,
			Sync:       labelNoGitHub,
			Clean:      clean,
			Dev:        devLabel(devExists),
		}
	}
	meta, ok := parseGhMeta(ghRaw)
	if !ok {
		return &repoRow{
			Repo:       dirName,
			Visibility: labelNA,
			Sync:       labelNoGitHub,
			Clean:      cleanLabel(statusText),
			Dev:        devLabel(devExists),
		}
	}
	if meta.IsArchived && !opts.showArchived {
		return nil
	}

	syncRaw := "0\t0"
	if statusErr == nil {
		if raw, err := h.capture(fullPath, "git", "rev-list", "--left-right", "--count", "HEAD...@{u}"); err == nil {
			syncRaw = raw
		}
	}
	if statusErr != nil {
		return &repoRow{
			Repo:       dirName,
			Visibility: visibilityLabel(meta.IsPrivate, meta.IsArchived),
			Sync:       syncLabel(parseAheadBehind(syncRaw)),
			Clean:      labelNA,
			Dev:        devLabel(devExists),
			archived:   meta.IsArchived,
		}
	}

	ahead, behind := parseAheadBehind(syncRaw)
	return &repoRow{
		Repo:       dirName,
		Visibility: visibilityLabel(meta.IsPrivate, meta.IsArchived),
		Sync:       syncLabel(ahead, behind),
		Clean:      cleanLabel(statusText),
		Dev:        devLabel(devExists),
		archived:   meta.IsArchived,
	}
}
