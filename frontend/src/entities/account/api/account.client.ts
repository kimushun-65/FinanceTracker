import { apiClient } from '@/shared/api/client';
import type {
  Account,
  CreateAccountPayload,
  UpdateAccountPayload,
  AccountListResponse,
  BackendAccountResponse,
  BackendAccountListResponse,
} from '../model';
import { accountEndpoints } from './account.endpoints';
import { transformBackendAccountToFrontend } from '../lib/account.transformers';

export const accountApi = {
  list: async (): Promise<AccountListResponse> => {
    const response = await apiClient.get<BackendAccountListResponse>(
      accountEndpoints.list,
    );
    const transformedAccounts = response.data.accounts.map(
      transformBackendAccountToFrontend,
    );
    return {
      accounts: transformedAccounts,
      total: response.data.total_count,
      totalAssets: response.data.total_balance
        ? parseFloat(response.data.total_balance)
        : undefined,
      totalDebt: response.data.total_debt
        ? parseFloat(response.data.total_debt)
        : undefined,
      netWorth: response.data.net_worth
        ? parseFloat(response.data.net_worth)
        : undefined,
    };
  },

  get: async (id: string): Promise<Account> => {
    const response = await apiClient.get<BackendAccountResponse>(
      accountEndpoints.get(id),
    );
    return transformBackendAccountToFrontend(response.data);
  },

  create: async (payload: CreateAccountPayload): Promise<Account> => {
    // Transform the payload to match backend expectations
    const backendPayload = {
      name: payload.name,
      account_type: payload.accountType,
      initial_balance: payload.initialBalance?.amount.toString() || '0',
      currency: payload.initialBalance?.currency || 'JPY',
    };

    const response = await apiClient.post<BackendAccountResponse>(
      accountEndpoints.create,
      backendPayload,
    );
    return transformBackendAccountToFrontend(response.data);
  },

  update: async (
    id: string,
    payload: UpdateAccountPayload,
  ): Promise<Account> => {
    // Transform the payload to match backend expectations
    const backendPayload: any = {};
    if (payload.name !== undefined) {
      backendPayload.name = payload.name;
    }
    if (payload.balance !== undefined) {
      backendPayload.balance = payload.balance.amount.toString();
      backendPayload.currency = payload.balance.currency;
    }

    const response = await apiClient.put<BackendAccountResponse>(
      accountEndpoints.update(id),
      backendPayload,
    );
    return transformBackendAccountToFrontend(response.data);
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(accountEndpoints.delete(id));
  },

  getBalance: async (id: string): Promise<Account['balance']> => {
    const account = await accountApi.get(id);
    return account.balance;
  },
} as const;

export type AccountApi = typeof accountApi;
