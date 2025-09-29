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
    const transformedAccounts = response.data.accounts.map(transformBackendAccountToFrontend);
    return {
      accounts: transformedAccounts,
      total: response.data.total_count,
    };
  },

  get: async (id: string): Promise<Account> => {
    const response = await apiClient.get<BackendAccountResponse>(accountEndpoints.get(id));
    return transformBackendAccountToFrontend(response.data);
  },

  create: async (payload: CreateAccountPayload): Promise<Account> => {
    const response = await apiClient.post<BackendAccountResponse>(
      accountEndpoints.create,
      payload,
    );
    return transformBackendAccountToFrontend(response.data);
  },

  update: async (
    id: string,
    payload: UpdateAccountPayload,
  ): Promise<Account> => {
    const response = await apiClient.put<BackendAccountResponse>(
      accountEndpoints.update(id),
      payload,
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
