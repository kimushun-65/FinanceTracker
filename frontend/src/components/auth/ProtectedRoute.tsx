'use client';

import React from 'react';
import { withAuthenticationRequired } from '@auth0/auth0-react';

interface ProtectedRouteProps {
  component: React.ComponentType<any>;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component }) => {
  const ProtectedComponent = withAuthenticationRequired(Component, {
    onRedirecting: () => (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    ),
  });

  return <ProtectedComponent />;
};