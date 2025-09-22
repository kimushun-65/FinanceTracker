'use client';

import React, { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { useRouter } from 'next/navigation';
import { LoginCard } from '../../../widgets/auth';

export const HomePage: React.FC = () => {
  const { isAuthenticated, isLoading, error } = useAuth0();
  const router = useRouter();

  useEffect(() => {
    if (isAuthenticated) {
      console.log('Auth0 authentication detected, redirecting to dashboard');
      router.push('/dashboard');
    }
  }, [isAuthenticated, router]);

  if (isLoading) {
    return (
      <div className='flex min-h-screen items-center justify-center'>
        <div className='h-12 w-12 animate-spin rounded-full border-b-2 border-blue-600'></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className='flex min-h-screen items-center justify-center'>
        <div className='text-center'>
          <h2 className='mb-2 text-2xl font-bold text-red-600'>
            Authentication Error
          </h2>
          <p className='text-gray-600'>{error.message}</p>
          <pre className='mt-4 rounded bg-gray-100 p-4 text-left text-sm'>
            {JSON.stringify(error, null, 2)}
          </pre>
        </div>
      </div>
    );
  }

  // ログイン済みの場合はダッシュボードにリダイレクトするので、ここは未ログイン時のみ
  if (isAuthenticated) {
    return null;
  }

  return (
    <div className='min-h-screen bg-gray-50'>
      <div className='mx-auto max-w-7xl px-4'>
        <div className='flex items-center justify-center py-4'>
          <img
            src='/image/logo_2/logo_2.png'
            alt='FinSight Logo'
            className='h-60 w-60'
          />
        </div>
      </div>

      <main className='mx-auto max-w-4xl px-4'>
        <div className='text-center'>
          <h2 className='mb-6 text-4xl font-bold text-gray-900'>
            Personal Finance Management
          </h2>
        </div>

        <div className='mb-12 grid gap-8 md:grid-cols-3'>
          <div className='p-6 text-center'>
            <div className='mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-100'>
              <svg
                className='h-8 w-8 text-blue-600'
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'
                />
              </svg>
            </div>
            <h3 className='mb-2 text-lg font-semibold text-gray-900'>
              Track Expenses
            </h3>
            <p className='text-gray-600'>
              Monitor your spending patterns and categorize expenses
              automatically
            </p>
          </div>

          <div className='p-6 text-center'>
            <div className='mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100'>
              <svg
                className='h-8 w-8 text-green-600'
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1'
                />
              </svg>
            </div>
            <h3 className='mb-2 text-lg font-semibold text-gray-900'>
              Manage Budgets
            </h3>
            <p className='text-gray-600'>
              Set and track budgets to keep your spending on target
            </p>
          </div>

          <div className='p-6 text-center'>
            <div className='mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-purple-100'>
              <svg
                className='h-8 w-8 text-purple-600'
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
                />
              </svg>
            </div>
            <h3 className='mb-2 text-lg font-semibold text-gray-900'>
              Financial Reports
            </h3>
            <p className='text-gray-600'>
              Get detailed insights with comprehensive financial reports
            </p>
          </div>
        </div>

        <div className='flex justify-center'>
          <LoginCard />
        </div>
      </main>
    </div>
  );
};
