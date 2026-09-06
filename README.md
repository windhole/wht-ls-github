# wht-ls-github

macOS 上のローカルディレクトリにある GitHub リポジトリ群の状態を、一括で確認・表示するコマンドです。

## 主な機能

- **リポジトリ一覧の自動取得**: 既定では `~/Documents/GitHub` 直下をスキャンする
- **最終更新順ソート**: ディレクトリの mtime が新しい順
- **公開/非公開の判定**: `gh` で GitHub 上の Public / Private を表示
- **リモート同期状態の検知**:
  - `Pull needed`: リモートにある更新をローカルに反映していない
  - `Push needed`: ローカルにあるコミットをリモートに送っていない
- **ローカルのクリーン判定**: 未コミットの変更（Modified / Untracked）があるか
- **Dev Container 検知**: `.devcontainer` の有無

## 前提

PATH にあれば足りるもの:

- [GitHub CLI (`gh`)](https://cli.github.com/)（対象ホストにログイン済み）
- `git`

配布バイナリの対象は macOS（Apple Silicon）です。実行側に Go や bun は不要です。

## 入れ方

[Releases](https://github.com/windhole/wht-ls-github/releases) から `wht-ls-github` を落とし、実行権限を付けて PATH へ置きます。

```bash
chmod +x wht-ls-github
install -m 0755 wht-ls-github "$HOME/bin/wht-ls-github"
```

ソースからビルドする場合:

```bash
make
./dist/wht-ls-github --version
```

版数の確認と、パッチ / マイナー / メジャーを上げて GitHub Releases に載せる手順は [DEVELOPING.md](DEVELOPING.md) を見てください。


## 使い方

```bash
# 既定の ~/Documents/GitHub を走査する
wht-ls-github

# 走査先を指定する
wht-ls-github --dir /path/to/repos

# 版数
wht-ls-github --version
```

各リポジトリで `git fetch` します。リモートの更新を見るための副作用です。`gh repo view` が失敗したディレクトリは行から除外します。

## 出力例

※プライベートリポジトリ名はマスクしています

![wht-ls-githubの出力例](./wht-ls-github_sample.png)
