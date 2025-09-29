import { useState } from 'react';
import {
  useAccounts,
  useAccountAggregates,
} from '@/features/account-management';
import { TotalAssetsWidget } from '@/widgets/account/total-assets/ui/TotalAssetsWidget';
import { AccountListTable } from '@/widgets/account/account-list/ui/AccountListTable';
import { Button, Loading } from '@/shared/ui';
import type { Account } from '@/entities/account';

export const AccountContainer = () => {
  const { data: accounts, isLoading, error } = useAccounts();
  const { totalAssets, accountCount } = useAccountAggregates(accounts);

  const [editingAccount, setEditingAccount] = useState<Account | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);

  const handleEditAccount = (account: Account) => {
    setEditingAccount(account);
  };

  const handleDeleteAccount = (account: Account) => {
    // TODO: 削除確認ダイアログを実装
    console.log('Delete account:', account.id);
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

      {/* TODO: モーダル実装 */}
      {showCreateModal && <div>Create Account Modal (TODO)</div>}

      {editingAccount && <div>Edit Account Modal (TODO)</div>}
    </div>
  );
};
