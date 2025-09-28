import { useMutation, useQueryClient } from '@tanstack/react-query';
import { userApi, userKeys } from '@/entities/user';

export const useUpdateProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: userApi.updateProfile,
    onSuccess: (data) => {
      queryClient.setQueryData(userKeys.profile(), data);
    },
  });
};
