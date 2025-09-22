import React from 'react';
import { useLogout } from '../../features/auth';

export const LogoutButton: React.FC = () => {
  const { logout, isLoading } = useLogout();

  return (
    <button
      onClick={logout}
      disabled={isLoading}
      className='rounded bg-gray-600 px-4 py-2 text-white transition-colors hover:bg-gray-700 disabled:opacity-50'
    >
      {isLoading ? 'Logging out...' : 'Log Out'}
    </button>
  );
};
