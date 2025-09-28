export const endpoints = {
  categories: {
    list: '/api/categories',
    get: (id: string) => `/api/categories/${id}`,
    create: '/api/categories',
    update: (id: string) => `/api/categories/${id}`,
    delete: (id: string) => `/api/categories/${id}`,
  },
  categoryMasters: {
    list: '/api/category-masters',
    get: (id: string) => `/api/category-masters/${id}`,
  },
} as const;
