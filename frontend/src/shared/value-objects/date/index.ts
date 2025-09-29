/**
 * Date値オブジェクト公開API
 */

// 型定義
export type {
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
} from './date.types';

// エラークラス
export {
  DateValidationError,
  TimeValidationError,
  PeriodValidationError,
} from './date.types';

// 定数
export {
  DATE_REGEX,
  WEEKDAY_NAMES_JA,
  WEEKDAY_NAMES_SHORT_JA,
  MONTH_NAMES_JA,
  MONTH_NAMES_EN,
  MONTH_NAMES_SHORT_EN,
  DEFAULT_DATE_FORMAT,
  DEFAULT_TIME_FORMAT,
  DATE_BOUNDS,
  TIME_BOUNDS,
  DAYS_IN_MONTH,
  DAYS_IN_MONTH_LEAP,
  MILLISECONDS,
} from './date.constants';

// ユーティリティ関数
export {
  isLeapYear,
  validateDateString,
  validateTimeString,
  validateMonthString,
  createDateString,
  createTimeString,
  createMonthString,
  createISODateString,
  dateToDateString,
  dateToTimeString,
  dateToMonthString,
  dateToISOString,
  dateStringToDate,
  isoStringToDate,
  now,
  today,
  thisMonth,
  addDays,
  addMonths,
  addYears,
  addRelativeDate,
  diffInDays,
  compareDates,
  validateDateRange,
  isDateInRange,
  formatDate,
  formatTime,
  getFirstDayOfMonth,
  getLastDayOfMonth,
  monthRangeToDateRange,
  getPeriodDateRange,
  formatDateShort,
} from './date.utils';
