/**
 * このモジュールは以下の機能を提供します：
 * - 開発/本番環境の判定
 * - 環境に応じたAPI URLの自動切り替え
 * - 環境変数の型安全なアクセス
 */

// 環境タイプの定義
// 'development'（開発環境）または 'production'（本番環境）のみを許可
export type Environment = 'development' | 'production';

// 現在の環境を取得する関数
// NODE_ENVから環境を判定し、型安全な値を返す
export const getEnvironment = (): Environment => {
  const env = process.env.NODE_ENV as string;
  return env === 'production' ? 'production' : 'development';
};

// 環境設定オブジェクト
// アプリケーション全体で使用する環境変数を一元管理
export const env = {
  // 現在の環境
  environment: getEnvironment(),
  // 開発環境かどうかのフラグ（条件分岐で使いやすい）
  isDevelopment: getEnvironment() === 'development',
  // 本番環境かどうかのフラグ（条件分岐で使いやすい）
  isProduction: getEnvironment() === 'production',

  // API設定
  api: {
    // APIのベースURL
    // 本番環境：環境変数から取得、なければデフォルトの本番URL
    // 開発環境：環境変数から取得、なければローカルホスト
    baseUrl:
      getEnvironment() === 'production'
        ? process.env.NEXT_PUBLIC_API_BASE_URL ||
          'https://api.starup-eagleai.com/api/v1'
        : process.env.NEXT_PUBLIC_API_BASE_URL ||
          'http://localhost:8080/api/v1',
  },
} as const; // as constで読み取り専用として扱う
