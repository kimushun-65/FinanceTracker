'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useMemo } from 'react';
import { ApiClient } from './client';
import { API_BASE_URL } from '../config';

export const useApiClient = () => {
  const { getAccessTokenSilently } = useAuth0();

  const apiClient = useMemo(() => {
    return new ApiClient({
      baseUrl: API_BASE_URL,
      getAccessToken: async () => {
        try {
          const token = await getAccessTokenSilently();
          console.log('Got access token');
          return token;
        } catch (error) {
          console.error('Failed to get access token:', error);
          throw error;
        }
      },
    });
  }, [getAccessTokenSilently]);

  return apiClient;
};
