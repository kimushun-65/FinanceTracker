import { useMutation, useQueryClient } from '@tanstack/react-query';
import { categoryApi, categoryKeys } from '@/entities/category';

export const useDeleteCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: categoryApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
    },
  });
};
