import { useAuth0 } from '@auth0/auth0-react';

export const useLogin = () => {
  const { loginWithRedirect, isLoading } = useAuth0();

  const login = () => {
    loginWithRedirect();
  };

  return {
    login,
    isLoading,
  };
};
