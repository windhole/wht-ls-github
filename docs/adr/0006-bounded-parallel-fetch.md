# 0006. リポジトリ走査は同時実行数を制限する

Date: 2026-09-06
Status: Accepted

## Context

現行の bun 実装は、走査対象の全ディレクトリに対して `Promise.all` で同時に `git fetch` と `gh repo view` を投げる。件数が増えるとリモートへの接続が一気に増え、待ちや失敗の原因になる。Go でも並列は欲しいが、上限は付けたい。

同期判定は `git fetch` のあとに `git rev-list HEAD...@{u}` を見ないと古いままになる。現行は同一リポジトリ内でも両方を同時起動しており、fetch 完了前に差を読む競合がある。

並び順はディレクトリ mtime の新しい順であり、完了順に並べ替えてはいけない。

## Decision

- リポジトリ間はワーカープールにする。同時実行数は定数 8。フラグにはしない。
- 結果は入力（mtime 順）のインデックスに書き戻し、`gh` 失敗などで欠けた行だけ落とす。
- 同一リポジトリ内では `git fetch` の完了を待ってから `rev-list` する。`gh repo view`、`git status -s`、`.devcontainer` の確認は fetch と並列でよい。

## Consequences

- 件数が多いときの fetch 嵐を抑えられる。
- Pull / Push 判定が fetch 後の状態を見る。
- 表の行順は従来どおり更新日時順。
