# 0007. make の既定はホスト向けビルド、版数操作は grg と同じ

Date: 2026-09-06
Status: Accepted

## Context

ADR-0003 は `make` / `make build` も darwin/arm64 固定にしていた。開発や確認で `make` したあと `./dist/wht-ls-github --version` をすぐ走らせたいが、Linux や Intel Mac ではそのバイナリを実行できない。配布物はこれまでどおり Apple Silicon 向けでよい。

版数の確認と、パッチ / マイナー / メジャーを上げてタグ・ビルド・GitHub Releases まで進める操作は、姉妹リポジトリ `wht-get-ready-for-github`（`grg`）と同じにしておきたい。

## Decision

- 既定の `make` / `make build` は、そのマシンの `go env GOOS` / `GOARCH` 向けに `dist/wht-ls-github` を出す。`CGO_ENABLED=0` のまま。
- `make release` / `release-minor` / `release-major` の成果物だけ `GOOS=darwin GOARCH=arm64` にする。GitHub Releases の対象は変えない。
- 版数の正は `vMAJOR.MINOR.PATCH` の git タグ。ファイルには持たない。
- `make show-version`（別名 `make version`）で現在のタグと次の patch / minor / major を出す。
- `make release` は最新タグのパッチを +1（無ければ `v0.1.0`）。桁を上げるときは `make release-minor` / `make release-major`。流れは grg と同じ（きれいな作業ツリー、テスト、ビルド、タグ、`gh release create`）。

## Consequences

- `make` の直後に、そのマシンでバイナリを実行して版数やヘルプを確認できる。
- 配布用の darwin/arm64 は `make release*` のときだけ作る。
- ADR-0003 の「`make` も常に darwin/arm64」は、この ADR で改める。リリース対象は ADR-0003 のまま。
