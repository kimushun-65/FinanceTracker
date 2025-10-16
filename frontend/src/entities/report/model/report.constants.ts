export const COMPARISON_PERIODS = {
  MONTH: 'month',
  QUARTER: 'quarter',
  YEAR: 'year',
} as const;

export const COMPARISON_PERIOD_LABELS: Record<string, string> = {
  month: '月次',
  quarter: '四半期',
  year: '年次',
};
