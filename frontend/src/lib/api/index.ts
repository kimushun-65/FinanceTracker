'use client';

import { useAuth0 } from '@auth0/auth0-react';
import { useMemo } from 'react';
import ApiClient from './client';

export const useApiClient = () => {
  const { getAccessTokenSilently } = useAuth0();

  const apiClient = useMemo(() => {
    return new ApiClient({
      baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
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