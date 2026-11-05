# Flutter／Dart コーディング規約

## 命名規則

**ファイル名（File Names）**:
- すべてのファイル名は `snake_case` を使う
- クラス名とファイル名を一致させる（例: `LoginScreen` → `login_screen.dart`）
- **例**: `auth_provider.dart`、`login_screen.dart`、`user_model.dart`

**クラス名（Class Names）**:
- `PascalCase` を使う
- 役割が明確にわかる接尾辞を使う。
  - Screen: `LoginScreen`、`HomeScreen`
  - Widget: `BookCard`、`LoadingIndicator`
  - Provider／Notifier: `AuthNotifier`、`LibraryNotifier`
  - Repository: `AuthRepository`、`BookRepository`
  - Model: `User`、`Book`、`LoginRequest`
  - DataSource: `AuthRemoteDataSource`
- **例**: `AuthNotifier`、`LibraryScreen`、`BookRepository`

**変数名・関数名（Variables & Functions）**:
- `lowerCamelCase` を使う
- 関数は動詞で始める
- bool の変数／関数には `is`、`has`、`can` の接頭辞を使う
- **例**: `userName`、`isLoading`、`fetchBooks()`、`hasToken()`

**定数（Constants）**:
- トップレベル定数: `lowerCamelCase`（Dart 公式の推奨）
- クラス内の定数: `lowerCamelCase`
- **例**: `const defaultTimeout = 30;`、`static const baseUrl = '...';`

**Provider 名**:
- 接尾辞 `Provider` を使う
- 「機能 ＋ 役割 ＋ Provider」の形式にする
- **例**: `authStateProvider`、`libraryBooksProvider`、`dioProvider`

## エラー処理

- すべての API 呼び出しを try-catch で囲む
- カスタム例外クラスを使う（`ApiException`）
- エラー状態を State に含め、UI 側で処理する
- ユーザーにわかりやすいエラーメッセージを表示する

## ウィジェットの書き方

- ウィジェットは単一責任の原則を守る
- 大きなウィジェットは小さなウィジェットに分割する
- `const` コンストラクタを積極的に活用する
- ビジネスロジックはウィジェットから分離する（Provider に移す）

## import の整理

`import` 文は、次の順序で並べます。
1. Dart SDK（`dart:`）
2. Flutter SDK（`package:flutter/`）
3. 外部パッケージ（`package:`）
4. プロジェクト内部のパッケージ（`package:<app>/`）
5. 相対パス（`../`、`./`）

## ロギング

- `debugPrint()` を使う（リリースビルドでは自動的に除去される）
- 機密情報はログに出力しない（トークン、パスワードなど）

## パフォーマンス

- `const` ウィジェットを積極的に活用する
- 不要なリビルドを防ぐ（`ref.watch` と `ref.read` を適切に使い分ける）
- 画像キャッシュを適用する

## テスト

**テストのルール**:
- モックライブラリ: `mocktail` を使う
- テストファイル名: `<対象>_test.dart`（例: `auth_provider_test.dart`）
- グルーピング: `group()` で関連するテストをまとめる
- パターン: Arrange → Act → Assert

**テスト対象**:
- Provider／Notifier: 状態の変化とメソッドの動作
- Repository: データ変換とエラー処理
- 主要なウィジェット: ユーザーインタラクションと状態表示

## コード品質チェックリスト

生成されたコードが、次を満たしているか確認する。
- [ ] ディレクトリ構造に従っている（Feature-First）
- [ ] 命名規則に一致している（ファイルは snake_case、クラスは PascalCase）
- [ ] Riverpod 3.x のパターンを使っている（Notifier、AsyncNotifier）
- [ ] エラー処理が完了している
- [ ] テストコードが存在する
- [ ] import がグループ化され、整列されている
- [ ] const コンストラクタを活用している
- [ ] 単一責任の原則を守っている
- [ ] レイヤー分離を守っている（Presentation → Domain → Data）
