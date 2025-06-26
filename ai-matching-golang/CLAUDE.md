# AI-Matching Golang Backend - Development Guide

## プロジェクト概要

このプロジェクトは、マルチテナント対応のクリニック管理システムのバックエンドAPIです。Clean Architecture/Hexagonal Architectureパターンに従い、保守性・テスタビリティ・拡張性を重視した設計になっています。

## 技術スタック

- **言語**: Go 1.24.4
- **Webフレームワーク**: Fiber v2 + Huma v2
- **データベース**: PostgreSQL
- **ORM/クエリビルダー**: SQLC (型安全なSQL)
- **認証**: JWT (現在はモック実装)
- **開発ツール**: Air (ホットリロード), golang-migrate (マイグレーション)

## プロジェクト構成

```
ai-matching-golang/
├── db/                    # データベース関連
│   ├── migrations/        # SQLマイグレーションファイル
│   ├── query/            # SQLCクエリ定義
│   └── sqlc/             # SQLC生成コード（自動生成）
├── src/                   # アプリケーションソースコード
│   ├── api/              # API層（プレゼンテーション層）
│   │   ├── auth/         # 認証必須エンドポイント
│   │   └── public/       # パブリックエンドポイント
│   ├── di/               # 依存性注入
│   ├── domain/           # ドメイン層（インターフェース定義）
│   └── infrastructure/   # インフラ層（実装）
├── docs/                  # ドキュメント
├── etc/                   # その他のリソース
└── tmp/                   # 一時ファイル（Air用）
```

## アーキテクチャ設計

### 1. Clean Architecture の実装

このプロジェクトはClean Architectureパターンに従い、以下の層に分離されています：

#### レイヤー構成
```
[API層] → [ユースケース層] → [リポジトリ層] → [データベース]
   ↓            ↓                ↓
Controller   Usecase         Repository
   ↓            ↓                ↓
Huma/Fiber  ビジネスロジック   SQLC生成コード
```

#### 各層の責務

1. **API層** (`src/api/`)
   - HTTPリクエスト/レスポンスの処理
   - 入力値の検証
   - エラーレスポンスの生成
   - OpenAPIドキュメントの定義

2. **ユースケース層** (`src/api/*/usecase/`)
   - ビジネスロジックの実装
   - トランザクション管理
   - データ変換（ドメインモデル ↔ DTOs）

3. **ドメイン層** (`src/domain/`)
   - インターフェース定義
   - ビジネスルールの定義
   - エンティティの定義

4. **インフラ層** (`src/infrastructure/`)
   - リポジトリインターフェースの実装
   - 外部サービスとの通信
   - データベースアクセス

### 2. 依存性注入パターン

```go
// di/container.go
type Container struct {
    DB      *sqlx.DB
    Queries db.Querier  // インターフェースを使用（テスト可能）
}

// 使用例
container := di.NewContainer()
userRepo := repository.NewUserRepository(container.Queries)
userUsecase := usecase.NewUserUsecase(userRepo)
userController := controller.NewUserController(userUsecase)
```

## SQLC の使用方法

### 1. 設定ファイル (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/query/"           # SQLクエリファイルの場所
    schema: "db/migrations/"       # スキーマ定義の場所
    gen:
      go:
        package: "db"
        out: "db/sqlc"             # 生成コードの出力先
        emit_json_tags: true       # JSON タグを生成
        emit_interface: true       # Querier インターフェースを生成
        emit_empty_slices: true    # 空の結果を nil ではなく [] で返す
```

### 2. クエリの書き方

#### 基本的なCRUD操作

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET first_name = $1, last_name = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

#### クエリタイプ
- `:one` - 単一レコードを返す
- `:many` - 複数レコードを返す
- `:exec` - 実行のみ（戻り値なし）
- `:execrows` - 影響を受けた行数を返す

### 3. NULL値の扱い

```go
// NULL可能なフィールドの作成
params := db.CreateUserParams{
    Email:         email,
    PasswordHash:  hash,
    FirstName:     sql.NullString{String: firstName, Valid: firstName != ""},
    LastName:      sql.NullString{String: lastName, Valid: lastName != ""},
}

// NULL値の読み取り
if user.FirstName.Valid {
    response.FirstName = &user.FirstName.String
}
```

### 4. トランザクション処理

```go
// トランザクション開始
tx, err := container.DB.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// トランザクション内でクエリ実行
qtx := container.Queries.WithTx(tx)
user, err := qtx.CreateUser(ctx, params)
if err != nil {
    return err
}

// コミット
return tx.Commit()
```

## API設計パターン

### 1. エンドポイント構成

```
/api/v1/public/     # 認証不要
  - /health         # ヘルスチェック
  - /auth/login     # ログイン
  - /auth/register  # ユーザー登録
  - /auth/refresh   # トークンリフレッシュ

/api/v1/auth/       # 認証必須
  - /users          # ユーザー管理
  - /organizations  # 組織管理
  - /tenants        # テナント管理
```

### 2. API機能モジュールの構成

各APIエンドポイントは、以下の5つのコンポーネントで構成されています：

#### フォルダ構成
```
src/api/
├── auth/           # 認証が必要なエンドポイント
│   ├── users/
│   │   ├── controller/    # HTTPリクエスト処理
│   │   ├── requests/      # リクエストDTO定義
│   │   ├── response/      # レスポンスDTO定義
│   │   ├── router/        # ルーティング設定
│   │   └── usecase/       # ビジネスロジック
│   ├── organizations/
│   └── tenants/
└── public/         # 認証不要のエンドポイント
    ├── auth/
    │   ├── controller/
    │   ├── requests/
    │   ├── response/
    │   ├── router/
    │   └── usecase/
    └── health/
```

#### 各コンポーネントの役割

1. **Controller** (`controller/`)
   - HTTPリクエストとレスポンスの処理
   - Humaフレームワークのハンドラー実装
   - ユースケースの呼び出し
   - エラーレスポンスの生成
   ```go
   type UserController struct {
       userUsecase *usecase.UserUsecase
   }
   
   func (c *UserController) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
       // リクエストをユースケースに渡し、レスポンスを返す
   }
   ```

2. **Requests** (`requests/`)
   - APIリクエストのデータ構造定義
   - バリデーションルールの設定
   - Humaの自動バリデーション用タグ
   ```go
   type CreateUserRequest struct {
       Email     string  `json:"email" validate:"required,email" doc:"User email"`
       Password  string  `json:"password" validate:"required,min=8" doc:"User password"`
       FirstName *string `json:"first_name" doc:"User first name"`
   }
   ```

3. **Response** (`response/`)
   - APIレスポンスのデータ構造定義
   - データベースモデルからAPIレスポンスへの変換
   - クライアントに返すデータ形式
   ```go
   type UserResponse struct {
       ID        int64   `json:"id"`
       Email     string  `json:"email"`
       FirstName *string `json:"first_name"`
       CreatedAt string  `json:"created_at"`
   }
   ```

4. **Router** (`router/`)
   - エンドポイントのルーティング設定
   - HTTPメソッドとパスの定義
   - OpenAPI仕様の設定
   - 認証要件の指定
   ```go
   func RegisterRoutes(api huma.API, controller *controller.UserController) {
       huma.Register(api, huma.Operation{
           OperationID: "create-user",
           Method:      "POST",
           Path:        "/api/v1/auth/users",
           Summary:     "Create user",
           Tags:        []string{"Users"},
           Security:    []map[string][]string{{"bearer": {}}},
       }, controller.CreateUser)
   }
   ```

5. **Usecase** (`usecase/`)
   - ビジネスロジックの実装
   - リポジトリの呼び出し
   - トランザクション管理
   - データ変換とビジネスルールの適用
   ```go
   type UserUsecase struct {
       userRepo repository.UserRepository
   }
   
   func (u *UserUsecase) CreateUser(ctx context.Context, req requests.CreateUserRequest) (*response.UserResponse, error) {
       // パスワードのハッシュ化
       // リポジトリを使用してユーザー作成
       // レスポンスの生成
   }
   ```

#### 新しい機能を追加する際の手順

1. 該当するディレクトリ（`auth/` または `public/`）に新しい機能フォルダを作成
2. 5つのサブフォルダ（controller, requests, response, router, usecase）を作成
3. 各コンポーネントを実装
4. メインのルーターファイル（`src/api/router.go`）に新しいルートを登録

### 3. Humaフレームワークの使用

#### リクエスト/レスポンス定義

```go
// リクエスト
type CreateUserInput struct {
    Body requests.CreateUserRequest `doc:"User creation request"`
}

type CreateUserRequest struct {
    Email     string  `json:"email" validate:"required,email" doc:"User email"`
    Password  string  `json:"password" validate:"required,min=8" doc:"User password"`
    FirstName *string `json:"first_name" doc:"User first name"`
    LastName  *string `json:"last_name" doc:"User last name"`
}

// レスポンス
type CreateUserOutput struct {
    Body response.UserResponse
}

type UserResponse struct {
    ID        int64   `json:"id"`
    Email     string  `json:"email"`
    FirstName *string `json:"first_name"`
    LastName  *string `json:"last_name"`
}
```

#### エンドポイント登録

```go
huma.Register(api, huma.Operation{
    OperationID: "create-user",
    Method:      "POST",
    Path:        "/api/v1/auth/users",
    Summary:     "Create user",
    Description: "Create a new user",
    Tags:        []string{"Users"},
    Security:    []map[string][]string{{"bearer": {}}},
}, userController.CreateUser)
```

### 4. エラーハンドリング

```go
// グローバルエラーハンドラー
app.Use(func(c *fiber.Ctx) error {
    err := c.Next()
    if err != nil {
        code := fiber.StatusInternalServerError
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
        }
        return c.Status(code).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return nil
})
```

## 開発時の注意事項

### 1. コーディング規約

#### 命名規則
- パッケージ名: 小文字、単数形 (`repository`, `controller`)
- インターフェース名: 能力を表す名前 (`UserRepository`, `Querier`)
- 構造体名: PascalCase (`UserController`, `CreateUserRequest`)
- メソッド名: PascalCase (`GetUser`, `CreateOrganization`)

#### ファイル構成
- 1ファイル1型を基本とする
- 関連する機能はパッケージでグループ化
- インポートは標準ライブラリ、外部ライブラリ、内部パッケージの順

### 2. リポジトリパターンの実装

```go
// インターフェース定義 (domain/interface/repository/)
type UserRepository interface {
    GetUser(ctx context.Context, id int64) (db.User, error)
    CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
}

// 実装 (infrastructure/repository/)
type userRepository struct {
    queries db.Querier  // インターフェースを使用
}

func NewUserRepository(queries db.Querier) repository.UserRepository {
    return &userRepository{queries: queries}
}
```

### 3. ユースケースの実装

```go
type UserUsecase struct {
    userRepo repository.UserRepository
}

func (u *UserUsecase) CreateUser(ctx context.Context, req requests.CreateUserRequest) (*response.UserResponse, error) {
    // 1. バリデーション
    // 2. ビジネスロジック（パスワードハッシュ化など）
    // 3. リポジトリ呼び出し
    // 4. レスポンス変換
    return response, nil
}
```

### 4. コンテキストの伝播

すべてのデータベース操作とビジネスロジックでコンテキストを伝播させること：

```go
func (r *userRepository) GetUser(ctx context.Context, id int64) (db.User, error) {
    return r.queries.GetUser(ctx, id)  // 必ずctxを渡す
}
```

### 5. データベースマイグレーション

```bash
# マイグレーションファイルの作成
make migrate-create name=add_new_table

# マイグレーションの実行
make migrate-up

# ロールバック
make migrate-down
```

マイグレーションファイルの命名規則：
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

### 6. SQLCコード生成

```bash
# SQLCコードの再生成
make sqlc-generate

# 注意: db/sqlc/ 内のファイルは自動生成なので編集しない
```

## セキュリティ考慮事項

### 1. 認証・認可

現在の実装ではモックトークンを使用していますが、本番環境では以下を実装する必要があります：

- 適切なJWTライブラリの使用
- トークンの署名と検証
- リフレッシュトークンの安全な管理
- ミドルウェアでの認証チェック

### 2. 入力検証

- Humaのバリデーションタグを活用
- カスタムバリデータの実装
- SQLインジェクション対策（SQLCが自動的に対応）

### 3. CORS設定

本番環境では適切なOriginを設定：

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://your-domain.com",  // 本番では具体的なドメインを指定
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
}))
```

### 4. エラーメッセージ

- 本番環境では詳細なエラーメッセージを露出しない
- ログには詳細を記録し、レスポンスは一般的なメッセージを返す

## パフォーマンス最適化

### 1. データベース接続プール

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 2. クエリ最適化

- 必要なカラムのみをSELECT
- 適切なインデックスの使用
- N+1問題の回避（JOINを使用）

### 3. ページネーション

すべてのリスト系APIでページネーションを実装：

```sql
-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;
```

## テスト戦略

### 1. 単体テスト

```go
// モックを使用したリポジトリテスト
type mockQuerier struct {
    mock.Mock
}

func (m *mockQuerier) GetUser(ctx context.Context, id int64) (db.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(db.User), args.Error(1)
}
```

### 2. 統合テスト

- Dockerを使用したPostgreSQLでのテスト
- マイグレーションの実行確認
- エンドツーエンドのAPIテスト

## 開発フロー

### 1. 新機能の追加手順

1. データベーススキーマの設計とマイグレーション作成
2. SQLCクエリの作成
3. `make sqlc-generate` でコード生成
4. リポジトリインターフェースの定義
5. リポジトリ実装の作成
6. ユースケースの実装
7. コントローラーの実装
8. ルーターへの登録
9. テストの作成

### 2. 開発コマンド

```bash
# 開発サーバーの起動（ホットリロード付き）
make run

# ビルド
make build

# テスト実行
make test

# リント実行
make lint

# フォーマット
make fmt
```

## マルチテナント対応

### 1. データモデル

```
Organization (組織)
    ↓
Tenant (テナント/クリニック)
    ↓
User (ユーザー)
```

### 2. テナント分離の実装

- すべてのクエリでテナントIDによるフィルタリング
- ミドルウェアでテナントコンテキストの設定
- サブドメインベースのテナント識別

## トラブルシューティング

### 1. SQLC生成エラー

```bash
# SQLCのバージョン確認
sqlc version

# クエリ構文の検証
sqlc compile
```

### 2. マイグレーションエラー

```bash
# 現在のマイグレーションステータス確認
migrate -database $DATABASE_URL -path db/migrations version

# 強制的に特定バージョンに設定
migrate -database $DATABASE_URL -path db/migrations force VERSION
```

### 3. 依存関係の問題

```bash
# 依存関係の整理
go mod tidy

# キャッシュクリア
go clean -modcache
```

## ベストプラクティスまとめ

1. **Clean Architectureの維持**: 各層の責務を明確に分離
2. **インターフェースの活用**: テスタビリティと柔軟性の確保
3. **コンテキストの伝播**: タイムアウトとキャンセレーションの対応
4. **エラーハンドリング**: 適切なエラーメッセージとログ記録
5. **型安全性**: SQLCによる型安全なデータベースアクセス
6. **ドキュメント**: OpenAPI仕様の自動生成を活用
7. **セキュリティ**: 入力検証、認証・認可の適切な実装
8. **パフォーマンス**: インデックス、ページネーション、接続プールの最適化