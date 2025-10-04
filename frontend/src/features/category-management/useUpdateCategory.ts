import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  categoryApi,
  categoryKeys,
  type CategoryWithMaster,
  type UpdateCategoryPayload,
} from '@/entities/category';

export const useUpdateCategory = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, ...payload }: { id: string } & UpdateCategoryPayload) =>
      categoryApi.update(id, payload),

    onMutate: async ({ id, ...newData }) => {
      await queryClient.cancelQueries({ queryKey: categoryKeys.detail(id) });

      const previousCategory = queryClient.getQueryData<CategoryWithMaster>(
        categoryKeys.detail(id),
      );

      queryClient.setQueryData<CategoryWithMaster>(
        categoryKeys.detail(id),
        (old) => ({
          ...old!,
          ...newData,
          updatedAt: new Date().toISOString(),
        }),
      );

      return { previousCategory };
    },

    onError: (_, { id }, context) => {
      if (context?.previousCategory) {
        queryClient.setQueryData(
          categoryKeys.detail(id),
          context.previousCategory,
        );
      }
    },

    onSuccess: (data, { id }) => {
      queryClient.setQueryData(categoryKeys.detail(id), data);
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
    },
  });
};
