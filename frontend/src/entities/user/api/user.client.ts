import { apiClient } from '@/shared/api/client';
import type {
  User,
  UserProfile,
  UpdateUserPayload,
  CreateUserPayload,
} from '../model';
import { endpoints } from './user.endpoints';

export const userApi = {
  getProfile: async (): Promise<UserProfile> => {
    const response = await apiClient.get<UserProfile>(endpoints.profile);
    return response.data;
  },

  updateProfile: async (payload: UpdateUserPayload): Promise<User> => {
    const response = await apiClient.put<User>(endpoints.update, payload);
    return response.data;
  },

  create: async (payload: CreateUserPayload): Promise<User> => {
    const response = await apiClient.post<User>(endpoints.create, payload);
    return response.data;
  },

  deleteProfile: async (): Promise<void> => {
    await apiClient.delete(endpoints.delete);
  },
};
