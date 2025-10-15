'use client';

import React from 'react';
import { AppLayout } from '../../../widgets/layout';
import { useAuthWithCookie, useLogout } from '../../../features/auth';
import {
  FinancialSummaryCards,
  RecentTransactionsList,
} from '../../../widgets/dashboard';

export const DashboardContainer: React.FC = () => {
  // Auth0認証とCookieの同期（ダッシュボードでトークン管理）
  useAuthWithCookie();

  // ログアウト時のクリーンアップも含む
  useLogout();

  return (
    <AppLayout title='Dashboard'>
      <div className='mb-8'>
        <h2 className='mb-2 text-3xl font-bold text-gray-900'>Dashboard</h2>
      </div>

      {/* 財務サマリーカード */}
      <div className='mb-8'>
        <FinancialSummaryCards />
      </div>

      {/* 最近の取引リスト */}
      <div className='mb-8'>
        <RecentTransactionsList />
      </div>
    </AppLayout>
  );
};
