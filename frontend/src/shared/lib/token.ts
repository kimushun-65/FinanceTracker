'use client';

const TOKEN_COOKIE_NAME = 'access_token';

export const tokenManager = {
  setToken: (token: string) => {
    document.cookie = `${TOKEN_COOKIE_NAME}=${token}; path=/; secure; samesite=strict`;
  },

  getToken: (): string | null => {
    if (typeof document === 'undefined') return null;
    
    const cookies = document.cookie.split(';');
    const tokenCookie = cookies.find(cookie => 
      cookie.trim().startsWith(`${TOKEN_COOKIE_NAME}=`)
    );
    
    if (!tokenCookie) return null;
    
    return tokenCookie.split('=')[1];
  },

  removeToken: () => {
    document.cookie = `${TOKEN_COOKIE_NAME}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`;
  },

  hasToken: (): boolean => {
    return !!tokenManager.getToken();
  }
};