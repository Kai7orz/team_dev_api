# team_dev_api
Web API のチーム開発

# Quick Start

## サーバーの起動

```bash
make run
```


# 開発者向け

## 静的解析
コミット時にlintを実行してエラーがあれば修正してください
```bash
make lint
```

## コードフォーマット
コミット時にコードフォーマットを実行してください
```bash
make fmt
```


## APIドキュメンテーションの生成
エンドポイントの追加・変更時にAPIドキュメンテーションを生成・更新してください
1. `swag` CLIのインストール（初回のみ）

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. go.modのパッケージをインストール（初回のみ）

```bash
go mod tidy
```

3. APIドキュメンテーションを生成・更新

```bash
swag init -g cmd/server/main.go -o ./docs
```

4. サーバーを起動

```bash
go run cmd/server/main.go
```

5. ドキュメントを確認

以下にアクセス
```bash
http://localhost:8080/swagger/index.html
```