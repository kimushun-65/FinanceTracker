'use client';

import { env } from '../config/env';

export const tokenManager = {
  setToken: async (token: string): Promise<boolean> => {
    try {
      const response = await fetch(`${env.api.baseUrl}/api/v1/auth/token`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ token }),
      });
      return response.ok;
    } catch (error) {
      console.error('Failed to set token:', error);
      return false;
    }
  },

  checkAuth: async (): Promise<{ authenticated: boolean; user?: any }> => {
    try {
      const response = await fetch(`${env.api.baseUrl}/api/v1/auth/check`, {
        method: 'GET',
        credentials: 'include',
      });
      
      if (!response.ok) {
        return { authenticated: false };
      }
      
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Failed to check auth:', error);
      return { authenticated: false };
    }
  },

  removeToken: async (): Promise<boolean> => {
    try {
      const response = await fetch(`${env.api.baseUrl}/api/v1/auth/token`, {
        method: 'DELETE',
        credentials: 'include',
      });
      return response.ok;
    } catch (error) {
      console.error('Failed to remove token:', error);
      return false;
    }
  },

  // Deprecated: Use checkAuth instead
  getToken: (): null => {
    console.warn('getToken is deprecated. Use checkAuth instead.');
    return null;
  },

  // Deprecated: Use checkAuth instead
  hasToken: (): boolean => {
    console.warn('hasToken is deprecated. Use checkAuth instead.');
    return false;
  }
};