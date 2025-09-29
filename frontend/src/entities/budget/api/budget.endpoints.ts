export const endpoints = {
  list: '/budgets',
  get: (id: string) => `/budgets/${id}`,
  create: '/budgets',
  update: (id: string) => `/budgets/${id}`,
  delete: (id: string) => `/budgets/${id}`,
  current: '/budgets/current',
  summary: '/budgets/summary',
} as const;
