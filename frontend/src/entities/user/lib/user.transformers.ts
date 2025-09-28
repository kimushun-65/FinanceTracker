import type { User, AuthUser, UserProfile } from '../model';
import { createEmail } from '@/shared/value-objects';

export const transformAuthUserToUser = (authUser: AuthUser): Partial<User> => {
  return {
    auth0UserId: authUser.sub || '',
    email: authUser.email ? createEmail(authUser.email) : undefined,
    name: authUser.name || '',
    emailVerified: authUser.email_verified || false,
  };
};

export const transformApiUserResponse = (apiUser: any): User => {
  return {
    id: apiUser.id,
    createdAt: apiUser.createdAt,
    updatedAt: apiUser.updatedAt,
    auth0UserId: apiUser.auth0UserId,
    email: createEmail(apiUser.email),
    name: apiUser.name,
    emailVerified: apiUser.emailVerified,
  };
};

export const transformApiUserProfileResponse = (
  apiUserProfile: any,
): UserProfile => {
  return {
    id: apiUserProfile.id,
    createdAt: apiUserProfile.createdAt,
    updatedAt: apiUserProfile.updatedAt,
    auth0UserId: apiUserProfile.auth0UserId,
    email: apiUserProfile.email,
    name: apiUserProfile.name,
    emailVerified: apiUserProfile.emailVerified,
    picture: apiUserProfile.picture,
    lastLoginAt: apiUserProfile.lastLoginAt,
    loginCount: apiUserProfile.loginCount || 0,
  };
};

export const getUserDisplayName = (
  user: User | UserProfile | AuthUser,
): string => {
  if ('name' in user && user.name) {
    return user.name;
  }

  if ('email' in user && user.email) {
    // Email型はstring & brandなので、stringとして扱える
    const emailString = String(user.email);
    return emailString.split('@')[0];
  }

  return 'ユーザー';
};

export const getUserInitials = (
  user: User | UserProfile | AuthUser,
): string => {
  const displayName = getUserDisplayName(user);

  if (displayName === 'ユーザー') {
    return 'U';
  }

  const names = displayName.split(' ');
  if (names.length >= 2) {
    return (names[0][0] + names[1][0]).toUpperCase();
  }

  return displayName.slice(0, 2).toUpperCase();
};
