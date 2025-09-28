import { baseEntitySchema, moneySchema } from '@/shared/schemas';
import { z } from 'zod';

export const transactionTypeSchema = z.enum(['income', 'expense']);

export const transactionSchema = baseEntitySchema.extend({
  userId: z.string().uuid(),
  accountId: z.string().uuid(),
  categoryId: z.string().uuid(),
  amount: moneySchema,
  type: transactionTypeSchema,
  description: z.string().min(0).max(500),
  date: z.string().datetime(),
});

export const createTransactionPayloadSchema = z.object({
  accountId: z.string().uuid(),
  categoryId: z.string().uuid(),
  amount: moneySchema,
  type: transactionTypeSchema,
  description: z.string().min(0).max(500),
  date: z.string().datetime().optional(),
});

export const updateTransactionPayloadSchema = z.object({
  amount: moneySchema.optional(),
  description: z.string().min(0).max(500).optional(),
  date: z.string().datetime().optional(),
  categoryId: z.string().uuid().optional(),
  accountId: z.string().uuid().optional(),
});

export const transactionListParamsSchema = z.object({
  accountId: z.string().uuid().optional(),
  categoryId: z.string().uuid().optional(),
  type: transactionTypeSchema.optional(),
  startDate: z.string().datetime().optional(),
  endDate: z.string().datetime().optional(),
  limit: z.number().min(1).max(100).default(20).optional(),
  offset: z.number().min(0).default(0).optional(),
  sortBy: z.enum(['date', 'amount']).default('date').optional(),
  sortOrder: z.enum(['asc', 'desc']).default('desc').optional(),
});
