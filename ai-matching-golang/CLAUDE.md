# AI-Matching Golang Backend - Development Guide

## プロジェクト概要

このプロジェクトは、マルチテナント対応のクリニック管理システムのバックエンドAPIです。Clean Architecture/Hexagonal Architectureパターンに従い、保守性・テスタビリティ・拡張性を重視した設計になっています。

## 技術スタック

- **言語**: Go 1.24.4
- **Webフレームワーク**: Fiber v2 + Huma v2
- **データベース**: PostgreSQL
- **ORM/クエリビルダー**: SQLC (型安全なSQL)
- **認証**: AWS Cognito (JWT)
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
│   │   └── interface/
│   │       ├── repository/  # リポジトリインターフェース
│   │       └── external/    # 外部サービスインターフェース
│   └── infrastructure/   # インフラ層（実装）
│       ├── repository/   # リポジトリ実装
│       ├── middleware/   # 認証・認可ミドルウェア
│       └── external/     # 外部サービス実装（AWS SDK等）
├── docs/                  # ドキュメント
├── etc/                   # その他のリソース
└── tmp/                   # 一時ファイル（Air用）
```

## 最近の重要な変更

### 1. UUID型の導入とパラメータ命名規則の統一

すべてのエンドポイントのパスパラメータが統一された命名規則に変更されました：

- `id` → `organizationId` (組織関連エンドポイント)
- `id` → `userId` (ユーザー関連エンドポイント)
- `id` → `tenantId` (テナント関連エンドポイント)

**例：**
```go
// 旧: /api/v1/organizations/{id}/tenants/{id}
// 新: /api/v1/organizations/{organizationId}/tenants/{tenantId}

type GetTenantInput struct {
    OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
    TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}
```

### 2. 認証・認可ミドルウェアの強化

`FiberMiddleware`が組織・テナントの所属確認を自動的に行うように改善されました：

- URLパスから`organizationId`と`tenantId`を自動抽出
- ユーザーの組織所属を自動検証
- ユーザーのテナント所属を自動検証
- **重要**: コントローラーでの所属確認は不要になりました

### 3. システム管理者機能の追加

ユーザーテーブルに`is_system_admin`フィールドが追加され、システム全体の管理者を識別できるようになりました。

## API設計の重要な規約

### エンドポイント形式

**必須**: やむを得ない事情がない限り、エンドポイントは以下の形式で記述してください：

```
/organizations/{organizationId}/tenants/{tenantId}
/organizations/{organizationId}/users/{userId}
/organizations/{organizationId}/tenants/{tenantId}/users/{userId}
```

### 認証・認可ミドルウェアについて

**重要**: 組織やテナントへのユーザー所属確認は`FiberMiddleware`で完了しています。

```go
// middleware/auth_middleware.go での処理内容：
// 1. JWTトークンの検証
// 2. ユーザー情報の取得とコンテキストへの設定
// 3. URLパスからorganizationId/tenantIdを抽出
// 4. ユーザーの組織所属を検証（organizationIdがある場合）
// 5. ユーザーのテナント所属を検証（tenantIdがある場合）
```

**したがって、コントローラーでは以下の確認は不要です：**
- ユーザーが組織に所属しているかの確認
- ユーザーがテナントに所属しているかの確認

コントローラーは純粋にビジネスロジックの実行に集中できます。

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

[API層] → [ユースケース層] → [外部サービス層] → [外部API]
   ↓            ↓                ↓
Controller   Usecase         External Client
   ↓            ↓                ↓
Huma/Fiber  ビジネスロジック   AWS SDK/外部ライブラリ
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
   - ミドルウェアの実装
   
   **重要**: 外部ライブラリに依存する実装は以下のディレクトリ構造に配置：
   - `src/infrastructure/external/` - 外部サービスの実装（AWS SDK、外部APIクライアントなど）
   - `src/domain/interface/external/` - 外部サービスのインターフェース定義

### 2. 依存性注入パターン

依存性注入は `ai-matching-golang/src/di/container.go` で一元管理されています。

```go
// di/container.go
type Container struct {
    DB            *sqlx.DB
    Queries       db.Querier
    CognitoClient external.CognitoClient
    
    // Repositories
    UserRepository         repository.UserRepository
    AuthRepository         repository.AuthRepository
    OrganizationRepository repository.OrganizationRepository
    TenantRepository       repository.TenantRepository
    TenantUserRepository   repository.TenantUserRepository
    
    // Usecases
    AuthUsecase         *publicAuthUsecase.AuthUsecase
    UserUsecase         *userUsecase.UserUsecase
    OrganizationUsecase *organizationUsecase.OrganizationUsecase
    TenantUsecase       *tenantUsecase.TenantUsecase
    TenantUserUsecase   *tenantUserUsecase.TenantUserUsecase
    
    // Controllers
    AuthController         *publicAuthController.AuthController
    UserController         *userController.UserController
    OrganizationController *authController.OrganizationController
    TenantController       *tenantController.TenantController
    TenantUserController   *tenantUserController.TenantUserController
    HealthController       *healthController.HealthController
}

// 使用例
container := di.NewContainer()
// すべての依存関係はコンテナ内で初期化される
// ルーターはコンテナからコントローラーを取得して使用
authRouter.RegisterAuthRoutes(api, publicAPI, container.AuthController)
```

**重要**: すべてのユースケース、リポジトリ、外部サービスの注入は `ai-matching-golang/src/di/container.go` で行われます。新しい依存関係を追加する際は、このファイルを更新してください。

## ミドルウェアの詳細

### AuthMiddleware の動作

```go
// src/infrastructure/middleware/auth_middleware.go

func (m *AuthMiddleware) FiberMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. JWTトークンの検証
        token, err := m.jwtValidator.ValidateToken(tokenString)
        
        // 2. ユーザー情報の取得とコンテキスト設定
        c.Locals("user_id", userID)
        c.Locals("organization_id", orgID)
        c.Locals("tenant_id", tenantID)
        
        // 3. URLパスから organizationId を抽出して検証
        organizationId := extractOrganizationIdFromPath(path)
        if organizationId != "" {
            // ユーザーの組織所属を自動検証
        }
        
        // 4. URLパスから tenantId を抽出して検証
        tenantId := extractTenantIdFromPath(path)
        if tenantId != "" {
            // ユーザーのテナント所属を自動検証
            _, err := m.tenantUserRepo.GetTenantUser(ctx, tenantID, userID)
        }
        
        return c.Next()
    }
}
```

### コンテキストからのユーザー情報取得

```go
// GetUserFromContext でユーザー情報を取得
userContext, err := middleware.GetUserFromContext(ctx)

// UserContext構造体
type UserContext struct {
    UserID         uuid.UUID
    Email          string
    Token          string
    OrganizationID uuid.UUID
    Tenant         *Tenant
}
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
/api/v1/public/                           # 認証不要
  - /health                               # ヘルスチェック
  - /auth/login                           # ログイン
  - /auth/register                        # ユーザー登録
  - /auth/refresh                         # トークンリフレッシュ

/api/v1/                             # 認証必須（ミドルウェア適用）
  - /organizations/{organizationId}/*     # 組織関連エンドポイント
  - /organizations/{organizationId}/tenants/{tenantId}/*  # テナント関連
  - /organizations/{organizationId}/users/{userId}/*      # ユーザー関連
  - /tenants/subdomain/{subdomain}       # サブドメイン検索（特殊ケース）
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
   type TenantController struct {
       usecase *usecase.TenantUsecase
   }
   
   func (c *TenantController) CreateTenantInOrganization(ctx context.Context, input *CreateTenantInput) (*CreateTenantOutput, error) {
       // パスパラメータからorganizationIdを設定
       input.Body.OrganizationID = input.OrganizationID
       // ユースケースに処理を委譲
       resp, err := c.usecase.CreateTenant(ctx, input.Body)
       return &CreateTenantOutput{Body: *resp}, nil
   }
   ```

2. **Requests** (`requests/`)
   - APIリクエストのデータ構造定義
   - バリデーションルールの設定
   - Humaの自動バリデーション用タグ
   ```go
   type CreateTenantRequest struct {
       Name           string    `json:"name" validate:"required" doc:"Tenant name"`
       Subdomain      string    `json:"subdomain" validate:"required,min=3" doc:"Tenant subdomain"`
       OrganizationID uuid.UUID `json:"-"` // パスパラメータから設定
   }
   ```

3. **Response** (`response/`)
   - APIレスポンスのデータ構造定義
   - データベースモデルからAPIレスポンスへの変換
   - クライアントに返すデータ形式
   ```go
   type TenantResponse struct {
       ID             uuid.UUID `json:"id"`
       Name           string    `json:"name"`
       Subdomain      string    `json:"subdomain"`
       OrganizationID uuid.UUID `json:"organization_id"`
       IsActive       bool      `json:"is_active"`
       CreatedAt      string    `json:"created_at"`
   }
   ```

4. **Router** (`router/`)
   - エンドポイントのルーティング設定
   - HTTPメソッドとパスの定義
   - OpenAPI仕様の設定
   - 認証要件の指定
   ```go
   func RegisterTenantRoutes(api huma.API, router fiber.Router, controller *controller.TenantController) {
       huma.Register(api, huma.Operation{
           OperationID: "create-tenant-in-organization",
           Method:      "POST",
           Path:        "/api/v1/organizations/{organizationId}/tenants",
           Summary:     "Create tenant in organization",
           Tags:        []string{"Tenants"},
           Security:    []map[string][]string{{"bearer": {}}},
       }, controller.CreateTenantInOrganization)
   }
   ```

5. **Usecase** (`usecase/`)
   - ビジネスロジックの実装
   - リポジトリの呼び出し
   - トランザクション管理
   - データ変換とビジネスルールの適用
   ```go
   type TenantUsecase struct {
       tenantRepo repository.TenantRepository
   }
   
   func (u *TenantUsecase) CreateTenant(ctx context.Context, req requests.CreateTenantRequest) (*response.TenantResponse, error) {
       // サブドメインの重複チェック
       // テナント作成
       // レスポンスの生成
   }
   ```

#### 新しい機能を追加する際の手順

1. 該当するディレクトリ（`auth/` または `public/`）に新しい機能フォルダを作成
2. 5つのサブフォルダ（controller, requests, response, router, usecase）を作成
3. 各コンポーネントを実装
4. DIコンテナ（`src/di/container.go`）に追加
5. メインのルーターファイル（`src/di/router.go`）に新しいルートを登録

### 3. Humaフレームワークの使用

#### リクエスト/レスポンス定義

```go
// リクエスト（パスパラメータ付き）
type GetTenantInput struct {
    OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
    TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

// リクエスト（ボディ付き）
type CreateTenantInput struct {
    OrganizationID uuid.UUID                   `path:"organizationId" doc:"Organization ID"`
    Body           requests.CreateTenantRequest `doc:"Tenant creation request"`
}

// レスポンス
type CreateTenantOutput struct {
    Body response.TenantResponse
}
```

#### エンドポイント登録

```go
huma.Register(api, huma.Operation{
    OperationID: "get-tenant-in-organization",
    Method:      "GET",
    Path:        "/api/v1/organizations/{organizationId}/tenants/{tenantId}",
    Summary:     "Get tenant in organization",
    Description: "Get tenant by ID within an organization",
    Tags:        []string{"Tenants"},
    Security:    []map[string][]string{{"bearer": {}}},
}, tenantController.GetTenantInOrganization)
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
- パスパラメータ: camelCase (`organizationId`, `tenantId`, `userId`)

#### ファイル構成
- 1ファイル1型を基本とする
- 関連する機能はパッケージでグループ化
- インポートは標準ライブラリ、外部ライブラリ、内部パッケージの順

### 2. リポジトリパターンの実装

```go
// インターフェース定義 (domain/interface/repository/)
type TenantUserRepository interface {
    GetTenantUser(ctx context.Context, tenantID, userID uuid.UUID) (*db.TenantUser, error)
    ListTenantUsers(ctx context.Context, tenantID uuid.UUID) ([]db.TenantUser, error)
}

// 実装 (infrastructure/repository/)
type tenantUserRepository struct {
    queries db.Querier  // インターフェースを使用
}

func NewTenantUserRepository(queries db.Querier) repository.TenantUserRepository {
    return &tenantUserRepository{queries: queries}
}
```

### 2.1. 外部サービスパターンの実装

```go
// インターフェース定義 (domain/interface/external/)
type CognitoClient interface {
    SignUp(ctx context.Context, email, password string, attributes map[string]string) (*cognitoidentityprovider.SignUpOutput, error)
    InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
}

// 実装 (infrastructure/external/cognito/)
type cognitoClient struct {
    client *cognitoidentityprovider.Client
}

func NewCognitoClient() (external.CognitoClient, error) {
    // AWS SDK設定
    return &cognitoClient{client: client}, nil
}
```

### 3. ユースケースの実装

```go
type TenantUserUsecase struct {
    tenantUserRepo repository.TenantUserRepository
    userRepo       repository.UserRepository
}

func (u *TenantUserUsecase) AddUserToTenant(ctx context.Context, req requests.AddUserToTenantRequest) (*response.TenantUserResponse, error) {
    // 1. バリデーション（ユーザーの存在確認など）
    // 2. ビジネスロジック（重複チェックなど）
    // 3. リポジトリ呼び出し
    // 4. レスポンス変換
    return response, nil
}
```

### 4. コンテキストの伝播

すべてのデータベース操作とビジネスロジックでコンテキストを伝播させること：

```go
func (r *tenantUserRepository) GetTenantUser(ctx context.Context, tenantID, userID uuid.UUID) (*db.TenantUser, error) {
    return r.queries.GetTenantUser(ctx, db.GetTenantUserParams{
        TenantID: tenantID,
        UserID:   userID,
    })
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

AWS Cognitoを使用したJWT認証が実装されています：

- JWTトークンの署名と検証
- リフレッシュトークンの管理
- ミドルウェアでの自動認証チェック
- 組織・テナントレベルでのアクセス制御

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
8. DIコンテナへの登録
9. ルーターへの登録
10. テストの作成

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
User (ユーザー) - tenant_users テーブルで関連付け
```

### 2. テナント分離の実装

- すべてのクエリでテナントIDによるフィルタリング
- ミドルウェアでテナントコンテキストの設定
- サブドメインベースのテナント識別
- URLパスベースのアクセス制御

### 3. アクセス制御の流れ

1. ユーザーがAPIにアクセス
2. AuthMiddlewareがJWTトークンを検証
3. URLパスから organizationId/tenantId を抽出
4. ユーザーの所属を自動検証
5. 検証成功後、コントローラーに処理を渡す
6. コントローラーは所属確認なしでビジネスロジックを実行


## ベストプラクティスまとめ

1. **Clean Architectureの維持**: 各層の責務を明確に分離
2. **インターフェースの活用**: テスタビリティと柔軟性の確保
3. **コンテキストの伝播**: タイムアウトとキャンセレーションの対応
4. **エラーハンドリング**: 適切なエラーメッセージとログ記録
5. **型安全性**: SQLCによる型安全なデータベースアクセス、UUID型の使用
6. **ドキュメント**: OpenAPI仕様の自動生成を活用
7. **セキュリティ**: 入力検証、認証・認可の適切な実装
8. **パフォーマンス**: インデックス、ページネーション、接続プールの最適化
9. **命名規則の統一**: パスパラメータは必ず具体的な名前（organizationId, tenantId, userId）を使用
10. **ミドルウェアの活用**: 共通処理（認証・認可）はミドルウェアで一元管理