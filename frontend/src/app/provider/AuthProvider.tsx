'use client';

import React from 'react';
import { Auth0Provider } from '@auth0/auth0-react';
import { getAuthConfig } from '@/entities';

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const config = getAuthConfig();

  return (
    <Auth0Provider
      domain={config.domain}
      clientId={config.clientId}
      authorizationParams={{
        redirect_uri: config.redirectUri,
        audience: config.audience,
        scope: config.scope,
      }}
    >
      {children}
    </Auth0Provider>
  );
};
