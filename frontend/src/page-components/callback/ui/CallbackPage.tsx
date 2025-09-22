'use client';

import { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { useRouter } from 'next/navigation';

export const CallbackPage: React.FC = () => {
  const { isLoading, error } = useAuth0();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !error) {
      router.push('/');
    }
  }, [isLoading, error, router]);

  if (error) {
    return (
      <div className='flex min-h-screen items-center justify-center'>
        <div className='text-center'>
          <h2 className='mb-2 text-2xl font-bold text-red-600'>
            Authentication Error
          </h2>
          <p className='text-gray-600'>{error.message}</p>
          <button
            onClick={() => router.push('/')}
            className='mt-4 rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700'
          >
            Return Home
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className='flex min-h-screen items-center justify-center'>
      <div className='text-center'>
        <div className='mx-auto h-12 w-12 animate-spin rounded-full border-b-2 border-blue-600'></div>
        <p className='mt-4 text-gray-600'>Completing authentication...</p>
      </div>
    </div>
  );
};
