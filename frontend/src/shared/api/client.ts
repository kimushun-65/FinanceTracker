/**
 * APIクライアント設定
 * インターセプターとエラーハンドリングを備えた集中管理されたHTTPクライアント
 */

import { env } from '../config';
import { AUTH_TOKEN_KEY, getCookie } from '../utils';

export interface ApiError {
  code: string;
  message: string;
  details?: unknown;
}

export interface ApiResponse<T = unknown> {
  data: T;
}

export interface RequestConfig extends RequestInit {
  params?: Record<string, string | number | boolean>;
  timeout?: number;
}

class ApiClient {
  private baseUrl: string;
  private defaultTimeout: number;

  constructor() {
    this.baseUrl = env.api.baseUrl;
    this.defaultTimeout = 30000; // 30秒
  }

  /**
   * クエリパラメータ付きのURLを構築
   */
  private buildUrl(
    endpoint: string,
    params?: Record<string, string | number | boolean>,
  ): string {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    return url.toString();
  }

  /**
   * 認証ヘッダーを取得
   */
  private getAuthHeaders(): HeadersInit {
    if (typeof window !== 'undefined') {
      // Cookieからトークンを取得
      const token = getCookie(AUTH_TOKEN_KEY);
      if (token) {
        return { Authorization: `Bearer ${token}` };
      }
    }
    return {};
  }

  /**
   * タイムアウト付きのfetchを作成
   */
  private fetchWithTimeout(
    url: string,
    options: RequestInit,
    timeout: number,
  ): Promise<Response> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    return fetch(url, {
      ...options,
      signal: controller.signal,
    }).finally(() => clearTimeout(timeoutId));
  }

  /**
   * APIエラーを処理
   */
  private async handleError(response: Response): Promise<never> {
    let error: ApiError;

    try {
      const data = await response.json();
      error = {
        code: data.code || 'INTERNAL_ERROR',
        message: data.error || data.message || 'An error occurred',
        details: data.details,
      };
    } catch {
      error = {
        code: 'INTERNAL_ERROR',
        message: response.statusText || 'An error occurred',
      };
    }

    // 特定のステータスコードを処理
    switch (response.status) {
      case 401:
        error.code = 'UNAUTHORIZED';
        break;
      case 403:
        error.code = 'FORBIDDEN';
        break;
      case 404:
        error.code = 'NOT_FOUND';
        break;
      case 422:
        error.code = 'VALIDATION_ERROR';
        break;
    }

    throw error;
  }

  /**
   * HTTPリクエストを実行
   */
  private async request<T>(
    method: string,
    endpoint: string,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    const { params, timeout = this.defaultTimeout, ...options } = config || {};
    const url = this.buildUrl(endpoint, params);

    const requestOptions: RequestInit = {
      ...options,
      method,
      credentials: 'include', // Cookieを含める
      headers: {
        'Content-Type': 'application/json',
        ...this.getAuthHeaders(),
        ...options.headers,
      },
    };

    try {
      const response = await this.fetchWithTimeout(
        url,
        requestOptions,
        timeout,
      );

      if (!response.ok) {
        await this.handleError(response);
      }

      // 204 No Content の場合は空のレスポンスを返す
      if (response.status === 204) {
        return { data: null } as ApiResponse<T>;
      }

      const data = await response.json();
      return { data };
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        throw {
          code: 'NETWORK_ERROR',
          message: 'Request timeout',
        } as ApiError;
      }
      throw error;
    }
  }

  /**
   * GETリクエスト
   */
  async get<T>(
    endpoint: string,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    return this.request<T>('GET', endpoint, config);
  }

  /**
   * POSTリクエスト
   */
  async post<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    return this.request<T>('POST', endpoint, {
      ...config,
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * PUTリクエスト
   */
  async put<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    return this.request<T>('PUT', endpoint, {
      ...config,
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * PATCHリクエスト
   */
  async patch<T>(
    endpoint: string,
    data?: unknown,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    return this.request<T>('PATCH', endpoint, {
      ...config,
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * DELETEリクエスト
   */
  async delete<T>(
    endpoint: string,
    config?: RequestConfig,
  ): Promise<ApiResponse<T>> {
    return this.request<T>('DELETE', endpoint, config);
  }
}

// シングルトンインスタンスをエクスポート
export const apiClient = new ApiClient();

// 便利なメソッドをエクスポート
export const api = {
  get: apiClient.get.bind(apiClient),
  post: apiClient.post.bind(apiClient),
  put: apiClient.put.bind(apiClient),
  patch: apiClient.patch.bind(apiClient),
  delete: apiClient.delete.bind(apiClient),
} as const;
