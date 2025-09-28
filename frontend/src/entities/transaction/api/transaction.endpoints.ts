export const endpoints = {
  list: '/api/transactions',
  get: (id: string) => `/api/transactions/${id}`,
  create: '/api/transactions',
  update: (id: string) => `/api/transactions/${id}`,
  delete: (id: string) => `/api/transactions/${id}`,
} as const;
