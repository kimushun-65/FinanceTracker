'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';
import { tokenManager } from './token';

/**
 * Auth0認証とHttpOnlyクッキー管理を統合するカスタムフック
 */
export const useAuthWithCookie = () => {
  const { isAuthenticated, getAccessTokenSilently, user, isLoading } =
    useAuth0();
  const [backendAuthStatus, setBackendAuthStatus] = useState<{
    authenticated: boolean;
    user?: any;
  } | null>(null);

  // ログイン後にトークンをHttpOnlyクッキーに保存
  useEffect(() => {
    const saveTokenToCookie = async () => {
      if (isAuthenticated && !isLoading) {
        try {
          console.log('Auth0 authenticated, getting access token...');
          const token = await getAccessTokenSilently();
          console.log('Got access token, saving to HttpOnly cookie...');
          const success = await tokenManager.setToken(token);
          if (success) {
            console.log('Token saved to HttpOnly cookie successfully');
            // バックエンドでの認証状態を確認
            const authStatus = await tokenManager.checkAuth();
            console.log('Backend auth status:', authStatus);
            setBackendAuthStatus(authStatus);
          } else {
            console.error('Failed to save token to HttpOnly cookie');
          }
        } catch (error) {
          console.error('Failed to save token to cookie:', error);
        }
      }
    };

    saveTokenToCookie();
  }, [isAuthenticated, isLoading, getAccessTokenSilently]);

  // ログアウト時にHttpOnlyクッキーからトークンを削除
  useEffect(() => {
    const handleLogout = async () => {
      if (!isAuthenticated && !isLoading) {
        await tokenManager.removeToken();
        setBackendAuthStatus(null);
        console.log('Token removed from HttpOnly cookie');
        // ログアウト後はホームページにリダイレクト
        if (typeof window !== 'undefined') {
          window.location.href = '/';
        }
      }
    };

    handleLogout();
  }, [isAuthenticated, isLoading]);

  // ページロード時にバックエンドの認証状態を確認
  useEffect(() => {
    const checkBackendAuth = async () => {
      if (!isLoading) {
        const authStatus = await tokenManager.checkAuth();
        setBackendAuthStatus(authStatus);
      }
    };

    checkBackendAuth();
  }, [isLoading]);

  return {
    isAuthenticated,
    user,
    isLoading,
    backendAuthStatus,
    isFullyAuthenticated: isAuthenticated && backendAuthStatus?.authenticated,
  };
};
