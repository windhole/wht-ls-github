package main

import "testing"

func TestParseGhPrivate(t *testing.T) {
	t.Parallel()

	priv, ok := parseGhPrivate(`{"isPrivate":true}`)
	if !ok || !priv {
		t.Fatalf("private: got %v %v", priv, ok)
	}
	pub, ok := parseGhPrivate(`{"isPrivate":false}`)
	if !ok || pub {
		t.Fatalf("public: got %v %v", pub, ok)
	}
	if _, ok := parseGhPrivate(""); ok {
		t.Fatal("empty should fail")
	}
	if _, ok := parseGhPrivate("not-json"); ok {
		t.Fatal("invalid json should fail")
	}
}

func TestParseAheadBehind(t *testing.T) {
	t.Parallel()

	ahead, behind := parseAheadBehind("1\t2")
	if ahead != 1 || behind != 2 {
		t.Fatalf("got %d %d", ahead, behind)
	}
	ahead, behind = parseAheadBehind(" 3\t0\n")
	if ahead != 3 || behind != 0 {
		t.Fatalf("got %d %d", ahead, behind)
	}
	ahead, behind = parseAheadBehind("x\t1")
	if ahead != 0 || behind != 0 {
		t.Fatalf("invalid should be 0 0, got %d %d", ahead, behind)
	}
	ahead, behind = parseAheadBehind("1")
	if ahead != 0 || behind != 0 {
		t.Fatalf("short should be 0 0, got %d %d", ahead, behind)
	}
}

func TestSyncLabel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		ahead, behind int
		want          string
	}{
		{0, 0, "✅ Synced"},
		{2, 0, "📤 Push needed (↑2)"},
		{0, 3, "📥 Pull needed (↓3)"},
		{1, 4, "🔄 Diverged (↑1 ↓4)"},
	}
	for _, tc := range cases {
		got := syncLabel(tc.ahead, tc.behind)
		if got != tc.want {
			t.Fatalf("ahead=%d behind=%d: got %q want %q", tc.ahead, tc.behind, got, tc.want)
		}
	}
}

func TestVisibilityCleanDevLabels(t *testing.T) {
	t.Parallel()

	if visibilityLabel(true) != "🔒 Priv" {
		t.Fatal(visibilityLabel(true))
	}
	if visibilityLabel(false) != "🌐 Pub" {
		t.Fatal(visibilityLabel(false))
	}
	if cleanLabel("") != "✅" {
		t.Fatal(cleanLabel(""))
	}
	if cleanLabel(" M file.go") != "⚠️ Mod" {
		t.Fatal(cleanLabel(" M file.go"))
	}
	if devLabel(true) != "📦" {
		t.Fatal(devLabel(true))
	}
	if devLabel(false) != "-" {
		t.Fatal(devLabel(false))
	}
}
