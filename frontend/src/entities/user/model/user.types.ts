import type { BaseEntity } from '@/shared/types';
import type { Email } from '@/shared/value-objects';

export type User = BaseEntity & {
  auth0UserId: string;
  email: Email;
  name: string;
  emailVerified: boolean;
};

export type UpdateUserPayload = {
  name?: string;
  email?: string;
};

export type AuthUser = {
  id?: string;
  name?: string;
  email?: string;
  picture?: string;
  email_verified?: boolean;
  sub?: string;
};

export type UserProfile = {
  id: string;
  auth0UserId: string;
  email: string;
  name: string;
  emailVerified: boolean;
  picture?: string;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
  loginCount: number;
};

export type CreateUserPayload = {
  auth0UserId: string;
  email: string;
  name: string;
  emailVerified?: boolean;
};