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

// バックエンドAPIのレスポンス型（snake_case、フラット構造）
interface BackendCategoryResponse {
  id: string;
  user_id: string;
  category_master_id: string;
  name: string;
  custom_name?: string;
  is_active: boolean;
  icon: string;
  color?: string;
  category_type: string;
  display_order: number;
  created_at: string;
  updated_at: string;
}

interface BackendCategoryListResponse {
  categories: BackendCategoryResponse[];
  total_count: number;
}

// バックエンドのフラット構造をフロントエンドのネスト構造に変換
const transformBackendCategory = (
  backendCategory: BackendCategoryResponse,
): CategoryWithMaster => {
  return {
    id: backendCategory.id,
    userId: backendCategory.user_id,
    categoryMasterId: backendCategory.category_master_id,
    customName: backendCategory.custom_name,
    isActive: backendCategory.is_active,
    createdAt: backendCategory.created_at,
    updatedAt: backendCategory.updated_at,
    master: {
      id: backendCategory.category_master_id,
      name: backendCategory.name,
      type: backendCategory.category_type as CategoryType,
      icon: backendCategory.icon,
      color: backendCategory.color || '',
      displayOrder: backendCategory.display_order,
      createdAt: backendCategory.created_at,
      updatedAt: backendCategory.updated_at,
    },
    displayName: backendCategory.custom_name || backendCategory.name,
  };
};

export const categoryApi = {
  list: async (type?: CategoryType): Promise<CategoryListResponse> => {
    const url = type
      ? `${endpoints.categories.list}?type=${type}`
      : endpoints.categories.list;
    const response = await apiClient.get<BackendCategoryListResponse>(url);

    // バックエンドのフラット構造をネスト構造に変換
    const categories = response.data.categories.map(transformBackendCategory);

    return {
      categories,
      total: response.data.total_count,
    };
  },

  get: async (id: string): Promise<CategoryWithMaster> => {
    const response = await apiClient.get<BackendCategoryResponse>(
      endpoints.categories.get(id),
    );
    return transformBackendCategory(response.data);
  },

  create: async (
    payload: CreateCategoryPayload,
  ): Promise<CategoryWithMaster> => {
    // フロントエンドのcamelCaseをバックエンドのsnake_caseに変換
    const backendPayload = {
      category_master_id: payload.categoryMasterId,
      custom_name: payload.customName,
    };

    const response = await apiClient.post<BackendCategoryResponse>(
      endpoints.categories.create,
      backendPayload,
    );
    return transformBackendCategory(response.data);
  },

  update: async (
    id: string,
    payload: UpdateCategoryPayload,
  ): Promise<CategoryWithMaster> => {
    // フロントエンドのcamelCaseをバックエンドのsnake_caseに変換
    const backendPayload = {
      custom_name: payload.customName,
      is_active: payload.isActive,
    };

    const response = await apiClient.put<BackendCategoryResponse>(
      endpoints.categories.update(id),
      backendPayload,
    );
    return transformBackendCategory(response.data);
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
