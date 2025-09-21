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
import { MultiSelectOption } from "@/components/multi-select";
import { mapUsersToOptions } from "@/lib/utils";
import { User } from "@/types/User";

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
	const [shareDialogOptions, setShareDialogOptions] = useState<MultiSelectOption[]>([]);
	const [loading, setLoading] = useState(false)

	const fetchUsers = async () => {
		try {
			const res = await api.get("/users",
				{ withCredentials: true },
			)
			const users: User[] = await res.data
			setShareDialogOptions(mapUsersToOptions(users));
		} catch (error) {
			console.log("error while fetching users:", error)
		}
	}

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

	const handleSizeFilterChange = (sizeFilter: { min_size: number | null, max_size: number | null }) => {
		setActiveFilters(prevFilters => {
			// Remove any old size filters
			let updatedFilters = prevFilters.filter(
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

		// Convert the filters array to an object for the API.
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);

		// Debounce the fetch call
		const debounceTimer = setTimeout(() => {
			fetchContents(currentFolder.id, filtersObject);
			fetchUsers();
		}, 750);

		return () => clearTimeout(debounceTimer);
	}, [path, activeFilters, fetchContents]);

	const refreshContents = () => {
		const currentFolder = path[path.length - 1];
		const filtersObject = activeFilters.reduce((acc, filter) => {
			acc[filter.column] = filter.value;
			return acc;
		}, {} as Record<string, string>);
		fetchContents(currentFolder.id, filtersObject);
		fetchUsers();
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
						onSizeFilterChange={handleSizeFilterChange}
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
							<FilesTable contents={contents} onDataChange={refreshContents} shareDialogOptions={shareDialogOptions} />
						</Card>
					</>
				}
			</div>
		</>

	)
}
export default DashboardPage;
