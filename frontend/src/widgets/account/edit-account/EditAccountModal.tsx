'use client';

import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Button, Input, Label, useToast } from '@/shared/ui';
import {
  ACCOUNT_TYPE_LABELS,
  type UpdateAccountPayload,
  type Account,
} from '@/entities/account';
import { createMoney } from '@/shared/value-objects/money';
import { useUpdateAccount } from '@/features/account-management';

interface EditAccountModalProps {
  account: Account | null;
  isOpen: boolean;
  onClose: () => void;
}

export function EditAccountModal({
  account,
  isOpen,
  onClose,
}: EditAccountModalProps) {
  const { toast } = useToast();
  const updateAccountMutation = useUpdateAccount();

  type FormData = {
    name: string;
    currentBalance: number;
  };

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<FormData>();

  useEffect(() => {
    if (account) {
      reset({
        name: account.name,
        currentBalance: account.balance.current.amount,
      });
    }
  }, [account, reset]);

  const onSubmit = async (data: FormData) => {
    if (!account) return;

    try {
      const payload: UpdateAccountPayload = {
        name: data.name,
        // Only send balance if it changed
        ...(data.currentBalance !== account.balance.current.amount && {
          balance: createMoney(
            data.currentBalance,
            account.balance.current.currency,
          ),
        }),
      };

      await updateAccountMutation.mutateAsync({
        id: account.id,
        ...payload,
      });
      toast({
        title: 'Success',
        description: 'Account updated successfully',
      });
      onClose();
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to update account',
        variant: 'destructive',
      });
    }
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  if (!isOpen || !account) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center'>
      <div
        className='fixed inset-0 bg-black bg-opacity-50'
        onClick={handleClose}
      />

      <div className='relative w-full max-w-md rounded-lg bg-white shadow-xl'>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className='flex items-center justify-between border-b p-6'>
            <h2 className='text-lg font-semibold'>Edit Account</h2>
            <button
              type='button'
              onClick={handleClose}
              className='text-2xl font-light text-gray-400 hover:text-gray-600'
            >
              ×
            </button>
          </div>

          <div className='space-y-4 p-6'>
            <div>
              <Label htmlFor='name'>Account Name</Label>
              <Input
                id='name'
                {...register('name')}
                placeholder='Enter account name'
                className={errors.name ? 'border-red-500' : ''}
              />
              {errors.name && (
                <p className='mt-1 text-sm text-red-500'>
                  {errors.name.message}
                </p>
              )}
            </div>

            <div>
              <Label htmlFor='accountType'>Account Type</Label>
              <div className='rounded-md bg-gray-50 p-3 text-lg font-medium'>
                {ACCOUNT_TYPE_LABELS[account.accountType]}
              </div>
              <p className='mt-1 text-sm text-gray-500'>
                Account type cannot be changed
              </p>
            </div>

            <div>
              <Label htmlFor='currentBalance'>Current Balance</Label>
              <div className='relative'>
                <span className='absolute left-3 top-1/2 -translate-y-1/2 text-gray-500'>
                  ¥
                </span>
                <Input
                  id='currentBalance'
                  type='number'
                  step='0.01'
                  className='pl-8'
                  {...register('currentBalance', {
                    valueAsNumber: true,
                    required: 'Current balance is required',
                    ...(account.accountType !== 'credit_card' && {
                      min: { value: 0, message: 'Balance must be positive' },
                    }),
                  })}
                />
              </div>
              {errors.currentBalance && (
                <p className='mt-1 text-sm text-red-500'>
                  {errors.currentBalance.message}
                </p>
              )}
            </div>
          </div>

          <div className='flex justify-end gap-3 border-t bg-gray-50 p-6'>
            <Button type='button' variant='outline' onClick={handleClose}>
              Cancel
            </Button>
            <Button type='submit' disabled={isSubmitting}>
              {isSubmitting ? 'Updating...' : 'Save'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
