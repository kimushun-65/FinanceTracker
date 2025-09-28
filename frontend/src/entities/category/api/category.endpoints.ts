export const endpoints = {
  categories: {
    list: '/categories',
    get: (id: string) => `/categories/${id}`,
    create: '/categories',
    update: (id: string) => `/categories/${id}`,
    delete: (id: string) => `/categories/${id}`,
  },
  categoryMasters: {
    list: '/categories/master',
    get: (id: string) => `/categories/master/${id}`,
  },
} as const;
