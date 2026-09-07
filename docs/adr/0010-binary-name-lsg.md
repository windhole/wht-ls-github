# 0010. 配布バイナリの名前は lsg とする

Date: 2026-09-07
Status: Accepted

## Context

リポジトリ名は `wht-ls-github` のままが分かりやすいが、PATH に置く実行ファイル名としては長い。姉妹リポジトリの `grg` と同じく、短いコマンド名にしたい。

## Decision

- `make` / `make release` / GitHub Actions が出力するファイルは `dist/lsg` とする。GitHub Release の Asset 名も `lsg`。
- コマンドとして呼び出す名前、`--version` の表示は `lsg` とする。
- リポジトリ名・Go モジュール名・ドキュメントのリポジトリ URL は `wht-ls-github` のまま。
- ADR-0003 の「実行ファイル名はリポジトリ名と同じ」は、この ADR で改める。

## Consequences

- 入れ方は `install dist/lsg "$HOME/bin/lsg"`。以降は `lsg` で実行する。
- 古い名前のバイナリが手元に残っていても、`.gitignore` は `/lsg` と `/wht-ls-github` の両方を無視する。
