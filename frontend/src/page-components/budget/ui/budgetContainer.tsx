'use client';

import {
  AppLayout,
  BudgetOverviewCards,
  BudgetListTable,
  CreateBudgetModal,
} from '@/widgets';
import { Button, Loading } from '@/shared/ui';
import { useCurrentBudgets, useBudgetSummary, useCategories } from '@/features';

export default function BudgetContainer() {
  const { data: budgetData, isLoading: budgetIsLoading } = useCurrentBudgets();
  const { data: summary, isLoading: summaryIsLoading } =
    useBudgetSummary('monthly');
  const { data: categories } = useCategories();

  if (budgetIsLoading || summaryIsLoading) {
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
          <div className='p-6'>
            <BudgetListTable budgets={budgets} categories={categories} />
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
