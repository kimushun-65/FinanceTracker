import type { FC } from 'react';
import type { CategoryWithMaster } from '../../model';
import { getCategoryColor, getCategoryDisplayName } from '../../lib';

export type CategoryColorBadgeProps = {
  category: CategoryWithMaster;
  size?: 'sm' | 'default' | 'lg';
  showName?: boolean;
  className?: string;
};

export const CategoryColorBadge: FC<CategoryColorBadgeProps> = ({
  category,
  size = 'default',
  showName = true,
  className = '',
}) => {
  const color = getCategoryColor(category);
  const displayName = getCategoryDisplayName(category);

  const getSizeClasses = (badgeSize: string) => {
    const sizeMap = {
      sm: 'h-3 w-3 text-xs',
      default: 'h-4 w-4 text-sm',
      lg: 'h-5 w-5 text-base',
    };
    return sizeMap[badgeSize as keyof typeof sizeMap] || sizeMap.default;
  };

  return (
    <div className={`inline-flex items-center gap-2 ${className}`}>
      <div
        className={`rounded-full border border-gray-200 ${getSizeClasses(size)}`}
        style={{ backgroundColor: color }}
        title={displayName}
      />
      {showName && (
        <span
          className={`${size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-base' : 'text-sm'}`}
        >
          {displayName}
        </span>
      )}
    </div>
  );
};
