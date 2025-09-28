'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useEffect } from 'react';
import { tokenManager } from '../../../../entities/auth/api';
import { api, ApiError } from '../../../../shared/api/client';

/**
 * Auth0認証とHttpOnlyクッキー管理を統合するカスタムフック
 */
export const useAuthWithCookie = () => {
  const { isAuthenticated, getAccessTokenSilently, isLoading } = useAuth0();

  // ログイン後にトークンをHttpOnlyクッキーに保存
  useEffect(() => {
    const saveTokenToCookie = async () => {
      if (isAuthenticated && !isLoading) {
        try {
          const token = await getAccessTokenSilently();
          const success = await tokenManager.setToken(token);
          if (success) {
            // ユーザー情報をDBに同期
            try {
              await api.get('/api/v1/auth/user');
            } catch (error) {
              if (error && typeof error === 'object' && 'code' in error) {
                console.error(
                  'Failed to sync user:',
                  (error as ApiError).message,
                );
              } else {
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

  // このフックは副作用のみで値を返さない
  // DashboardContainerでトークン管理の初期化のために呼ばれる
};
