import { useState } from 'react';
import {
  useAccounts,
  useAccountAggregates,
  useDeleteAccount,
} from '@/features/account-management';
import { TotalAssetsWidget } from '@/widgets/account/total-assets/ui/TotalAssetsWidget';
import { AccountListTable } from '@/widgets/account/account-list/ui/AccountListTable';
import { CreateAccountModal } from '@/widgets/account/create-account/CreateAccountModal';
import { EditAccountModal } from '@/widgets/account/edit-account/EditAccountModal';
import { DeleteAccountModal } from '@/widgets/account/delete-account';
import { Button, Loading, useToast } from '@/shared/ui';
import type { Account } from '@/entities/account';

export const AccountContainer = () => {
  const { data: accounts, isLoading } = useAccounts();
  const { totalAssets, accountCount } = useAccountAggregates(accounts);
  const deleteAccountMutation = useDeleteAccount();
  const { toast } = useToast();

  const [editingAccount, setEditingAccount] = useState<Account | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [deletingAccount, setDeletingAccount] = useState<Account | null>(null);

  const handleEditAccount = (account: Account) => {
    setEditingAccount(account);
  };

  const handleDeleteAccount = (account: Account) => {
    setDeletingAccount(account);
  };

  const confirmDelete = async () => {
    if (!deletingAccount) return;

    try {
      await deleteAccountMutation.mutateAsync(deletingAccount.id);
      toast({
        title: 'Success',
        description: 'Account deleted successfully',
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete account',
        variant: 'destructive',
      });
    } finally {
      setDeletingAccount(null);
    }
  };

  const handleCreateAccount = () => {
    setShowCreateModal(true);
  };

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div className='mx-auto max-w-7xl space-y-6'>
      {/* ヘッダー */}
      <div className='flex items-center justify-between'>
        <h1 className='text-2xl font-bold text-gray-900'></h1>
        <Button onClick={handleCreateAccount}>+ Add Account</Button>
      </div>

      {/* 総資産ウィジェット */}
      <TotalAssetsWidget
        totalAssets={totalAssets}
        accountCount={accountCount}
      />

      {/* 口座一覧テーブル */}
      <AccountListTable
        accounts={accounts || []}
        totalAssets={totalAssets}
        onEditAccount={handleEditAccount}
        onDeleteAccount={handleDeleteAccount}
      />

      {/* Modals */}
      <CreateAccountModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
      />

      <EditAccountModal
        account={editingAccount}
        isOpen={!!editingAccount}
        onClose={() => setEditingAccount(null)}
      />

      {/* Delete Confirmation Modal */}
      <DeleteAccountModal
        account={deletingAccount}
        isOpen={!!deletingAccount}
        onClose={() => setDeletingAccount(null)}
        onConfirm={confirmDelete}
        isDeleting={deleteAccountMutation.isPending}
      />
    </div>
  );
};
