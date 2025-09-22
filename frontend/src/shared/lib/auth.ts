'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useEffect } from 'react';
import { tokenManager } from './token';

/**
 * Auth0認証とCookie管理を統合するカスタムフック
 */
export const useAuthWithCookie = () => {
  const { isAuthenticated, getAccessTokenSilently, user, isLoading } = useAuth0();

  // ログイン後にトークンをCookieに保存
  useEffect(() => {
    const saveTokenToCookie = async () => {
      if (isAuthenticated && !isLoading) {
        try {
          const token = await getAccessTokenSilently();
          tokenManager.setToken(token);
          console.log('Token saved to cookie successfully');
        } catch (error) {
          console.error('Failed to save token to cookie:', error);
        }
      }
    };

    saveTokenToCookie();
  }, [isAuthenticated, isLoading, getAccessTokenSilently]);

  // ログアウト時にCookieからトークンを削除
  useEffect(() => {
    if (!isAuthenticated && !isLoading) {
      tokenManager.removeToken();
      console.log('Token removed from cookie');
      // ログアウト後はホームページにリダイレクト
      if (typeof window !== 'undefined') {
        window.location.href = '/';
      }
    }
  }, [isAuthenticated, isLoading]);

  return {
    isAuthenticated,
    user,
    isLoading,
    hasTokenInCookie: tokenManager.hasToken(),
  };
};