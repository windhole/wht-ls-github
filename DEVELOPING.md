# 開発者向けメモ

このリポジトリをビルドしたり、生成物の元データを変えたりする人向けです。コマンドの使い方は [README.md](README.md) を見てください。

Go の標準ライブラリだけで実装しています。第三者パッケージは使いません。実行時に呼ぶ外部コマンドは `git` と `gh` だけです。

## レイアウト

- `src/` … Go のソース
- `dist/` … `make` の成果物（gitignore 対象）
- `data/` … ローカル作業用（gitignore 対象。リポジトリには含めない）
- `docs/adr/` … 設計判断
- `docs/devlog/` … 作業ログ

## 前提

- [Go](https://go.dev/) 1.22 以降（Homebrew なら `brew install go`）
- `make`
- リリースするときだけ [GitHub CLI (`gh`)](https://cli.github.com/)（ログイン済み）

## ビルド

成果物は `dist/wht-ls-github`（darwin/arm64）です。できたバイナリを Mac に置くときは Go は不要です。版数は `-ldflags` で `main.version` に埋め込みます。未指定時は `git describe`（無ければ `dev`）です。

```bash
make
make build VERSION=v0.1.0   # 版を明示する場合
./dist/wht-ls-github --version
```

PATH の通った場所へ置いて使います。

```bash
install -m 0755 dist/wht-ls-github "$HOME/bin/wht-ls-github"
```

ソースから直接走らせる場合:

```bash
go run ./src
go run ./src --dir /path/to/repos
```

## リリース

GitHub Releases へ載せるときは `make release` だけ実行します。最新の `vX.Y.Z` タグのパッチを 1 つ上げます。タグがまだ無ければ `v0.1.0` です。

```bash
make release
```

桁を上げたいとき:

```bash
make release-minor
make release-major
```

現在の版と、次の patch / minor / major を見る場合:

```bash
make show-version
```

作業ツリーがきれいな状態で、タグ作成・darwin/arm64 のビルド・`gh release create` まで行います。アップロードするファイルは `dist/wht-ls-github` で、Release 上の Asset 名は `wht-ls-github` です。
