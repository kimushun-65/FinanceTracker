export const accountEndpoints = {
  list: '/api/accounts',
  get: (id: string) => `/api/accounts/${id}` as const,
  create: '/api/accounts',
  update: (id: string) => `/api/accounts/${id}` as const,
  delete: (id: string) => `/api/accounts/${id}` as const,
  getBalance: (id: string) => `/api/accounts/${id}/balance` as const,
  getTransactions: (id: string) => `/api/accounts/${id}/transactions` as const,
} as const;

export type AccountEndpoints = typeof accountEndpoints;
