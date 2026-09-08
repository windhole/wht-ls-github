# 0002. Go 標準ライブラリのみで CLI を実装する

Date: 2026-09-06
Status: Accepted (table rendering amended by ADR-0011)

## Context

bun の TypeScript 単一ファイルでも第三者パッケージは使っていなかったが、実行には bun ランタイムが要る。方針を「サードパーティのライブラリに依存しない」「言語標準と official のみ」「macOS では単一バイナリで配りたい」に寄せる。機能は ADR-0001 と同じ（走査、git / gh の呼び出し、表出力）で、実行基盤だけを置き換える。姉妹リポジトリ `wht-get-ready-for-github`（`grg`）と同じ制約にする。

## Decision

- 実装言語は Go。依存モジュールは置かず、標準ライブラリだけを使う（`os/exec`、`encoding/json`、`flag`、`os`、`path/filepath`、`sync`）。表の罫線は ADR-0011。
- 配布物は `go build` した単一バイナリ。実行時に必要な外部コマンドは従来どおり PATH 上の `git` と `gh` のみ。
- bun / TypeScript のソースは削除する。
- Cobra などの CLI 枠、golang.org/x 配下の拡張パッケージ、go-git も使わない。
- 既定の走査先は ADR-0001 と同じ `~/Documents/GitHub`。テストと再利用のため `--dir` を足す。`--help` / `--version` も足す。
- 同期ラベルと凡例の文言・絵文字は現行互換を維持する。

## Consequences

- 実行側に bun は不要。Mac にはバイナリを 1 つ置けば足りる。
- ビルドには Go ツールチェーンが要る。対象は darwin を想定するが、ソースは Unix 系なら動く。
- git / GitHub API 自体は自前実装せず、既存の `git` / `gh` に任せる。
- ADR-0001 の走査ルールと除外条件は維持する。
