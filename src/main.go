package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type cliArgs struct {
	help    bool
	version bool
	dir     string
	all     bool
	proj    bool
	wh      bool
}

const commandName = "lsg"

func parseArgs(argv []string) (cliArgs, error) {
	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "")
	fs.BoolVar(help, "h", false, "")
	showVersion := fs.Bool("version", false, "")
	dir := fs.String("dir", "", "")
	all := fs.Bool("all", false, "")
	proj := fs.Bool("proj", false, "")
	wh := fs.Bool("wh", false, "")

	if err := fs.Parse(argv); err != nil {
		return cliArgs{}, fmt.Errorf("不明な引数があります")
	}
	if fs.NArg() > 0 {
		return cliArgs{}, fmt.Errorf("不明な引数: %s", fs.Arg(0))
	}

	out := cliArgs{
		help:    *help,
		version: *showVersion,
		dir:     strings.TrimSpace(*dir),
		all:     *all,
		proj:    *proj,
		wh:      *wh,
	}
	if out.help {
		out.version = false
	}
	return out, nil
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", commandName, versionString())
}

func printUsage(w io.Writer, cmd string) {
	fmt.Fprintf(w, "%s — ローカルの GitHub リポジトリ一覧を表示する\n", commandName)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "使い方")
	fmt.Fprintf(w, "  %s\n", cmd)
	fmt.Fprintln(w, "    既定の ~/Documents/GitHub を走査する")
	fmt.Fprintf(w, "  %s --dir PATH\n", cmd)
	fmt.Fprintln(w, "    指定したディレクトリ直下を走査する")
	fmt.Fprintf(w, "  %s --all\n", cmd)
	fmt.Fprintln(w, "    GitHub で archive されたリポジトリも表示する")
	fmt.Fprintf(w, "  %s --proj\n", cmd)
	fmt.Fprintln(w, "    名前が proj- で始まるディレクトリだけ表示する")
	fmt.Fprintf(w, "  %s --wh\n", cmd)
	fmt.Fprintln(w, "    名前が wh または windhole- で始まるディレクトリだけ表示する")
	fmt.Fprintf(w, "  %s --version\n", cmd)
	fmt.Fprintln(w, "    バージョンを表示する")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "前提")
	fmt.Fprintln(w, "  ・git / gh が PATH にある")
	fmt.Fprintln(w, "  ・gh は対象ホストにログイン済み")
	fmt.Fprintln(w, "  ・各リポジトリで git fetch する（git / GitHub 連携があるもの）")
	fmt.Fprintln(w, "  ・--proj と --wh を両方付けると、どちらかに合うものを出す")
}

func resolveScanDir(h *host, specified string) (string, error) {
	if specified != "" {
		return specified, nil
	}
	home, err := h.userHome()
	if err != nil || home == "" {
		return "", fmt.Errorf("ホームディレクトリが分かりません")
	}
	return defaultScanDir(home), nil
}

func run(h *host) int {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		h.errorf("%s\n", err.Error())
		h.errorf("使い方を見る: %s --help\n", h.invocation())
		return 2
	}

	if args.help {
		printUsage(h.stdout, h.invocation())
		return 0
	}
	if args.version {
		printVersion(h.stdout)
		return 0
	}

	if !h.hasCommand("git") {
		h.errorf("git が PATH にありません\n")
		return 1
	}
	if !h.hasCommand("gh") {
		h.errorf("gh が PATH にありません\n")
		return 1
	}

	root, err := resolveScanDir(h, args.dir)
	if err != nil {
		h.errorf("%s\n", err.Error())
		return 1
	}
	if _, err := h.stat(root); err != nil {
		h.errorf("走査先がありません: %s\n", root)
		return 1
	}

	h.printf("Checking GitHub repositories (fetching remotes)...\n")

	dirs, err := listRepoDirs(h, root)
	if err != nil {
		h.errorf("ディレクトリを読めません: %s\n", err.Error())
		return 1
	}
	dirs = filterDirs(dirs, args.proj, args.wh)

	rows := inspectAll(h, root, dirs, inspectOpts{showArchived: args.all})
	writeTable(h.stdout, rows)
	writeLegend(h.stdout)
	return 0
}

func main() {
	os.Exit(run(defaultHost()))
}
