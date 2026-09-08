package main

import (
	"io"
	"strconv"
	"strings"
	"unicode"
)

func writeTable(w io.Writer, rows []repoRow) {
	headers := []string{"", "Repo", "Visibility", "Sync", "Clean", "Dev"}
	cells := make([][]string, len(rows))
	for i, row := range rows {
		cells[i] = []string{
			strconv.Itoa(i),
			row.Repo,
			row.Visibility,
			row.Sync,
			row.Clean,
			row.Dev,
		}
	}
	writeBoxTable(w, headers, cells)
}

func writeLegend(w io.Writer) {
	_, _ = io.WriteString(w, "\n--- Legend ---\n")
	_, _ = io.WriteString(w, "📥 Pull needed : Remote has changes you don't have.\n")
	_, _ = io.WriteString(w, "📤 Push needed : You have local commits not on remote.\n")
	_, _ = io.WriteString(w, "⚠️ Mod         : You have uncommitted files (modified/untracked).\n")
	_, _ = io.WriteString(w, "📦             : .devcontainer environment detected.\n")
}

func writeBoxTable(w io.Writer, headers []string, rows [][]string) {
	cols := len(headers)
	widths := make([]int, cols)
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}
	for _, row := range rows {
		for i := 0; i < cols && i < len(row); i++ {
			if w := displayWidth(row[i]); w > widths[i] {
				widths[i] = w
			}
		}
	}

	_, _ = io.WriteString(w, ruleLine(widths, "┌", "┬", "┐")+"\n")
	_, _ = io.WriteString(w, dataLine(headers, widths, false)+"\n")
	if len(rows) == 0 {
		_, _ = io.WriteString(w, ruleLine(widths, "└", "┴", "┘")+"\n")
		return
	}
	_, _ = io.WriteString(w, ruleLine(widths, "├", "┼", "┤")+"\n")
	for i, row := range rows {
		_, _ = io.WriteString(w, dataLine(row, widths, true)+"\n")
		if i+1 < len(rows) {
			_, _ = io.WriteString(w, ruleLine(widths, "├", "┼", "┤")+"\n")
		}
	}
	_, _ = io.WriteString(w, ruleLine(widths, "└", "┴", "┘")+"\n")
}

func ruleLine(widths []int, left, mid, right string) string {
	var b strings.Builder
	b.WriteString(left)
	for i, w := range widths {
		if i > 0 {
			b.WriteString(mid)
		}
		b.WriteString(strings.Repeat("─", w+2))
	}
	b.WriteString(right)
	return b.String()
}

func dataLine(cells []string, widths []int, indexRight bool) string {
	var b strings.Builder
	b.WriteString("│")
	for i, w := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		b.WriteByte(' ')
		if indexRight && i == 0 {
			b.WriteString(padLeft(cell, w))
		} else {
			b.WriteString(padRight(cell, w))
		}
		b.WriteByte(' ')
		b.WriteString("│")
	}
	return b.String()
}

func padRight(s string, width int) string {
	return s + strings.Repeat(" ", padSpaces(s, width))
}

func padLeft(s string, width int) string {
	return strings.Repeat(" ", padSpaces(s, width)) + s
}

func padSpaces(s string, width int) int {
	n := width - displayWidth(s)
	if n < 0 {
		return 0
	}
	return n
}

func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

func runeWidth(r rune) int {
	switch {
	case r == 0 || r == '\n' || r == '\r' || r == '\t':
		return 0
	case r == 0x200D || r == 0xFE0E || r == 0xFE0F:
		return 0
	case unicode.Is(unicode.Mn, r), unicode.Is(unicode.Me, r), unicode.Is(unicode.Cf, r):
		return 0
	case r < 0x7F:
		if r >= 0x20 {
			return 1
		}
		return 0
	case r >= 0x1100 && isWideRune(r):
		return 2
	default:
		return 1
	}
}

func isWideRune(r rune) bool {
	switch {
	case r <= 0x115F:
		return true
	case r == 0x2329 || r == 0x232A:
		return true
	case r >= 0x2E80 && r <= 0xA4CF && r != 0x303F:
		return true
	case r >= 0xAC00 && r <= 0xD7A3:
		return true
	case r >= 0xF900 && r <= 0xFAFF:
		return true
	case r >= 0xFE10 && r <= 0xFE19:
		return true
	case r >= 0xFE30 && r <= 0xFE6F:
		return true
	case r >= 0xFF00 && r <= 0xFF60:
		return true
	case r >= 0xFFE0 && r <= 0xFFE6:
		return true
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return true
	case r >= 0x1F300 && r <= 0x1FAFF:
		return true
	case r >= 0x1F900 && r <= 0x1F9FF:
		return true
	case r >= 0x2190 && r <= 0x21FF:
		return false
	case r >= 0x2600 && r <= 0x27BF:
		return true
	case r >= 0x2B00 && r <= 0x2BFF:
		return true
	default:
		return false
	}
}
