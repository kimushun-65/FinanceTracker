/**
 * Email値オブジェクトのユーティリティ関数
 * メールアドレスのバリデーション、正規化、解析処理
 */

import {
  Email,
  EmailValidationError,
  EmailValidationOptions,
  ParsedEmail,
} from './email.types';
import {
  EMAIL_VALIDATION,
  EMAIL_REGEX,
  STRICT_EMAIL_REGEX,
  LOCAL_PART_REGEX,
  DOMAIN_PART_REGEX,
  DISPOSABLE_EMAIL_DOMAINS,
  FORBIDDEN_LOCAL_PATTERNS,
  DEFAULT_VALIDATION_OPTIONS,
} from './email.constants';

/**
 * 文字列がメールアドレスの形式として有効かチェック
 */
export const isValidEmailFormat = (value: string): boolean => {
  return EMAIL_REGEX.test(value);
};

/**
 * より厳密なメールアドレス形式のチェック
 */
export const isStrictValidEmailFormat = (value: string): boolean => {
  return STRICT_EMAIL_REGEX.test(value);
};

/**
 * メールアドレスを解析
 */
export const parseEmail = (email: string): ParsedEmail => {
  const atIndex = email.lastIndexOf('@');
  if (atIndex === -1) {
    throw new EmailValidationError(email, 'Missing @ symbol');
  }

  const local = email.substring(0, atIndex);
  const domain = email.substring(atIndex + 1);

  return {
    local,
    domain,
    normalized: `${local.toLowerCase()}@${domain.toLowerCase()}`,
  };
};

/**
 * ローカル部のバリデーション
 */
const validateLocalPart = (
  local: string,
  options: EmailValidationOptions,
): void => {
  if (local.length === 0) {
    throw new EmailValidationError(local, 'Local part is empty');
  }

  if (
    local.length > (options.maxLocalLength ?? EMAIL_VALIDATION.MAX_LOCAL_LENGTH)
  ) {
    throw new EmailValidationError(local, 'Local part is too long');
  }

  if (!LOCAL_PART_REGEX.test(local)) {
    throw new EmailValidationError(
      local,
      'Local part contains invalid characters',
    );
  }

  // 禁止されたパターンのチェック
  for (const pattern of FORBIDDEN_LOCAL_PATTERNS) {
    if (pattern.test(local)) {
      throw new EmailValidationError(
        local,
        'Local part contains forbidden pattern',
      );
    }
  }
};

/**
 * ドメイン部のバリデーション
 */
const validateDomainPart = (
  domain: string,
  options: EmailValidationOptions,
): void => {
  if (domain.length === 0) {
    throw new EmailValidationError(domain, 'Domain part is empty');
  }

  if (
    domain.length >
    (options.maxDomainLength ?? EMAIL_VALIDATION.MAX_DOMAIN_LENGTH)
  ) {
    throw new EmailValidationError(domain, 'Domain part is too long');
  }

  if (!DOMAIN_PART_REGEX.test(domain)) {
    throw new EmailValidationError(
      domain,
      'Domain part contains invalid characters',
    );
  }

  // ドメインラベルの長さチェック
  const labels = domain.split('.');
  for (const label of labels) {
    if (label.length > EMAIL_VALIDATION.MAX_DOMAIN_LABEL_LENGTH) {
      throw new EmailValidationError(domain, 'Domain label is too long');
    }
  }

  // 使い捨てメールアドレスのチェック
  if (
    options.checkDisposable &&
    DISPOSABLE_EMAIL_DOMAINS.has(domain.toLowerCase())
  ) {
    throw new EmailValidationError(
      domain,
      'Disposable email domain is not allowed',
    );
  }
};

/**
 * メールアドレスの包括的なバリデーション
 */
export const validateEmail = (
  value: string,
  options: EmailValidationOptions = DEFAULT_VALIDATION_OPTIONS,
): void => {
  // 長さチェック
  if (value.length > (options.maxLength ?? EMAIL_VALIDATION.MAX_LENGTH)) {
    throw new EmailValidationError(value, 'Email address is too long');
  }

  // 基本的な形式チェック
  if (!isValidEmailFormat(value)) {
    throw new EmailValidationError(value, 'Invalid email format');
  }

  // 解析してローカル部とドメイン部を個別にバリデーション
  const parsed = parseEmail(value);
  validateLocalPart(parsed.local, options);
  validateDomainPart(parsed.domain, options);
};

/**
 * 安全なEmail値オブジェクトを作成
 */
export const createEmail = (
  value: string,
  options: EmailValidationOptions = DEFAULT_VALIDATION_OPTIONS,
): Email => {
  validateEmail(value, options);
  return value as Email;
};

/**
 * メールアドレスの正規化
 * 小文字化、トリミングなど
 */
export const normalizeEmail = (value: string): string => {
  const trimmed = value.trim();
  const parsed = parseEmail(trimmed);
  return parsed.normalized;
};

/**
 * Email値オブジェクトから文字列を取得
 */
export const emailToString = (email: Email): string => {
  return email as string;
};

/**
 * メールアドレスのドメイン部を取得
 */
export const getEmailDomain = (email: Email): string => {
  const parsed = parseEmail(email);
  return parsed.domain;
};

/**
 * メールアドレスのローカル部を取得
 */
export const getEmailLocal = (email: Email): string => {
  const parsed = parseEmail(email);
  return parsed.local;
};

/**
 * メールアドレスのマスク処理
 * プライバシー保護のため一部を隠す
 */
export const maskEmail = (email: Email, maskChar: string = '*'): string => {
  const parsed = parseEmail(email);
  const local = parsed.local;
  const domain = parsed.domain;

  if (local.length <= 2) {
    return `${maskChar}@${domain}`;
  }

  const visibleStart = Math.min(2, Math.floor(local.length / 3));
  const visibleEnd = Math.min(1, Math.floor(local.length / 3));
  const maskedLength = local.length - visibleStart - visibleEnd;

  const maskedLocal =
    local.substring(0, visibleStart) +
    maskChar.repeat(maskedLength) +
    local.substring(local.length - visibleEnd);

  return `${maskedLocal}@${domain}`;
};

/**
 * 使い捨てメールアドレスかチェック
 */
export const isDisposableEmail = (email: Email): boolean => {
  const domain = getEmailDomain(email).toLowerCase();
  return DISPOSABLE_EMAIL_DOMAINS.has(domain);
};

/**
 * 複数のメールアドレスが同じかチェック（正規化後）
 */
export const areEmailsEqual = (email1: Email, email2: Email): boolean => {
  return normalizeEmail(email1) === normalizeEmail(email2);
};
