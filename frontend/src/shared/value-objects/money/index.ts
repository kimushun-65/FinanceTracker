/**
 * Money値オブジェクト公開API
 */

// 型定義
export type {
  Money,
  Currency,
  MoneyComparison,
  MoneyFormatOptions,
} from './money.types';

// エラークラス
export {
  MoneyError,
  CurrencyMismatchError,
  InvalidAmountError,
} from './money.types';

// 定数
export {
  CURRENCY_SYMBOLS,
  CURRENCY_NAMES_JA,
  CURRENCY_NAMES_EN,
  CURRENCY_DECIMAL_PLACES,
  DEFAULT_CURRENCY,
  ZERO_MONEY,
  MAX_SAFE_AMOUNT,
  MIN_SAFE_AMOUNT,
} from './money.constants';

// ユーティリティ関数
export {
  createMoney,
  validateAmount,
  assertSameCurrency,
  addMoney,
  subtractMoney,
  multiplyMoney,
  divideMoney,
  compareMoney,
  isEqualMoney,
  isZero,
  isPositive,
  isNegative,
  absoluteMoney,
  negateMoney,
  sumMoney,
  formatMoney,
  formatMoneyCompact,
  parseMoney,
} from './money.utils';
