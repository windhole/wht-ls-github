# 0012. 未初期化・未連携を表示し、archive / 名前で絞る

Date: 2026-09-17
Status: Accepted

## Context

これまで `gh repo view` が失敗したディレクトリは行から外していた。`git init` していない場所や、git はあるが GitHub と繋がっていない場所も、一覧で区別して見たい。

GitHub 上で archive したリポジトリは普段の作業対象ではないので、既定では出したくない。名前のプレフィックスで自分のプロジェクト群だけ見たいこともある。

## Decision

- git リポジトリでないディレクトリは除外せず、Sync を `🚫 no git`、Visibility / Clean を `-` として出す。`.devcontainer` の有無は従来どおり。
- git リポジトリだが `gh repo view` が失敗したものは、Sync を `☁️ no GitHub`、Visibility を `-` として出す。`git status` と Dev は取る。fetch / 同期差は見ない。
- GitHub 上で archive されたものは既定では出さない。`--all` のときだけ出し、Visibility は `🗄️ Priv` / `🗄️ Pub`。
- `--proj` は名前が `proj-` で始まるものだけ。`--wh` は名前が `wh` で始まるものだけ。両方付けたときは和集合。`--all` と組み合わせてよい。
- archive 判定は `gh repo view --json isPrivate,isArchived`。git かどうかは `git rev-parse --is-inside-work-tree`。

## Consequences

- 走査先直下の非隠しディレクトリは、GitHub 連携の有無にかかわらず行になり得る。
- archive を見たいときは `--all` が必要。
- `--proj` と `--wh` は同時指定できる（和集合）。排他にはしない。
