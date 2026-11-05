# Flutter 設計パターン

## Riverpod パターン（3.x）

- **Notifier パターンを使う**（StateNotifier は使用禁止）
- 同期的な状態: `Notifier` ＋ `NotifierProvider`
- 非同期的な状態: `AsyncNotifier` ＋ `AsyncNotifierProvider`
- 状態クラスは immutable とし、`copyWith()` メソッドを提供する

## Repository パターン

- インターフェース: `domain/repositories/` に定義する
- 実装: `data/repositories/` に実装する
- DataSource と Storage を組み合わせて、データアクセスを抽象化する

## 依存性の注入

- Riverpod の Provider を通じて依存性を注入する
- DataSource → Repository → Notifier の順に Provider を定義する

## レイヤーの役割

- **Presentation**: UI の描画、ユーザー入力の処理
- **Provider**: 状態管理、UI とビジネスロジックの接続
- **Repository**: データソースの抽象化、キャッシュ戦略
- **DataSource**: 実際の API 呼び出し、ローカル DB へのアクセス
- **Model**: データ構造の定義、シリアライズ／デシリアライズ
