'use client';

import { useMemo } from 'react';
import {
  AppLayout,
  BudgetOverviewCards,
  BudgetListTable,
  CreateBudgetModal,
} from '@/widgets';
import { Button, Loading } from '@/shared/ui';
import { useBudgets, useCategories } from '@/features';
import type { BudgetSummary } from '@/entities/budget';
import {
  calculateTotalBudgetAmount,
  calculateTotalUsedAmount,
} from '@/entities/budget';
import { subtractMoney } from '@/shared/value-objects';

export default function BudgetContainer() {
  const { data: budgetData, isLoading: budgetIsLoading } = useBudgets();
  const { data: categories } = useCategories();

  const summary: BudgetSummary | null = useMemo(() => {
    if (!budgetData?.budgets || budgetData.budgets.length === 0) {
      return {
        totalBudget: { amount: 0, currency: 'JPY' },
        totalUsed: { amount: 0, currency: 'JPY' },
        totalRemaining: { amount: 0, currency: 'JPY' },
        overBudgetCount: 0,
        activeCount: 0,
      };
    }

    const budgets = budgetData.budgets;
    const totalBudget = calculateTotalBudgetAmount(budgets);
    const totalUsed = calculateTotalUsedAmount(budgets);
    const totalRemaining = subtractMoney(totalBudget, totalUsed);
    const overBudgetCount = budgets.filter(
      (b) => b.status === 'exceeded',
    ).length;
    const activeCount = budgets.filter((b) => b.isActive).length;

    return {
      totalBudget,
      totalUsed,
      totalRemaining,
      overBudgetCount,
      activeCount,
    };
  }, [budgetData]);

  if (budgetIsLoading) {
    return <Loading />;
  }

  if (!budgetData || !summary || !categories) {
    return <Loading />;
  }

  const budgets = budgetData?.budgets || [];

  return (
    <AppLayout>
      <div className='container mx-auto px-4 py-6'>
        {/* Page Header */}
        <div className='mb-6 flex items-center justify-between'>
          <h1 className='text-3xl font-bold text-gray-900'>
            Budget Management
          </h1>
          <CreateBudgetModal
            trigger={<Button size='lg'>+ Create Budget</Button>}
            categories={categories}
          />
        </div>

        {/* Budget Summary Cards */}
        <BudgetOverviewCards summary={summary} />

        {/* Budget List */}
        <div className='rounded-lg bg-white shadow'>
          <div className='border-b border-gray-200 p-6'>
            <h2 className='text-xl font-semibold text-gray-900'>
              Current Month Budget
            </h2>
            <p className='mt-1 text-sm text-gray-500'>
              Manage your category budgets and track spending
            </p>
          </div>
          <div className='p-6'>
            <BudgetListTable budgets={budgets} categories={categories} />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
