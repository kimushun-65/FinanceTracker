import { apiClient } from '@/shared/api/client';
import type {
  Account,
  CreateAccountPayload,
  UpdateAccountPayload,
  AccountListResponse,
} from '../model';
import { accountEndpoints } from './account.endpoints';

export const accountApi = {
  list: async (): Promise<AccountListResponse> => {
    const response = await apiClient.get<AccountListResponse>(
      accountEndpoints.list,
    );
    return response.data;
  },

  get: async (id: string): Promise<Account> => {
    const response = await apiClient.get<Account>(accountEndpoints.get(id));
    return response.data;
  },

  create: async (payload: CreateAccountPayload): Promise<Account> => {
    const response = await apiClient.post<Account>(
      accountEndpoints.create,
      payload,
    );
    return response.data;
  },

  update: async (
    id: string,
    payload: UpdateAccountPayload,
  ): Promise<Account> => {
    const response = await apiClient.put<Account>(
      accountEndpoints.update(id),
      payload,
    );
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(accountEndpoints.delete(id));
  },

  getBalance: async (id: string): Promise<Account['balance']> => {
    const response = await apiClient.get<Account['balance']>(
      accountEndpoints.getBalance(id),
    );
    return response.data;
  },
} as const;

export type AccountApi = typeof accountApi;
