import type { FC } from 'react';
import type { BudgetWithUsage } from '../../model';
import { BUDGET_STATUS_COLORS } from '../../model';
import { calculateBudgetProgress } from '../../lib';

export type BudgetProgressBarProps = {
  budget: BudgetWithUsage;
  size?: 'sm' | 'default' | 'lg';
  showPercentage?: boolean;
  className?: string;
};

export const BudgetProgressBar: FC<BudgetProgressBarProps> = ({
  budget,
  size = 'default',
  showPercentage = true,
  className = '',
}) => {
  const progress = calculateBudgetProgress(budget.amount, budget.used);
  const statusColor = BUDGET_STATUS_COLORS[budget.status];

  const getSizeClasses = (progressSize: string) => {
    const sizeMap = {
      sm: 'h-2',
      default: 'h-3',
      lg: 'h-4',
    };
    return sizeMap[progressSize as keyof typeof sizeMap] || sizeMap.default;
  };

  const getTextSize = (progressSize: string) => {
    const sizeMap = {
      sm: 'text-xs',
      default: 'text-sm',
      lg: 'text-base',
    };
    return sizeMap[progressSize as keyof typeof sizeMap] || sizeMap.default;
  };

  return (
    <div className={`space-y-1 ${className}`}>
      {showPercentage && (
        <div className='flex items-center justify-between'>
          <span className={`font-medium ${getTextSize(size)}`}>予算進捗</span>
          <span
            className={`font-semibold ${getTextSize(size)}`}
            style={{ color: statusColor }}
          >
            {budget.usagePercentage.toFixed(1)}%
          </span>
        </div>
      )}

      <div
        className={`w-full overflow-hidden rounded-full bg-gray-200 ${getSizeClasses(size)}`}
      >
        <div
          className='h-full rounded-full transition-all duration-300 ease-in-out'
          style={{
            width: `${Math.min(progress, 100)}%`,
            backgroundColor: statusColor,
          }}
        />
        {progress > 100 && (
          <div
            className='h-full transition-all duration-300 ease-in-out'
            style={{
              width: `${Math.min(progress - 100, 100)}%`,
              backgroundColor: BUDGET_STATUS_COLORS.exceeded,
              marginLeft: '-100%',
              opacity: 0.7,
            }}
          />
        )}
      </div>
    </div>
  );
};
