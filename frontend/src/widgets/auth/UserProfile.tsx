import React from 'react';
import { useProfile } from '../../features/auth';

export const UserProfile: React.FC = () => {
  const { user, isAuthenticated, isLoading } = useProfile();

  if (isLoading) {
    return <div className='p-4'>Loading...</div>;
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className='rounded-lg border bg-white p-4 shadow-sm'>
      <div className='flex items-center space-x-4'>
        {user.picture && (
          <img
            src={user.picture}
            alt={user.name || 'User'}
            className='h-12 w-12 rounded-full'
          />
        )}
        <div>
          <h3 className='font-semibold'>{user.name}</h3>
          <p className='text-sm text-gray-600'>{user.email}</p>
        </div>
      </div>
    </div>
  );
};
