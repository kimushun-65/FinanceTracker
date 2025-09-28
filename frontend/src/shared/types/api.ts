/**
 * API関連の型定義
 * リクエスト・レスポンス、エラーハンドリングに関する型
 */

import { ApiError, PaginationParams, PaginatedResponse } from './base';

/**
 * HTTPメソッドの型
 */
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

/**
 * APIエンドポイントの設定
 */
export type ApiEndpoint = {
  method: HttpMethod;
  path: string;
};

/**
 * クエリパラメータの型
 */
export type QueryParams = Record<string, string | number | boolean | undefined>;

/**
 * リクエストヘッダーの型
 */
export type RequestHeaders = Record<string, string>;

/**
 * APIリクエストの設定
 */
export type ApiRequestConfig = {
  params?: QueryParams;
  headers?: RequestHeaders;
  timeout?: number;
};

/**
 * APIエラーの詳細型
 */
export interface ValidationError extends ApiError {
  code: 'VALIDATION_ERROR';
  details: {
    field: string;
    message: string;
  }[];
}

export interface AuthenticationError extends ApiError {
  code: 'UNAUTHORIZED';
}

export interface AuthorizationError extends ApiError {
  code: 'FORBIDDEN';
}

export interface NotFoundError extends ApiError {
  code: 'NOT_FOUND';
}

export interface NetworkError extends ApiError {
  code: 'NETWORK_ERROR';
}

export interface InternalError extends ApiError {
  code: 'INTERNAL_ERROR';
}

/**
 * 可能なAPIエラーの型の合併型
 */
export type KnownApiError =
  | ValidationError
  | AuthenticationError
  | AuthorizationError
  | NotFoundError
  | NetworkError
  | InternalError;

/**
 * リスト取得用のクエリパラメータ
 */
export type ListQueryParams = PaginationParams & {
  /** 並び順 */
  sort?: string;
  /** 昇順・降順 */
  order?: 'asc' | 'desc';
  /** 検索クエリ */
  search?: string;
};

/**
 * 日付範囲フィルタ用のパラメータ
 */
export type DateRangeParams = {
  /** 開始日 (YYYY-MM-DD) */
  startDate?: string;
  /** 終了日 (YYYY-MM-DD) */
  endDate?: string;
};

// Re-export base types
export type {
  ApiResponse,
  ApiError,
  PaginationParams,
  PaginatedResponse,
} from './base';
