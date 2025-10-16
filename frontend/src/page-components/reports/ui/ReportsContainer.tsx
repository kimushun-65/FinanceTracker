'use client';

import React, { useState } from 'react';
import { startOfMonth, endOfMonth } from 'date-fns';
import { AppLayout } from '../../../widgets/layout';
import { ReportFilters } from '../../../widgets/reports/report-filters';
import { MonthlyReportCard } from '../../../widgets/reports/monthly-report';
import {
  useTransactions,
  useTransactionCategorySummary,
} from '../../../features/transaction-management';
import { useCategories } from '../../../features/category-management';
import { useMonthlyReport } from '../../../features/report-management';
import type { DateRange, ComparisonPeriod } from '../../../entities/report';

export const ReportsContainer: React.FC = () => {
  // フィルター状態
  const [dateRange, setDateRange] = useState<DateRange>({
    from: startOfMonth(new Date()),
    to: endOfMonth(new Date()),
  });
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [comparisonPeriod, setComparisonPeriod] =
    useState<ComparisonPeriod>('month');

  // データ取得
  const { data: categories } = useCategories();
  const { data: transactionsData, isLoading: isTransactionsLoading } =
    useTransactions({
      startDate: dateRange.from.toISOString().split('T')[0],
      endDate: dateRange.to.toISOString().split('T')[0],
    });
  const { data: categorySummary, isLoading: isCategorySummaryLoading } =
    useTransactionCategorySummary({
      from: dateRange.from.toISOString().split('T')[0],
      to: dateRange.to.toISOString().split('T')[0],
      type: 'expense',
    });

  // 計算処理
  const monthlyReport = useMonthlyReport(
    transactionsData?.transactions,
    categorySummary,
  );

  const isLoading = isTransactionsLoading || isCategorySummaryLoading;

  return (
    <AppLayout title='Reports'>
      <div className='mb-8'>
        <h2 className='mb-2 text-3xl font-bold text-gray-900'>Reports</h2>
      </div>

      <div className='space-y-6'>
        {/* フィルター */}
        <ReportFilters
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
          selectedCategories={selectedCategories}
          onCategoriesChange={setSelectedCategories}
          categories={categories || []}
          comparisonPeriod={comparisonPeriod}
          onComparisonPeriodChange={setComparisonPeriod}
        />

        {/* 月次レポートカード */}
        <MonthlyReportCard
          report={monthlyReport}
          dateRange={dateRange}
          isLoading={isLoading}
        />

        {/* 他のウィジェットは後で追加 */}
      </div>
    </AppLayout>
  );
};
