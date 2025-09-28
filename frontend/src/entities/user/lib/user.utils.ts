import type { User, UserProfile } from '../model';

export const isEmailVerified = (user: User | UserProfile): boolean => {
  return user.emailVerified;
};

export const getUserAge = (user: User | UserProfile): number | null => {
  // If we had birthDate field, we could calculate age here
  // For now, return null as birthDate is not in our current schema
  return null;
};

export const formatLastLoginAt = (lastLoginAt?: string): string => {
  if (!lastLoginAt) {
    return '未ログイン';
  }

  const date = new Date(lastLoginAt);
  const now = new Date();
  const diffInDays = Math.floor(
    (now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (diffInDays === 0) {
    return '今日';
  } else if (diffInDays === 1) {
    return '昨日';
  } else if (diffInDays < 7) {
    return `${diffInDays}日前`;
  } else if (diffInDays < 30) {
    return `${Math.floor(diffInDays / 7)}週間前`;
  } else if (diffInDays < 365) {
    return `${Math.floor(diffInDays / 30)}ヶ月前`;
  } else {
    return `${Math.floor(diffInDays / 365)}年前`;
  }
};

export const getLoginFrequency = (
  loginCount: number,
  createdAt: string,
): string => {
  const created = new Date(createdAt);
  const now = new Date();
  const daysSinceCreation = Math.floor(
    (now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24),
  );

  if (daysSinceCreation === 0) {
    return '初回ログイン';
  }

  const frequency = loginCount / daysSinceCreation;

  if (frequency >= 1) {
    return '毎日';
  } else if (frequency >= 0.5) {
    return '頻繁';
  } else if (frequency >= 0.2) {
    return '定期的';
  } else {
    return '稀に';
  }
};

export const isNewUser = (user: UserProfile): boolean => {
  const created = new Date(user.createdAt);
  const now = new Date();
  const daysSinceCreation = Math.floor(
    (now.getTime() - created.getTime()) / (1000 * 60 * 60 * 24),
  );

  return daysSinceCreation <= 7; // 7日以内に作成されたユーザーを新規とみなす
};

export const maskEmail = (email: string): string => {
  const [localPart, domain] = email.split('@');
  if (localPart.length <= 2) {
    return `${localPart}***@${domain}`;
  }
  return `${localPart.slice(0, 2)}***@${domain}`;
};
