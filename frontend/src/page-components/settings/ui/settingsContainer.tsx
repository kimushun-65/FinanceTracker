'use client';

import { AppLayout } from '@/widgets';
import { Loading } from '@/shared/ui';
import { useUserProfile, useCategories } from '@/features';
import { CategorySelectionSection, ProfileSection } from '@/widgets/settings';

export default function SettingsContainer() {
  const { data: userProfile, isLoading: profileLoading } = useUserProfile();
  const { data: categories, isLoading: categoriesLoading } = useCategories();

  if (profileLoading || categoriesLoading) {
    return <Loading />;
  }

  if (!userProfile || !categories) {
    return <Loading />;
  }

  return (
    <AppLayout>
      <div className='container mx-auto px-4 py-6'>
        {/* Page Header */}

        <div className='space-y-6'>
          {/* Profile Section */}
          <ProfileSection userProfile={userProfile} />

          {/* Category Selection Section */}
          <CategorySelectionSection categories={categories} />
        </div>
      </div>
    </AppLayout>
  );
}
