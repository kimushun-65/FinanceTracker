import { z } from 'zod';
import { baseEntitySchema, moneySchema } from '@/shared/schemas';

export const accountTypeSchema = z.enum(['checking', 'investment', 'cash']);

export const balanceStatusSchema = z.enum(['normal', 'zero', 'negative']);

export const balanceSchema = z.object({
  current: moneySchema,
  initial: moneySchema,
  gain: moneySchema,
  status: balanceStatusSchema,
});

export const accountSchema = baseEntitySchema.extend({
  userId: z.string().uuid({ version: 'v4' }),
  name: z.string().min(1).max(100),
  accountType: accountTypeSchema,
  balance: balanceSchema,
});

export const createAccountPayloadSchema = z.object({
  name: z.string().min(1).max(100),
  accountType: accountTypeSchema,
  initialBalance: moneySchema.optional(),
});

export const updateAccountPayloadSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  accountType: accountTypeSchema.optional(),
});
