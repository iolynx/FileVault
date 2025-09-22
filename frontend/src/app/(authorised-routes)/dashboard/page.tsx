"use client";

import FilesTable from "@/components/FilesTable";
import { Card } from "@/components/ui/card";
import { useEffect, useState } from "react";
import { FileUploadMenu } from "@/components/FileUploadMenu";
import { ActiveFilter } from "@/types/Filter";
import { SearchAndFilterComponent } from "@/components/SearchAndFilter";
import { useContentStore } from "@/stores/useContentStore";
import { Breadcrumbs } from "@/components/Breadcrumbs";
import { DropzoneOverlay } from "@/components/DropzoneOverlay";
import { useDragAndDrop } from "@/hooks/useDragAndDrop";
import { useFileUploader } from "@/hooks/useFileUploader";
import { SortConfig } from "@/types/Sort";
import { DataTablePagination } from '@/components/DataTablePagination';
import FilesSkeleton from "@/components/FilesSkeleton";

/**
 * DashboardPage component
 * 
 * Displays the user's files and folders, allows filtering, sorting, pagination, 
 * uploading new files, and drag-and-drop functionality.
 * 
 * Uses several custom hooks:
 * - `useContentStore` for fetching and managing file/folder data
 * - `useFileUploader` for handling file uploads
 * - `useDragAndDrop` for drag-and-drop support
 * 
 * Internal Features:
 * - Filter by column values and file sizes
 * - Sort files by a specific key and direction
 * - Pagination support
 * - Refresh contents after actions
 * - Dropzone overlay when dragging files
 * 
 * @component
 */
const DashboardPage = () => {
	const { contents, path, isLoading, setLoading, totalCount, fetchContents } = useContentStore();

	const [activeFilters, setActiveFilters] = useState<ActiveFilter[]>([]);
	const [sortConfig, setSortConfig] = useState<SortConfig>({ key: 'uploaded_at', direction: 'desc' });
	const [pagination, setPagination] = useState({
		pageIndex: 0,
		pageSize: 10,
	});

	const currentFolder = path[path.length - 1];
	const currentFolderId = currentFolder ? currentFolder.id : null;

	/**
	 * Updates active filters when a user changes a filter for a column.
	 * @param {string} column - The column being filtered
	 * @param {string | Date | undefined} value - The value to filter by
	 */
	const handleFilterChange = (column: string, value: string | Date | undefined) => {
		setLoading(true);
		setActiveFilters(prevFilters => {
			const existingFilterIndex = prevFilters.findIndex(f => f.column === column);

			let stringValue: string;
			if (value instanceof Date) {
				stringValue = value.toISOString();
			} else {
				stringValue = value || '';
			}

			// If the value is empty, remove the filter from the array
			if (stringValue === '') {
				return prevFilters.filter(f => f.column !== column);
			}

			// If the filter already exists, update its value
			if (existingFilterIndex > -1) {
				const newFilters = [...prevFilters];
				newFilters[existingFilterIndex].value = stringValue;
				return newFilters;
			}

			// Otherwise, add the new filter to the array
			return [...prevFilters, { column, value: stringValue }];
		});
	};

	/**
	 * Updates size-related filters
	 * @param {{ min_size: number | null, max_size: number | null }} sizeFilter 
	 */
	const handleSizeFilterChange = (sizeFilter: { min_size: number | null, max_size: number | null }) => {
		setLoading(true);
		setActiveFilters(prevFilters => {
			// Remove any old size filters
			const updatedFilters = prevFilters.filter(
				f => f.column !== 'min_size' && f.column !== 'max_size'
			);

			// If the new size filters aren't null, add to activeFilters
			if (sizeFilter.min_size !== null) {
				updatedFilters.push({ column: 'min_size', value: String(sizeFilter.min_size) });
			}
			if (sizeFilter.max_size !== null) {
				updatedFilters.push({ column: 'max_size', value: String(sizeFilter.max_size) });
			}

			return updatedFilters;
		});
	};

	useEffect(() => {
		const currentFolder = path[path.length - 1];

		// Convert the filters array to an object for the API
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);

		// Add sorting to the filters object
		if (sortConfig) {
			filtersObject.sort_by = sortConfig.key;
			filtersObject.sort_order = sortConfig.direction;
		}

		// Debounce the fetch call
		// TOOD: do this using useDebounce
		const debounceTimer = setTimeout(() => {
			fetchContents(currentFolder.id, filtersObject, pagination);
		}, 750);

		return () => clearTimeout(debounceTimer);
	}, [path, activeFilters, sortConfig, pagination, fetchContents]);

	/**
	 * Refreshes file contents using the current folder and active filters
	 * Uses the fetchContents hook from useContentStore
	 */
	const refreshContents = () => {
		const currentFolder = path[path.length - 1];
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);
		fetchContents(currentFolder.id, filtersObject, pagination);
	};

	/**
	 * Handles file uploads
	 * @param {File[]} files - Array of files to upload
	 */
	const handleUpload = async (files: File[]) => {
		const currentFolderId = path[path.length - 1]?.id || null;
		uploadFiles(files, currentFolderId)
		console.log("Uploading files from drop:", files);
	};

	const { uploadFiles } = useFileUploader({
		onUploadComplete: refreshContents,
	})

	/**
	 * Handles sorting by a specific key
	 * @param {string} key - Column key to sort by
	 */
	const handleSort = (key: string) => {
		setLoading(true)
		setSortConfig(prevConfig => {
			// If clicking the same key, toggle direction; otherwise, reset to 'asc'
			const direction = prevConfig.key === key && prevConfig.direction === 'asc' ? 'desc' : 'asc';
			return { key, direction };
		});
	};

	const { isDragging } = useDragAndDrop({ onDrop: (files) => handleUpload(files) });

	return (
		<>
			<DropzoneOverlay isDragging={isDragging} />

			<div className="flex flex-col items-center">
				<div className="flex flex-col items-center my-6">
					<h1 className="text-3xl font-bold"> Dashboard </h1>
					<p> View and Manage your files here.</p>
				</div>
				<div>
					<SearchAndFilterComponent
						activeFilters={activeFilters}
						onFilterChange={handleFilterChange}
						onSizeFilterChange={handleSizeFilterChange}
					/>
				</div>
				<div className="m-4">
					<FileUploadMenu
						onActionComplete={refreshContents}
						currentFolderID={currentFolderId}
					/>
				</div>

				{isLoading ?
					<FilesSkeleton />
					:
					<div className="flex flex-col p-0 m-0">
						<Breadcrumbs />
						<Card className="rounded-2xl border shadow-sm justify-around overflow-hidden max-w-7xl min-w-6xl mt-4 pt-1 pb-1">
							<FilesTable
								contents={contents}
								onDataChange={refreshContents}
								sortConfig={sortConfig}
								onSort={handleSort}
							/>
							<DataTablePagination
								pageIndex={pagination.pageIndex}
								pageSize={pagination.pageSize}
								setPageSize={(size) => setPagination({ pageIndex: 0, pageSize: size })}
								setPageIndex={(page) => setPagination(p => ({ ...p, pageIndex: page }))}
								totalCount={totalCount}
							/>
						</Card>
					</div>
				}
			</div>
		</>

	)
}
export default DashboardPage;
