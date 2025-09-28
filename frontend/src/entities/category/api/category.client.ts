import { apiClient } from '@/shared/api/client';
import type {
  Category,
  CategoryMaster,
  CategoryWithMaster,
  CreateCategoryPayload,
  UpdateCategoryPayload,
  CategoryListResponse,
  CategoryMasterListResponse,
  CategoryType,
} from '../model';
import { endpoints } from './category.endpoints';

export const categoryApi = {
  list: async (type?: CategoryType): Promise<CategoryListResponse> => {
    const url = type
      ? `${endpoints.categories.list}?type=${type}`
      : endpoints.categories.list;
    const response = await apiClient.get<CategoryListResponse>(url);
    return response.data;
  },

  get: async (id: string): Promise<CategoryWithMaster> => {
    const response = await apiClient.get<CategoryWithMaster>(
      endpoints.categories.get(id),
    );
    return response.data;
  },

  create: async (payload: CreateCategoryPayload): Promise<Category> => {
    const response = await apiClient.post<Category>(
      endpoints.categories.create,
      payload,
    );
    return response.data;
  },

  update: async (
    id: string,
    payload: UpdateCategoryPayload,
  ): Promise<Category> => {
    const response = await apiClient.put<Category>(
      endpoints.categories.update(id),
      payload,
    );
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(endpoints.categories.delete(id));
  },
};

export const categoryMasterApi = {
  list: async (type?: CategoryType): Promise<CategoryMasterListResponse> => {
    const url = type
      ? `${endpoints.categoryMasters.list}?type=${type}`
      : endpoints.categoryMasters.list;
    const response = await apiClient.get<CategoryMasterListResponse>(url);
    return response.data;
  },

  get: async (id: string): Promise<CategoryMaster> => {
    const response = await apiClient.get<CategoryMaster>(
      endpoints.categoryMasters.get(id),
    );
    return response.data;
  },
};
