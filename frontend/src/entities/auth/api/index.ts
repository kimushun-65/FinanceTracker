import { api, ApiError } from '../../../shared/api/client';

interface AuthCheckResponse {
  authenticated: boolean;
  user?: any;
}

interface TokenResponse {
  success: boolean;
}

export const tokenManager = {
  setToken: async (token: string): Promise<boolean> => {
    try {
      const response = await api.post<TokenResponse>('/api/v1/auth/token', {
        token,
      });
      return response.data.success;
    } catch (error) {
      if (error && typeof error === 'object' && 'code' in error) {
        console.error('Failed to set token:', (error as ApiError).message);
      } else {
        console.error('Failed to set token:', error);
      }
      return false;
    }
  },

  checkAuth: async (): Promise<{ authenticated: boolean; user?: any }> => {
    try {
      const response = await api.get<AuthCheckResponse>('/api/v1/auth/check');
      return response.data;
    } catch (error) {
      return { authenticated: false };
    }
  },

  removeToken: async (): Promise<boolean> => {
    try {
      const response = await api.delete<TokenResponse>('/api/v1/auth/token');
      return response.data.success;
    } catch (error) {
      if (error && typeof error === 'object' && 'code' in error) {
        console.error('Failed to remove token:', (error as ApiError).message);
      } else {
        console.error('Failed to remove token:', error);
      }
      return false;
    }
  },
};
