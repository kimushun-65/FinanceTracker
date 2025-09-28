import { useMutation, useQueryClient } from '@tanstack/react-query';
import { categoryApi, categoryKeys, type CreateCategoryPayload } from '@/entities/category';

export const useCreateCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: categoryApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
    },
  });
};