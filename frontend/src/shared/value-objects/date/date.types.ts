/**
 * Date関連の値オブジェクトと型定義
 * 日付・時刻の型安全な処理のための型定義
 */

/**
 * ISO 8601形式の日付文字列型
 */
export type ISODateString = string & { readonly __isoDateBrand: unique symbol };

/**
 * YYYY-MM-DD形式の日付文字列型
 */
export type DateString = string & { readonly __dateBrand: unique symbol };

/**
 * HH:MM形式の時間文字列型
 */
export type TimeString = string & { readonly __timeBrand: unique symbol };

/**
 * YYYY-MM形式の月文字列型
 */
export type MonthString = string & { readonly __monthBrand: unique symbol };

/**
 * 日付のフォーマット設定
 */
export type DateFormatOptions = {
  /** ロケール */
  locale?: string;
  /** 年の表示 */
  year?: 'numeric' | '2-digit';
  /** 月の表示 */
  month?: 'numeric' | '2-digit' | 'long' | 'short' | 'narrow';
  /** 日の表示 */
  day?: 'numeric' | '2-digit';
  /** 曜日の表示 */
  weekday?: 'long' | 'short' | 'narrow';
  /** タイムゾーンの表示 */
  timeZone?: string;
};

/**
 * 時間のフォーマット設定
 */
export type TimeFormatOptions = {
  /** ロケール */
  locale?: string;
  /** 時間の表示 */
  hour?: 'numeric' | '2-digit';
  /** 分の表示 */
  minute?: 'numeric' | '2-digit';
  /** 秒の表示 */
  second?: 'numeric' | '2-digit';
  /** 12時間制か24時間制か */
  hour12?: boolean;
  /** タイムゾーンの表示 */
  timeZone?: string;
};

/**
 * 日付の範囲
 */
export type DateRange = {
  /** 開始日 */
  start: DateString;
  /** 終了日 */
  end: DateString;
};

/**
 * 月の範囲
 */
export type MonthRange = {
  /** 開始月 */
  start: MonthString;
  /** 終了月 */
  end: MonthString;
};

/**
 * 期間の型
 */
export type Period = {
  /** 期間の単位 */
  unit: 'day' | 'week' | 'month' | 'quarter' | 'year';
  /** 期間の長さ */
  length: number;
};

/**
 * 相対的な日付を表す型
 */
export type RelativeDate = {
  /** 基準日からの相対的な日数 */
  days?: number;
  /** 基準日からの相対的な週数 */
  weeks?: number;
  /** 基準日からの相対的な月数 */
  months?: number;
  /** 基準日からの相対的な年数 */
  years?: number;
};

/**
 * 日付のバリデーションエラー
 */
export class DateValidationError extends Error {
  constructor(value: string, reason: string) {
    super(`Invalid date "${value}": ${reason}`);
    this.name = 'DateValidationError';
  }
}

/**
 * 時間のバリデーションエラー
 */
export class TimeValidationError extends Error {
  constructor(value: string, reason: string) {
    super(`Invalid time "${value}": ${reason}`);
    this.name = 'TimeValidationError';
  }
}

/**
 * 期間のバリデーションエラー
 */
export class PeriodValidationError extends Error {
  constructor(reason: string) {
    super(`Invalid period: ${reason}`);
    this.name = 'PeriodValidationError';
  }
}
