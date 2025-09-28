import { z } from 'zod';

export const baseEntitySchema = z.object({
  id: z.string().uuid({ version: 'v4' }),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});
