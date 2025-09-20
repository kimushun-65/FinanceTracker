# ドメインエンティティ詳細設計

## 1. ユーザー管理コンテキスト（User Management Context）

### 集約: User（ユーザー集約）

#### 集約ルート: User
```
エンティティ: User (users)
責務: ユーザーアカウントの管理とAuth0統合

属性:
- id: UUID - ユーザー一意識別子
- auth0UserId: string - Auth0ユーザーID
- email: Email - メールアドレス（値オブジェクト）
- name: string - ユーザー名
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

値オブジェクト:
- Email: メールアドレス

ビジネスルール:
- メールアドレスはシステム全体で一意
- Auth0ユーザーIDはシステム全体で一意
- メールアドレスは有効な形式である必要がある
- 認証はAuth0を通じて行われる
```

#### 値オブジェクト: Email
```
値オブジェクト: Email
責務: 有効なメールアドレスを表現

属性:
- value: string - メールアドレス

振る舞い:
- validate(): boolean - メールアドレスの形式検証
- getDomain(): string - ドメイン部分を取得
- equals(other: Email): boolean - メールアドレスの同一性判定
- toString(): string - 文字列表現

ビジネスルール:
- RFC 5322に準拠した形式
- 最大255文字
- 一度設定された値は不変
```

#### 値オブジェクト: HexColor
```
値オブジェクト: HexColor
責務: HEX形式の色を表現

属性:
- value: string - HEX形式の色（#RRGGBB）

振る舞い:
- validate(): boolean - HEX形式の検証
- toRgb(): {r: number, g: number, b: number} - RGB値に変換
- equals(other: HexColor): boolean - 色の同一性判定
- toString(): string - 文字列表現

ビジネスルール:
- #で始まる6桁の16進数（例: #3B82F6）
- 大文字小文字は区別しない
- 一度設定された値は不変
```

---

## 2. 金融口座管理コンテキスト（Account Management Context）

### 集約: Account（口座集約）

#### 集約ルート: Account
```
エンティティ: Account (accounts)
責務: 金融口座の管理と残高追跡

属性:
- id: UUID - 口座一意識別子
- userId: UUID - ユーザーID（FK）
- name: AccountName - 口座名（値オブジェクト）
- type: AccountType - 口座種別（値オブジェクト）
- balance: Balance - 残高情報（値オブジェクト）
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

含まれるエンティティ:
- AccountMovement: 口座の入出金履歴

値オブジェクト:
- AccountName: 口座名称
- AccountType: 口座種別
- Balance: 残高（初期残高・現在残高）

ビジネスルール:
- 口座名は必須
- 初期残高は0以上
- 残高の変更は必ずAccountMovementを通じて行う
```

#### 値オブジェクト: AccountType
```
値オブジェクト: AccountType
責務: 口座の種別を表現

属性:
- value: enum (checking, savings, investment, credit_card, loan)

振る舞い:
- isAsset(): boolean - 資産口座か判定
- isLiability(): boolean - 負債口座か判定
- canGoNegative(): boolean - マイナス残高を許可するか
- toString(): string - 口座種別の文字列表現
- getDisplayName(): string - 表示用名称（普通預金、投資信託等）

ビジネスルール:
- checking（普通預金）、savings（定期預金）、investment（投資）は資産
- credit_card（クレジットカード）、loan（ローン）は負債
- 一度設定された種別は不変
```

#### 値オブジェクト: Balance
```
値オブジェクト: Balance
責務: 口座の残高情報を表現

属性:
- initialBalance: Money - 初期残高
- currentBalance: Money - 現在残高

振る舞い:
- add(amount: Money): Balance - 入金
- subtract(amount: Money): Balance - 出金
- getDifference(): Money - 初期残高との差額
- isPositive(): boolean - プラス残高か判定
- canWithdraw(amount: Money): boolean - 出金可能か判定

ビジネスルール:
- 金額は整数で管理（最小単位：円）
- 資産口座の残高は0以上
- 変更時は新しいインスタンスを返す（不変性）
```

#### エンティティ: AccountMovement
```
エンティティ: AccountMovement (account_movements)
責務: 口座の入出金履歴管理

属性:
- id: UUID - 移動一意識別子
- userId: UUID - ユーザーID（FK）
- accountId: UUID - 口座ID（FK）
- amount: Money - 金額（値オブジェクト）
- occurredAt: timestamp - 発生日時
- note: string - メモ

値オブジェクト:
- Money: 金額

ビジネスルール:
- 金額は0以外の値
- プラスは入金、マイナスは出金
- 一度記録された履歴は変更不可（イミュータブル）
```

---

## 3. 取引管理コンテキスト（Transaction Management Context）

### 集約: Transaction（取引集約）

#### 集約ルート: Transaction
```
エンティティ: Transaction (transactions)
責務: 収支取引の記録と管理

属性:
- id: UUID - 取引一意識別子
- userId: UUID - ユーザーID（FK）
- categoryId: UUID - カテゴリID（FK）
- type: TransactionType - 取引種別（値オブジェクト）
- amount: Money - 金額（値オブジェクト）
- transactionDate: date - 取引日
- description: string - 説明・メモ
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

値オブジェクト:
- TransactionType: 取引種別（収入/支出）
- Money: 金額

ビジネスルール:
- 金額は正の値のみ
- 取引日は必須
- カテゴリは必須
- 過去1年以上前の取引は編集不可
```

#### 値オブジェクト: TransactionType
```
値オブジェクト: TransactionType
責務: 取引の種別を表現

属性:
- value: enum (income, expense)

振る舞い:
- isIncome(): boolean - 収入か判定
- isExpense(): boolean - 支出か判定
- getSign(): number - 符号を取得（収入:+1、支出:-1）
- toString(): string - 種別の文字列表現
- getDisplayName(): string - 表示用名称（収入、支出）

ビジネスルール:
- income（収入）またはexpense（支出）のみ
- 一度設定された種別は不変
```

#### 値オブジェクト: Money
```
値オブジェクト: Money
責務: 金額を表現

属性:
- amount: integer - 金額（円）
- currency: string - 通貨（デフォルト: JPY）

振る舞い:
- add(other: Money): Money - 加算
- subtract(other: Money): Money - 減算
- multiply(factor: number): Money - 乗算
- isPositive(): boolean - 正の値か判定
- isZero(): boolean - ゼロか判定
- isNegative(): boolean - 負の値か判定
- format(): string - フォーマット済み文字列（¥1,000）
- equals(other: Money): boolean - 金額の同一性判定

ビジネスルール:
- 整数で管理（小数点以下なし）
- 同一通貨のみ演算可能
- 不変オブジェクト（演算結果は新しいインスタンス）
```

---

## 4. カテゴリ管理コンテキスト（Category Management Context）

### 集約: CategoryMaster（カテゴリマスター集約）

#### 集約ルート: CategoryMaster
```
エンティティ: CategoryMaster (category_master)
責務: システム標準カテゴリの管理

属性:
- id: UUID - カテゴリマスター一意識別子
- parentId: UUID - 親カテゴリID（FK、自己参照）
- name: CategoryName - カテゴリ名（値オブジェクト）
- sortOrder: integer - 並び順
- isActive: boolean - 有効フラグ
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

値オブジェクト:
- CategoryName: カテゴリ名称

ビジネスルール:
- カテゴリは2階層まで（親・子）
- 並び順は同一親内で一意
- 無効化されたカテゴリは新規選択不可
```

### 集約: Category（ユーザーカテゴリ集約）

#### 集約ルート: Category
```
エンティティ: Category (categories)
責務: ユーザー固有のカテゴリ管理

属性:
- id: UUID - カテゴリ一意識別子
- userId: UUID - ユーザーID（FK）
- masterId: UUID - マスターカテゴリID（FK、オプション）
- parentId: UUID - 親カテゴリID（FK、自己参照）
- name: CategoryName - カテゴリ名（値オブジェクト）
- icon: string - カテゴリアイコン（絵文字またはアイコンID）
- color: HexColor - カテゴリ色（値オブジェクト）
- isCustom: boolean - カスタムカテゴリフラグ
- sortOrder: integer - 並び順
- isActive: boolean - 有効フラグ
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

値オブジェクト:
- CategoryName: カテゴリ名称
- HexColor: HEX形式の色

ビジネスルール:
- マスターカテゴリはユーザーごとにコピーされる
- カスタムカテゴリは自由に追加可能
- 使用中のカテゴリは削除不可（非活性化のみ）
```

#### 値オブジェクト: CategoryName
```
値オブジェクト: CategoryName
責務: カテゴリの名称を表現

属性:
- value: string - カテゴリ名

振る舞い:
- validate(): boolean - カテゴリ名の検証
- equals(other: CategoryName): boolean - カテゴリ名の同一性判定
- toString(): string - 文字列表現

ビジネスルール:
- 1文字以上50文字以下
- 空白のみは不可
- 一度設定された値は不変
```

---

## 5. 予算管理コンテキスト（Budget Management Context）

### 集約: Budget（予算集約）

#### 集約ルート: Budget
```
エンティティ: Budget (budgets)
責務: 月次予算の設定と管理

属性:
- id: UUID - 予算一意識別子
- userId: UUID - ユーザーID（FK）
- categoryId: UUID - カテゴリID（FK）
- amount: Money - 予算額（値オブジェクト）
- month: string - 対象年月（YYYY-MM形式）
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

値オブジェクト:
- Money: 金額

ビジネスルール:
- ユーザー、カテゴリ、年月の組み合わせは一意
- 予算額は0より大きい値
- 過去の予算もコピー元として参照可能
- monthはYYYY-MM形式（例: 2024-02）
```

### 集約: BudgetSuggestion（予算提案集約）

#### 集約ルート: BudgetSuggestion
```
エンティティ: BudgetSuggestion (budget_suggestions)
責務: AIによる予算提案の管理

属性:
- id: UUID - 提案一意識別子
- userId: UUID - ユーザーID（FK）
- month: string - 対象年月（YYYY-MM形式）
- categoryId: UUID - カテゴリID（FK）
- suggestedAmount: Money - 提案金額（値オブジェクト）
- currentBudget: Money - 現在の予算額（値オブジェクト）
- lastMonthActual: Money - 前月実績（値オブジェクト）
- threeMonthAverage: Money - 3ヶ月平均（値オブジェクト）
- status: SuggestionStatus - ステータス（値オブジェクト）
- reason: string - 提案理由
- confidence: float - 信頼度（0-1）
- createdAt: timestamp - 作成日時

値オブジェクト:
- SuggestionStatus: 提案ステータス

ビジネスルール:
- 月ごとに最新の提案のみ有効
- 信頼度は0以上1以下
- 採用/却下後は変更不可
```

#### 値オブジェクト: SuggestionStatus
```
値オブジェクト: SuggestionStatus
責務: 予算提案の状態を表現

属性:
- value: enum (pending, accepted, rejected)

振る舞い:
- canTransitionTo(newStatus: SuggestionStatus): boolean - 状態遷移の可否判定
- isFinal(): boolean - 最終状態か判定
- toString(): string - ステータスの文字列表現
- getDisplayName(): string - 表示用名称（検討中、採用、却下）

ビジネスルール:
- pendingからaccepted/rejectedへの遷移のみ
- accepted/rejectedは最終状態
- 一度設定されたステータスは不変
```

---

## 6. 資産管理コンテキスト（Asset Management Context）

### 集約: AssetSnapshot（資産スナップショット集約）

#### 集約ルート: AssetSnapshot
```
エンティティ: AssetSnapshot (asset_snapshots)
責務: 特定時点の総資産記録

属性:
- id: UUID - スナップショット一意識別子
- userId: UUID - ユーザーID（FK）
- snapshotDate: date - スナップショット日付
- totalAssets: Money - 総資産額（値オブジェクト）
- accountBreakdown: JSON - 口座別内訳（JSONBフィールド）
- createdAt: timestamp - 作成日時

値オブジェクト:
- Money: 金額

ビジネスルール:
- ユーザーと日付の組み合わせは一意
- 日次で自動作成
- 過去のスナップショットは変更不可
```

### 集約: AssetForecast（資産予測集約）

#### 集約ルート: AssetForecast
```
エンティティ: AssetForecast (asset_forecasts)
責務: 将来の資産予測管理

属性:
- id: UUID - 予測一意識別子
- userId: UUID - ユーザーID（FK）
- horizonMonths: integer - 予測期間（月数）
- forecastDate: date - 予測基準日
- predictedAssets: Money - 予測資産額（値オブジェクト）
- method: ForecastMethod - 予測手法（値オブジェクト）
- assumptions: Assumptions - 前提条件（値オブジェクト）
- confidence: float - 信頼区間
- createdAt: timestamp - 作成日時

値オブジェクト:
- ForecastMethod: 予測手法
- Assumptions: 予測前提条件

ビジネスルール:
- 予測期間は1-60ヶ月
- 前提条件は必須
- 最新の予測が有効
```

#### 値オブジェクト: ForecastMethod
```
値オブジェクト: ForecastMethod
責務: 資産予測の手法を表現

属性:
- value: enum (linear_regression, moving_average, monte_carlo, ml_based)

振る舞い:
- getDescription(): string - 手法の説明
- getRequiredDataMonths(): integer - 必要なデータ期間
- toString(): string - 手法名の文字列表現

ビジネスルール:
- 各手法には最低限必要なデータ期間がある
- 一度設定された手法は不変
```

#### 値オブジェクト: Assumptions
```
値オブジェクト: Assumptions
責務: 予測の前提条件を表現

属性:
- monthlyIncome: Money - 月収予測
- monthlyExpense: Money - 月支出予測
- savingsRate: float - 貯蓄率
- investmentReturn: float - 投資リターン率
- inflationRate: float - インフレ率
- majorExpenses: MajorExpense[] - 大型支出予定

振る舞い:
- toJson(): string - JSON形式で出力
- validate(): boolean - 前提条件の妥当性検証
- calculateMonthlySavings(): Money - 月間貯蓄額を計算

ビジネスルール:
- 各率は-1から1の範囲
- 月収・支出は0以上
- 一度設定された前提条件は不変
```

---

## 7. 通知管理コンテキスト（Notification Management Context）

### 集約: NotificationSettings（通知設定集約）

#### 集約ルート: NotificationSettings
```
エンティティ: NotificationSettings (notification_settings)
責務: ユーザーの通知設定管理

属性:
- id: UUID - 通知設定一意識別子
- userId: UUID - ユーザーID（FK）
- monthlyReportEnabled: boolean - 月次レポート有効フラグ
- monthlyReportSendDay: integer - 送信日（1-31）
- monthlyReportSendTime: Time - 送信時刻
- monthlyReportEmail: Email - 送信先メールアドレス
- budgetExceededAlertEnabled: boolean - 予算超過アラート有効フラグ
- budgetExceededAlertEmail: Email - アラート送信先メールアドレス
- createdAt: timestamp - 作成日時
- updatedAt: timestamp - 更新日時

ビジネスルール:
- ユーザーごとに1つのみ存在
- 送信日は1-31の範囲
- 送信時刻はHH:mm形式
```

---

## ドメイン間の連携

### イベント駆動による連携
各集約間は以下のドメインイベントで連携：

1. **UserCreated** → 各コンテキストでユーザー参照を作成、通知設定を初期化
2. **TransactionCreated** → 予算消化計算、口座残高更新
3. **BudgetUpdated** → 予算達成率の再計算
4. **AccountBalanceChanged** → 資産スナップショットの更新
5. **AssetSnapshotCreated** → 資産予測の再計算トリガー
6. **BudgetExceeded** → 予算超過通知のトリガー

### 整合性の保証
- 各集約内：強整合性（トランザクション保証）
- 集約間：結果整合性（イベントによる同期）
- 参照整合性：UUIDによる疎結合

### クエリモデル
読み取り専用の以下のビューを提供：
- **MonthlyFinancialSummary**: 月次収支サマリー
- **BudgetPerformance**: 予算対実績
- **AssetTrend**: 資産推移
- **CategorySpending**: カテゴリ別支出分析