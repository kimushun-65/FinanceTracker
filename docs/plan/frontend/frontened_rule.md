# フロントエンド コーディングルール

## 基本原則

### 1. アーキテクチャ

- **FSD (Feature-Sliced Design)** に準拠したディレクトリ構造を維持する
- レイヤー間の依存関係を厳守する（上位レイヤーは下位レイヤーに依存しない）
- 同一レイヤー間の依存を禁止（widgets間、features間、entities間での相互参照は不可）
- 各スライスは独立して動作可能な設計にする
- 共通で使用するコンポーネントは**shared**レイヤーに配置する

#### ディレクトリ構造とファイル配置

各スライス内では以下の構造を厳守する：

```
feature-name/
├── api/          # エンドポイント呼び出し
├── lib/          # ビジネスロジック・カスタムフック
├── model/        # 型定義・状態管理
└── ui/           # コンポーネント・見た目
```

**配置ルール**:

- **api/**: APIエンドポイントの呼び出し処理のみ
- **lib/**: ビジネスロジック、カスタムフック、ヘルパー関数
- **model/**: 型定義、インターフェース、状態管理（store/slice）
- **ui/**: UIコンポーネント、スタイリング、見た目の実装

### 2. コード品質の原則

- **DRY原則**: 同じコードを繰り返さない。共通ロジックは適切に抽出する
- **KISS原則**: シンプルで理解しやすいコードを書く。過度な抽象化を避ける
- **単一責任の原則**: 1つの関数/コンポーネントは1つの責任のみを持つ

### 3. TypeScript規約

```typescript
// ❌ 悪い例
const getData = (params: any): any => { ... }

// ✅ 良い例
const getData = (params: DataParams): DataResponse => { ... }
```

- `any`型の使用を禁止。必ず適切な型定義を行う
- 型定義は各スライスの`model/types.ts`に集約する
- インターフェースには`I`プレフィックスを付けない

### 4. 命名規則

- **変数・関数**: camelCase（例: `getUserData`）
- **定数**: UPPER_SNAKE_CASE（例: `MAX_RETRY_COUNT`）
- **型・インターフェース**: PascalCase（例: `UserProfile`）
- **ファイル名**: kebab-case（例: `user-profile.ts`）
- **コンポーネント**: PascalCase（例: `UserProfileCard.tsx`）

### 5. データ取得パターン

#### `use` hookを優先する（React 19+）

```typescript
// ✅ 推奨: useフックを使用
const userPromise = useMemo(
  () => fetch(`/api/users/${userId}`).then((res) => res.json()),
  [userId],
);
const user = use(userPromise);
```

#### 従来の`useEffect`パターン（必要な場合のみ）

```typescript
// ⚠️ 必要な場合のみ使用
const [data, setData] = useState<User | null>(null);
const [loading, setLoading] = useState(true);

useEffect(() => {
  // 複雑な副作用処理が必要な場合のみ
}, [dependency]);
```

#### `use` hookの利点

- よりシンプルなコード
- Suspenseとの自動統合
- 競合状態の回避
- 型安全性の向上

### 6. エラーハンドリング

- 全てのAPI呼び出しにはエラーハンドリングを実装する
- ユーザーに分かりやすいエラーメッセージを表示する
- `client.ts`を用いる
- Error Boundaryを適切に配置する

### 7. コンポーネント設計

- Propsは必ず型定義する
- デフォルトエクスポートは避け、名前付きエクスポートを使用する
- ビジネスロジックはカスタムフックに切り出す

#### Widgets層のコンポーネント:

- Props駆動型で設計し、外部から動作を制御可能にする
- UIロジックは`lib/`内のカスタムフックに分離する
- ページ固有のビジネスロジックは含めない

#### Features層のコンポーネント:

- フォームロジックやAPI呼び出しはカスタムフックに分離する
- UIとビジネスロジックを明確に分離する

### 8. FSD特有の命名規則

各レイヤーでのカスタムフックの命名を統一する：

- **page-components**: `use[PageName]Page.ts`（例: `useItemMasterPage`）
- **widgets**: `use[WidgetName]Widget.ts`（例: `useStoreTableWidget`）
- **features**: `use[FeatureName].ts`（例: `useStoreEdit`）
- **entities**: `[entity]Slice.ts`, `[entity]Api.ts`（例: `itemSlice.ts`, `itemApi.ts`）

## よくある違反パターンと対処法

開発時に陥りやすいアンチパターンと正しい実装方法：

### レイヤー違反

```typescript
// ❌ 悪い例: widgets間でインポート
import { ExportButton } from '@/widgets/export-button';

// ✅ 良い例: sharedに共通化
import { ExportButton } from '@/shared/ui';
```

### 上位レイヤー参照

```typescript
// ❌ 悪い例: featuresから上位レイヤー参照
import { useItemMasterPage } from '@/page-components/item-master';

// ✅ 良い例: propsで受け取る
interface Props {
  onSuccess: (data: Item) => void;
}
```

### 適切でないデータ取得パターン

```typescript
// ❌ 悪い例: 復雑なuseEffectでのデータ取得
useEffect(() => {
  const fetchData = async () => {
    // 50行以上の処理...
  };
  fetchData();
}, []);

// ✅ 良い例: useフックまたはカスタムフックを使用
const data = use(dataPromise); // または
const { data, isLoading } = useItemData(itemId);
```

## PRチェックリスト

プルリクエスト作成前に確認すべき項目：

- [ ] ESLintエラーなし (`npm run lint`)
- [ ] TypeScript型エラーなし (`npm run typecheck`)
- [ ] レイヤー依存関係の確認（同一レイヤー間の依存なし）
- [ ] `any`型を使用していない
- [ ] 適切なエラーハンドリング（try-catch、エラートースト）
- [ ] 不要な`console.log`を削除
- [ ] コンポーネントにPropsの型定義がある
- [ ] カスタムフックでロジックを分離している

## まとめ

これらのルールを守ることで、保守性が高く、チーム開発に適したコードベースを維持できます。

---

_最終更新: 2025-08-14_