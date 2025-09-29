import { cn } from '@/shared/utils';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  className,
}: PaginationProps) {
  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  return (
    <div className={cn('flex items-center justify-center gap-2', className)}>
      <button
        onClick={handlePrevious}
        disabled={currentPage <= 1}
        className={cn(
          'rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50',
          'disabled:cursor-not-allowed disabled:opacity-50',
        )}
      >
        ◀ Previous
      </button>

      <span className='px-3 py-2 text-sm text-gray-700'>
        Page {currentPage} of {totalPages}
      </span>

      <button
        onClick={handleNext}
        disabled={currentPage >= totalPages}
        className={cn(
          'rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50',
          'disabled:cursor-not-allowed disabled:opacity-50',
        )}
      >
        Next ▶
      </button>
    </div>
  );
}
