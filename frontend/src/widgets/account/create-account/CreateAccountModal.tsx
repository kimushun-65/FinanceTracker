'use client';

import { useForm } from 'react-hook-form';
import { Button, Input, Label, Select, useToast } from '@/shared/ui';
import {
  ACCOUNT_TYPE_LABELS,
  type CreateAccountPayload,
  type AccountType,
} from '@/entities/account';
import { useCreateAccount } from '@/features/account-management';
import { createMoney } from '@/shared/value-objects/money';

type FormData = {
  name: string;
  accountType: AccountType;
  initialBalance: number;
};

interface CreateAccountModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CreateAccountModal({
  isOpen,
  onClose,
}: CreateAccountModalProps) {
  const { toast } = useToast();
  const createAccountMutation = useCreateAccount();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    setValue,
    watch,
  } = useForm<FormData>({
    mode: 'onChange',
    defaultValues: {
      name: '',
      accountType: 'checking',
      initialBalance: 0,
    },
  });

  const accountType = watch('accountType');

  const onSubmit = async (data: FormData) => {
    try {
      const payload: CreateAccountPayload = {
        name: data.name,
        accountType: data.accountType,
        initialBalance: createMoney(data.initialBalance || 0, 'JPY'),
      };
      await createAccountMutation.mutateAsync(payload);
      toast({
        title: 'Success',
        description: 'Account created successfully',
      });
      reset();
      onClose();
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to create account',
        variant: 'destructive',
      });
    }
  };

  const handleClose = () => {
    reset({
      name: '',
      accountType: 'checking',
      initialBalance: 0,
    });
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center'>
      <div
        className='fixed inset-0 bg-black bg-opacity-50'
        onClick={handleClose}
      />

      <div className='relative w-full max-w-md rounded-lg bg-white shadow-xl'>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className='flex items-center justify-between border-b p-6'>
            <h2 className='text-lg font-semibold'>Add Account</h2>
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
              <Select
                value={accountType}
                onValueChange={(value) =>
                  setValue('accountType', value as AccountType)
                }
                options={Object.entries(ACCOUNT_TYPE_LABELS).map(
                  ([value, label]) => ({
                    value,
                    label,
                  }),
                )}
                placeholder='Select account type'
              />
              {errors.accountType && (
                <p className='mt-1 text-sm text-red-500'>
                  {errors.accountType.message}
                </p>
              )}
            </div>

            <div>
              <Label htmlFor='initialBalance'>Initial Balance</Label>
              <div className='relative'>
                <span className='absolute left-3 top-1/2 -translate-y-1/2 text-gray-500'>
                  ¥
                </span>
                <Input
                  id='initialBalance'
                  type='number'
                  placeholder='0'
                  className='pl-8'
                  {...register('initialBalance', {
                    valueAsNumber: true,
                  })}
                />
              </div>
              {errors.initialBalance && (
                <p className='mt-1 text-sm text-red-500'>
                  {errors.initialBalance.message}
                </p>
              )}
            </div>
          </div>

          <div className='flex justify-end gap-3 border-t bg-gray-50 p-6'>
            <Button type='button' variant='outline' onClick={handleClose}>
              Cancel
            </Button>
            <Button type='submit' disabled={isSubmitting}>
              {isSubmitting ? 'Creating...' : 'Add'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
