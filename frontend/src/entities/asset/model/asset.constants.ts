export const ASSET_PERIOD_OPTIONS = [
  { value: '1m', label: '1 Month' },
  { value: '3m', label: '3 Months' },
  { value: '6m', label: '6 Months' },
  { value: '1y', label: '1 Year' },
  { value: 'all', label: 'All Time' },
] as const;

export type AssetPeriod = '1m' | '3m' | '6m' | '1y' | 'all';
