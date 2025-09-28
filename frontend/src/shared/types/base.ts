/**
 * 共通の基底型定義
 * 全エンティティで使用される基本的な型を定義
 */

/**
 * 全エンティティの基底型
 * バックエンドのBaseEntityと対応
 */
export type BaseEntity = {
  /** エンティティのUUID */
  id: string;
  /** 作成日時 (ISO 8601) */
  createdAt: string;
  /** 更新日時 (ISO 8601) */
  updatedAt: string;
};

/**
 * API レスポンスの共通型
 * 既存のapiClientと整合性を保つ
 */
export type ApiResponse<T = unknown> = {
  data: T;
};

/**
 * API エラーの型定義
 * 既存のapiClientのエラーハンドリングと整合性を保つ
 */
export type ApiError = {
  /** エラーコード */
  code: string;
  /** エラーメッセージ */
  message: string;
  /** エラーの詳細情報（オプション） */
  details?: unknown;
};

/**
 * ページネーション用の共通型
 */
export type PaginationParams = {
  /** 取得件数 */
  limit?: number;
  /** オフセット */
  offset?: number;
};

export type PaginatedResponse<T> = {
  data: T[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
};

/**
 * 一般的なID型
 */
export type ID = string;

/**
 * タイムスタンプの型
 */
export type Timestamp = string;
