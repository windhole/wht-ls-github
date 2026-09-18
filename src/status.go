package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type repoRow struct {
	Repo       string
	Visibility string
	Sync       string
	Clean      string
	Dev        string
	archived   bool
}

type ghMeta struct {
	IsPrivate  bool `json:"isPrivate"`
	IsArchived bool `json:"isArchived"`
}

const (
	labelNoGit    = "🚫 no git"
	labelNoGitHub = "☁️ no GitHub"
	labelNA       = "-"
	labelArchPriv = "🗄️ Priv"
	labelArchPub  = "🗄️ Pub"
	prefixProj    = "proj-"
)

var prefixWH = []string{"wh", "windhole-"}

func parseGhMeta(raw string) (ghMeta, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ghMeta{}, false
	}
	var payload ghMeta
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return ghMeta{}, false
	}
	return payload, true
}

func parseAheadBehind(raw string) (int, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0
	}
	parts := strings.Split(raw, "\t")
	if len(parts) < 2 {
		return 0, 0
	}
	ahead, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	behind, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return 0, 0
	}
	return ahead, behind
}

func syncLabel(ahead, behind int) string {
	switch {
	case behind > 0 && ahead > 0:
		return "🔄 Diverged (↑" + strconv.Itoa(ahead) + " ↓" + strconv.Itoa(behind) + ")"
	case behind > 0:
		return "📥 Pull needed (↓" + strconv.Itoa(behind) + ")"
	case ahead > 0:
		return "📤 Push needed (↑" + strconv.Itoa(ahead) + ")"
	default:
		return "✅ Synced"
	}
}

func visibilityLabel(isPrivate, archived bool) string {
	if archived {
		if isPrivate {
			return labelArchPriv
		}
		return labelArchPub
	}
	if isPrivate {
		return "🔒 Priv"
	}
	return "🌐 Pub"
}

func cleanLabel(statusText string) string {
	if strings.TrimSpace(statusText) == "" {
		return "✅"
	}
	return "⚠️ Mod"
}

func devLabel(exists bool) string {
	if exists {
		return "📦"
	}
	return "-"
}

func matchDirName(name string, proj, wh bool) bool {
	if !proj && !wh {
		return true
	}
	ok := false
	if proj && strings.HasPrefix(name, prefixProj) {
		ok = true
	}
	if wh && hasAnyPrefix(name, prefixWH...) {
		ok = true
	}
	return ok
}

func hasAnyPrefix(name string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func isGitWorkTreeOutput(out string, err error) bool {
	if err != nil {
		return false
	}
	s := strings.TrimSpace(out)
	return s == "true" || s == ""
}
