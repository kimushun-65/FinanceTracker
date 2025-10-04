'use client';

import { Button } from '@/shared/ui';
import type { Account } from '@/entities/account';

interface DeleteAccountModalProps {
  account: Account | null;
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  isDeleting?: boolean;
}

export function DeleteAccountModal({
  account,
  isOpen,
  onClose,
  onConfirm,
  isDeleting = false,
}: DeleteAccountModalProps) {
  if (!isOpen || !account) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center'>
      <div className='fixed inset-0 bg-black bg-opacity-50' onClick={onClose} />
      <div className='relative w-full max-w-md rounded-lg bg-white p-6 shadow-xl'>
        <h3 className='mb-4 text-lg font-semibold'>Delete Account</h3>
        <p className='mb-6 text-gray-600'>
          Are you sure you want to delete &quot;{account.name}&quot;?
        </p>
        <div className='flex justify-end gap-3'>
          <Button variant='outline' onClick={onClose}>
            Cancel
          </Button>
          <Button
            variant='destructive'
            onClick={onConfirm}
            disabled={isDeleting}
          >
            {isDeleting ? 'Deleting...' : 'Delete'}
          </Button>
        </div>
      </div>
    </div>
  );
}
