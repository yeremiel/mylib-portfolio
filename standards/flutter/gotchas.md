# Flutter／Dart の注意点

## `copyWith()` で nullable フィールドを `null` に戻せない

Dart の `copyWith()` パターンでは、`field: field ?? this.field` は `null` を「変更なし」として扱います。
nullable なフィールドを `null` に初期化したい場合は、`copyWith()` ではなく `const State()` で新しいインスタンスを生成します。

```dart
// ❌ meta が null にリセットされない
state = state.copyWith(meta: null, books: []);

// ✅ 新しいインスタンスを生成して解決
state = const BookSearchState();
```

## 生成された API モデルとプロジェクトのモデルで import が衝突する

Swagger の生成コードには、`User` など、プロジェクトのモデルと同じ名前のクラスが存在することがあります。
衝突した場合は、生成コードの export ファイルで `hide` を使って除外します。

```dart
// ❌ User クラスが衝突する
import 'package:mylib_app/core/api/generated/export.dart';

// ✅ hide で衝突を解決
import 'package:mylib_app/core/api/generated/export.dart' hide User;
```

## Riverpod の Notifier 内で `state` の名前が衝突する

Notifier クラスの内部では、`state` は Riverpod の状態プロパティです。
OAuth などの外部パラメータに `state` という名前がある場合は衝突するため、別の名前を使います。

```dart
// ❌ Riverpod の state プロパティと衝突する
final state = generateOAuthState();

// ✅ 明示的な名前で区別する
final oauthState = generateOAuthState();
```

## `.g.dart` ファイルは git に含めない

`build_runner` が生成する `.g.dart` ファイルは、`.gitignore` に登録されています。
`git add` 時には自動的に除外されますが、手動で追加しないよう注意します。
