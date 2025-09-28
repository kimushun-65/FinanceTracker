import type { FC } from 'react';
import type { CategoryWithMaster } from '../../model';
import { getCategoryIcon } from '../../lib';

export type CategoryIconProps = {
  category: CategoryWithMaster;
  size?: number;
  className?: string;
};

export const CategoryIcon: FC<CategoryIconProps> = ({
  category,
  size = 20,
  className = '',
}) => {
  const icon = getCategoryIcon(category);

  return (
    <span
      className={`inline-flex items-center justify-center ${className}`}
      style={{ fontSize: `${size}px` }}
      title={category.displayName}
    >
      {icon}
    </span>
  );
};
