# API レスポンス・エラー標準（v1）

このドキュメントは、**バックエンド API のレスポンス形式（Response Envelope）** と **エラーコード体系（Error Code System）** を、ひとつの標準として定義します。
言語やフレームワーク（Go、Java、Node など）に依存せず、**すべてのバックエンド API** に同じように適用します。

> 目標: クライアント（Flutter／Web／SDK）が `status / code / error / message` だけで、**一貫した分岐処理**を行えるようにすること。

---

## 1. 基本原則（Principles）

1. **HTTP ステータス ＝ code**
   - レスポンス本文の `code` は、**常に HTTP ステータスと同じ値**にします。
2. **機械可読と人間可読の分離**
   - `error`: 機械が分岐するための**標準化されたコード**
   - `message`: ユーザー表示・ログ用の**メッセージ**（i18n キー、またはプレーンなメッセージ）
3. **後方互換性を基本とする**
   - フィールドの拡張は可能ですが、意味の変更や破壊的変更（breaking change）は最小限にします。
4. **2xx にエラーを載せない**
   - `status="error"` なのに HTTP 200 を返すレスポンスは禁止です。

---

## 2. Response Envelope（レスポンス構造）

### 2.1 成功レスポンス

```json
{
  "status": "success",
  "code": 200,
  "message": "ok",
  "data": { },
  "metadata": {
    "timestamp": "2026-02-12T10:30:00Z",
    "traceId": "abc123"
  }
}
```

### 2.2 失敗レスポンス

```json
{
  "status": "error",
  "code": 401,
  "error": "AUTH_TOKEN_EXPIRED",
  "message": "auth.token_expired",
  "metadata": {
    "timestamp": "2026-02-12T10:30:00Z",
    "traceId": "abc123"
  }
}
```

---

## 3. フィールド定義（Field Definitions）

| フィールド | 型 | 必須 | 目的 | ルール |
|---|---|---:|---|---|
| status | string | ○ | 成功／失敗の区別 | `"success"` または `"error"` |
| code | number | ○ | HTTP ステータスのミラー | **常に HTTP ステータスと同じ** |
| data | object/array/primitive | △ | 成功時のデータ | **成功時のみ使用**。データがない場合は省略（omit）を推奨 |
| error | string | △ | 分岐用のエラーコード | **失敗時に推奨**。標準の一覧に従う |
| message | string | ○ | 表示・説明用メッセージ | i18n 使用時はキー、未使用時はプレーンなメッセージも可 |
| errors | array | △ | フィールド単位の検証エラー | **422 で推奨**（第7章を参照） |
| metadata | object | ○ | 共通メタ情報 | timestamp は必須。traceId／requestId などで拡張可能 |

### 3.1 `message` のルール（i18n／非 i18n）

- i18n を使うサービス:
  - `message` は **i18n キー**（例: `auth.token_expired`）
- i18n を使わないサービス:
  - `message` は**プレーンな人間向けメッセージ**も可（例: `Token expired`）

> クライアントでの表示方針はプロジェクトの性質によって異なるため、`message` は**表示可能な値**であることを保証します。

### 3.2 `data` のルール（省略か null か）

- 返すデータがない場合、`data` は**省略（omit）を推奨**します。
- `null` も可能ですが、クライアント側の分岐処理が複雑になるため推奨しません。

---

## 4. Metadata 標準（Observability）

`metadata` は、観測・追跡（Observability）のために使用します。

### 4.1 必須フィールド

- `metadata.timestamp`: RFC3339／ISO8601（例: `2026-02-12T10:30:00Z`）

### 4.2 推奨する拡張

- `metadata.traceId`: 分散トレーシング／リクエスト追跡
- `metadata.requestId`: リクエスト識別子（ゲートウェイまたはサーバーで生成）
- `metadata.version`: レスポンススキーマ／サーバーのバージョン（任意）

例:

```json
"metadata": {
  "timestamp": "2026-02-12T10:30:00Z",
  "traceId": "abc123",
  "requestId": "req-9f8e",
  "version": "v1"
}
```

---

## 5. Error Code System（エラーコード体系）

### 5.1 命名規則（Naming）

- 形式: `UPPER_SNAKE_CASE`
- できる限り**具体的に**（原因 ＋ 対象 ＋ 状況）
- ドメイン／関心事の prefix を推奨します。

| Prefix | 意味 |
|---|---|
| REQUEST_* | リクエストの形式／パラメータ／検証 |
| AUTH_* | 認証・認可（Authentication／Authorization） |
| INTERNAL_* | サーバー内部／依存先の障害 |
| <DOMAIN>_* | ドメイン別（User／Book／Order など） |

例:
- `REQUEST_MALFORMED_JSON`
- `AUTH_INVALID_TOKEN`
- `USER_EMAIL_ALREADY_EXISTS`
- `BOOK_NOT_FOUND`

### 5.2 HTTP ステータス ↔ エラーコードの対応ルール

> HTTP ステータスは「失敗の種類」、`error` は「失敗の正確な意味」を表します。

| HTTP | 意味 | 代表的な error コードの例 |
|---|---|---|
| 400 Bad Request | 不正なリクエスト（形式／パース） | `REQUEST_INVALID`、`REQUEST_MALFORMED_JSON` |
| 401 Unauthorized | 認証失敗 | `AUTH_UNAUTHORIZED`、`AUTH_INVALID_TOKEN`、`AUTH_TOKEN_EXPIRED` |
| 403 Forbidden | 権限なし | `AUTH_FORBIDDEN`、`AUTH_SCOPE_INSUFFICIENT` |
| 404 Not Found | リソースなし | `RESOURCE_NOT_FOUND`、`USER_NOT_FOUND`、`BOOK_NOT_FOUND` |
| 409 Conflict | 競合／重複 | `RESOURCE_CONFLICT`、`USER_EMAIL_ALREADY_EXISTS` |
| 422 Unprocessable Entity | 意味・ドメインの検証失敗 | `REQUEST_VALIDATION_FAILED`、`BOOK_INVALID_ISBN` |
| 429 Too Many Requests | レートリミット | `REQUEST_RATE_LIMITED` |
| 500 Internal Server Error | サーバー内部エラー | `INTERNAL_ERROR` |
| 503 Service Unavailable | 依存先の障害／メンテナンス | `INTERNAL_DEPENDENCY_UNAVAILABLE` |

### 5.3 禁止・注意事項

- `status="error"` なのに `code=200` は禁止です。
- 4xx は、クライアントの修正や再試行で解決できるエラー。5xx は、サーバー／インフラに原因があるエラーとして区別します。
- 同じ条件では、**常に同じ `error`** を返します（一貫性）。

---

## 6. 標準エラーコード一覧（v1 最小セット）

> 最初からすべてのドメインを完璧に作るのではなく、**共通 ＋ 認証・認可**から始めます。

### 6.1 Common

- `REQUEST_INVALID` (400)
- `REQUEST_MALFORMED_JSON` (400)
- `REQUEST_RATE_LIMITED` (429)
- `INTERNAL_ERROR` (500)
- `INTERNAL_DEPENDENCY_UNAVAILABLE` (503)

### 6.2 Validation (422)

- `REQUEST_VALIDATION_FAILED` (422)

### 6.3 Auth (401/403)

- `AUTH_UNAUTHORIZED` (401) — 認証が必要
- `AUTH_INVALID_CREDENTIALS` (401) — ログイン失敗（資格情報の不一致）
- `AUTH_INVALID_TOKEN` (401) — トークンの形式／署名が不正
- `AUTH_TOKEN_EXPIRED` (401) — トークンの期限切れ
- `AUTH_FORBIDDEN` (403) — 権限なし
- `AUTH_SCOPE_INSUFFICIENT` (403) — scope／role が不足

> ドメイン別（`USER_*`、`BOOK_*` など）は、実際に必要になった時点で追加します。

---

## 7. Validation Errors (422) の詳細仕様

422 では、ユーザー入力の検証結果を UI に正確に対応付けられるよう、`errors[]` による拡張を推奨します。

```json
{
  "status": "error",
  "code": 422,
  "error": "REQUEST_VALIDATION_FAILED",
  "message": "request.validation_failed",
  "errors": [
    { "field": "email", "reason": "invalid_format", "message": "user.email.invalid" },
    { "field": "password", "reason": "too_short", "message": "auth.password.too_short" }
  ],
  "metadata": { "timestamp": "2026-02-12T10:30:00Z", "traceId": "abc123" }
}
```

ルール:
- `errors[]` は **422 でのみ使用（推奨）**
- `field` は、リクエストペイロードのキーと同じ名前にする
- `reason` は、UI の分岐用の短いトークンにする（`invalid_format`、`too_short` など）
- `message` は i18n キーを推奨。未使用の場合はプレーンなメッセージも可

---

## 8. Client Handling Rule（クライアントの処理ルール）

クライアントは、次の優先順位で処理することを推奨します。

1. `status` で成功／失敗を分岐する
2. 失敗時:
   - `error` があれば、**エラーコードを見て分岐**する（例: ログイン画面へ遷移、再試行など）
   - `message` は、ユーザーへの表示やログに使う
3. `error` がない場合:
   - `message` にフォールバックする（互換性の維持）

---

## 9. Implementation Guideline（サーバー実装ガイド）

### 9.1 Single Source of Truth

- エラーコードは enum／const で管理します。
- ハンドラー／コントローラーは、例外やエラーを受け取り、**(httpStatus, errorCode, message)** に変換する Mapper を使います。

### 9.2 Logging

- エラーログには、最低限次を含めます。
  - `code`、`error`、`traceId`（または requestId）、原因のスタック／コンテキスト

---

## 10. Change Policy（変更ポリシー）

- このドキュメントの変更は API 契約の変更にあたるため、変更時には必ず次を行います。
  - サーバー／クライアントへの影響範囲の確認
  - バージョンタグ／マイグレーション方針の記録（必要に応じて）
  - テスト（サーバー／クライアント）の更新

---

## 付録 A. 例（Examples）

### A.1 Not Found

```json
{
  "status": "error",
  "code": 404,
  "error": "BOOK_NOT_FOUND",
  "message": "book.not_found",
  "metadata": { "timestamp": "2026-02-12T10:30:00Z", "traceId": "abc123" }
}
```

### A.2 Conflict

```json
{
  "status": "error",
  "code": 409,
  "error": "USER_EMAIL_ALREADY_EXISTS",
  "message": "user.email_already_exists",
  "metadata": { "timestamp": "2026-02-12T10:30:00Z", "traceId": "abc123" }
}
```
