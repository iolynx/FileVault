"use client";

import FilesTable from "@/components/FilesTable";
import { Card } from "@/components/ui/card";
import api from "@/lib/axios";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { FileUploadMenu } from "@/components/FileUploadMenu";
import { Skeleton } from "@/components/ui/skeleton";
import { useRouter } from "next/navigation";
import { ActiveFilter, FilterOption } from "@/types/Filter";
import { SearchAndFilterComponent } from "@/components/SearchAndFilter";
import { useContentStore } from "@/stores/useContentStore";
import { Breadcrumbs } from "@/components/Breadcrumbs";
import { DropzoneOverlay } from "@/components/DropzoneOverlay";
import { useDragAndDrop } from "@/hooks/useDragAndDrop";
import { useFileUploader } from "@/hooks/useFileUploader";

const filterOptions: FilterOption[] = [
	{ value: 'content_type', label: 'MIME Type' },
	{ value: 'user_owns_file', label: 'Location' },
	{ value: 'uploaded_before', label: 'Before:' },
	{ value: 'uploaded_after', label: 'After:' }
];

const DashboardPage = () => {
	const router = useRouter();
	const { contents, path, isLoading, fetchContents } = useContentStore();
	const [activeFilters, setActiveFilters] = useState<ActiveFilter[]>([]);
	const currentFolder = path[path.length - 1];
	const currentFolderId = currentFolder ? currentFolder.id : null;

	const handleFilterChange = (column: string, value: string | Date | undefined) => {
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

	useEffect(() => {
		const currentFolder = path[path.length - 1];

		// Convert the filters array to an object for the API.
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);

		// Debounce the fetch call
		const debounceTimer = setTimeout(() => {
			fetchContents(currentFolder.id, filtersObject);
		}, 250);

		return () => clearTimeout(debounceTimer);
	}, [path, activeFilters, fetchContents]);

	const refreshContents = () => {
		const currentFolder = path[path.length - 1];
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);
		fetchContents(currentFolder.id, filtersObject);
	};

	const handleUpload = async (files: File[]) => {
		const currentFolderId = path[path.length - 1]?.id || null;
		uploadFiles(files, currentFolderId)
		console.log("Uploading files from drop:", files);
	};

	const { uploadFiles } = useFileUploader({
		onUploadComplete: refreshContents,
	})

	const { isDragging } = useDragAndDrop({ onDrop: (files) => handleUpload(files) });

	return (
		<>
			<DropzoneOverlay isDragging={isDragging} />

			<div className="flex flex-col items-center">
				<div className="flex flex-col items-center my-10">
					<h1 className="text-3xl font-bold"> Dashboard </h1>
					<p> View and Manage your files here.</p>
				</div>
				<div>
					<SearchAndFilterComponent
						activeFilters={activeFilters}
						onFilterChange={handleFilterChange}
					/>
				</div>
				<div className="m-4">
					<FileUploadMenu onActionComplete={refreshContents} currentFolderID={currentFolderId} />
				</div>

				{isLoading ?
					<Skeleton className="h-[500px] w-[86%] rounded-2xl max-w-7xl mt-4" /> :
					<>
						<Breadcrumbs />
						<Card className="rounded-2xl border shadow-sm overflow-hidden w-full max-w-7xl mt-4 pt-1 pb-1">
							<FilesTable contents={contents} onDataChange={refreshContents} />
						</Card>
					</>
				}
			</div>
		</>

	)
}
export default DashboardPage;
