package main

import (
	"strings"
	"testing"
)

func TestWriteTableAndLegend(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	writeTable(&b, []repoRow{
		{Repo: "alpha", Visibility: "🌐 Pub", Sync: "✅ Synced", Clean: "✅", Dev: "-"},
		{Repo: "beta", Visibility: "🔒 Priv", Sync: "📤 Push needed (↑1)", Clean: "⚠️ Mod", Dev: "📦"},
	})
	writeLegend(&b)
	out := b.String()

	if !strings.Contains(out, "Repo") || !strings.Contains(out, "Visibility") {
		t.Fatalf("missing header: %s", out)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Fatalf("missing rows: %s", out)
	}
	if !strings.Contains(out, "┌") || !strings.Contains(out, "┼") || !strings.Contains(out, "└") {
		t.Fatalf("missing box drawing: %s", out)
	}
	if !strings.Contains(out, "│  0 │") && !strings.Contains(out, "│ 0 │") {
		t.Fatalf("missing index column: %s", out)
	}
	if !strings.Contains(out, "--- Legend ---") {
		t.Fatalf("missing legend: %s", out)
	}
	if !strings.Contains(out, "📥 Pull needed") {
		t.Fatalf("missing pull legend: %s", out)
	}

	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	var tableLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "┌") || strings.HasPrefix(line, "├") || strings.HasPrefix(line, "└") || strings.HasPrefix(line, "│") {
			tableLines = append(tableLines, line)
		}
	}
	if len(tableLines) < 5 {
		t.Fatalf("too few table lines: %#v", tableLines)
	}
	width := displayWidth(tableLines[0])
	for _, line := range tableLines {
		if got := displayWidth(line); got != width {
			t.Fatalf("uneven line width %d vs %d\n%s\n%s", width, got, tableLines[0], line)
		}
	}
}

func TestWriteTableEmpty(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	writeTable(&b, nil)
	out := b.String()
	if !strings.Contains(out, "Repo") {
		t.Fatalf("missing header: %s", out)
	}
	if !strings.Contains(out, "┌") || !strings.Contains(out, "└") {
		t.Fatalf("missing box: %s", out)
	}
	if strings.Count(out, "├") != 0 {
		t.Fatalf("empty table should not have mid rules: %s", out)
	}
}

func TestDisplayWidth(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"Repo", 4},
		{"-", 1},
		{"✅", 2},
		{"📦", 2},
		{"🔒 Priv", 7},
		{"🌐 Pub", 6},
		{"⚠️ Mod", 6},
		{"✅ Synced", 9},
	}
	for _, tc := range cases {
		if got := displayWidth(tc.in); got != tc.want {
			t.Fatalf("%q: got %d want %d", tc.in, got, tc.want)
		}
	}
}
