export const accountEndpoints = {
  list: '/accounts',
  get: (id: string) => `/accounts/${id}` as const,
  create: '/accounts',
  update: (id: string) => `/accounts/${id}` as const,
  delete: (id: string) => `/accounts/${id}` as const,
  getBalance: (id: string) => `/accounts/${id}/balance` as const,
  getTransactions: (id: string) => `/accounts/${id}/transactions` as const,
} as const;

export type AccountEndpoints = typeof accountEndpoints;
