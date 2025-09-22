import { AuthConfig } from '../model';

export const getAuthConfig = (): AuthConfig => ({
  domain:
    process.env.NEXT_PUBLIC_AUTH0_DOMAIN || 'dev-kimushun3765.jp.auth0.com',
  clientId:
    process.env.NEXT_PUBLIC_AUTH0_CLIENT_ID ||
    'N19XPSt4tzz3Fa7reXz26yRgLfoeuoJ9',
  redirectUri:
    process.env.NEXT_PUBLIC_AUTH0_REDIRECT_URI ||
    'http://localhost:3000/callback',
  audience:
    process.env.NEXT_PUBLIC_AUTH0_AUDIENCE ||
    'https://api.financetracker.local',
  scope: 'openid profile email',
});
