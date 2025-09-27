import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useCallback } from 'react';
import { tokenManager } from '../../../../entities/auth/api';

export const useLogout = () => {
  const { logout, isLoading, isAuthenticated } = useAuth0();

  // Auth0の認証状態を監視し、ログアウト時にクリーンアップを実行
  useEffect(() => {
    const handlePostLogout = async () => {
      if (!isAuthenticated && !isLoading) {
        // HttpOnlyクッキーからトークンを削除
        await tokenManager.removeToken();

        // ホームページにリダイレクト（Auth0のreturnToが効かない場合のフォールバック）
        if (typeof window !== 'undefined' && window.location.pathname !== '/') {
          window.location.href = '/';
        }
      }
    };

    handlePostLogout();
  }, [isAuthenticated, isLoading]);

  const logoutUser = useCallback(async () => {
    // 即座にトークンを削除（Auth0のログアウト完了を待たずに）
    await tokenManager.removeToken();

    // Auth0ログアウトを実行
    logout({
      logoutParams: {
        returnTo: window.location.origin,
      },
    });
  }, [logout]);

  return {
    logout: logoutUser,
    isLoading,
  };
};
