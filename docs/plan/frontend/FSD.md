# Frontend Architecture - Feature-Sliced Design (FSD)

## **概要**

Feature-Sliced Design (FSD) を Next.js App Router の制約に合わせて実践的に適応したアーキテクチャパターンです。理論的な純粋性よりも、開発効率と保守性を重視しています。

## **基本構造**

```
src/
├── app/                    # Next.js ルーティング + グローバル設定
├── page-components/        # ページのオーケストレーション
├── widgets/                # 再利用可能な複合UIブロック
├── features/               # ビジネス機能の実装
├── entities/               # ビジネスエンティティ
└── shared/                 # 共有リソース
```

## **重要な原則**

### **1. 依存関係は一方向のみ（上から下へ）**

```
app → page-components → widgets → features → entities → shared
```

### **2. 同一レイヤー間の依存は禁止**

- ❌ widget から別の widget をインポート
- ❌ feature から別の feature をインポート
- ✅ widget から feature をインポート

### **3. 各モジュールは index.ts で Public API を定義**

```tsx
// features/hr/skill/skill-get/index.ts
export { useSkills } from './lib/useSkills';
```

## **実践的な適応ポイント**

### **Next.js との統合**

- `src/app/` に routing と providers を共存させる（理論的には分離すべきだが、Next.js の制約により統合）

### **ドメイン駆動の組織化**

- widgets と features をドメイン（HR、プロジェクト等）ごとに整理
- 例：`widgets/hr/`, `features/project/`

### **細粒度の機能分割**

- features を「エンティティ/操作」単位で分割
- 例：`features/hr/skill/skill-get/`
- 命名規則：`use[Entity][Action].ts`

### **実用的なディレクトリ構造**

```
features/
└── hr/                     # ドメイン
    └── skill/              # エンティティ
        ├── skill-get/      # 操作単位
        │   ├── lib/
        │   │   └── useSkills.ts
        │   └── index.ts
        └── skill-create/
```

## **各レイヤーの設計指針**

### **App レイヤー**

#### **App に書くべきもの**

1. **ルーティング定義**: page.tsx、layout.tsx、loading.tsx、error.tsx
2. **グローバルプロバイダー**: Redux Provider、React Query Provider
3. **グローバル状態の定義**: Redux store、認証コンテキスト
4. **グローバルスタイル**: globals.css、テーマ設定

#### **App に書いてはいけないもの**

1. **ビジネスロジック**: features や entities に委譲
2. **複雑な UI コンポーネント**: widgets で実装
3. **ページ固有の状態管理**: page-components で実装

### **Page-Components レイヤー**

#### **Page-Components に書くべきもの**

1. **features の呼び出しと組み合わせ**: 複数の features を組み合わせてページを構成
2. **ページ全体の状態管理**: ページレベルのローディング、エラー処理
3. **データの取得と配布**: API データを取得し、widgets に配布
4. **ルーティングロジック**: ページ遷移、クエリパラメータの管理

#### **Page-Components に書いてはいけないもの**

1. **UI コンポーネント**: 表示は widgets に委譲
2. **細かいビジネスロジック**: features や entities に委譲
3. **グローバル状態の定義**: app レイヤーで定義

### **Widgets レイヤー**

#### **Widgets に書くべきもの**

1. **Widget内の状態管理**: フィルター条件、ソート順、開閉状態、選択状態など
2. **表示用のデータ整形**: フィルタリング、ソート、グルーピング、ページネーション
3. **複数のUIコンポーネントの組み合わせ**: Card + Badge + Button など
4. **Widget固有のカスタムフック**: `useHrListFilter`, `useProjectTable` など

#### **Widgets に書いてはいけないもの**

1. **データ取得**: API 呼び出しは features か page-components で
2. **ビジネスロジック**: 計算や変換は entities で
3. **グローバル状態の更新**: Redux の dispatch など

### **Features レイヤー**

#### **Features に書くべきもの**

1. **API 通信の実装**: GET/POST/UPDATE/DELETE のカスタムフック
2. **機能固有の状態管理**: フォーム入力、一時的なUI状態
3. **機能単位のビジネスロジック**: バリデーション、エラーハンドリング
4. **React Query/SWR の設定**: キャッシュ、リトライ、楽観的更新

#### **Features に書いてはいけないもの**

1. **複雑な UI コンポーネント**: widgets に委譲
2. **ページ全体の状態**: page-components で管理
3. **他の feature への依存**: 同一レイヤー間の依存は禁止
4. **エンティティの型定義**: entities で定義

### **Entities レイヤー**

#### **Entities に書くべきもの**

1. **ビジネスエンティティの型定義**: インターフェース、型、定数
2. **API クライアント**: エンドポイント定義、レスポンス変換
3. **ビジネスロジック**: 計算、変換、バリデーションルール
4. **クエリキーの管理**: React Query のキー定義
5. **UI コンポーネント**: entityが関わる共通UI（ステータスバッジなど）

#### **Entities に書いてはいけないもの**

1. **ページ固有のロジック**: より汎用的な実装を心がける
2. **フレームワーク固有の実装**: 可能な限り純粋な TypeScript で

### **Shared レイヤー**

#### **Shared に書くべきもの**

1. **基本的な UI コンポーネント**: Button、Input、Card など（shadcn/ui）
2. **汎用ユーティリティ**: 日付処理、フォーマット、共通バリデーション
3. **共通設定**: API クライアント設定、定数、環境変数
4. **共通フック**: useDebounce、useLocalStorage など

#### **Shared に書いてはいけないもの**

1. **ドメイン固有のロジック**: features や entities で実装
2. **ビジネスエンティティ**: entities レイヤーで定義
3. **複雑な状態管理**: 上位レイヤーで実装

## **データフローの原則**

```
page-components (features を呼び出し、データと関数を管理)
    ↓ props でデータと関数を流し込む
widgets (表示 + Widget内の状態管理)
    ↓ 受け取った関数を実行
features (GET/POST/UPDATE/DELETE などの機能を提供)
```

## **状態管理の分担**

- **app**: グローバル状態（Redux store、認証状態、テーマ設定）
- **page-components**: ページレベルの状態、features の orchestration
- **widgets**: Widget 内で完結する UI 状態（フィルター、ソート、開閉など）
- **features**: 機能固有の一時的な状態（フォーム入力、送信状態）
- **entities**: ビジネスデータの型定義（状態そのものは持たない）

## **なぜこれが効果的か**

1. **発見しやすさ**: ドメインとエンティティで整理されているため、コードの場所が予測可能
2. **独立性**: 細粒度の分割により、機能同士の結合度が低い
3. **スケーラビリティ**: チームが大きくなっても、レイヤールールにより秩序を保てる
4. **実用性**: Next.js の制約を受け入れつつ、FSD の利点を最大限活用