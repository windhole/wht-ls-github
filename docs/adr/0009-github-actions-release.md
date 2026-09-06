# 0009. GitHub Actions から darwin/arm64 をビルドし Releases に載せる

Date: 2026-09-06
Status: Accepted

## Context

手元に Go や make がある PC が無くても、GitHub の Web UI でソースを直して、コンパイルと GitHub Releases 登録まで済ませたい。版数の上げ方と成果物（darwin/arm64 の `wht-ls-github`）は ADR-0003 / ADR-0007 のままにしたい。GoReleaser などの第三者リリースツールは使わない。

GitHub ホステッドランナーは Linux だが、`CGO_ENABLED=0` のクロスビルドで darwin/arm64 はすでに作れている。

## Decision

- `.github/workflows/release.yml` を置く。起動は `workflow_dispatch` のみ。入力は `patch` / `minor` / `major`（既定は `patch`）。
- main への push では自動リリースしない。Web で編集・コミットしたあと、Actions 画面から明示的に走らせる。
- ランナーは `ubuntu-latest`。使う Action は GitHub 公式の `actions/checkout` と `actions/setup-go` だけ。Go の版は `go.mod` に従う。
- 版数計算・テスト・ビルド・タグ・`gh release create` は既存の Makefile を呼ぶ。CI ではすでにそのコミットにいるので `git push origin HEAD` はしない。
- 手元の `make release*` は残す。どちらの入口も同じタグ規則。

## Consequences

- ブラウザだけで「編集 → Actions で Run workflow → Release にバイナリ」ができる。
- 誤って毎回の push で版が上がることはない。
- 手元と Actions を同時に走らせるとタグがぶつかることがある。片方だけ使う。
