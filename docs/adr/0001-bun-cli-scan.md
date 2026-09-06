# 0001. bun 単一ファイル CLI でローカル GitHub リポジトリを一覧する

Date: 2026-09-06
Status: Superseded by ADR-0002

## Context

macOS の `~/Documents/GitHub` 配下に置いた複数リポジトリの公開範囲・リモート同期・未コミット変更・Dev Container 有無を、毎回個別に `git` / `gh` で確認するのは手間が大きい。手元の作業ディレクトリを一括で眺めるコマンドが欲しい。

Git プロトコルや GitHub API を自前で扱う必要はなく、すでに PATH にある `git` と GitHub CLI（`gh`）で足りる。

## Decision

- ランタイムは bun。依存パッケージは置かず、TypeScript 単一ファイルとする。
- 走査先は `~/Documents/GitHub` 直下の非隠しディレクトリ。ディレクトリ自身の mtime が新しい順。
- 各ディレクトリで `gh repo view --json isPrivate`、`git fetch`、`git status -s`、`git rev-list --left-right --count HEAD...@{u}`、`.devcontainer` の存在確認を行う。
- `gh repo view` が失敗したディレクトリは一覧から除外する。
- 結果は表と凡例として標準出力へ出す。

## Consequences

- bun と git / gh があればスクリプトを置いて実行できる。
- 一覧のたびに各リポジトリで `git fetch` するため、リモート追従の副作用がある。
- bun ランタイムが実行側に必要。単一バイナリでは配れない。
