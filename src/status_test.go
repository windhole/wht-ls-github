package main

import "testing"

func TestParseGhMeta(t *testing.T) {
	t.Parallel()

	priv, ok := parseGhMeta(`{"isPrivate":true}`)
	if !ok || !priv.IsPrivate || priv.IsArchived {
		t.Fatalf("private: %+v %v", priv, ok)
	}
	pub, ok := parseGhMeta(`{"isPrivate":false,"isArchived":true}`)
	if !ok || pub.IsPrivate || !pub.IsArchived {
		t.Fatalf("archived public: %+v %v", pub, ok)
	}
	if _, ok := parseGhMeta(""); ok {
		t.Fatal("empty should fail")
	}
	if _, ok := parseGhMeta("not-json"); ok {
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

	if visibilityLabel(true, false) != "🔒 Priv" {
		t.Fatal(visibilityLabel(true, false))
	}
	if visibilityLabel(false, false) != "🌐 Pub" {
		t.Fatal(visibilityLabel(false, false))
	}
	if visibilityLabel(true, true) != "🗄️ Priv" {
		t.Fatal(visibilityLabel(true, true))
	}
	if visibilityLabel(false, true) != "🗄️ Pub" {
		t.Fatal(visibilityLabel(false, true))
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

func TestMatchDirName(t *testing.T) {
	t.Parallel()

	if !matchDirName("anything", false, false) {
		t.Fatal("no filter should match all")
	}
	if !matchDirName("proj-foo", true, false) || matchDirName("wht-ls", true, false) {
		t.Fatal("--proj")
	}
	if !matchDirName("wht-ls", false, true) || !matchDirName("whc-ai", false, true) {
		t.Fatal("--wh should match wh*")
	}
	if !matchDirName("windhole-note", false, true) {
		t.Fatal("--wh should match windhole-")
	}
	if matchDirName("proj-foo", false, true) {
		t.Fatal("proj- should not match --wh")
	}
	if !matchDirName("proj-foo", true, true) || !matchDirName("wht-ls", true, true) {
		t.Fatal("both filters are a union")
	}
	if matchDirName("other", true, true) {
		t.Fatal("union should still exclude others")
	}
}

func TestIsGitWorkTreeOutput(t *testing.T) {
	t.Parallel()

	if !isGitWorkTreeOutput("true", nil) {
		t.Fatal("true")
	}
	if !isGitWorkTreeOutput("", nil) {
		t.Fatal("empty success is treated as git (tests)")
	}
	if isGitWorkTreeOutput("true", errSentinel{}) {
		t.Fatal("error")
	}
}

type errSentinel struct{}

func (errSentinel) Error() string { return "no" }
