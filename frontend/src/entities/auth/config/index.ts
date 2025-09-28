import { AuthConfig } from '../model';

export const getAuthConfig = (): AuthConfig => {
  // 静的エクスポートの場合、ビルド時の環境変数ではなく
  // 実行時のURLを使用する必要がある
  const getRedirectUri = () => {
    if (typeof window !== 'undefined') {
      return `${window.location.origin}/callback`;
    }
    // サーバーサイドレンダリング中はデフォルト値を使用
    return (
      process.env.NEXT_PUBLIC_AUTH0_REDIRECT_URI ||
      'http://localhost:3000/callback'
    );
  };

  return {
    domain:
      process.env.NEXT_PUBLIC_AUTH0_DOMAIN || 'dev-kimushun3765.jp.auth0.com',
    clientId:
      process.env.NEXT_PUBLIC_AUTH0_CLIENT_ID ||
      'N19XPSt4tzz3Fa7reXz26yRgLfoeuoJ9',
    redirectUri: getRedirectUri(),
    audience:
      process.env.NEXT_PUBLIC_AUTH0_AUDIENCE ||
      'https://api.financetracker.local',
    scope: 'openid profile email',
  };
};
