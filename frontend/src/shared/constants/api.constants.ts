/**
 * API関連の定数定義
 * エンドポイント、タイムアウト、その他API設定
 */

/**
 * APIバージョン
 */
export const API_VERSION = 'v1' as const;

/**
 * APIのベースパス
 */
export const API_BASE_PATH = '/api' as const;

/**
 * タイムアウト設定（ミリ秒）
 */
export const TIMEOUT = {
  /** デフォルトのリクエストタイムアウト */
  DEFAULT: 30000, // 30秒
  /** 短いリクエスト用（認証チェックなど） */
  SHORT: 10000, // 10秒
  /** 長いリクエスト用（ファイルアップロードなど） */
  LONG: 60000, // 60秒
  /** 非常に長いリクエスト用（レポート生成など） */
  VERY_LONG: 300000, // 5分
} as const;

/**
 * HTTP ステータスコード
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const;

/**
 * APIエンドポイントのパス
 */
export const API_ENDPOINTS = {
  // 認証関連
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    PROFILE: '/auth/profile',
  },

  // ユーザー関連
  USERS: {
    BASE: '/users',
    PROFILE: '/users/profile',
    UPDATE: '/users/profile',
  },

  // アカウント関連
  ACCOUNTS: {
    BASE: '/accounts',
    BY_ID: (id: string) => `/accounts/${id}`,
  },

  // 取引関連
  TRANSACTIONS: {
    BASE: '/transactions',
    BY_ID: (id: string) => `/transactions/${id}`,
    BY_ACCOUNT: (accountId: string) => `/transactions?accountId=${accountId}`,
  },

  // カテゴリ関連
  CATEGORIES: {
    BASE: '/categories',
    BY_ID: (id: string) => `/categories/${id}`,
    MASTERS: '/categories/masters',
  },

  // 予算関連
  BUDGETS: {
    BASE: '/budgets',
    BY_ID: (id: string) => `/budgets/${id}`,
    BY_CATEGORY: (categoryId: string) => `/budgets?categoryId=${categoryId}`,
    SUGGESTIONS: '/budgets/suggestions',
  },

  // 資産関連
  ASSETS: {
    SNAPSHOTS: '/assets/snapshots',
    FORECASTS: '/assets/forecasts',
  },

  // 口座移動関連
  ACCOUNT_MOVEMENTS: {
    BASE: '/account-movements',
    BY_ID: (id: string) => `/account-movements/${id}`,
    BY_ACCOUNT: (accountId: string) =>
      `/account-movements?accountId=${accountId}`,
  },

  // 通知設定関連
  NOTIFICATION_SETTINGS: {
    BASE: '/notification-settings',
  },
} as const;

/**
 * ページネーションのデフォルト値
 */
export const PAGINATION = {
  /** デフォルトの1ページあたりの件数 */
  DEFAULT_LIMIT: 20,
  /** 最大の1ページあたりの件数 */
  MAX_LIMIT: 100,
  /** 最小の1ページあたりの件数 */
  MIN_LIMIT: 1,
  /** デフォルトのオフセット */
  DEFAULT_OFFSET: 0,
} as const;

/**
 * リクエストヘッダー
 */
export const REQUEST_HEADERS = {
  CONTENT_TYPE: 'Content-Type',
  AUTHORIZATION: 'Authorization',
  ACCEPT: 'Accept',
  USER_AGENT: 'User-Agent',
  X_REQUESTED_WITH: 'X-Requested-With',
} as const;

/**
 * コンテンツタイプ
 */
export const CONTENT_TYPES = {
  JSON: 'application/json',
  FORM_DATA: 'multipart/form-data',
  URL_ENCODED: 'application/x-www-form-urlencoded',
  TEXT: 'text/plain',
} as const;

/**
 * キャッシュ関連の設定
 */
export const CACHE = {
  /** React Queryのキャッシュ時間（ミリ秒） */
  STALE_TIME: 5 * 60 * 1000, // 5分
  /** React Queryのガベージコレクション時間（ミリ秒） */
  CACHE_TIME: 10 * 60 * 1000, // 10分
  /** リフェッチ間隔（ミリ秒） */
  REFETCH_INTERVAL: 30 * 60 * 1000, // 30分
} as const;

/**
 * リトライ設定
 */
export const RETRY = {
  /** デフォルトのリトライ回数 */
  DEFAULT_COUNT: 3,
  /** リトライ間隔（ミリ秒） */
  DELAY: 1000,
  /** 指数バックオフの倍率 */
  BACKOFF_MULTIPLIER: 2,
} as const;

/**
 * APIエラーコード
 */
export const API_ERROR_CODES = {
  // ネットワーク関連
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT_ERROR: 'TIMEOUT_ERROR',

  // 認証・認可関連
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  TOKEN_EXPIRED: 'TOKEN_EXPIRED',

  // データ関連
  NOT_FOUND: 'NOT_FOUND',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  CONFLICT: 'CONFLICT',

  // サーバー関連
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',

  // ビジネスロジック関連
  INSUFFICIENT_BALANCE: 'INSUFFICIENT_BALANCE',
  INVALID_TRANSACTION: 'INVALID_TRANSACTION',
  BUDGET_EXCEEDED: 'BUDGET_EXCEEDED',
} as const;

/**
 * 並び順の設定
 */
export const SORT_ORDER = {
  ASC: 'asc',
  DESC: 'desc',
} as const;

/**
 * 日付範囲のプリセット
 */
export const DATE_RANGE_PRESETS = {
  TODAY: 'today',
  THIS_WEEK: 'this_week',
  THIS_MONTH: 'this_month',
  THIS_QUARTER: 'this_quarter',
  THIS_YEAR: 'this_year',
  LAST_7_DAYS: 'last_7_days',
  LAST_30_DAYS: 'last_30_days',
  LAST_90_DAYS: 'last_90_days',
  LAST_YEAR: 'last_year',
} as const;
