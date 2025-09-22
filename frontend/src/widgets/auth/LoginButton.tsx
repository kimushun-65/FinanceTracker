import React from 'react';
import { useLogin } from '../../features/auth';

export const LoginButton: React.FC = () => {
  const { login, isLoading } = useLogin();

  return (
    <button
      onClick={login}
      disabled={isLoading}
      className='rounded bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700 disabled:opacity-50'
    >
      {isLoading ? 'Logging in...' : 'Log In'}
    </button>
  );
};
