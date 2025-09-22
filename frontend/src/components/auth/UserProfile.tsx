import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';

export const UserProfile: React.FC = () => {
  const { user, isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className="p-4 border rounded-lg bg-white shadow-sm">
      <div className="flex items-center space-x-4">
        {user.picture && (
          <img 
            src={user.picture} 
            alt={user.name} 
            className="w-12 h-12 rounded-full"
          />
        )}
        <div>
          <h3 className="font-semibold">{user.name}</h3>
          <p className="text-sm text-gray-600">{user.email}</p>
        </div>
      </div>
    </div>
  );
};