'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';
import { tokenManager } from './token';
import { env } from '../config/env';

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
  const [isCheckingAuth, setIsCheckingAuth] = useState(false);
  const [isTokenReady, setIsTokenReady] = useState(false);

  // ログイン後にトークンをHttpOnlyクッキーに保存
  useEffect(() => {
    const saveTokenToCookie = async () => {
      if (isAuthenticated && !isLoading) {
        try {
          const token = await getAccessTokenSilently();
          const success = await tokenManager.setToken(token);
          if (success) {
            
            // トークンが準備完了したことをマーク
            setIsTokenReady(true);
            
            // クッキーの保存が確実に完了するまで待機
            await new Promise((resolve) => setTimeout(resolve, 500));

            // バックエンドでの認証状態を確認（リトライ付き）
            let authStatus = { authenticated: false };
            let retryCount = 0;
            const maxRetries = 3;
            
            while (retryCount < maxRetries) {
              try {
                authStatus = await tokenManager.checkAuth();
                if (authStatus.authenticated) {
                  break;
                }
                retryCount++;
                if (retryCount < maxRetries) {
                  await new Promise((resolve) => setTimeout(resolve, 500));
                }
              } catch (error) {
                retryCount++;
                if (retryCount < maxRetries) {
                  await new Promise((resolve) => setTimeout(resolve, 500));
                }
              }
            }
            
            setBackendAuthStatus(authStatus);

            // 認証が成功した場合のみユーザー情報をDBに同期
            if (authStatus?.authenticated) {
              try {
                const response = await fetch(`${env.api.baseUrl}/api/v1/auth/user`, {
                  method: 'GET',
                  credentials: 'include',
                });

                if (!response.ok) {
                  console.error('Failed to sync user:', response.statusText);
                }
              } catch (error) {
                console.error('Error syncing user:', error);
              }
            }
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
        setIsTokenReady(false);
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
    // Auth0の認証中、既に確認中、既に認証状態が確定している、またはトークンが準備できていない場合はスキップ
    if (isLoading || isCheckingAuth || backendAuthStatus !== null || !isTokenReady) {
      return;
    }
    
    const checkBackendAuth = async () => {
      setIsCheckingAuth(true);
      try {
        const authStatus = await tokenManager.checkAuth();
        setBackendAuthStatus(authStatus);
      } catch (error) {
        // 初回ロード時の401エラーは想定内なので、エラーログは出さない
        setBackendAuthStatus({ authenticated: false });
      } finally {
        setIsCheckingAuth(false);
      }
    };

    checkBackendAuth();
  }, [isLoading, isCheckingAuth, backendAuthStatus, isTokenReady]);

  return {
    isAuthenticated,
    user,
    isLoading,
    backendAuthStatus,
    isFullyAuthenticated: isAuthenticated && backendAuthStatus?.authenticated,
  };
};
