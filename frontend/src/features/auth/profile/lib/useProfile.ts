import { useAuth0 } from '@auth0/auth0-react';
import type { User } from '../../../../entities/auth';

export const useProfile = () => {
  const { user, isAuthenticated, isLoading } = useAuth0();

  return {
    user: user as User | undefined,
    isAuthenticated,
    isLoading,
  };
};
