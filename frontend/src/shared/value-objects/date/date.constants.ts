/**
 * Date関連の定数定義
 */

/**
 * 日付フォーマットの正規表現
 */
export const DATE_REGEX = {
  /** YYYY-MM-DD形式 */
  DATE: /^\d{4}-\d{2}-\d{2}$/,
  /** ISO 8601形式（基本） */
  ISO_DATE: /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z?$/,
  /** HH:MM形式 */
  TIME: /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/,
  /** HH:MM:SS形式 */
  TIME_WITH_SECONDS: /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/,
  /** YYYY-MM形式 */
  MONTH: /^\d{4}-\d{2}$/,
} as const;

/**
 * 日本語の曜日名
 */
export const WEEKDAY_NAMES_JA = [
  '日曜日',
  '月曜日',
  '火曜日',
  '水曜日',
  '木曜日',
  '金曜日',
  '土曜日',
] as const;

/**
 * 日本語の曜日名（短縮形）
 */
export const WEEKDAY_NAMES_SHORT_JA = [
  '日',
  '月',
  '火',
  '水',
  '木',
  '金',
  '土',
] as const;

/**
 * 日本語の月名
 */
export const MONTH_NAMES_JA = [
  '1月',
  '2月',
  '3月',
  '4月',
  '5月',
  '6月',
  '7月',
  '8月',
  '9月',
  '10月',
  '11月',
  '12月',
] as const;

/**
 * 英語の月名
 */
export const MONTH_NAMES_EN = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
] as const;

/**
 * 英語の月名（短縮形）
 */
export const MONTH_NAMES_SHORT_EN = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
] as const;

/**
 * デフォルトの日付フォーマット設定
 */
export const DEFAULT_DATE_FORMAT = {
  locale: 'ja-JP',
  year: 'numeric' as const,
  month: 'long' as const,
  day: 'numeric' as const,
  timeZone: 'Asia/Tokyo',
};

/**
 * デフォルトの時間フォーマット設定
 */
export const DEFAULT_TIME_FORMAT = {
  locale: 'ja-JP',
  hour: '2-digit' as const,
  minute: '2-digit' as const,
  hour12: false,
  timeZone: 'Asia/Tokyo',
};

/**
 * 日付の境界値
 */
export const DATE_BOUNDS = {
  /** 最小年 */
  MIN_YEAR: 1900,
  /** 最大年 */
  MAX_YEAR: 2100,
} as const;

/**
 * 時間の境界値
 */
export const TIME_BOUNDS = {
  /** 最小時間 */
  MIN_HOUR: 0,
  /** 最大時間 */
  MAX_HOUR: 23,
  /** 最小分 */
  MIN_MINUTE: 0,
  /** 最大分 */
  MAX_MINUTE: 59,
  /** 最小秒 */
  MIN_SECOND: 0,
  /** 最大秒 */
  MAX_SECOND: 59,
} as const;

/**
 * 月の日数（平年）
 */
export const DAYS_IN_MONTH = [
  31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
] as const;

/**
 * 月の日数（うるう年）
 */
export const DAYS_IN_MONTH_LEAP = [
  31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
] as const;

/**
 * ミリ秒の換算値
 */
export const MILLISECONDS = {
  SECOND: 1000,
  MINUTE: 60 * 1000,
  HOUR: 60 * 60 * 1000,
  DAY: 24 * 60 * 60 * 1000,
} as const;
