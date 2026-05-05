import { $ } from "bun";
import { readdir, stat } from "node:fs/promises";
import { join } from "node:path";

// 探索対象のルートディレクトリ
const GITHUB_DIR = join(process.env.HOME!, "Documents/GitHub");

/**
 * 個別のリポジトリの状態を解析する関数
 * @param dirName ディレクトリ名
 */
async function getRepoStatus(dirName: string) {
  const fullPath = join(GITHUB_DIR, dirName);
  
  try {
    // --- [並列実行セクション] ---
    // 複数の非同期処理を同時に開始して、トータルの待ち時間を短縮します。
    const [ghDataRaw, _fetch, statusText, devContainerExists, syncStatus] = await Promise.all([
      // 1. GitHub CLIを使って、リポジトリのメタデータ（非公開かどうか）を取得
      $`cd ${fullPath} && gh repo view --json isPrivate`.quiet().text().catch(() => null),
      
      // 2. リモートの最新情報を取得（これを行わないと、リモートの更新を検知できない）
      $`cd ${fullPath} && git fetch`.quiet().catch(() => null),
      
      // 3. ローカルの未コミット変更（modified）があるか確認 (-s は短縮形式)
      $`cd ${fullPath} && git status -s`.quiet().text(),
      
      // 4. .devcontainer ディレクトリが存在するかチェック
      stat(join(fullPath, ".devcontainer")).then(() => true).catch(() => false),
      
      // 5. リモートとローカルの「コミット数の差」を計算
      // --left-right: HEAD(左)とリモート(右)の差を出す
      // --count: 数値で出力
      // HEAD...@{u}: 現在のブランチと、その追跡対象(upstream)を比較
      $`cd ${fullPath} && git rev-list --left-right --count HEAD...@{u}`.quiet().text().catch(() => "0\t0")
    ]);

    // ghDataRawがnull = GitHubにリポジトリが存在しない、または権限がない場合は一覧から除外
    if (!ghDataRaw) return null;
    const ghData = JSON.parse(ghDataRaw);

    // --- [同期状態の判定ロジック] ---
    // syncStatus は "1\t2" (Ahead 1, Behind 2) のような文字列で返ってくるので分解
    const [ahead, behind] = syncStatus.trim().split("\t").map(Number);
    
    let syncLabel = "✅ Synced";
    if (behind > 0 && ahead > 0) syncLabel = `🔄 Diverged (↑${ahead} ↓${behind})`; // 両方に独自の更新がある
    else if (behind > 0) syncLabel = `📥 Pull needed (↓${behind})`;               // リモートが先行（Pullが必要）
    else if (ahead > 0) syncLabel = `📤 Push needed (↑${ahead})`;               // ローカルが先行（Pushが必要）

    // 表示用オブジェクトの生成
    return {
      Repo: dirName,
      Visibility: ghData.isPrivate ? "🔒 Priv" : "🌐 Pub",
      Sync: syncLabel,
      Clean: statusText.trim() === "" ? "✅" : "⚠️ Mod", // 未コミット変更の有無
      Dev: devContainerExists ? "📦" : "-",
    };
  } catch (error) {
    // 予期せぬエラー（ディレクトリが削除された等）が発生した場合はその行をスキップ
    return null;
  }
}

// --- [メイン処理セクション] ---

console.log("Checking GitHub repositories (fetching remotes)...");

// 1. ルートディレクトリ内のファイル・フォルダ一覧を取得
const entries = await readdir(GITHUB_DIR, { withFileTypes: true });

// 2. ディレクトリの最終更新日時(mtime)を取得し、ソートの準備をする
const dirInfos = await Promise.all(
  entries
    .filter(e => e.isDirectory() && !e.name.startsWith(".")) // 隠しフォルダを除外
    .map(async (e) => {
      const fullPath = join(GITHUB_DIR, e.name);
      const s = await stat(fullPath).catch(() => null);
      return s ? { name: e.name, mtime: s.mtime } : null;
    })
);

// 3. 更新日時が新しい順に名前をソート
const sortedDirs = dirInfos
  .filter((d): d is NonNullable<typeof d> => d !== null)
  .sort((a, b) => b.mtime.getTime() - a.mtime.getTime())
  .map(d => d.name);

// 4. ソート済みのディレクトリに対して、ステータス取得処理を並列実行
const results = (await Promise.all(sortedDirs.map(getRepoStatus)))
  .filter((r): r is NonNullable<typeof r> => r !== null);

// 5. テーブル形式でコンソールに出力
console.table(results);

// 6. ユーザー向けの凡例表示
console.log("\n--- Legend ---");
console.log("📥 Pull needed : Remote has changes you don't have.");
console.log("📤 Push needed : You have local commits not on remote.");
console.log("⚠️ Mod         : You have uncommitted files (modified/untracked).");
console.log("📦             : .devcontainer environment detected.");
