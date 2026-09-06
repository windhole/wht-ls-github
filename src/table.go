package main

import (
	"io"
	"text/tabwriter"
)

func writeTable(w io.Writer, rows []repoRow) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = io.WriteString(tw, "Repo\tVisibility\tSync\tClean\tDev\n")
	for _, row := range rows {
		_, _ = io.WriteString(tw, row.Repo+"\t"+row.Visibility+"\t"+row.Sync+"\t"+row.Clean+"\t"+row.Dev+"\n")
	}
	_ = tw.Flush()
}

func writeLegend(w io.Writer) {
	_, _ = io.WriteString(w, "\n--- Legend ---\n")
	_, _ = io.WriteString(w, "📥 Pull needed : Remote has changes you don't have.\n")
	_, _ = io.WriteString(w, "📤 Push needed : You have local commits not on remote.\n")
	_, _ = io.WriteString(w, "⚠️ Mod         : You have uncommitted files (modified/untracked).\n")
	_, _ = io.WriteString(w, "📦             : .devcontainer environment detected.\n")
}
