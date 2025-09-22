import { useAuth0 } from '@auth0/auth0-react';

export const useLogout = () => {
  const { logout, isLoading } = useAuth0();

  const logoutUser = () => {
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
