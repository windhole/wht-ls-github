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
}

func parseGhPrivate(raw string) (bool, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, false
	}
	var payload struct {
		IsPrivate bool `json:"isPrivate"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return false, false
	}
	return payload.IsPrivate, true
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

func visibilityLabel(isPrivate bool) string {
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
