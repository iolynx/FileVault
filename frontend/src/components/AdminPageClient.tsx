'use client';

import { useState, useEffect, useCallback } from 'react';
import api from '@/lib/axios';
import { AdminFileItem } from '@/types/AdminFileItem';
import AdminFilesTable from '@/components/AdminFilesTable';
import { SortConfig } from '@/types/Sort';
import { useDebounce } from '@/hooks/useDebounce';
import { DataTablePagination } from '@/components/DataTablePagination';

export function AdminPageClient() {
	const [files, setFiles] = useState<AdminFileItem[]>([]);
	const [isLoading, setIsLoading] = useState(true);
	const [totalCount, setTotalCount] = useState(0);
	const [sortConfig, setSortConfig] = useState<SortConfig>({ key: 'uploaded_at', direction: 'desc' });
	const [pagination, setPagination] = useState({
		pageIndex: 0,
		pageSize: 10,
	});


	// Debounce the sort config to prevent rapid API calls
	const debouncedSortConfig = useDebounce(sortConfig, 300);

	const fetchAdminFiles = useCallback(async () => {
		setIsLoading(true);
		try {
			const params = new URLSearchParams({
				limit: String(pagination.pageSize),
				offset: String(pagination.pageIndex * pagination.pageSize),
				sort_by: debouncedSortConfig.key,
				sort_order: debouncedSortConfig.direction,
			});
			const response = await api.get(`/admin/files?${params.toString()}`);
			setFiles(response.data.data);
			setTotalCount(response.data.totalCount);
		} catch (error) {
			console.error("Failed to fetch admin files:", error);
			// Add toast notification for error
		} finally {
			setIsLoading(false);
		}
	}, [debouncedSortConfig, pagination]);

	useEffect(() => {
		fetchAdminFiles();
	}, [fetchAdminFiles]);

	const handleSort = (key: string) => {
		setSortConfig(prevConfig => {
			const direction = prevConfig.key === key && prevConfig.direction === 'asc' ? 'desc' : 'asc';
			return { key, direction };
		});
	};

	return (
		<div>
			<AdminFilesTable
				files={files}
				sortConfig={sortConfig}
				onSort={handleSort}
				onDataChange={fetchAdminFiles}
			/>

			<DataTablePagination
				pageIndex={pagination.pageIndex}
				pageSize={pagination.pageSize}
				setPageSize={(size) => setPagination({ pageIndex: 0, pageSize: size })}
				setPageIndex={(page) => setPagination(p => ({ ...p, pageIndex: page }))}
				totalCount={totalCount}
			/>
		</div>
	);
}
