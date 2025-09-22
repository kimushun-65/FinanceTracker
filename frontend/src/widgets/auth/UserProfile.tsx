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
    <div className='flex items-center space-x-2'>
      {user.picture && (
        <img
          src={user.picture}
          alt={user.name || 'User'}
          className='h-8 w-8 rounded-full'
        />
      )}
      <div className='hidden sm:block'>
        <p className='text-sm font-medium text-gray-900'>{user.name}</p>
      </div>
    </div>
  );
};
