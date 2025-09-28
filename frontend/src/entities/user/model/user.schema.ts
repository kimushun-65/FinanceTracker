import { z } from 'zod';
import { baseEntitySchema } from '@/shared/schemas';

export const userSchema = baseEntitySchema.extend({
  auth0UserId: z.string().min(1),
  email: z.string().email(),
  name: z.string().min(1).max(100),
  emailVerified: z.boolean(),
});

export const updateUserPayloadSchema = z.object({
  name: z.string().min(1).max(100).optional(),
  email: z.string().email().optional(),
});

export const createUserPayloadSchema = z.object({
  auth0UserId: z.string().min(1),
  email: z.string().email(),
  name: z.string().min(1).max(100),
  emailVerified: z.boolean().optional().default(false),
});

export const authUserSchema = z.object({
  id: z.string().optional(),
  name: z.string().optional(),
  email: z.string().email().optional(),
  picture: z.string().url().optional(),
  email_verified: z.boolean().optional(),
  sub: z.string().optional(),
});

export const userProfileSchema = userSchema.extend({
  picture: z.string().url().optional(),
  lastLoginAt: z.string().datetime().optional(),
  loginCount: z.number().int().min(0),
});
