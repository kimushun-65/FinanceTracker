/**
 * Email関連の定数定義
 */

/**
 * メールアドレスのバリデーション定数
 */
export const EMAIL_VALIDATION = {
  /** メールアドレスの最大長 */
  MAX_LENGTH: 254,
  /** ローカル部の最大長 */
  MAX_LOCAL_LENGTH: 64,
  /** ドメイン部の最大長 */
  MAX_DOMAIN_LENGTH: 253,
  /** ドメインラベルの最大長 */
  MAX_DOMAIN_LABEL_LENGTH: 63,
} as const;

/**
 * 基本的なメールアドレスの正規表現
 * RFC 5322に基づく簡略版
 */
export const EMAIL_REGEX =
  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

/**
 * より厳密なメールアドレスの正規表現
 * 一般的なビジネス用途に適した形式
 */
export const STRICT_EMAIL_REGEX =
  /^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?@[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?\.[a-zA-Z]{2,}$/;

/**
 * ローカル部の正規表現
 */
export const LOCAL_PART_REGEX = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+$/;

/**
 * ドメイン部の正規表現
 */
export const DOMAIN_PART_REGEX =
  /^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

/**
 * 使い捨てメールアドレスのドメインリスト（一部）
 * 実際の運用では外部APIやより包括的なリストを使用することを推奨
 */
export const DISPOSABLE_EMAIL_DOMAINS = new Set<string>([
  '10minutemail.com',
  'guerrillamail.com',
  'mailinator.com',
  'tempmail.org',
  'yopmail.com',
  '0-mail.com',
  '10mail.org',
  '20mail.it',
  'getnada.com',
  'mohmal.com',
  'sharklasers.com',
  'temp-mail.org',
  'throwaway.email',
]);

/**
 * 禁止されたローカル部のパターン
 */
export const FORBIDDEN_LOCAL_PATTERNS = [
  /^\./, // ピリオドで始まる
  /\.$/, // ピリオドで終わる
  /\.\./, // 連続するピリオド
] as const;

/**
 * デフォルトのバリデーション設定
 */
export const DEFAULT_VALIDATION_OPTIONS = {
  maxLength: EMAIL_VALIDATION.MAX_LENGTH,
  maxLocalLength: EMAIL_VALIDATION.MAX_LOCAL_LENGTH,
  maxDomainLength: EMAIL_VALIDATION.MAX_DOMAIN_LENGTH,
  allowInternationalDomains: false,
  checkDisposable: true,
} as const;
