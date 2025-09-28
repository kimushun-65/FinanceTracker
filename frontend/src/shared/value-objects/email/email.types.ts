/**
 * Email値オブジェクトの型定義
 * 型安全なメールアドレス処理のための branded type
 */

/**
 * 型安全なメールアドレス型
 * branded typeを使用してプリミティブ型との区別を明確にする
 */
export type Email = string & { readonly __brand: unique symbol };

/**
 * メールアドレスのバリデーションエラー
 */
export class EmailValidationError extends Error {
  constructor(email: string, reason: string) {
    super(`Invalid email address "${email}": ${reason}`);
    this.name = 'EmailValidationError';
  }
}

/**
 * メールアドレスのバリデーション設定
 */
export type EmailValidationOptions = {
  /** 最大長 */
  maxLength?: number;
  /** ローカル部の最大長 */
  maxLocalLength?: number;
  /** ドメイン部の最大長 */
  maxDomainLength?: number;
  /** 国際化ドメイン名を許可するか */
  allowInternationalDomains?: boolean;
  /** 使い捨てメールアドレスをチェックするか */
  checkDisposable?: boolean;
};

/**
 * メールアドレスの解析結果
 */
export type ParsedEmail = {
  /** ローカル部（@より前） */
  local: string;
  /** ドメイン部（@より後） */
  domain: string;
  /** 正規化されたメールアドレス */
  normalized: string;
};
