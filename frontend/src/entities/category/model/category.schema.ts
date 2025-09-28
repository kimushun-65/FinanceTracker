import { z } from 'zod';
import { baseEntitySchema } from '@/shared/schemas';

export const categoryTypeSchema = z.enum(['income', 'expense']);

export const categorySchema = baseEntitySchema.extend({
  userId: z.string().uuid({ version: 'v4' }),
  categoryMasterId: z.string().uuid({ version: 'v4' }),
  customName: z.string().min(1).max(50).optional(),
  isActive: z.boolean(),
});

export const categoryMasterSchema = baseEntitySchema.extend({
  name: z.string().min(1).max(50),
  type: categoryTypeSchema,
  icon: z.string().min(1).max(50),
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/),
  displayOrder: z.number().int().min(0),
});

export const createCategoryPayloadSchema = z.object({
  categoryMasterId: z.string().uuid({ version: 'v4' }),
  customName: z.string().min(1).max(50).optional(),
});

export const updateCategoryPayloadSchema = z.object({
  customName: z.string().min(1).max(50).optional(),
  isActive: z.boolean().optional(),
});
