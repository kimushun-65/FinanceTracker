'use client';

import Image from 'next/image';
import { useAuth0 } from '@auth0/auth0-react';
import React from 'react';
import { LoginButton } from '@/components/auth/LoginButton';
import { LogoutButton } from '@/components/auth/LogoutButton';
import { UserProfile } from '@/components/auth/UserProfile';

export default function Home() {
  const { isAuthenticated, isLoading, error, getAccessTokenSilently } = useAuth0();

  // デバッグ用：認証後にトークン取得を試みる
  React.useEffect(() => {
    if (isAuthenticated) {
      getAccessTokenSilently()
        .then(token => {
          console.log('Successfully got token:', token.substring(0, 20) + '...');
        })
        .catch(err => {
          console.error('Failed to get token:', err);
        });
    }
  }, [isAuthenticated, getAccessTokenSilently]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600 mb-2">Authentication Error</h2>
          <p className="text-gray-600">{error.message}</p>
          <pre className="mt-4 p-4 bg-gray-100 rounded text-sm text-left">
            {JSON.stringify(error, null, 2)}
          </pre>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-8">
      <header className="max-w-7xl mx-auto mb-8">
        <div className="flex justify-between items-center">
          <div className="flex items-center gap-4">
            <Image
              src="/image/logo_1/logo_1.png"
              alt="FinSight Logo"
              width={50}
              height={50}
              priority
            />
            <h1 className="text-2xl font-bold">FinSight</h1>
          </div>
          
          <div className="flex items-center gap-4">
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

      <main className="max-w-7xl mx-auto">
        {isAuthenticated ? (
          <div className="grid gap-6">
            <div className="bg-white p-6 rounded-lg shadow-sm border">
              <h2 className="text-xl font-semibold mb-4">Welcome to FinSight</h2>
              <p className="text-gray-600">
                You are now authenticated. Start managing your finances with our powerful tools.
              </p>
            </div>
            
            <div className="grid md:grid-cols-3 gap-6">
              <div className="bg-white p-6 rounded-lg shadow-sm border">
                <h3 className="font-semibold mb-2">Track Expenses</h3>
                <p className="text-sm text-gray-600">Monitor your spending patterns and categorize expenses</p>
              </div>
              
              <div className="bg-white p-6 rounded-lg shadow-sm border">
                <h3 className="font-semibold mb-2">Manage Accounts</h3>
                <p className="text-sm text-gray-600">Connect and manage all your financial accounts in one place</p>
              </div>
              
              <div className="bg-white p-6 rounded-lg shadow-sm border">
                <h3 className="font-semibold mb-2">View Reports</h3>
                <p className="text-sm text-gray-600">Get insights with detailed financial reports and analytics</p>
              </div>
            </div>
          </div>
        ) : (
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-4">Personal Finance Management</h2>
            <p className="text-xl text-gray-600 mb-8">
              Track expenses, manage budgets, and achieve your financial goals
            </p>
            <div className="flex justify-center">
              <LoginButton />
            </div>
          </div>
        )}
      </main>
    </div>
  );
}