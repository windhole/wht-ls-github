# lsg

ローカルディレクトリにある GitHub リポジトリ群の状態を、一括で確認・表示するコマンドです。

既定では `~/Documents/GitHub` 直下を、ディレクトリの mtime が新しい順に見ます。各リポジトリで `git fetch` したうえで、公開範囲・リモートとの差・未コミット変更・`.devcontainer` の有無を表に出します。

| フラグ | 動き |
|--------|------|
| （なし） | 既定の `~/Documents/GitHub` を走査する |
| `--dir PATH` | 指定したディレクトリ直下を走査する |
| `--version` | バージョンだけ表示する |
| `--help` / `-h` | 使い方を表示する |

## 前提

PATH にあれば足りるもの:

- git
- [GitHub CLI (`gh`)](https://cli.github.com/)（対象ホストにログイン済み）

ホストやユーザーは `gh` の現在の認証先に従います。URL は埋め込んでいません。

配布バイナリの対象は macOS（Apple Silicon）です。実行側に Go は不要です。

## 入れ方

[Releases](https://github.com/windhole/wht-ls-github/releases) から `lsg` を落とし、実行権限を付けて PATH へ置きます。

```bash
chmod +x lsg
install -m 0755 lsg "$HOME/bin/lsg"
```

ソースからビルドする場合は [development.md](development.md) を見てください。

## 使い方

```bash
# 既定の ~/Documents/GitHub を走査する
lsg

# 走査先を指定する
lsg --dir /path/to/repos

# バージョンだけ
lsg --version
```

各リポジトリで `git fetch` します。リモートの更新を見るための副作用です。`gh repo view` が失敗したディレクトリは行から除外します。

## 表示すること

- **Repo**: ディレクトリ名
- **Visibility**: GitHub 上の Public / Private（`gh repo view --json isPrivate`）
- **Sync**: 追跡ブランチとの差（Synced / Pull needed / Push needed / Diverged）
- **Clean**: 未コミットの変更（modified / untracked）があるか
- **Dev**: `.devcontainer` があるか

隠しディレクトリ（`.` で始まるもの）は走査しません。

## やらないこと

- リポジトリの作成、push、pull の実行
- Organization 配下や、走査先より深い階層の再帰走査
- Intel Mac / Linux 向けバイナリの配布

## 出力例

※プライベートリポジトリ名はマスクしています

![lsg の出力例](./wht-ls-github_sample.png)
