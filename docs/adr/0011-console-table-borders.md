# 0011. 表出力は TypeScript 版の console.table に合わせる

Date: 2026-09-08
Status: Accepted

## Context

Go 移植の初版は `text/tabwriter` で空白揃えだけしていた。TypeScript 版は bun の `console.table` で、番号列とセルを囲む罫線がある。罫線がないと行の対応が追いづらい、というフィードバックがあった。第三者の表ライブラリは使わない。

## Decision

- 表は標準ライブラリだけで箱線（`┌─┬┐` など）を描く。行のあいだにも横線を入れる。
- 先頭に 0 始まりの番号列を置く（ヘッダは空）。列は TypeScript 版と同じ `Repo` / `Visibility` / `Sync` / `Clean` / `Dev`。
- 絵文字の表示幅は端末向けに自前で数える。`golang.org/x/text` は使わない。
- 凡例の文言は変えない。

## Consequences

- `text/tabwriter` は使わなくなる。ADR-0002 の「表は tabwriter」は、この ADR で改める。
- 絵文字の幅は端末実装差で 1 桁ずれることがある。使うラベルが崩れたら幅判定を直す。
