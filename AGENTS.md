# Repository Guidelines

## プロジェクト構成とモジュール構成
リポジトリは Go モノリポ構成で、主要コードはルート直下の `main.go` と各 AWS リソース向けモジュール（`instances.go`、`security_groups.go`、`loadbalancers.go`、`confluence.go`）に分割されています。Confluence 投稿テンプレートは `template.go` にまとめています。共通の依存関係は `go.mod` で管理し、追加の設定ファイルは `renovate.json` に配置しています。将来テストを追加する場合は、対象モジュールと同じ階層に `_test.go` を配置してください。

## ビルド・テスト・開発コマンド
`go mod tidy` を実行して依存関係を同期し、未使用モジュールを削除します。`go build ./...` で全パッケージをビルドし、出力バイナリは必要に応じて `bin/` にまとめてください。動作確認には `go run ./main.go --help` を使い、Kingpin CLI 引数の挙動を確認します。ユニットテストは `go test ./...`、詳細ログが必要な場合は `go test -v ./...` を利用します。静的解析は `go vet ./...` を基準とし、ローカルで追加ツールを使う場合は README に追記してください。

## コーディングスタイルと命名規約
Go 言語標準に従い、`gofmt` および `goimports` で整形します。パッケージ内の公開関数は AWS リソース名と動詞を組み合わせた UpperCamelCase（例: `RenderSecurityGroups`）、プライベート関数は lowerCamelCase を推奨します。構造体フィールドは Confluence JSON や AWS SDK モデルに合わせ、略語は大文字 (`ID`, `URL`) を維持してください。複雑な処理は短い関数に分割し、テンプレート操作は `template.go` に集約します。

## テストガイドライン
テストは Go の標準 testing パッケージを使用し、テーブル駆動でケースを増やしやすくしてください。ファイル名はターゲットと同一で `_test.go` サフィックスを付けます。外部 API コールは `aws` または `confluence` パッケージをモックするインターフェースを導入し、`TestMain` で環境変数を設定します。境界値とエラー処理を最低 1 ケースずつ含め、`go test -cover` で主要ロジック（特にテンプレート生成と一覧取得）が 80% 以上のカバレッジになるよう目指します。

## コミットとプルリクエストガイドライン
コミットメッセージは Conventional Commits に従い、`<type>(<scope>): <subject>` 形式で記述します。`type` は `feat`（機能追加）、`fix`（バグ修正）、`docs`（ドキュメント）、`chore`（雑務）、`build`（依存更新）などから選択し、必要に応じて `aws` や `confluence` などの `scope` を指定してください。`subject` は英語で現在形・動詞始まりの一文（例: `fix(aws): handle session creation errors`）とし、詳細は本文に記述します。プルリクエストでは目的、主要変更点、検証手順、影響範囲を箇条書きで記載し、関連 Issue や Jira を `Closes #123` 形式でリンクします。レビュー前に `go test ./...` と `go vet ./...` を実行し、結果を PR に貼付してください。

## セキュリティと設定のヒント
Confluence と AWS の資格情報は環境変数でのみ渡し、`.env` やコード内にハードコーディングしないでください。ローカル検証には AWS IAM の最小権限ロールを作成し、`CONFLUENCE_` 系変数は一時トークンを利用します。機密値を共有する場合はチームのシークレットマネージャーを経由し、ログに URL やトークンが残らないよう注意してください。
