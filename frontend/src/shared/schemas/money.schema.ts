import { z } from 'zod';

export const currencySchema = z.enum([
  'USD',
  'EUR',
  'JPY',
  'GBP',
  'CNY',
  'KRW',
]);

export const moneySchema = z.object({
  amount: z.number(),
  currency: currencySchema,
});
