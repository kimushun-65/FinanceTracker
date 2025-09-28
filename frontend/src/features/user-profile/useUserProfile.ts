import { useQuery } from '@tanstack/react-query';
import { userApi, userKeys } from '@/entities/user';

export const useUserProfile = () => {
  return useQuery({
    queryKey: userKeys.profile(),
    queryFn: () => userApi.getProfile(),
    staleTime: 10 * 60 * 1000, // 10分
  });
};
