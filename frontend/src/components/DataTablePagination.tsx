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

export function DataTablePagination({
	pageIndex,
	pageSize,
	totalCount,
	setPageIndex,
	setPageSize,
	siblingCount = 1,
}: DataTablePaginationProps) {
	const currentPage = pageIndex + 1;
	const totalPageCount = Math.ceil(totalCount / pageSize);

	const paginationRange = usePagination({
		currentPage,
		totalCount,
		siblingCount,
		pageSize,
	});

	const onNext = () => {
		if (currentPage < totalPageCount) {
			setPageIndex(pageIndex + 1);
		}
	};

	const onPrevious = () => {
		if (currentPage > 1) {
			setPageIndex(pageIndex - 1);
		}
	};

	if (totalPageCount <= 1) {
		return null; // Do not render pagination if there's only one page or less
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
							<PaginationItem key={pageNumber}>
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
