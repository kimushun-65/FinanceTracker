import type { Money } from '@/shared/value-objects';
import { formatMoney as formatMoneyUtil } from '@/shared/value-objects/money';

export const formatAssetAmount = (money: Money): string => {
  return formatMoneyUtil(money);
};

export const formatPercentage = (value: number): string => {
  return `${value >= 0 ? '+' : ''}${value.toFixed(1)}%`;
};

export const formatSnapshotDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
};
