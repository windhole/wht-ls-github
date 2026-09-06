# 0008. 利用者向けは README.md、開発作業は development.md

Date: 2026-09-06
Status: Accepted

## Context

コマンドの概要・使い方と、ビルド / リリース手順が同じファイルに混ざると、使う人と作る人のどちらにも長く見える。姉妹リポジトリ `wht-get-ready-for-github`（`grg`）は README に使い方、別ファイルに開発手順を置いている。開発メモのファイル名は grg の `DEVELOPING.md` ではなく、このリポジトリでは `development.md` にする。

## Decision

- トップの [README.md](../../README.md) には、コマンドの概要、前提、入れ方、使い方、表示する内容、やらないことだけを書く。
- ビルド、版数の見方、リリースはトップの [development.md](../../development.md) に書く。README からはリンクするだけにする。
- `DEVELOPING.md` は置かない。

## Consequences

- 利用者が clone や Releases から入れたあと、README だけで使える。
- 開発者は `development.md` を見れば `make` と `make release*` が分かる。
- grg の `DEVELOPING.md` とはファイル名が異なる。
