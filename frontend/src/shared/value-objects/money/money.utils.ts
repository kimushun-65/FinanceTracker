/**
 * Money値オブジェクトのユーティリティ関数
 * 金額の計算、比較、フォーマット処理
 */

import {
  Money,
  Currency,
  MoneyComparison,
  MoneyFormatOptions,
  CurrencyMismatchError,
  InvalidAmountError,
} from './money.types';
import {
  CURRENCY_SYMBOLS,
  CURRENCY_DECIMAL_PLACES,
  DEFAULT_CURRENCY,
  ZERO_MONEY,
  MAX_SAFE_AMOUNT,
  MIN_SAFE_AMOUNT,
} from './money.constants';

/**
 * Moneyオブジェクトを作成
 */
export const createMoney = (
  amount: number,
  currency: Currency = DEFAULT_CURRENCY,
): Money => {
  validateAmount(amount);
  return { amount, currency };
};

/**
 * 金額の妥当性を検証
 */
export const validateAmount = (amount: number): void => {
  if (!Number.isFinite(amount)) {
    throw new InvalidAmountError(amount);
  }
  if (amount > MAX_SAFE_AMOUNT || amount < MIN_SAFE_AMOUNT) {
    throw new InvalidAmountError(amount);
  }
};

/**
 * 通貨が同一かチェック
 */
export const assertSameCurrency = (money1: Money, money2: Money): void => {
  if (money1.currency !== money2.currency) {
    throw new CurrencyMismatchError(money1.currency, money2.currency);
  }
};

/**
 * 金額の加算
 */
export const addMoney = (money1: Money, money2: Money): Money => {
  assertSameCurrency(money1, money2);
  const result = money1.amount + money2.amount;
  validateAmount(result);
  return createMoney(result, money1.currency);
};

/**
 * 金額の減算
 */
export const subtractMoney = (money1: Money, money2: Money): Money => {
  assertSameCurrency(money1, money2);
  const result = money1.amount - money2.amount;
  validateAmount(result);
  return createMoney(result, money1.currency);
};

/**
 * 金額の乗算
 */
export const multiplyMoney = (money: Money, multiplier: number): Money => {
  validateAmount(multiplier);
  const result = money.amount * multiplier;
  validateAmount(result);
  return createMoney(Math.round(result), money.currency);
};

/**
 * 金額の除算
 */
export const divideMoney = (money: Money, divisor: number): Money => {
  if (divisor === 0) {
    throw new InvalidAmountError(divisor);
  }
  validateAmount(divisor);
  const result = money.amount / divisor;
  validateAmount(result);
  return createMoney(Math.round(result), money.currency);
};

/**
 * 金額の比較
 * @returns -1: money1 < money2, 0: money1 = money2, 1: money1 > money2
 */
export const compareMoney = (money1: Money, money2: Money): MoneyComparison => {
  assertSameCurrency(money1, money2);
  if (money1.amount < money2.amount) return -1;
  if (money1.amount > money2.amount) return 1;
  return 0;
};

/**
 * 金額が等しいかチェック
 */
export const isEqualMoney = (money1: Money, money2: Money): boolean => {
  return money1.currency === money2.currency && money1.amount === money2.amount;
};

/**
 * 金額がゼロかチェック
 */
export const isZero = (money: Money): boolean => {
  return money.amount === 0;
};

/**
 * 金額が正の値かチェック
 */
export const isPositive = (money: Money): boolean => {
  return money.amount > 0;
};

/**
 * 金額が負の値かチェック
 */
export const isNegative = (money: Money): boolean => {
  return money.amount < 0;
};

/**
 * 金額の絶対値を取得
 */
export const absoluteMoney = (money: Money): Money => {
  return createMoney(Math.abs(money.amount), money.currency);
};

/**
 * 金額の符号を反転
 */
export const negateMoney = (money: Money): Money => {
  return createMoney(-money.amount, money.currency);
};

/**
 * 複数の金額の合計を計算
 */
export const sumMoney = (moneyList: Money[]): Money => {
  if (moneyList.length === 0) {
    return { ...ZERO_MONEY };
  }

  return moneyList.reduce((sum, money) => addMoney(sum, money));
};

/**
 * 金額をフォーマットして文字列として表示
 */
export const formatMoney = (
  money: Money,
  options: MoneyFormatOptions = {},
): string => {
  const {
    showCurrencySymbol = true,
    useGrouping = true,
    minimumFractionDigits = CURRENCY_DECIMAL_PLACES[money.currency],
    maximumFractionDigits = CURRENCY_DECIMAL_PLACES[money.currency],
  } = options;

  const formatter = new Intl.NumberFormat('ja-JP', {
    useGrouping,
    minimumFractionDigits,
    maximumFractionDigits,
  });

  const formattedAmount = formatter.format(money.amount);

  if (showCurrencySymbol) {
    const symbol = CURRENCY_SYMBOLS[money.currency];
    return `${symbol}${formattedAmount}`;
  }

  return formattedAmount;
};

/**
 * 金額をコンパクト形式でフォーマット（例: ¥1.2K, ¥3.4M）
 */
export const formatMoneyCompact = (money: Money): string => {
  const formatter = new Intl.NumberFormat('ja-JP', {
    notation: 'compact',
    compactDisplay: 'short',
    minimumFractionDigits: 0,
    maximumFractionDigits: 1,
  });

  const formattedAmount = formatter.format(money.amount);
  const symbol = CURRENCY_SYMBOLS[money.currency];
  return `${symbol}${formattedAmount}`;
};

/**
 * 文字列から金額を解析
 */
export const parseMoney = (
  amountString: string,
  currency: Currency = DEFAULT_CURRENCY,
): Money => {
  // 数字以外の文字を除去
  const cleanString = amountString.replace(/[^\d.-]/g, '');
  const amount = parseFloat(cleanString);

  if (isNaN(amount)) {
    throw new InvalidAmountError(amount);
  }

  return createMoney(Math.round(amount), currency);
};
