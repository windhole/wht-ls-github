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
	if !strings.Contains(out, "--- Legend ---") {
		t.Fatalf("missing legend: %s", out)
	}
	if !strings.Contains(out, "📥 Pull needed") {
		t.Fatalf("missing pull legend: %s", out)
	}
}
