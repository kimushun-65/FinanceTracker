/**
 * Money関連の定数定義
 */

import { Currency } from './money.types';

/**
 * 通貨記号のマッピング
 */
export const CURRENCY_SYMBOLS: Record<Currency, string> = {
  JPY: '¥',
  USD: '$',
  EUR: '€',
  GBP: '£',
  CNY: '¥',
  KRW: '₩',
} as const;

/**
 * 通貨名のマッピング（日本語）
 */
export const CURRENCY_NAMES_JA: Record<Currency, string> = {
  JPY: '日本円',
  USD: '米ドル',
  EUR: 'ユーロ',
  GBP: '英ポンド',
  CNY: '中国元',
  KRW: '韓国ウォン',
} as const;

/**
 * 通貨名のマッピング（英語）
 */
export const CURRENCY_NAMES_EN: Record<Currency, string> = {
  JPY: 'Japanese Yen',
  USD: 'US Dollar',
  EUR: 'Euro',
  GBP: 'British Pound',
  CNY: 'Chinese Yuan',
  KRW: 'South Korean Won',
} as const;

/**
 * 通貨ごとの小数点以下桁数
 */
export const CURRENCY_DECIMAL_PLACES: Record<Currency, number> = {
  JPY: 0, // 日本円は小数点なし
  USD: 2,
  EUR: 2,
  GBP: 2,
  CNY: 2,
  KRW: 0, // 韓国ウォンは小数点なし
} as const;

/**
 * デフォルト通貨
 */
export const DEFAULT_CURRENCY: Currency = 'JPY';

/**
 * ゼロ金額の定数
 */
export const ZERO_MONEY = {
  amount: 0,
  currency: DEFAULT_CURRENCY,
} as const;

/**
 * 金額の最大値（安全な整数の最大値）
 */
export const MAX_SAFE_AMOUNT = Number.MAX_SAFE_INTEGER;

/**
 * 金額の最小値
 */
export const MIN_SAFE_AMOUNT = Number.MIN_SAFE_INTEGER;
