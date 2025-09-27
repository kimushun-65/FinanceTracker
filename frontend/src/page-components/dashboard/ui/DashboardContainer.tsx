'use client';

import React from 'react';
import { AppLayout } from '../../../widgets/layout';
import { Card, CardContent, CardHeader, CardTitle } from '../../../shared/ui';
import { useAuthWithCookie, useLogout } from '../../../features/auth';

export const DashboardContainer: React.FC = () => {
  // Auth0認証とCookieの同期（ダッシュボードでトークン管理）
  useAuthWithCookie();
  
  // ログアウト時のクリーンアップも含む
  useLogout();

  return (
    <AppLayout title='Dashboard'>
      <div className='mb-8'>
        <h2 className='mb-2 text-3xl font-bold text-gray-900'>Dashboard</h2>
        <p className='text-gray-600'>
          Welcome back! Here&apos;s your financial overview.
        </p>
      </div>

      <div className='mb-8 grid gap-6 md:grid-cols-2 lg:grid-cols-4'>
        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>Total Balance</CardTitle>
            <svg
              className='h-4 w-4 text-muted-foreground'
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
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>¥1,234,567</div>
            <p className='text-xs text-muted-foreground'>
              +2.5% from last month
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Monthly Expenses
            </CardTitle>
            <svg
              className='h-4 w-4 text-muted-foreground'
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
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>¥345,678</div>
            <p className='text-xs text-muted-foreground'>
              -5.1% from last month
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Budget Remaining
            </CardTitle>
            <svg
              className='h-4 w-4 text-muted-foreground'
              fill='none'
              stroke='currentColor'
              viewBox='0 0 24 24'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M13 7h8m0 0v8m0-8l-8 8-4-4-6 6'
              />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>¥154,322</div>
            <p className='text-xs text-muted-foreground'>
              68% of monthly budget
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>Savings Goal</CardTitle>
            <svg
              className='h-4 w-4 text-muted-foreground'
              fill='none'
              stroke='currentColor'
              viewBox='0 0 24 24'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>82%</div>
            <p className='text-xs text-muted-foreground'>
              ¥820,000 of ¥1,000,000
            </p>
          </CardContent>
        </Card>
      </div>

      <div className='grid gap-6 md:grid-cols-2'>
        <Card>
          <CardHeader>
            <CardTitle>Recent Transactions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className='space-y-4'>
              <div className='flex items-center justify-between'>
                <div>
                  <p className='font-medium'>Grocery Store</p>
                  <p className='text-sm text-muted-foreground'>Food & Dining</p>
                </div>
                <p className='font-medium text-red-600'>-¥8,420</p>
              </div>
              <div className='flex items-center justify-between'>
                <div>
                  <p className='font-medium'>Salary</p>
                  <p className='text-sm text-muted-foreground'>Income</p>
                </div>
                <p className='font-medium text-green-600'>+¥350,000</p>
              </div>
              <div className='flex items-center justify-between'>
                <div>
                  <p className='font-medium'>Gas Station</p>
                  <p className='text-sm text-muted-foreground'>
                    Transportation
                  </p>
                </div>
                <p className='font-medium text-red-600'>-¥6,500</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Budget Overview</CardTitle>
          </CardHeader>
          <CardContent>
            <div className='space-y-4'>
              <div>
                <div className='mb-2 flex items-center justify-between'>
                  <span className='text-sm font-medium'>Food & Dining</span>
                  <span className='text-sm text-muted-foreground'>
                    ¥45,000 / ¥60,000
                  </span>
                </div>
                <div className='h-2 w-full rounded-full bg-gray-200'>
                  <div
                    className='h-2 rounded-full bg-blue-600'
                    style={{ width: '75%' }}
                  ></div>
                </div>
              </div>
              <div>
                <div className='mb-2 flex items-center justify-between'>
                  <span className='text-sm font-medium'>Transportation</span>
                  <span className='text-sm text-muted-foreground'>
                    ¥25,000 / ¥40,000
                  </span>
                </div>
                <div className='h-2 w-full rounded-full bg-gray-200'>
                  <div
                    className='h-2 rounded-full bg-green-600'
                    style={{ width: '62.5%' }}
                  ></div>
                </div>
              </div>
              <div>
                <div className='mb-2 flex items-center justify-between'>
                  <span className='text-sm font-medium'>Entertainment</span>
                  <span className='text-sm text-muted-foreground'>
                    ¥18,000 / ¥30,000
                  </span>
                </div>
                <div className='h-2 w-full rounded-full bg-gray-200'>
                  <div
                    className='h-2 rounded-full bg-yellow-600'
                    style={{ width: '60%' }}
                  ></div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  );
};
