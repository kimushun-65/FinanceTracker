'use client';

import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui';
import { formatMoney } from '@/shared/value-objects/money';
import type { AssetTrendData, AssetPeriod } from '@/entities/asset';

interface AssetTrendChartProps {
  data: AssetTrendData[];
  period: AssetPeriod;
  onPeriodChange: (period: AssetPeriod) => void;
  isLoading?: boolean;
}

const PERIOD_OPTIONS: { value: AssetPeriod; label: string }[] = [
  { value: '1m', label: '1ヶ月' },
  { value: '3m', label: '3ヶ月' },
  { value: '6m', label: '6ヶ月' },
  { value: '1y', label: '1年' },
  { value: 'all', label: 'すべて' },
];

export const AssetTrendChart: React.FC<AssetTrendChartProps> = ({
  data,
  period,
  onPeriodChange,
  isLoading = false,
}) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
    });
  };

  const formatYAxis = (value: number) => {
    if (value >= 1000000) {
      return `¥${(value / 1000000).toFixed(1)}M`;
    }
    if (value >= 1000) {
      return `¥${(value / 1000).toFixed(0)}K`;
    }
    return `¥${value}`;
  };

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className='rounded-lg border bg-white p-3 shadow-lg'>
          <p className='mb-2 font-medium'>{formatDate(label)}</p>
          <p className='text-lg font-bold text-blue-600'>
            {formatMoney({ amount: payload[0].value, currency: 'JPY' })}
          </p>
        </div>
      );
    }
    return null;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>資産推移</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='flex h-[400px] items-center justify-center'>
            <div className='animate-pulse text-muted-foreground'>
              データを読み込み中...
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card>
        <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
          <CardTitle>資産推移</CardTitle>
          <div className='flex gap-2'>
            {PERIOD_OPTIONS.map((option) => (
              <button
                key={option.value}
                onClick={() => onPeriodChange(option.value)}
                className={`rounded px-3 py-1 text-sm ${
                  period === option.value
                    ? 'bg-blue-600 text-white'
                    : 'bg-muted text-muted-foreground hover:bg-muted/80'
                }`}
              >
                {option.label}
              </button>
            ))}
          </div>
        </CardHeader>
        <CardContent>
          <div className='flex h-[400px] items-center justify-center'>
            <div className='text-center text-muted-foreground'>
              <p className='mb-2 text-lg font-medium'>
                スナップショットデータがありません
              </p>
              <p className='text-sm'>
                「スナップショット作成」ボタンから資産を記録してください
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
        <CardTitle>資産推移</CardTitle>
        <div className='flex gap-2'>
          {PERIOD_OPTIONS.map((option) => (
            <button
              key={option.value}
              onClick={() => onPeriodChange(option.value)}
              className={`rounded px-3 py-1 text-sm transition-colors ${
                period === option.value
                  ? 'bg-blue-600 text-white'
                  : 'bg-muted text-muted-foreground hover:bg-muted/80'
              }`}
            >
              {option.label}
            </button>
          ))}
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width='100%' height={400}>
          <LineChart
            data={data}
            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray='3 3' stroke='#e0e0e0' />
            <XAxis
              dataKey='date'
              tickFormatter={formatDate}
              stroke='#666'
              style={{ fontSize: '12px' }}
            />
            <YAxis
              tickFormatter={formatYAxis}
              stroke='#666'
              style={{ fontSize: '12px' }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend
              wrapperStyle={{ paddingTop: '20px' }}
              formatter={() => '総資産'}
            />
            <Line
              type='monotone'
              dataKey='totalAssets'
              stroke='#3b82f6'
              strokeWidth={2}
              dot={{ fill: '#3b82f6', r: 4 }}
              activeDot={{ r: 6 }}
              name='総資産'
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
};
