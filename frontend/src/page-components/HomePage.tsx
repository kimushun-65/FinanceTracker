'use client';

import Image from 'next/image';
import React, { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { LoginButton, LogoutButton, UserProfile } from '../widgets/auth';

export const HomePage: React.FC = () => {
  const { isAuthenticated, isLoading, error, getAccessTokenSilently } =
    useAuth0();

  useEffect(() => {
    if (isAuthenticated) {
      getAccessTokenSilently()
        .then((token) => {
          console.log(
            'Successfully got token:',
            token.substring(0, 20) + '...',
          );
        })
        .catch((err) => {
          console.error('Failed to get token:', err);
        });
    }
  }, [isAuthenticated, getAccessTokenSilently]);

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

  return (
    <div className='min-h-screen p-8'>
      <header className='mx-auto mb-8 max-w-7xl'>
        <div className='flex items-center justify-between'>
          <div className='flex items-center gap-4'>
            <Image
              src='/image/logo_1/logo_1.png'
              alt='FinSight Logo'
              width={50}
              height={50}
              priority
            />
            <h1 className='text-2xl font-bold'>FinSight</h1>
          </div>

          <div className='flex items-center gap-4'>
            {isAuthenticated ? (
              <>
                <UserProfile />
                <LogoutButton />
              </>
            ) : (
              <LoginButton />
            )}
          </div>
        </div>
      </header>

      <main className='mx-auto max-w-7xl'>
        {isAuthenticated ? (
          <div className='grid gap-6'>
            <div className='rounded-lg border bg-white p-6 shadow-sm'>
              <h2 className='mb-4 text-xl font-semibold'>
                Welcome to FinSight
              </h2>
              <p className='text-gray-600'>
                You are now authenticated. Start managing your finances with our
                powerful tools.
              </p>
            </div>

            <div className='grid gap-6 md:grid-cols-3'>
              <div className='rounded-lg border bg-white p-6 shadow-sm'>
                <h3 className='mb-2 font-semibold'>Track Expenses</h3>
                <p className='text-sm text-gray-600'>
                  Monitor your spending patterns and categorize expenses
                </p>
              </div>

              <div className='rounded-lg border bg-white p-6 shadow-sm'>
                <h3 className='mb-2 font-semibold'>Manage Accounts</h3>
                <p className='text-sm text-gray-600'>
                  Connect and manage all your financial accounts in one place
                </p>
              </div>

              <div className='rounded-lg border bg-white p-6 shadow-sm'>
                <h3 className='mb-2 font-semibold'>View Reports</h3>
                <p className='text-sm text-gray-600'>
                  Get insights with detailed financial reports and analytics
                </p>
              </div>
            </div>
          </div>
        ) : (
          <div className='text-center'>
            <h2 className='mb-4 text-3xl font-bold'>
              Personal Finance Management
            </h2>
            <p className='mb-8 text-xl text-gray-600'>
              Track expenses, manage budgets, and achieve your financial goals
            </p>
            <div className='flex justify-center'>
              <LoginButton />
            </div>
          </div>
        )}
      </main>
    </div>
  );
};
