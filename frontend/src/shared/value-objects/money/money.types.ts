/**
 * Money値オブジェクトの型定義
 * バックエンドのMoney構造体と整合性を保つ
 */

/**
 * サポートされる通貨
 * バックエンドのCurrency列挙型と対応
 */
export type Currency = 'JPY' | 'USD' | 'EUR' | 'GBP' | 'CNY' | 'KRW';

/**
 * 金額を表す値オブジェクト
 * バックエンドの int64 ベースの金額処理に対応
 */
export type Money = {
  /** 金額（整数値、日本円の場合は円単位） */
  amount: number;
  /** 通貨コード */
  currency: Currency;
};

/**
 * 金額の比較結果
 */
export type MoneyComparison = -1 | 0 | 1;

/**
 * 金額のフォーマット設定
 */
export type MoneyFormatOptions = {
  /** 通貨記号を表示するか */
  showCurrencySymbol?: boolean;
  /** 桁区切り文字を使用するか */
  useGrouping?: boolean;
  /** 小数点以下の桁数（通貨によって自動設定） */
  minimumFractionDigits?: number;
  /** 最大小数点以下桁数 */
  maximumFractionDigits?: number;
};

/**
 * 金額の計算エラー
 */
export class MoneyError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'MoneyError';
  }
}

/**
 * 通貨不一致エラー
 */
export class CurrencyMismatchError extends MoneyError {
  constructor(currency1: Currency, currency2: Currency) {
    super(`Currency mismatch: ${currency1} and ${currency2}`);
    this.name = 'CurrencyMismatchError';
  }
}

/**
 * 無効な金額エラー
 */
export class InvalidAmountError extends MoneyError {
  constructor(amount: number) {
    super(`Invalid amount: ${amount}`);
    this.name = 'InvalidAmountError';
  }
}
