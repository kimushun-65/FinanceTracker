export const endpoints = {
  list: '/transactions',
  get: (id: string) => `/transactions/${id}`,
  create: '/transactions',
  update: (id: string) => `/transactions/${id}`,
  delete: (id: string) => `/transactions/${id}`,
  monthlySummary: '/transactions/summary/monthly',
} as const;
