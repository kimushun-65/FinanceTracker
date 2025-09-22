import { useAuth0 } from '@auth0/auth0-react';
import { tokenManager } from '../../../../shared/lib/token';

export const useLogout = () => {
  const { logout, isLoading } = useAuth0();

  const logoutUser = () => {
    // Cookieからトークンを削除
    tokenManager.removeToken();
    
    logout({
      logoutParams: {
        returnTo: window.location.origin,
      },
    });
  };

  return {
    logout: logoutUser,
    isLoading,
  };
};
