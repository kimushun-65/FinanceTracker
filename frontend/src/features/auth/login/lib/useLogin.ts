import { useAuth0 } from '@auth0/auth0-react';

export const useLogin = () => {
  const { loginWithRedirect, isLoading } = useAuth0();

  const login = () => {
    loginWithRedirect({
      authorizationParams: {
        connection: 'google-oauth2', // Googleログインを強制
      },
    });
  };

  return {
    login,
    isLoading,
  };
};
