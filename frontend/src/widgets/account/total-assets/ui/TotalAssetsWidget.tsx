import { Card } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { Money } from '@/shared/value-objects';

interface TotalAssetsWidgetProps {
  totalAssets: Money;
  accountCount: number;
}

export const TotalAssetsWidget = ({
  totalAssets,
  accountCount,
}: TotalAssetsWidgetProps) => {

  return (
    <Card className='p-6'>
      <div className='space-y-4'>
        <div>
          <h2 className='mb-2 text-lg font-medium text-gray-600'>
            Total Assets
          </h2>
          <div className='text-4xl font-bold text-gray-900'>
            {formatMoney(totalAssets)}
          </div>
        </div>

        <div className='text-sm text-gray-500'>{accountCount} accounts</div>

        <div className='text-xs text-gray-400'>
          Updated: {new Date().toLocaleString('ja-JP')}
        </div>
      </div>
    </Card>
  );
};
