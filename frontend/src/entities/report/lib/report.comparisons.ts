import { subMonths, subQuarters, subYears, format } from 'date-fns';
import type {
  DateRange,
  ComparisonPeriod,
  TrendData,
  CategorySummary,
} from '../model';

// CategoryWithMasterの最小限の型定義
interface CategoryWithMaster {
  id: string;
  displayName: string;
}

/**
 * 比較期間の計算
 */
export const calculateComparisonPeriod = (
  currentRange: DateRange,
  period: ComparisonPeriod,
): DateRange => {
  switch (period) {
    case 'month':
      return {
        from: subMonths(currentRange.from, 1),
        to: subMonths(currentRange.to, 1),
      };
    case 'quarter':
      return {
        from: subQuarters(currentRange.from, 1),
        to: subQuarters(currentRange.to, 1),
      };
    case 'year':
      return {
        from: subYears(currentRange.from, 1),
        to: subYears(currentRange.to, 1),
      };
  }
};

/**
 * トレンドデータ変換
 */
export const transformToTrendData = (
  categorySummary: CategorySummary[],
  categories: CategoryWithMaster[],
): TrendData[] => {
  // グループ化された日付別データを作成
  const dateMap = new Map<string, Record<string, number>>();

  categorySummary.forEach((summary) => {
    // ここでは簡易実装として、カテゴリごとのデータを返す
    const date = format(new Date(), 'yyyy-MM-dd');

    if (!dateMap.has(date)) {
      dateMap.set(date, {});
    }

    const dayData = dateMap.get(date)!;
    dayData[summary.categoryId] = summary.totalAmount.amount;
  });

  // TrendData配列に変換
  return Array.from(dateMap.entries()).map(([date, data]) => ({
    date,
    ...data,
  }));
};
