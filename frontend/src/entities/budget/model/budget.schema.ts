import { z } from 'zod';
import { baseEntitySchema, moneySchema } from '@/shared/schemas';

export const periodTypeSchema = z.enum(['monthly', 'yearly']);

export const budgetStatusSchema = z.enum(['normal', 'warning', 'exceeded']);

export const budgetSchema = baseEntitySchema.extend({
  userId: z.string().uuid({ version: 'v4' }),
  categoryId: z.string().uuid({ version: 'v4' }),
  amount: moneySchema,
  periodType: periodTypeSchema,
  startDate: z.string().date(),
  endDate: z.string().date().optional(),
  isActive: z.boolean(),
});

export const createBudgetPayloadSchema = z.object({
  categoryId: z.string().uuid({ version: 'v4' }),
  amount: moneySchema,
  periodType: periodTypeSchema,
  startDate: z.string().date(),
  endDate: z.string().date().optional(),
});

export const updateBudgetPayloadSchema = z.object({
  amount: moneySchema.optional(),
  startDate: z.string().date().optional(),
  endDate: z.string().date().optional(),
  isActive: z.boolean().optional(),
});

export const budgetWithUsageSchema = budgetSchema.extend({
  used: moneySchema,
  remaining: moneySchema,
  usagePercentage: z.number().min(0).max(200),
  status: budgetStatusSchema,
});
