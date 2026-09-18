# 0013. --wh は wh と windhole- の両方を対象にする

Date: 2026-09-18
Status: Accepted

## Context

ADR-0012 の `--wh` は名前が `wh` で始まるものだけだった。`wht-` / `whc-` / `whex-` はそれに入る。`windhole-` で始まるディレクトリも同じグループとして見たい。

`windhole-` は文字列としては `wh` で始まるので、実装上はすでに含まれていた。意図をプレフィックスの列として明示し、ドキュメントもそれに合わせる。

## Decision

- `--wh` の対象は、名前が `wh` で始まるもの、または `windhole-` で始まるもの。判定はプレフィックスの列で行う。
- `--proj` は従来どおり `proj-` のみ。
- README と `--help` に `windhole-` を書く。

## Consequences

- `wht-ls-github` も `windhole-note` も `--wh` で出る。
- あとから `wh` の意味を狭めても、`windhole-` は列に残せる。
