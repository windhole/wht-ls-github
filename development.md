# 開発者向けメモ

このリポジトリをビルドしたり、リリースしたりする人向けです。コマンドの使い方は [README.md](README.md) を見てください。

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

`make`（または `make build`）は、いま動かしているマシン向けに `dist/wht-ls-github` を出します。できたバイナリを置く先のマシンでは Go は不要です。版数は `-ldflags` で `main.version` に埋め込みます。未指定時は `git describe`（無ければ `dev`）です。

```bash
make
make help
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

テスト:

```bash
make test
```

## 版数

版数の正は `vMAJOR.MINOR.PATCH` の git タグです。ファイルには持ちません。現在のタグと、次の patch / minor / major は次で見ます。

```bash
make show-version
make version          # 同じ
```

## リリース

GitHub Releases へ載せるときは `make release` だけ実行します。最新の `vX.Y.Z` タグのパッチを 1 つ上げてから、darwin/arm64 のバイナリをビルドして載せます。タグがまだ無ければ `v0.1.0` です。

```bash
make release          # パッチ +1（例: v0.1.0 → v0.1.1）
```

桁を上げたいとき:

```bash
make release-minor    # マイナー +1（例: v0.1.1 → v0.2.0）
make release-major    # メジャー +1（例: v0.2.0 → v1.0.0）
```

作業ツリーがきれいな状態で、テスト・darwin/arm64 のビルド・タグ作成・`gh release create` まで行います。アップロードするファイルは `dist/wht-ls-github` で、Release 上の Asset 名は `wht-ls-github` です。
