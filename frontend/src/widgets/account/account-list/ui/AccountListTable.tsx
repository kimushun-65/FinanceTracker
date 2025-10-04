import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
  Button,
} from '@/shared/ui';
import { AccountTypeBadge, AccountBalanceDisplay } from '@/entities/account/ui';
import type { Account } from '@/entities/account';
import type { Money } from '@/shared/value-objects';

interface AccountListTableProps {
  accounts: Account[];
  totalAssets: Money;
  onEditAccount?: (account: Account) => void;
  onDeleteAccount?: (account: Account) => void;
}

export const AccountListTable = ({
  accounts,
  totalAssets,
  onEditAccount,
  onDeleteAccount,
}: AccountListTableProps) => {
  const calculatePercentage = (balance: Money): number => {
    if (totalAssets.amount === 0) return 0;
    return (balance.amount / totalAssets.amount) * 100;
  };

  return (
    <div className='rounded-lg bg-white shadow'>
      <div className='border-b border-gray-200 px-6 py-4'>
        <h3 className='text-lg font-medium text-gray-900'>Account List</h3>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Account Name</TableHead>
            <TableHead>Type</TableHead>
            <TableHead className='text-right'>Balance</TableHead>
            <TableHead className='text-right'>Allocation</TableHead>
            <TableHead className='text-right'>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {accounts.map((account) => (
            <TableRow key={account.id}>
              <TableCell className='font-medium'>{account.name}</TableCell>
              <TableCell>
                <AccountTypeBadge type={account.accountType} />
              </TableCell>
              <TableCell className='text-right'>
                <AccountBalanceDisplay
                  balance={account.balance.current}
                  status={account.balance.status}
                  accountType={account.accountType}
                />
              </TableCell>
              <TableCell className='text-right'>
                {calculatePercentage(account.balance.current).toFixed(1)}%
              </TableCell>
              <TableCell className='text-right'>
                <div className='flex justify-end space-x-2'>
                  <Button
                    variant='outline'
                    size='sm'
                    onClick={() => onEditAccount?.(account)}
                  >
                    Edit
                  </Button>
                  <Button
                    variant='outline'
                    size='sm'
                    onClick={() => onDeleteAccount?.(account)}
                  >
                    Delete
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};
