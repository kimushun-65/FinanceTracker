'use client';

import { useState } from 'react';
import {
  Modal,
  ModalTrigger,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Input,
  Label,
  useModalContext,
} from '@/shared/ui';
import { useUpdateProfile } from '@/features';
import { useToast } from '@/shared/lib/hooks/use-toast';
import type { UserProfile } from '@/entities/user';

interface ProfileEditModalProps {
  trigger?: React.ReactNode;
  userProfile: UserProfile;
  onSuccess?: () => void;
}

function ProfileEditForm({
  userProfile,
  onSuccess,
}: {
  userProfile: UserProfile;
  onSuccess?: () => void;
}) {
  const { closeModal } = useModalContext();
  const updateProfileMutation = useUpdateProfile();
  const { toast } = useToast();

  const [formData, setFormData] = useState({
    name: userProfile.name,
    email: userProfile.email,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name || formData.name.trim().length === 0) {
      newErrors.name = 'Name is required';
    } else if (formData.name.length > 100) {
      newErrors.name = 'Name must be less than 100 characters';
    }

    if (!formData.email || formData.email.trim().length === 0) {
      newErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Invalid email format';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    // Check if anything changed
    const hasChanges =
      formData.name !== userProfile.name ||
      formData.email !== userProfile.email;

    if (!hasChanges) {
      closeModal();
      return;
    }

    try {
      await updateProfileMutation.mutateAsync({
        name: formData.name !== userProfile.name ? formData.name : undefined,
        email:
          formData.email !== userProfile.email ? formData.email : undefined,
      });

      toast({
        title: 'Success',
        description: 'Profile updated successfully',
      });

      closeModal();
      onSuccess?.();
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to update profile',
        variant: 'destructive',
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className='space-y-4'>
      {/* Name */}
      <div>
        <Label htmlFor='name'>Name</Label>
        <Input
          id='name'
          type='text'
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder='Enter your name'
          className={errors.name ? 'border-red-500' : ''}
        />
        {errors.name && (
          <p className='mt-1 text-sm text-red-600'>{errors.name}</p>
        )}
      </div>

      {/* Email */}
      <div>
        <Label htmlFor='email'>Email</Label>
        <Input
          id='email'
          type='email'
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder='Enter your email'
          className={errors.email ? 'border-red-500' : ''}
        />
        {errors.email && (
          <p className='mt-1 text-sm text-red-600'>{errors.email}</p>
        )}
      </div>

      <ModalFooter>
        <Button
          type='button'
          variant='outline'
          onClick={closeModal}
          disabled={updateProfileMutation.isPending}
        >
          Cancel
        </Button>
        <Button type='submit' disabled={updateProfileMutation.isPending}>
          {updateProfileMutation.isPending ? 'Saving...' : 'Save Changes'}
        </Button>
      </ModalFooter>
    </form>
  );
}

export function ProfileEditModal({
  trigger,
  userProfile,
  onSuccess,
}: ProfileEditModalProps) {
  return (
    <Modal>
      <ModalTrigger asChild>
        {trigger || <Button variant='outline'>Edit Profile</Button>}
      </ModalTrigger>
      <ModalContent>
        <ModalHeader>Edit Profile</ModalHeader>
        <ModalBody>
          <ProfileEditForm userProfile={userProfile} onSuccess={onSuccess} />
        </ModalBody>
      </ModalContent>
    </Modal>
  );
}
