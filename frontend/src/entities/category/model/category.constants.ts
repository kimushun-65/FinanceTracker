import type { CategoryType } from './category.types';

export const CATEGORY_TYPES: Record<CategoryType, string> = {
  income: '収入',
  expense: '支出',
} as const;

export const DEFAULT_CATEGORY_ICONS = {
  income: {
    salary: '💰',
    bonus: '🎁',
    investment: '📈',
    business: '💼',
    other: '💵',
  },
  expense: {
    food: '🍽️',
    transport: '🚗',
    entertainment: '🎮',
    shopping: '🛒',
    utilities: '💡',
    healthcare: '🏥',
    education: '📚',
    other: '💸',
  },
} as const;

export const DEFAULT_CATEGORY_COLORS = {
  income: ['#10b981', '#059669', '#047857', '#065f46'],
  expense: ['#ef4444', '#dc2626', '#b91c1c', '#991b1b'],
} as const;

export const CATEGORY_VALIDATION = {
  maxNameLength: 50,
  maxCustomNameLength: 50,
  colorRegex: /^#[0-9A-Fa-f]{6}$/,
} as const;
