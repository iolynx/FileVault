'use-client';

import {
	Pagination,
	PaginationContent,
	PaginationEllipsis,
	PaginationItem,
	PaginationLink,
	PaginationNext,
	PaginationPrevious,
} from '@/components/ui/pagination';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { usePagination, DOTS } from '@/hooks/usePagination';
import { cn } from '@/lib/utils';

interface DataTablePaginationProps {
	pageIndex: number;
	pageSize: number;
	totalCount: number;
	setPageIndex: (index: number) => void;
	setPageSize: (size: number) => void;
	siblingCount?: number;
}

/**
 * Props for DataTablePagination
 * 
 * @typedef {Object} DataTablePaginationProps
 * @property {number} pageIndex - Current zero-based page index
 * @property {number} pageSize - Number of items per page
 * @property {number} totalCount - Total number of items in the dataset
 * @property {(index: number) => void} setPageIndex - Callback to update the current page index
 * @property {(size: number) => void} setPageSize - Callback to update the number of items per page
 * @property {number} [siblingCount=1] - Number of sibling page numbers to show on each side of the current page
 */

/**
 * DataTablePagination component
 * 
 * Renders pagination controls for a table of data.
 * Supports:
 * - Navigating to the previous/next page
 * - Dynamically computing page numbers using `usePagination`
 * - Conditionally rendering pagination only when total pages > 1
 * 
 * @param {DataTablePaginationProps} props - Component props
 * @returns {JSX.Element | null} JSX element rendering the pagination controls, or null if pagination is unnecessary
 * 
 * @component
 */
export function DataTablePagination({
	pageIndex,
	pageSize,
	totalCount,
	setPageIndex,
	setPageSize,
	siblingCount = 1,
}: DataTablePaginationProps) {
	/** Current 1-based page number */
	const currentPage = pageIndex + 1;
	const totalPageCount = Math.ceil(totalCount / pageSize);

	const paginationRange = usePagination({
		currentPage,
		totalCount,
		siblingCount,
		pageSize,
	});

	/**
	 * Navigate to the next page
	 */
	const onNext = () => {
		if (currentPage < totalPageCount) {
			setPageIndex(pageIndex + 1);
		}
	};

	/**
	 * Navigate to the previous page
	 */
	const onPrevious = () => {
		if (currentPage > 1) {
			setPageIndex(pageIndex - 1);
		}
	};

	// Do not render pagination if the number of records 
	// are lesser than the smallest pageSize option (10)
	if (Math.ceil(totalCount / 10) <= 1) {
		return null;
	}

	return (
		<div className="flex items-center justify-between px-3 py-2 mt-4">
			<div className="flex items-center space-x-2">
				<p className="text-sm font-medium text-nowrap">Rows per page</p>
				<Select
					value={`${pageSize}`}
					onValueChange={(value) => setPageSize(Number(value))}
				>
					<SelectTrigger className="h-8 w-[70px]">
						<SelectValue placeholder={pageSize} />
					</SelectTrigger>
					<SelectContent side="top">
						{[10, 20, 30, 40, 50].map((size) => (
							<SelectItem key={size} value={`${size}`}>
								{size}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>

			<Pagination>
				<PaginationContent>
					<PaginationItem>
						<PaginationPrevious
							onClick={onPrevious}
							className={cn({
								'cursor-not-allowed text-muted-foreground': currentPage === 1,
								'cursor-pointer': currentPage != 1
							})}
						/>
					</PaginationItem>
					{paginationRange.map((pageNumber, index) => {
						if (pageNumber === DOTS) {
							return (
								<PaginationItem key={`dots-${index}`}>
									<PaginationEllipsis />
								</PaginationItem>
							);
						}
						return (
							<PaginationItem key={pageNumber} className='cursor-pointer'>
								<PaginationLink
									onClick={() => setPageIndex(Number(pageNumber) - 1)}
									isActive={currentPage === pageNumber}
								>
									{pageNumber}
								</PaginationLink>
							</PaginationItem>
						);
					})}
					<PaginationItem>
						<PaginationNext
							onClick={onNext}
							className={cn({
								'cursor-not-allowed text-muted-foreground': currentPage === totalPageCount,
								'cursor-pointer': currentPage != totalPageCount
							})}
						/>
					</PaginationItem>
				</PaginationContent>
			</Pagination>

			<div className="flex w-[100px] items-center justify-end text-sm font-medium text-nowrap">
				Total {totalCount} rows
			</div>
		</div>
	);
}
