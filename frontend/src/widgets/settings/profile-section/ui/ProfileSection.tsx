'use client';

import { Card } from '@/shared/ui';
import type { UserProfile } from '@/entities/user';
import { ProfileEditModal } from './ProfileEditModal';

interface ProfileSectionProps {
  userProfile: UserProfile;
}

export function ProfileSection({ userProfile }: ProfileSectionProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <Card className='p-6'>
      <div className='mb-4 flex items-center justify-between'>
        <h2 className='text-xl font-semibold text-gray-900'>Profile</h2>
        <ProfileEditModal userProfile={userProfile} />
      </div>

      <div className='space-y-4'>
        {/* User Name */}
        <div>
          <label className='text-sm font-medium text-gray-500'>Name</label>
          <p className='mt-1 text-base text-gray-900'>{userProfile.name}</p>
        </div>

        {/* Email */}
        <div>
          <label className='text-sm font-medium text-gray-500'>Email</label>
          <p className='mt-1 text-base text-gray-900'>{userProfile.email}</p>
          {userProfile.emailVerified && (
            <p className='mt-1 text-xs text-green-600'>✓ Verified</p>
          )}
        </div>

        {/* Registration Date */}
        <div>
          <label className='text-sm font-medium text-gray-500'>
            Member Since
          </label>
          <p className='mt-1 text-base text-gray-900'>
            {formatDate(userProfile.createdAt)}
          </p>
        </div>

        {/* Last Login */}
        {userProfile.lastLoginAt && (
          <div>
            <label className='text-sm font-medium text-gray-500'>
              Last Login
            </label>
            <p className='mt-1 text-base text-gray-900'>
              {formatDate(userProfile.lastLoginAt)}
            </p>
          </div>
        )}
      </div>
    </Card>
  );
}
