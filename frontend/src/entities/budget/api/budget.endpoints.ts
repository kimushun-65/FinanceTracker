export const endpoints = {
  list: '/api/budgets',
  get: (id: string) => `/api/budgets/${id}`,
  create: '/api/budgets',
  update: (id: string) => `/api/budgets/${id}`,
  delete: (id: string) => `/api/budgets/${id}`,
  summary: '/api/budgets/summary',
} as const;
