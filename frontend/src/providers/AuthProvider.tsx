'use client';

import React from 'react';
import { Auth0Provider } from '@auth0/auth0-react';

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  return (
    <Auth0Provider
      domain={process.env.NEXT_PUBLIC_AUTH0_DOMAIN || 'dev-kimushun3765.jp.auth0.com'}
      clientId={process.env.NEXT_PUBLIC_AUTH0_CLIENT_ID || 'aYKkTbJq4WG2L80JOVPEI2MmBSa9VtB7'}
      authorizationParams={{
        redirect_uri: process.env.NEXT_PUBLIC_AUTH0_REDIRECT_URI || 'http://localhost:3000/callback',
        audience: process.env.NEXT_PUBLIC_AUTH0_AUDIENCE || 'https://api.financetracker.local',
        scope: "openid profile email"
      }}
    >
      {children}
    </Auth0Provider>
  );
};