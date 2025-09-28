/**
 * Email値オブジェクト公開API
 */

// 型定義
export type { Email, EmailValidationOptions, ParsedEmail } from './email.types';

// エラークラス
export { EmailValidationError } from './email.types';

// 定数
export {
  EMAIL_VALIDATION,
  EMAIL_REGEX,
  STRICT_EMAIL_REGEX,
  LOCAL_PART_REGEX,
  DOMAIN_PART_REGEX,
  DISPOSABLE_EMAIL_DOMAINS,
  FORBIDDEN_LOCAL_PATTERNS,
  DEFAULT_VALIDATION_OPTIONS,
} from './email.constants';

// ユーティリティ関数
export {
  isValidEmailFormat,
  isStrictValidEmailFormat,
  parseEmail,
  validateEmail,
  createEmail,
  normalizeEmail,
  emailToString,
  getEmailDomain,
  getEmailLocal,
  maskEmail,
  isDisposableEmail,
  areEmailsEqual,
} from './email.utils';
