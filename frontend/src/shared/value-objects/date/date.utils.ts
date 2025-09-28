/**
 * Date関連のユーティリティ関数
 * 日付・時刻の変換、フォーマット、計算処理
 */

import {
  ISODateString,
  DateString,
  TimeString,
  MonthString,
  DateFormatOptions,
  TimeFormatOptions,
  DateRange,
  MonthRange,
  Period,
  RelativeDate,
  DateValidationError,
  TimeValidationError,
  PeriodValidationError,
} from './date.types';
import {
  DATE_REGEX,
  DEFAULT_DATE_FORMAT,
  DEFAULT_TIME_FORMAT,
  DATE_BOUNDS,
  TIME_BOUNDS,
  DAYS_IN_MONTH,
  DAYS_IN_MONTH_LEAP,
  MILLISECONDS,
} from './date.constants';

/**
 * うるう年かどうかを判定
 */
export const isLeapYear = (year: number): boolean => {
  return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0;
};

/**
 * 日付文字列のバリデーション
 */
export const validateDateString = (value: string): void => {
  if (!DATE_REGEX.DATE.test(value)) {
    throw new DateValidationError(
      value,
      'Invalid date format. Expected YYYY-MM-DD',
    );
  }

  const [year, month, day] = value.split('-').map(Number);

  if (year < DATE_BOUNDS.MIN_YEAR || year > DATE_BOUNDS.MAX_YEAR) {
    throw new DateValidationError(
      value,
      `Year must be between ${DATE_BOUNDS.MIN_YEAR} and ${DATE_BOUNDS.MAX_YEAR}`,
    );
  }

  if (month < 1 || month > 12) {
    throw new DateValidationError(value, 'Month must be between 1 and 12');
  }

  const daysInMonth = isLeapYear(year) ? DAYS_IN_MONTH_LEAP : DAYS_IN_MONTH;
  if (day < 1 || day > daysInMonth[month - 1]) {
    throw new DateValidationError(
      value,
      `Day must be between 1 and ${daysInMonth[month - 1]} for the given month`,
    );
  }
};

/**
 * 時間文字列のバリデーション
 */
export const validateTimeString = (value: string): void => {
  if (!DATE_REGEX.TIME.test(value)) {
    throw new TimeValidationError(value, 'Invalid time format. Expected HH:MM');
  }

  const [hour, minute] = value.split(':').map(Number);

  if (hour < TIME_BOUNDS.MIN_HOUR || hour > TIME_BOUNDS.MAX_HOUR) {
    throw new TimeValidationError(
      value,
      `Hour must be between ${TIME_BOUNDS.MIN_HOUR} and ${TIME_BOUNDS.MAX_HOUR}`,
    );
  }

  if (minute < TIME_BOUNDS.MIN_MINUTE || minute > TIME_BOUNDS.MAX_MINUTE) {
    throw new TimeValidationError(
      value,
      `Minute must be between ${TIME_BOUNDS.MIN_MINUTE} and ${TIME_BOUNDS.MAX_MINUTE}`,
    );
  }
};

/**
 * 月文字列のバリデーション
 */
export const validateMonthString = (value: string): void => {
  if (!DATE_REGEX.MONTH.test(value)) {
    throw new DateValidationError(
      value,
      'Invalid month format. Expected YYYY-MM',
    );
  }

  const [year, month] = value.split('-').map(Number);

  if (year < DATE_BOUNDS.MIN_YEAR || year > DATE_BOUNDS.MAX_YEAR) {
    throw new DateValidationError(
      value,
      `Year must be between ${DATE_BOUNDS.MIN_YEAR} and ${DATE_BOUNDS.MAX_YEAR}`,
    );
  }

  if (month < 1 || month > 12) {
    throw new DateValidationError(value, 'Month must be between 1 and 12');
  }
};

/**
 * DateStringを作成
 */
export const createDateString = (value: string): DateString => {
  validateDateString(value);
  return value as DateString;
};

/**
 * TimeStringを作成
 */
export const createTimeString = (value: string): TimeString => {
  validateTimeString(value);
  return value as TimeString;
};

/**
 * MonthStringを作成
 */
export const createMonthString = (value: string): MonthString => {
  validateMonthString(value);
  return value as MonthString;
};

/**
 * ISODateStringを作成
 */
export const createISODateString = (value: string): ISODateString => {
  // ISO形式のバリデーション
  if (!DATE_REGEX.ISO_DATE.test(value)) {
    throw new DateValidationError(value, 'Invalid ISO date format');
  }
  return value as ISODateString;
};

/**
 * DateオブジェクトからDateStringを作成
 */
export const dateToDateString = (date: Date): DateString => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return createDateString(`${year}-${month}-${day}`);
};

/**
 * DateオブジェクトからTimeStringを作成
 */
export const dateToTimeString = (date: Date): TimeString => {
  const hour = String(date.getHours()).padStart(2, '0');
  const minute = String(date.getMinutes()).padStart(2, '0');
  return createTimeString(`${hour}:${minute}`);
};

/**
 * DateオブジェクトからMonthStringを作成
 */
export const dateToMonthString = (date: Date): MonthString => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  return createMonthString(`${year}-${month}`);
};

/**
 * DateオブジェクトからISODateStringを作成
 */
export const dateToISOString = (date: Date): ISODateString => {
  return createISODateString(date.toISOString());
};

/**
 * DateStringからDateオブジェクトを作成
 */
export const dateStringToDate = (dateString: DateString): Date => {
  return new Date(dateString + 'T00:00:00.000Z');
};

/**
 * ISODateStringからDateオブジェクトを作成
 */
export const isoStringToDate = (isoString: ISODateString): Date => {
  return new Date(isoString);
};

/**
 * 現在日時を取得
 */
export const now = (): Date => {
  return new Date();
};

/**
 * 今日の日付を取得
 */
export const today = (): DateString => {
  return dateToDateString(now());
};

/**
 * 今月を取得
 */
export const thisMonth = (): MonthString => {
  return dateToMonthString(now());
};

/**
 * 日付の加算
 */
export const addDays = (dateString: DateString, days: number): DateString => {
  const date = dateStringToDate(dateString);
  date.setDate(date.getDate() + days);
  return dateToDateString(date);
};

/**
 * 月の加算
 */
export const addMonths = (
  dateString: DateString,
  months: number,
): DateString => {
  const date = dateStringToDate(dateString);
  date.setMonth(date.getMonth() + months);
  return dateToDateString(date);
};

/**
 * 年の加算
 */
export const addYears = (dateString: DateString, years: number): DateString => {
  const date = dateStringToDate(dateString);
  date.setFullYear(date.getFullYear() + years);
  return dateToDateString(date);
};

/**
 * 相対的な日付の計算
 */
export const addRelativeDate = (
  baseDate: DateString,
  relative: RelativeDate,
): DateString => {
  let result = baseDate;

  if (relative.days) {
    result = addDays(result, relative.days);
  }
  if (relative.months) {
    result = addMonths(result, relative.months);
  }
  if (relative.years) {
    result = addYears(result, relative.years);
  }

  return result;
};

/**
 * 日付の差分を日数で計算
 */
export const diffInDays = (date1: DateString, date2: DateString): number => {
  const d1 = dateStringToDate(date1);
  const d2 = dateStringToDate(date2);
  const diffTime = d2.getTime() - d1.getTime();
  return Math.floor(diffTime / MILLISECONDS.DAY);
};

/**
 * 日付の比較
 * @returns -1: date1 < date2, 0: date1 = date2, 1: date1 > date2
 */
export const compareDates = (
  date1: DateString,
  date2: DateString,
): -1 | 0 | 1 => {
  if (date1 < date2) return -1;
  if (date1 > date2) return 1;
  return 0;
};

/**
 * 日付範囲のバリデーション
 */
export const validateDateRange = (range: DateRange): void => {
  if (compareDates(range.start, range.end) > 0) {
    throw new PeriodValidationError(
      'Start date must be before or equal to end date',
    );
  }
};

/**
 * 日付が範囲内にあるかチェック
 */
export const isDateInRange = (date: DateString, range: DateRange): boolean => {
  validateDateRange(range);
  return (
    compareDates(range.start, date) <= 0 && compareDates(date, range.end) <= 0
  );
};

/**
 * 日付のフォーマット
 */
export const formatDate = (
  dateString: DateString,
  options: DateFormatOptions = DEFAULT_DATE_FORMAT,
): string => {
  const date = dateStringToDate(dateString);
  return new Intl.DateTimeFormat(options.locale, options).format(date);
};

/**
 * 時間のフォーマット
 */
export const formatTime = (
  timeString: TimeString,
  options: TimeFormatOptions = DEFAULT_TIME_FORMAT,
): string => {
  // 仮の日付を作成して時間をフォーマット
  const [hour, minute] = timeString.split(':').map(Number);
  const date = new Date(2000, 0, 1, hour, minute);
  return new Intl.DateTimeFormat(options.locale, options).format(date);
};

/**
 * 月の最初の日を取得
 */
export const getFirstDayOfMonth = (monthString: MonthString): DateString => {
  return createDateString(`${monthString}-01`);
};

/**
 * 月の最後の日を取得
 */
export const getLastDayOfMonth = (monthString: MonthString): DateString => {
  const [year, month] = monthString.split('-').map(Number);
  const daysInMonth = isLeapYear(year) ? DAYS_IN_MONTH_LEAP : DAYS_IN_MONTH;
  const lastDay = daysInMonth[month - 1];
  return createDateString(
    `${year}-${String(month).padStart(2, '0')}-${String(lastDay).padStart(2, '0')}`,
  );
};

/**
 * 月の範囲を日付範囲に変換
 */
export const monthRangeToDateRange = (monthRange: MonthRange): DateRange => {
  return {
    start: getFirstDayOfMonth(monthRange.start),
    end: getLastDayOfMonth(monthRange.end),
  };
};
